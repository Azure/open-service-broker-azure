// +build !unit

package lifecycle

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

// serviceLifecycleTestCase encapsulates all the required things for a lifecycle
// test case. A case should defines both createDependency and
// cleanUpDependency, or neither of them. And we assume that the dependency is
// in the same resource group with the service instance.
type serviceLifecycleTestCase struct {
	module                 service.Module
	description            string
	setup                  func() error
	serviceID              string
	planID                 string
	location               string
	provisioningParameters service.ProvisioningParameters
	bindingParameters      service.BindingParameters
	testCredentials        func(credentials service.Credentials) error
}

func (s serviceLifecycleTestCase) getName() string {
	base := fmt.Sprintf(
		"TestServices/lifecycle/%s",
		s.module.GetName(),
	)
	if s.description == "" {
		return base
	}
	return fmt.Sprintf(
		"%s/%s",
		base,
		strings.Replace(s.description, " ", "_", -1),
	)
}

func (s serviceLifecycleTestCase) execute(resourceGroup string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*20)
	defer cancel()

	name := s.getName()

	log.Printf("----> %s: starting\n", name)

	defer log.Printf("----> %s: completed\n", name)

	// This will periodically send status to stdout until the context is canceled.
	// THIS is what stops CI from timing out these tests!
	go s.showStatus(ctx)

	// Get the service and plan
	cat, err := s.module.GetCatalog()
	if err != nil {
		return fmt.Errorf(
			`error gettting catalog from module "%s"`,
			s.module.GetName(),
		)
	}
	svc, ok := cat.GetService(s.serviceID)
	if !ok {
		return fmt.Errorf(
			`service "%s" not found in module "%s" catalog`,
			s.serviceID,
			s.module.GetName(),
		)
	}
	plan, ok := svc.GetPlan(s.planID)
	if !ok {
		return fmt.Errorf(
			`plan "%s" not found for service "%s" in module "%s" catalog`,
			s.planID,
			s.serviceID,
			s.module.GetName(),
		)
	}
	serviceManager := svc.GetServiceManager()

	err = serviceManager.ValidateProvisioningParameters(s.provisioningParameters)
	if err != nil {
		return err
	}

	// Setup...
	if s.setup != nil {
		if err := s.setup(); err != nil {
			return err
		}
	}

	// Build an instance from test case details
	instance := service.Instance{
		ServiceID: s.serviceID,
		PlanID:    s.planID,
		Location:  s.location,
		// Force the resource group to be something known to this test executor
		// to ensure good cleanup
		ResourceGroup:          resourceGroup,
		ProvisioningContext:    serviceManager.GetEmptyProvisioningContext(),
		ProvisioningParameters: s.provisioningParameters,
	}

	// Provision...
	provisioner, err := serviceManager.GetProvisioner(plan)
	if err != nil {
		return err
	}
	stepName, ok := provisioner.GetFirstStepName()
	// There MUST be a first step
	if !ok {
		return fmt.Errorf(
			`Module "%s" provisioner has no steps`,
			s.module.GetName(),
		)
	}
	// Execute provisioning steps until there are none left
	for {
		var step service.ProvisioningStep
		step, ok = provisioner.GetStep(stepName)
		if !ok {
			return fmt.Errorf(
				`Module "%s" provisioning step "%s" not found`,
				s.module.GetName(),
				stepName,
			)
		}
		instance.ProvisioningContext, err = step.Execute(ctx, instance, plan)
		if err != nil {
			return err
		}
		stepName, ok = provisioner.GetNextStepName(stepName)
		// If there is no next step, we're done with provisioning
		if !ok {
			break
		}
	}

	// Bind
	bc, credentials, err := serviceManager.Bind(instance, s.bindingParameters)
	if err != nil {
		return err
	}

	// Test the credentials
	if s.testCredentials != nil {
		err = s.testCredentials(credentials)
		if err != nil {
			return err
		}
	}

	// Unbind
	err = serviceManager.Unbind(instance, bc)
	if err != nil {
		return err
	}

	// Deprovision...
	deprovisioner, err := serviceManager.GetDeprovisioner(plan)
	if err != nil {
		return nil
	}
	stepName, ok = deprovisioner.GetFirstStepName()
	// There MUST be a first step
	if !ok {
		return fmt.Errorf(
			`Module "%s" deprovisioner has no steps`,
			s.module.GetName(),
		)
	}
	// Execute deprovisioning steps until there are none left
	for {
		step, ok := deprovisioner.GetStep(stepName)
		if !ok {
			return fmt.Errorf(
				`Module "%s" deprovisioning step "%s" not found`,
				s.module.GetName(),
				stepName,
			)
		}
		// Assign results to temp variable in case they're nil. We don't want
		// pc to ever be nil, or we risk a nil pointer dereference in the
		// cleanup logic.
		instance.ProvisioningContext, err = step.Execute(ctx, instance, plan)
		if err != nil {
			return err
		}
		stepName, ok = deprovisioner.GetNextStepName(stepName)
		// If there is no next step, we're done with deprovisioning
		if !ok {
			break
		}
	}

	return nil
}

func (s serviceLifecycleTestCase) showStatus(ctx context.Context) {
	name := s.getName()
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Printf("----> %s: in progress\n", name)
		case <-ctx.Done():
			return
		}
	}
}
