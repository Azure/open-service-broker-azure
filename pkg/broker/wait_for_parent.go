package broker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/async/model"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
)

func (b *broker) doWaitForParentStep(
	ctx context.Context,
	args map[string]string,
) error {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	provisionFirstStep, ok := args["provisionFirstStep"]
	if !ok {
		return errors.New(`missing required argument "provisionFirstStep"`)
	}

	instanceID, ok := args["instanceID"]
	if !ok {
		return errors.New(`missing required argument "instanceID"`)
	}

	instance, ok, err := b.store.GetInstance(instanceID)
	if !ok {
		return b.handleProvisioningError(
			instanceID,
			"waitForParent",
			nil,
			"error loading persisted instance",
		)
	}
	if err != nil {
		return b.handleProvisioningError(
			instanceID,
			"waitForParent",
			err,
			"error loading persisted instance",
		)
	}
	waitForParent, err := b.waitForParent(instance)
	if err != nil {
		return b.handleProvisioningError(
			instance,
			"waitForParentStep",
			err,
			fmt.Sprintf("error: parent status invalid"),
		)
	}
	var task model.Task
	if waitForParent {
		task = model.NewDelayedTask(
			"waitForParentStep",
			map[string]string{
				"provisionFirstStep": provisionFirstStep,
				"instanceID":         instanceID,
			},
			time.Minute*5,
		)
		log.WithFields(log.Fields{
			"step":       "waitforParent",
			"instanceID": instanceID,
		}).Debug("parent not done, will wait again")
	} else {
		task = model.NewTask(
			"provisionStep",
			map[string]string{
				"stepName":   provisionFirstStep,
				"instanceID": instanceID,
			},
		)
		log.WithFields(log.Fields{
			"step":       "waitforParent",
			"instanceID": instanceID,
		}).Debug("parent done, sending start provision task")
	}
	if err = b.asyncEngine.SubmitTask(task); err != nil {
		return b.handleProvisioningError(
			instance,
			"waitForParentStep",
			err,
			fmt.Sprintf(`error starting next provision task step: "%s"`,
				provisionFirstStep),
		)
	}
	return nil

}

func (b *broker) waitForParent(instance service.Instance) (bool, error) {

	parent, parentFound, err := b.store.GetInstanceByAlias(instance.ParentAlias)

	//Parent has not been submitted yet, so wait for that
	if !parentFound {
		return true, nil
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":       "waitforParent",
			"instanceID":  instance.InstanceID,
			"parentAlias": instance.ParentAlias,
		}).Error(
			"bad provision request: unable to retrieve parent",
		)
		return false, err
	}

	//If parent failed, we should not even attempt to provision this
	if parent.Status == service.InstanceStateProvisioningFailed {
		log.WithFields(log.Fields{
			"error":       "waitforParent",
			"instanceID":  instance.InstanceID,
			"parentAlias": instance.ParentAlias,
		}).Error(
			"bad provision request: parent failed provisioning",
		)
		return false, fmt.Errorf("error provisioning: parent provision failed")
	}

	//If parent is deprovisioning, we should not even attempt to provision this
	if parent.Status == service.InstanceStateDeprovisioning {
		log.WithFields(log.Fields{
			"error":       "waitforParent",
			"instanceID":  instance.InstanceID,
			"parentAlias": instance.ParentAlias,
		}).Error(
			"bad provision request: parent is deprovisioning",
		)
		return false, fmt.Errorf("error provisioning: parent is deprovisioning")
	}

	if parent.Status == service.InstanceStateProvisioned {
		return false, nil
	}

	return true, nil
}
