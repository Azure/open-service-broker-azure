package keyvault

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (m *module) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"*keyvault.ProvisioningParameters",
		)
	}
	if pp.ObjectID == "" {
		return service.NewValidationError(
			"objectid",
			fmt.Sprintf(`invalid service principal objectid: "%s"`, pp.ObjectID),
		)
	}
	if pp.ClientID == "" {
		return service.NewValidationError(
			"clientId",
			fmt.Sprintf(`invalid service principal clientId: "%s"`, pp.ClientID),
		)
	}
	if pp.ClientSecret == "" {
		return service.NewValidationError(
			"clientSecret",
			fmt.Sprintf(`invalid service principal clientSecret: "%s"`, pp.ClientSecret),
		)
	}
	return nil
}

func (m *module) GetProvisioner(string, string) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", m.preProvision),
		service.NewProvisioningStep("deployARMTemplate", m.deployARMTemplate),
	)
}

func (m *module) preProvision(
	_ context.Context,
	_ string, // instanceID
	_ string, // serviceID
	_ string, // planID
	_ service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	_ service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*keyvaultProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *keyvaultProvisioningContext",
		)
	}
	pc.ARMDeploymentName = uuid.NewV4().String()
	pc.KeyVaultName = "sb" + uuid.NewV4().String()[:20]
	return pc, nil
}

func (m *module) deployARMTemplate(
	_ context.Context,
	_ string, // instanceID
	serviceID string,
	planID string,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	provisioningParameters service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*keyvaultProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *keyvaultProvisioningContext",
		)
	}
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters as " +
				"*keyvault.ProvisioningParameters",
		)
	}
	catalog, err := m.GetCatalog()
	if err != nil {
		return nil, fmt.Errorf("error retrieving catalog: %s", err)
	}
	service, ok := catalog.GetService(serviceID)
	if !ok {
		return nil, fmt.Errorf(
			`service "%s" not found in the "%s" module catalog`,
			serviceID,
			m.GetName(),
		)
	}
	plan, ok := service.GetPlan(planID)
	if !ok {
		return nil, fmt.Errorf(
			`plan "%s" not found for service "%s"`,
			planID,
			serviceID,
		)
	}

	outputs, err := m.armDeployer.Deploy(
		pc.ARMDeploymentName,
		standardProvisioningContext.ResourceGroup,
		standardProvisioningContext.Location,
		armTemplateBytes,
		nil, // Go template params
		map[string]interface{}{ // ARM template params
			"keyVaultName": pc.KeyVaultName,
			"vaultSku":     plan.GetProperties().Extended["vaultSku"],
			"tenantId":     m.keyvaultManager.GetTenantID(),
			"objectId":     pp.ObjectID,
		},
		standardProvisioningContext.Tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	vaultURI, ok := outputs["vaultUri"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving vaultUri from deployment: %s",
			err,
		)
	}
	pc.VaultURI = vaultURI
	pc.ClientID = pp.ClientID
	pc.ClientSecret = pp.ClientSecret

	return pc, nil
}
