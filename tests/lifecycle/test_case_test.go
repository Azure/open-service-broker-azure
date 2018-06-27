// +build !unit

package lifecycle

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/slice"
	"github.com/Azure/open-service-broker-azure/pkg/types"
	"github.com/stretchr/testify/assert"
)

// serviceLifecycleTestCase encapsulates all the required things for a lifecycle
// test case. A case should defines both createDependency and
// cleanUpDependency, or neither of them. And we assume that the dependency is
// in the same resource group with the service instance.
type serviceLifecycleTestCase struct {
	// To clarify-- this is a test grouping-- it is NOT a resource group
	group                  string
	name                   string
	serviceID              string
	planID                 string
	provisioningParameters map[string]interface{}
	updatingParameters     map[string]interface{}
	parentServiceInstance  *service.Instance
	bindingParameters      map[string]interface{}
	testCredentials        func(credentials map[string]interface{}) error
	childTestCases         []*serviceLifecycleTestCase
}

func (s serviceLifecycleTestCase) getName() string {
	return fmt.Sprintf("TestServices/lifecycle/%s/%s", s.group, s.name)
}

func (s serviceLifecycleTestCase) execute(
	t *testing.T,
	catalog service.Catalog,
	resourceGroup string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*40)
	defer cancel()

	name := s.getName()

	log.Printf("----> %s: starting\n", name)

	defer log.Printf("----> %s: completed\n", name)

	// This will periodically send status to stdout until the context is canceled.
	// THIS is what stops CI from timing out these tests!
	go s.showStatus(ctx)

	// Get the service and plan
	svc, ok := catalog.GetService(s.serviceID)
	if !ok {
		return fmt.Errorf(`service "%s" not found catalog`, s.serviceID)
	}
	plan, ok := svc.GetPlan(s.planID)
	if !ok {
		return fmt.Errorf(
			`plan "%s" not found for service "%s"`,
			s.planID,
			s.serviceID,
		)
	}

	// Force the resource group to be something known to this test executor
	// to ensure good cleanup
	if slice.ContainsString(
		plan.GetSchemas().ServiceInstances.ProvisioningParametersSchema.RequiredProperties, // nolint: lll
		"resourceGroup",
	) {
		s.provisioningParameters["resourceGroup"] = resourceGroup
	}

	if err :=
		plan.GetSchemas().ServiceInstances.ProvisioningParametersSchema.Validate(
			s.provisioningParameters,
		); err != nil {
		return err
	}

	serviceManager := svc.GetServiceManager()

	// Wrap the provisioning parameters with a "params" object that guides access
	// to the parameters using schema
	pp := &service.ProvisioningParameters{
		Parameters: service.Parameters{
			Schema: plan.GetSchemas().ServiceInstances.ProvisioningParametersSchema,
			Data:   s.provisioningParameters,
		},
	}

	// Build an instance from test case details
	instance := service.Instance{
		ServiceID: s.serviceID,
		Service:   svc,
		PlanID:    s.planID,
		Plan:      plan,
		ProvisioningParameters: pp,
		Parent:                 s.parentServiceInstance,
	}

	// Provision...
	provisioner, err := serviceManager.GetProvisioner(plan)
	if err != nil {
		return err
	}
	stepName, ok := provisioner.GetFirstStepName()
	// There MUST be a first step
	if !ok {
		return fmt.Errorf(`Provisioner for service "%s" has no steps`, s.serviceID)
	}
	// Execute provisioning steps until there are none left
	for {
		var step service.ProvisioningStep
		step, ok = provisioner.GetStep(stepName)
		if !ok {
			return fmt.Errorf(
				`Provisioner step "%s" for service "%s" not found`,
				stepName,
				s.serviceID,
			)
		}
		instance.Details, err = step.Execute(ctx, instance)
		if err != nil {
			return err
		}
		stepName, ok = provisioner.GetNextStepName(stepName)
		// If there is no next step, we're done with provisioning
		if !ok {
			break
		}
	}

	// Update...
	if len(s.updatingParameters) != 0 {
		// equivalent to func mergeUpdateParameters in api/update
		if len(s.provisioningParameters) != 0 {
			ppCopy := map[string]interface{}{}
			for key, value := range instance.ProvisioningParameters.Data {
				ppCopy[key] = value
			}
			for key, value := range s.updatingParameters {
				if !types.IsEmpty(value) {
					ppCopy[key] = value
				}
			}
			s.updatingParameters = ppCopy
		}

		// Wrap the updating parameters with a "params" object that guides access
		// to the parameters using schema
		pps := plan.GetSchemas().ServiceInstances.ProvisioningParametersSchema // nolint: lll
		up := &service.ProvisioningParameters{
			Parameters: service.Parameters{
				Schema: &pps,
				Data:   s.updatingParameters,
			},
		}
		instance.UpdatingParameters = up

		var updater service.Updater
		updater, err = serviceManager.GetUpdater(plan)
		if err != nil {
			return err
		}
		stepName, ok = updater.GetFirstStepName()
		// There MUST be a first step
		if !ok {
			return fmt.Errorf(`Updater for service "%s" has no steps`, s.serviceID)
		}
		// Execute updating steps until there are none left
		for {
			var step service.UpdatingStep
			step, ok = updater.GetStep(stepName)
			if !ok {
				return fmt.Errorf(
					`Updater step "%s" for service "%s" not found`,
					stepName,
					s.serviceID,
				)
			}
			instance.Details, err = step.Execute(ctx, instance)
			if err != nil {
				return err
			}
			stepName, ok = updater.GetNextStepName(stepName)
			// If there is no next step, we're done with provisioning
			if !ok {
				break
			}
		}
	}

	//Only test the binding operations if the service is bindable
	if svc.IsBindable() {
		// Wrap the binding parameters with a "params" object that guides access to
		// the parameters using schema
		bps := instance.Plan.GetSchemas().ServiceBindings.BindingParametersSchema
		bp := service.BindingParameters{
			Parameters: service.Parameters{
				Schema: &bps,
				Data:   s.bindingParameters,
			},
		}

		// Bind
		var bd service.BindingDetails
		bd, err = serviceManager.Bind(instance, bp)
		if err != nil {
			return err
		}

		binding := service.Binding{
			Details: bd,
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

	// Iterate through any child test cases, setting the instnace from this
	// test case as the parent.
	for _, childTestCase := range s.childTestCases {
		childTestCase.parentServiceInstance = &instance
		t.Run(childTestCase.getName(), func(t *testing.T) {
			tErr := childTestCase.execute(t, catalog, resourceGroup)
			// This will fail this subtest and also the parent lifecycle test
			assert.Nil(t, tErr)
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
			`DepProvisioner for service "%s" has no steps`,
			s.serviceID,
		)
	}
	// Execute deprovisioning steps until there are none left
	for {
		step, ok := deprovisioner.GetStep(stepName)
		if !ok {
			return fmt.Errorf(
				`Deprovisioner step "%s" for service "%s" not found`,
				stepName,
				s.serviceID,
			)
		}
		instance.Details, err = step.Execute(ctx, instance)
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
