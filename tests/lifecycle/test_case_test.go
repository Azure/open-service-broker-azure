// +build !unit

package lifecycle

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/stretchr/testify/assert"
)

// serviceLifecycleTestCase encapsulates all the required things for a lifecycle
// test case. A case should defines both createDependency and
// cleanUpDependency, or neither of them. And we assume that the dependency is
// in the same resource group with the service instance.
type serviceLifecycleTestCase struct {
	module                 service.Module
	description            string
	serviceID              string
	planID                 string
	location               string
	provisioningParameters service.CombinedProvisioningParameters
	parentServiceInstance  *service.Instance
	childTestCases         []*serviceLifecycleTestCase
	bindingParameters      service.CombinedBindingParameters
	testCredentials        func(credentials map[string]interface{}) error
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

func (s serviceLifecycleTestCase) execute(
	t *testing.T,
	resourceGroup string,
) error {
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

	pp, spp, err :=
		serviceManager.SplitProvisioningParameters(s.provisioningParameters)
	if err != nil {
		return err
	}

	if err = serviceManager.ValidateProvisioningParameters(pp, spp); err != nil {
		return err
	}

	// Build an instance from test case details
	instance := service.Instance{
		ServiceID: s.serviceID,
		Service:   svc,
		PlanID:    s.planID,
		Plan:      plan,
		Location:  s.location,
		// Force the resource group to be something known to this test executor
		// to ensure good cleanup
		ResourceGroup:                resourceGroup,
		ProvisioningParameters:       pp,
		SecureProvisioningParameters: spp,
		Parent: s.parentServiceInstance,
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
		instance.Details, instance.SecureDetails, err = step.Execute(ctx, instance)
		if err != nil {
			return err
		}
		stepName, ok = provisioner.GetNextStepName(stepName)
		// If there is no next step, we're done with provisioning
		if !ok {
			break
		}
	}

	//Only test the binding operations if the service is bindable
	if svc.IsBindable() {
		var bp service.BindingParameters
		var sbp service.SecureBindingParameters
		bp, sbp, err = serviceManager.SplitBindingParameters(s.bindingParameters)
		if err != nil {
			return err
		}

		// Bind
		var bd service.BindingDetails
		var sbd service.SecureBindingDetails
		bd, sbd, err = serviceManager.Bind(instance, bp, sbp)
		if err != nil {
			return err
		}

		binding := service.Binding{
			Details:       bd,
			SecureDetails: sbd,
		}

		var credentials service.Credentials
		credentials, err = serviceManager.GetCredentials(instance, binding)
		if err != nil {
			return err
		}

		// Convert the credentials to a map
		var credsMap map[string]interface{}
		credsMap, err = service.GetMapFromStruct(credentials)
		if err != nil {
			return err
		}

		// Test the credentials
		if s.testCredentials != nil {
			if err := s.testCredentials(credsMap); err != nil {
				return err
			}
		}

		// Unbind
		if err = serviceManager.Unbind(instance, binding); err != nil {
			return err
		}
	}

	//Iterate through any child test cases, setting the instnace from this
	//test case as the parent.
	for i, test := range s.childTestCases {
		test.parentServiceInstance = &instance
		var subTestName string
		if test.description == "" {
			subTestName = fmt.Sprintf("subtest-%v", i)
		} else {
			subTestName = strings.Replace(test.description, " ", "_", -1)
		}
		t.Run(subTestName, func(t *testing.T) {
			err = test.execute(t, resourceGroup)
			//This will fail this subtest, but also the parent lifecycle test
			assert.Nil(t, err)
		})
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
		instance.Details, instance.SecureDetails, err = step.Execute(ctx, instance)
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
