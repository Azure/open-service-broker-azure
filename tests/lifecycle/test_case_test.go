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
	setup                  func() (*service.Instance, error)
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
	svc, plan, err := s.getServiceAndPlan()
	if err != nil {
		return err
	}

	serviceManager := svc.GetServiceManager()

	err = serviceManager.ValidateProvisioningParameters(s.provisioningParameters)
	if err != nil {
		return err
	}

	// Setup...
	var parent *service.Instance
	if s.setup != nil {
		parent, err = s.setup()
		if err != nil {
			return fmt.Errorf("error running setup %s", err)
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
		Details:                serviceManager.GetEmptyInstanceDetails(),
		ProvisioningParameters: s.provisioningParameters,
		Parent:                 parent,
	}

	if _, err = s.provision(ctx, serviceManager, instance, plan); err != nil {
		return err
	}

	//Only test the binding operations if the service is bindable
	if svc.GetBindable() {
		// Bind
		bd, bErr := serviceManager.Bind(instance, s.bindingParameters)
		if bErr != nil {
			return bErr
		}

		binding := service.Binding{Details: bd}

		credentials, bErr := serviceManager.GetCredentials(instance, binding)
		if bErr != nil {
			return bErr
		}

		// Test the credentials
		if s.testCredentials != nil {
			bErr = s.testCredentials(credentials)
			if bErr != nil {
				return bErr
			}
		}

		// Unbind
		bErr = serviceManager.Unbind(instance, bd)
		if bErr != nil {
			return bErr
		}
	}

	// Deprovision...
	deprovisioner, err := serviceManager.GetDeprovisioner(plan)
	if err != nil {
		return nil
	}
	stepName, ok := deprovisioner.GetFirstStepName()
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
		instance.Details, err = step.Execute(ctx, instance, plan)
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

func (s serviceLifecycleTestCase) getServiceAndPlan() (service.Service, service.Plan, error) {
	// Get the service and plan
	cat, err := s.module.GetCatalog()
	if err != nil {
		return nil, nil, fmt.Errorf(
			`error gettting catalog from module "%s"`,
			s.module.GetName(),
		)
	}

	svc, ok := cat.GetService(s.serviceID)
	if !ok {
		return nil, nil, fmt.Errorf(
			`service "%s" not found in module "%s" catalog`,
			s.serviceID,
			s.module.GetName(),
		)
	}
	plan, ok := svc.GetPlan(s.planID)
	if !ok {
		return nil, nil, fmt.Errorf(
			`plan "%s" not found for service "%s" in module "%s" catalog`,
			s.planID,
			s.serviceID,
			s.module.GetName(),
		)
	}

	return svc, plan, nil
}

func (s serviceLifecycleTestCase) provision(
	ctx context.Context,
	serviceManager service.ServiceManager,
	instance service.Instance,
	plan service.Plan,
) (service.InstanceDetails, error) {
	// Provision...
	provisioner, err := serviceManager.GetProvisioner(plan)
	if err != nil {
		return nil, err
	}
	stepName, ok := provisioner.GetFirstStepName()
	// There MUST be a first step
	if !ok {
		return nil, fmt.Errorf(
			`Module "%s" provisioner has no steps`,
			s.module.GetName(),
		)
	}
	// Execute provisioning steps until there are none left
	for {
		var step service.ProvisioningStep
		step, ok = provisioner.GetStep(stepName)
		if !ok {
			return nil, fmt.Errorf(
				`Module "%s" provisioning step "%s" not found`,
				s.module.GetName(),
				stepName,
			)
		}
		instance.Details, err = step.Execute(ctx, instance, plan)
		if err != nil {
			return nil, err
		}
		stepName, ok = provisioner.GetNextStepName(stepName)
		// If there is no next step, we're done with provisioning
		if !ok {
			break
		}
	}
	return instance.Details, nil
}
