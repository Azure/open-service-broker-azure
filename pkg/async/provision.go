package async

import (
	"errors"
	"fmt"
	"log"

	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/RichardKnop/machinery/v1/tasks"
)

func (e *engine) Provision(instance *service.Instance) error {
	provisioner, ok := e.provisioners[instance.ServiceID]
	// If the provisioner wasn't found, something is seriously wrong
	if !ok {
		return fmt.Errorf("no provisioner found for service %s", instance.ServiceID)
	}
	firstStepName, ok := provisioner.GetFirstStepName()
	if !ok {
		return fmt.Errorf(
			"no steps are defined for provisioning the service %s",
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
		log.Println("error enqueing provisioning task")
		return err
	}
	return nil
}

func (e *engine) doProvisioningWork(
	instance *service.Instance,
	module service.Module,
	provisioningResult interface{},
	stepName string,
) (interface{}, string, bool, error) {
	provisioningParams := module.GetEmptyProvisioningParameters()
	err := instance.GetProvisioningParameters(provisioningParams)
	if err != nil {
		return nil, "", false, errors.New(
			"error decoding provisioning parameters from persisted instance",
		)
	}
	provisioner, ok := e.provisioners[instance.ServiceID]
	if !ok {
		return nil, "", false, errors.New(
			"no provisioner was found for handling the service",
		)
	}
	step, ok := provisioner.GetStep(stepName)
	if !ok {
		return nil, "", false, errors.New("provisioner does not know how to process step")
	}
	updatedProvisioningResult, err := step.Execute(
		provisioningResult,
		provisioningParams,
	)
	if err != nil {
		return nil, "", false, fmt.Errorf("error executing provisioning step: %s", err)
	}
	nextStepName, ok := provisioner.GetNextStepName(step.GetName())
	return updatedProvisioningResult, nextStepName, ok, nil
}
