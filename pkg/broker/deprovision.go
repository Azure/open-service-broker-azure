package broker

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/async/model"
	"github.com/Azure/azure-service-broker/pkg/service"
	log "github.com/Sirupsen/logrus"
)

func (b *broker) doDeprovisionStep(
	ctx context.Context,
	args map[string]string,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	stepName, ok := args["stepName"]
	if !ok {
		return errors.New(`missing required argument "stepName"`)
	}
	instanceID, ok := args["instanceID"]
	if !ok {
		return errors.New(`missing required argument "instanceID"`)
	}
	instance, ok, err := b.store.GetInstance(instanceID)
	if err != nil {
		return b.handleDeprovisioningError(
			instanceID,
			stepName,
			err,
			"error loading persisted instance",
		)
	}
	if !ok {
		return b.handleDeprovisioningError(
			instanceID,
			stepName,
			nil,
			"instance does not exist in the data store",
		)
	}
	log.WithFields(log.Fields{
		"step":       stepName,
		"instanceID": instance.InstanceID,
	}).Debug("executing deprovisioning step")
	svc, ok := b.catalog.GetService(instance.ServiceID)
	if !ok {
		return b.handleDeprovisioningError(
			instance,
			stepName,
			nil,
			fmt.Sprintf(
				`no service was found for handling serviceID "%s"`,
				instance.ServiceID,
			),
		)
	}
	plan, ok := svc.GetPlan(instance.PlanID)
	if !ok {
		return b.handleDeprovisioningError(
			instance,
			stepName,
			nil,
			fmt.Sprintf(
				`no plan was found for handling planID "%s"`,
				instance.ServiceID,
			),
		)
	}
	serviceManager := svc.GetServiceManager()
	provisioningContext := serviceManager.GetEmptyProvisioningContext()
	err = instance.GetProvisioningContext(provisioningContext, b.codec)
	if err != nil {
		return b.handleDeprovisioningError(
			instance,
			stepName,
			err,
			"error decoding provisioningContext from persisted instance",
		)
	}
	deprovisioner, err := serviceManager.GetDeprovisioner(plan)
	if err != nil {
		return b.handleDeprovisioningError(
			instance,
			stepName,
			err,
			fmt.Sprintf(
				`error retrieving deprovisioner for service "%s"`,
				instance.ServiceID,
			),
		)
	}
	step, ok := deprovisioner.GetStep(stepName)
	if !ok {
		return b.handleDeprovisioningError(
			instance,
			stepName,
			nil,
			`deprovisioner does not know how to process step "%s"`,
		)
	}
	updatedProvisioningContext, err := step.Execute(
		ctx,
		instanceID,
		plan,
		instance.StandardProvisioningContext,
		provisioningContext,
	)
	if err != nil {
		return b.handleDeprovisioningError(
			instance,
			stepName,
			err,
			"error executing deprovisioning step",
		)
	}
	err = instance.SetProvisioningContext(updatedProvisioningContext, b.codec)
	if err != nil {
		return b.handleDeprovisioningError(
			instance,
			stepName,
			err,
			"error encoding modified provisioningContext",
		)
	}
	if nextStepName, ok := deprovisioner.GetNextStepName(step.GetName()); ok {
		if err = b.store.WriteInstance(instance); err != nil {
			return b.handleDeprovisioningError(
				instance,
				stepName,
				err,
				"error persisting instance",
			)
		}
		task := model.NewTask(
			"deprovisionStep",
			map[string]string{
				"stepName":   nextStepName,
				"instanceID": instanceID,
			},
		)
		if err = b.asyncEngine.SubmitTask(task); err != nil {
			return b.handleDeprovisioningError(
				instance,
				stepName,
				err,
				fmt.Sprintf(`error enqueing next step: "%s"`, nextStepName),
			)
		}
	} else {
		// No next step-- we're done deprovisioning!
		_, err = b.store.DeleteInstance(instance.InstanceID)
		if err != nil {
			return b.handleDeprovisioningError(
				instance,
				stepName,
				err,
				"error deleting deprovisioned instance",
			)
		}
	}
	return nil
}

// handleDeprovisioningError tries to handle async deprovisioning errors. If an
// instance is passed in, its status is updated and an attempt is made to
// persist the instance with updated status. If this fails, we have a very
// serious problem on our hands, so we log that failure and kill the process.
// Barring such a failure, a nicely formatted error is returned to be, in-turn,
// returned by the caller of this function. If an instanceID is passed in
// (instead of an instance), only error formatting is handled.
func (b *broker) handleDeprovisioningError(
	instanceOrInstanceID interface{},
	stepName string,
	e error,
	msg string,
) error {
	instance, ok := instanceOrInstanceID.(*service.Instance)
	if !ok {
		instanceID := instanceOrInstanceID
		if e == nil {
			return fmt.Errorf(
				`error executing deprovisioning step "%s" for instance "%s": %s`,
				stepName,
				instanceID,
				msg,
			)
		}
		return fmt.Errorf(
			`error executing deprovisioning step "%s" for instance "%s": %s: %s`,
			stepName,
			instanceID,
			msg,
			e,
		)
	}
	// If we get to here, we have an instance (not just and instanceID)
	instance.Status = service.InstanceStateDeprovisioningFailed
	var ret error
	if e == nil {
		ret = fmt.Errorf(
			`error executing deprovisioning step "%s" for instance "%s": %s`,
			stepName,
			instance.InstanceID,
			msg,
		)
	} else {
		ret = fmt.Errorf(
			`error executing deprovisioning step "%s" for instance "%s": %s: %s`,
			stepName,
			instance.InstanceID,
			msg,
			e,
		)
	}
	instance.StatusReason = ret.Error()
	if err := b.store.WriteInstance(instance); err != nil {
		log.WithFields(log.Fields{
			"instanceID":       instance.InstanceID,
			"status":           instance.Status,
			"originalError":    ret,
			"persistenceError": err,
		}).Fatal("error persisting instance with updated status")
	}
	return ret
}
