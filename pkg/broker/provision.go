package broker

import (
	"context"
	"fmt"
	"log"

	"errors"

	"github.com/Azure/azure-service-broker/pkg/async/model"
	"github.com/Azure/azure-service-broker/pkg/service"
)

func (b *broker) doProvisionStep(ctx context.Context, args map[string]string) error {
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
		return b.handleProvisioningError(
			instanceID,
			stepName,
			fmt.Sprintf("error loading persisted instance: %s", err),
		)
	}
	if !ok {
		return b.handleProvisioningError(
			instanceID,
			stepName,
			"instance does not exist in the data store",
		)
	}
	log.Printf(
		`executing provisioning step "%s" for instance "%s"`,
		stepName,
		instance.InstanceID,
	)
	module, ok := b.modules[instance.ServiceID]
	if !ok {
		return b.handleProvisioningError(
			instanceID,
			stepName,
			fmt.Sprintf(
				`no module was found for handling service "%s"`,
				instance.ServiceID,
			),
		)
	}
	provisioningContext := module.GetEmptyProvisioningContext()
	err = instance.GetProvisioningContext(provisioningContext)
	if err != nil {
		return b.handleProvisioningError(
			instanceID,
			stepName,
			"error decoding provisioningContext from persisted instance",
		)
	}
	provisioningParams := module.GetEmptyProvisioningParameters()
	if err := instance.GetProvisioningParameters(provisioningParams); err != nil {
		return b.handleProvisioningError(
			instanceID,
			stepName,
			"error decoding provisioningParameters from persisted instance",
		)
	}
	provisioner, ok := b.provisioners[instance.ServiceID]
	if !ok {
		return b.handleProvisioningError(
			instanceID,
			stepName,
			fmt.Sprintf(
				`no provisioner was found for handling service "%s"`,
				instance.ServiceID,
			),
		)
	}
	step, ok := provisioner.GetStep(stepName)
	if !ok {
		return b.handleProvisioningError(
			instanceID,
			stepName,
			`provisioner does not know how to process step "%s"`,
		)
	}
	updatedProvisioningContext, err := step.Execute(
		ctx,
		provisioningContext,
		provisioningParams,
	)
	if err != nil {
		return b.handleProvisioningError(
			instanceID,
			stepName,
			fmt.Sprintf("error executing provisioning step: %s", err),
		)
	}
	err = instance.SetProvisioningContext(updatedProvisioningContext)
	if err != nil {
		return b.handleProvisioningError(
			instanceID,
			stepName,
			fmt.Sprintf("error encoding modified provisioningContext: %s", err),
		)
	}
	if nextStepName, ok := provisioner.GetNextStepName(step.GetName()); ok {
		err = b.store.WriteInstance(instance)
		if err != nil {
			return b.handleProvisioningError(
				instanceID,
				stepName,
				fmt.Sprintf("error persisting instance: %s", err),
			)
		}
		task := model.NewTask(
			"provisionStep",
			map[string]string{
				"stepName":   nextStepName,
				"instanceID": instanceID,
			},
		)
		if err := b.asyncEngine.SubmitTask(task); err != nil {
			return b.handleProvisioningError(
				instanceID,
				stepName,
				fmt.Sprintf(`error enqueing next step "%s"`, nextStepName),
			)
		}
	} else {
		// No next step-- we're done provisioning!
		instance.Status = service.InstanceStateProvisioned
		err = b.store.WriteInstance(instance)
		if err != nil {
			return b.handleProvisioningError(
				instanceID,
				stepName,
				fmt.Sprintf("error persisting instance: %s", err),
			)
		}
	}
	return nil
}

func (b *broker) handleProvisioningError(
	instanceOrInstanceID interface{},
	stepName string,
	msg string,
) error {
	instance, ok := instanceOrInstanceID.(*service.Instance)
	if ok {
		instance.Status = service.InstanceStateProvisioningFailed
		instance.StatusReason = fmt.Sprintf(
			`error executing provisioning step "%s" for instance "%s": %s`,
			stepName,
			instance.InstanceID,
			msg,
		)
		err := b.store.WriteInstance(instance)
		if err != nil {
			log.Fatalf(
				"error persisting instance %s with updated status: %s",
				instance.InstanceID,
				err,
			)
		}
		return errors.New(instance.StatusReason)
	}
	instanceID := instanceOrInstanceID
	return fmt.Errorf(
		`error executing provisioning step "%s" for instance "%s": %s`,
		stepName,
		instanceID,
		msg,
	)
}
