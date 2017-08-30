// +build !unit

package lifecycle

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-service-broker/pkg/service"
	uuid "github.com/satori/go.uuid"
)

type moduleLifecycleTestCase struct {
	module                 service.Module
	serviceID              string
	planID                 string
	provisioningParameters service.ProvisioningParameters
	bindingParameters      service.BindingParameters
}

func (m *moduleLifecycleTestCase) execute() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*20)
	defer cancel()
	// This will periodically send status to stdout until the context is canceled.
	// THIS is what stops Travis from timing out these tests!
	go m.showStatus(ctx)

	err := m.module.ValidateProvisioningParameters(m.provisioningParameters)
	if err != nil {
		return err
	}

	pc := m.module.GetEmptyProvisioningContext()
	var tempPC service.ProvisioningContext

	// Make sure we clean up after ourselves
	defer func() {
		resourceGroupName := pc.GetResourceGroupName()
		if resourceGroupName != "" {
			log.Printf("----> deleting resource group \"%s\"\n", resourceGroupName)
			if err := deleteResourceGroup(resourceGroupName); err != nil {
				log.Printf("----> error deleting resource group: %s", err)
			} else {
				log.Printf(
					"----> done deleting resource group \"%s\"\n",
					resourceGroupName,
				)
			}
		}
	}()

	// Provision...
	iid := uuid.NewV4().String()
	provisioner, err := m.module.GetProvisioner(m.serviceID, m.planID)
	if err != nil {
		return err
	}
	stepName, ok := provisioner.GetFirstStepName()
	// There MUST be a first step
	if !ok {
		return fmt.Errorf(
			`module "%s" provisioner has no steps`,
			m.module.GetName(),
		)
	}
	// Execute provisioning steps until there are none left
	for {
		var step service.ProvisioningStep
		step, ok = provisioner.GetStep(stepName)
		if !ok {
			return fmt.Errorf(
				`module "%s" provisioning step "%s" not found`,
				m.module.GetName(),
				stepName,
			)
		}
		// Assign results to temp variable in case they're nil. We don't want
		// pc to ever be nil, or we risk a nil pointer dereference in the
		// cleanup logic.
		tempPC, err = step.Execute(
			ctx,
			iid,
			m.serviceID,
			m.planID,
			pc,
			m.provisioningParameters,
		)
		if err != nil {
			return err
		}
		pc = tempPC
		stepName, ok = provisioner.GetNextStepName(stepName)
		// If there is no next step, we're done with provisioning
		if !ok {
			break
		}
	}

	// Bind
	bc, _, err := m.module.Bind(pc, m.bindingParameters)
	if err != nil {
		return err
	}

	// Unbind
	err = m.module.Unbind(pc, bc)
	if err != nil {
		return err
	}

	// Deprovision...
	deprovisioner, err := m.module.GetDeprovisioner(m.serviceID, m.planID)
	if err != nil {
		return nil
	}
	stepName, ok = deprovisioner.GetFirstStepName()
	// There MUST be a first step
	if !ok {
		return fmt.Errorf(
			`module "%s" deprovisioner has no steps`,
			m.module.GetName(),
		)
	}
	// Execute deprovisioning steps until there are none left
	for {
		step, ok := deprovisioner.GetStep(stepName)
		if !ok {
			return fmt.Errorf(
				`module "%s" deprovisioning step "%s" not found`,
				m.module.GetName(),
				stepName,
			)
		}
		// Assign results to temp variable in case they're nil. We don't want
		// pc to ever be nil, or we risk a nil pointer dereference in the
		// cleanup logic.
		tempPC, err = step.Execute(
			ctx,
			iid,
			m.serviceID,
			m.planID,
			pc,
		)
		if err != nil {
			return err
		}
		pc = tempPC
		stepName, ok = deprovisioner.GetNextStepName(stepName)
		// If there is no next step, we're done with deprovisioning
		if !ok {
			break
		}
	}

	return nil
}

func (m *moduleLifecycleTestCase) showStatus(ctx context.Context) {
	moduleName := m.module.GetName()
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Printf(
				"----> module \"%s\" lifecycle tests in progress\n",
				moduleName,
			)
		case <-ctx.Done():
			return
		}
	}
}
