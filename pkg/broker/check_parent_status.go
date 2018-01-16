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

func (b *broker) doCheckParentStatus(
	ctx context.Context,
	args map[string]string,
) error {

	instanceID, ok := args["instanceID"]
	if !ok {
		return errors.New(`missing required argument "instanceID"`)
	}

	instance, ok, err := b.store.GetInstance(instanceID)
	if !ok {
		return b.handleProvisioningError(
			instanceID,
			"checkParentStatus",
			nil,
			"error loading persisted instance",
		)
	}
	if err != nil {
		return b.handleProvisioningError(
			instanceID,
			"checkParentStatus",
			err,
			"error loading persisted instance",
		)
	}
	waitForParent, err := b.waitForParent(instance)
	if err != nil {
		return b.handleProvisioningError(
			instance,
			"checkParentStatus",
			err,
			fmt.Sprintf("error: parent status invalid"),
		)
	}
	var task model.Task
	if waitForParent {
		task = model.NewDelayedTask(
			"checkParentStatus",
			map[string]string{
				"instanceID": instanceID,
			},
			time.Minute*1,
		)
		log.WithFields(log.Fields{
			"instanceID": instanceID,
		}).Debug("parent not done, will wait again")
	} else {
		svc, ok := b.catalog.GetService(instance.ServiceID)
		if !ok {
			// If we don't find the Service in the catalog, something is really wrong.
			// (It should exist, because an instance with this serviceID exists.)
			log.WithFields(log.Fields{
				"instanceID": instanceID,
				"serviceID":  instance.ServiceID,
			}).Error(
				"provisioning error: no Service found for serviceID",
			)
			return b.handleProvisioningError(
				instanceID,
				"checkParentStatus",
				err,
				"error no service found for serviceID",
			)
		}
		plan, ok := svc.GetPlan(instance.PlanID)
		if !ok {
			// If we don't find the Service in the catalog, something is really wrong.
			// (It should exist, because an instance with this serviceID exists.)
			log.WithFields(log.Fields{
				"instanceID": instanceID,
				"serviceID":  instance.ServiceID,
				"planID":     instance.PlanID,
			}).Error(
				"provisioning error: no Plan found for planID in Service",
			)
			return b.handleProvisioningError(
				instanceID,
				"checkParentStatus",
				err,
				"error no plan found for planID",
			)
		}
		serviceManager := svc.GetServiceManager()
		var provisioner service.Provisioner
		provisioner, err = serviceManager.GetProvisioner(plan)
		if err != nil {
			log.WithFields(log.Fields{
				"instanceID": instanceID,
				"serviceID":  instance.ServiceID,
				"planID":     instance.PlanID,
				"error":      err,
			}).Error(
				"provisioning error: error retrieving provisioner for " +
					"service and plan",
			)
			return b.handleProvisioningError(
				instanceID,
				"checkParentStatus",
				err,
				"error retrieving provisioner for service and plan",
			)
		}
		provisionFirstStep, ok := provisioner.GetFirstStepName()
		if !ok {
			log.WithFields(log.Fields{
				"instanceID": instanceID,
				"serviceID":  instance.ServiceID,
				"planID":     instance.PlanID,
			}).Error(
				"provisioning error: no steps found for provisioning " +
					"service and plan",
			)
			return b.handleProvisioningError(
				instanceID,
				"checkParentStatus",
				err,
				"error: no steps found for provisioning service and plan",
			)
		}
		task = model.NewTask(
			"provisionStep",
			map[string]string{
				"stepName":   provisionFirstStep,
				"instanceID": instanceID,
			},
		)
		log.WithFields(log.Fields{
			"step":       "checkParentStatus",
			"instanceID": instanceID,
		}).Debug("parent done, sending start provision task")
	}
	if err = b.asyncEngine.SubmitTask(task); err != nil {
		return b.handleProvisioningError(
			instance,
			"checkParentStatus",
			err,
			fmt.Sprintf(
				`error submitting task %s from checkParentStatus`,
				task.GetJobName(),
			),
		)
	}
	return nil

}

func (b *broker) waitForParent(instance service.Instance) (bool, error) {
	//Parent has not been submitted yet, so wait for that
	if instance.Parent == nil {
		return true, nil
	}

	//If parent failed, we should not even attempt to provision this
	if instance.Parent.Status == service.InstanceStateProvisioningFailed {
		log.WithFields(log.Fields{
			"error":      "waitforParent",
			"instanceID": instance.InstanceID,
			"parentID":   instance.Parent.InstanceID,
		}).Info(
			"bad provision request: parent failed provisioning",
		)
		return false, fmt.Errorf("error provisioning: parent provision failed")
	}
	//If parent is deprovisioning, we should not even attempt to provision this
	if instance.Parent.Status == service.InstanceStateDeprovisioning {
		log.WithFields(log.Fields{
			"error":      "waitforParent",
			"instanceID": instance.InstanceID,
			"parentID":   instance.Parent.InstanceID,
		}).Info(
			"bad provision request: parent is deprovisioning",
		)
		return false, fmt.Errorf("error provisioning: parent is deprovisioning")
	}
	if instance.Parent.Status == service.InstanceStateProvisioned {
		return false, nil
	}
	return true, nil
}
