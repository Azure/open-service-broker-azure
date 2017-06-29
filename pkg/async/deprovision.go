package async

import (
	"errors"
	"fmt"
	"log"

	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/RichardKnop/machinery/v1/tasks"
)

func (e *engine) Deprovision(instance *service.Instance) error {
	deprovisioner, ok := e.deprovisioners[instance.ServiceID]
	// If the deprovisioner wasn't found, something is seriously wrong
	if !ok {
		return fmt.Errorf(
			"no deprovisioner found for service %s",
			instance.ServiceID,
		)
	}
	firstStepName, ok := deprovisioner.GetFirstStepName()
	if !ok {
		return fmt.Errorf(
			"no steps are defined for deprovisioning the service %s",
			instance.ServiceID,
		)
	}
	signature := &tasks.Signature{
		Name: "work",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: instance.InstanceID,
			},
			{
				Type:  "string",
				Value: firstStepName,
			},
		},
	}
	_, err := e.machineryServer.SendTask(signature)
	if err != nil {
		log.Println("error enqueing deprovisioning task")
		return err
	}
	return nil
}

func (e *engine) doDeprovisioningWork(
	instance *service.Instance,
	module service.Module,
	provisioningResult interface{},
	stepName string,
) (interface{}, string, bool, error) {
	deprovisioner, ok := e.deprovisioners[instance.ServiceID]
	if !ok {
		return nil, "", false, errors.New(
			"no deprovisioner was found for handling the service",
		)
	}
	step, ok := deprovisioner.GetStep(stepName)
	if !ok {
		return nil, "", false, errors.New(
			"deprovisioner does not know how to process step",
		)
	}
	updatedProvisioningResult, err := step.Execute(provisioningResult)
	if err != nil {
		return nil, "", false, fmt.Errorf("error executing deprovisioning step: %s", err)
	}
	nextStepName, ok := deprovisioner.GetNextStepName(step.GetName())
	return updatedProvisioningResult, nextStepName, ok, nil
}
