package broker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	"github.com/krancour/async"
)

func (b *broker) doCheckParentStatus(
	_ context.Context,
	task async.Task,
) ([]async.Task, error) {

	instanceID, ok := task.GetArgs()["instanceID"]
	if !ok {
		return nil, errors.New(`missing required argument "instanceID"`)
	}

	instance, ok, err := b.store.GetInstance(instanceID)
	if !ok {
		return nil, b.handleProvisioningError(
			instanceID,
			"checkParentStatus",
			nil,
			"error loading persisted instance",
		)
	}
	if err != nil {
		return nil, b.handleProvisioningError(
			instanceID,
			"checkParentStatus",
			err,
			"error loading persisted instance",
		)
	}
	waitForParent, err := b.waitForParent(instance)
	if err != nil {
		return nil, b.handleProvisioningError(
			instance,
			"checkParentStatus",
			err,
			fmt.Sprintf("error: parent status invalid"),
		)
	}
	if waitForParent {
		log.WithFields(log.Fields{
			"instanceID": instanceID,
		}).Debug("parent not done, will wait again")
		return []async.Task{
			async.NewDelayedTask(
				"checkParentStatus",
				map[string]string{
					"instanceID": instanceID,
				},
				time.Minute*1,
			),
		}, nil
	}
	serviceManager := instance.Service.GetServiceManager()
	var provisioner service.Provisioner
	provisioner, err = serviceManager.GetProvisioner(instance.Plan)
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
		return nil, b.handleProvisioningError(
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
		return nil, b.handleProvisioningError(
			instanceID,
			"checkParentStatus",
			err,
			"error: no steps found for provisioning service and plan",
		)
	}
	log.WithFields(log.Fields{
		"step":       "checkParentStatus",
		"instanceID": instanceID,
	}).Debug("parent done, sending start provision task")

	// Update the status
	instance.Status = service.InstanceStateProvisioning
	if err = b.store.WriteInstance(instance); err != nil {
		return nil, b.handleProvisioningError(
			instance,
			"checkParentStatus",
			err,
			"error: error updating instance status",
		)
	}

	// Put the real provision task into the queue
	return []async.Task{
		async.NewTask(
			"executeProvisioningStep",
			map[string]string{
				"stepName":   provisionFirstStep,
				"instanceID": instanceID,
			},
		),
	}, nil
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
