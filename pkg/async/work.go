package async

import (
	"fmt"
	"log"

	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/RichardKnop/machinery/v1/tasks"
)

func (e *engine) doWork(instanceID, stepName string) error {
	instance, ok, err := e.store.GetInstance(instanceID)
	if err != nil {
		e.handleWorkError(
			instanceID,
			stepName,
			fmt.Sprintf("error loading persisted instance: %s", err),
		)
		return nil
	}
	if !ok {
		e.handleWorkError(
			instanceID,
			stepName,
			"instance does not exist in the data store",
		)
		return nil
	}
	var operation string
	switch instance.Status {
	case service.InstanceStateProvisioning:
		operation = operationProvisioning
	case service.InstanceStateDeprovisioning:
		operation = operationDeprovisioning
	default:
		e.handleWorkError(
			instance,
			stepName,
			fmt.Sprintf(
				`operation cannot be inferred from instance status "%s"`,
				instance.Status,
			),
		)
		return nil
	}
	log.Printf(
		`executing %s step "%s" for instance "%s"`,
		operation,
		stepName,
		instance.InstanceID,
	)
	module, ok := e.modules[instance.ServiceID]
	if !ok {
		e.handleWorkError(
			instance,
			stepName,
			"no module was found for handling the service",
		)
		return nil
	}
	provisioningResult := module.GetEmptyProvisioningResult()
	err = instance.GetProvisioningResult(provisioningResult)
	if err != nil {
		e.handleWorkError(
			instance,
			stepName,
			"error decoding provisioning result from persisted instance",
		)
		return nil
	}
	var nextStepName string
	var nextStepExists bool
	switch instance.Status {
	case service.InstanceStateProvisioning:
		provisioningResult, nextStepName, nextStepExists, err = e.doProvisioningWork(
			instance,
			module,
			provisioningResult,
			stepName,
		)
	case service.InstanceStateDeprovisioning:
		provisioningResult, nextStepName, nextStepExists, err = e.doDeprovisioningWork(
			instance,
			module,
			provisioningResult,
			stepName,
		)
	}
	if err != nil {
		e.handleWorkError(
			instance,
			stepName,
			err.Error(),
		)
	}
	err = instance.SetProvisioningResult(provisioningResult)
	if err != nil {
		e.handleWorkError(
			instance,
			stepName,
			fmt.Sprintf("error encoding modified provisioning result: %s", err),
		)
		return nil
	}
	err = e.store.WriteInstance(instance)
	if err != nil {
		e.handleWorkError(
			instance,
			stepName,
			fmt.Sprintf("error persisting instance: %s", err),
		)
		return nil
	}
	if nextStepExists {
		// Enque the next task
		signature := &tasks.Signature{
			Name: "work",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: instance.InstanceID,
				},
				{
					Type:  "string",
					Value: nextStepName,
				},
			},
		}
		_, err := e.machineryServer.SendTask(signature)
		if err != nil {
			e.handleWorkError(
				instance,
				stepName,
				fmt.Sprintf("error enqueing the next provisioning step: %s", err),
			)
		}
	}
	return nil
}

// handleWorkError handles errors that occur within the async engine. It
// attempts to record the failure appropriately and terminates execution if it
// cannot.
func (e *engine) handleWorkError(
	instanceOrInstanceID interface{},
	stepName string,
	message string,
) {
	instance, ok := instanceOrInstanceID.(*service.Instance)
	if ok {
		var operation string
		if instance.StatusReason == service.InstanceStateProvisioning {
			operation = operationProvisioning
			instance.Status = service.InstanceStateProvisioningFailed
		} else if instance.StatusReason == service.InstanceStateDeprovisioning {
			operation = operationDeprovisioning
			instance.Status = service.InstanceStateDeprovisioningFailed
		} else {
			log.Printf(
				`error executing asynchronous work on instance "%s": %s`,
				instance.InstanceID,
				message,
			)
			return
		}
		instance.StatusReason = fmt.Sprintf(
			`error executing %s step "%s" on instance %#v: %s`,
			operation,
			stepName,
			instance,
			message,
		)
		log.Println(instance.StatusReason)
		err := e.store.WriteInstance(instance)
		if err != nil {
			log.Fatalf(
				"error persisting instance %s with updated status: %s",
				instance.InstanceID,
				err,
			)
		}
	} else {
		log.Printf(
			`error executing asynchronous work on instance "%s": %s`,
			instance,
			message,
		)
	}
}
