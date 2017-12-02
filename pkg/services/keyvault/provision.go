package keyvault

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (s *serviceManager) ValidateProvisioningParameters(
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

func (s *serviceManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", s.preProvision),
		service.NewProvisioningStep("deployARMTemplate", s.deployARMTemplate),
	)
}

func (s *serviceManager) preProvision(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
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

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	_ string, // instanceID
	plan service.Plan,
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

	outputs, err := s.armDeployer.Deploy(
		pc.ARMDeploymentName,
		standardProvisioningContext.ResourceGroup,
		standardProvisioningContext.Location,
		armTemplateBytes,
		nil, // Go template params
		map[string]interface{}{ // ARM template params
			"keyVaultName": pc.KeyVaultName,
			"vaultSku":     plan.GetProperties().Extended["vaultSku"],
			"tenantId":     s.keyvaultManager.GetTenantID(),
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
