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
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*keyvaultInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *keyvaultInstanceDetails",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.KeyVaultName = "sb" + uuid.NewV4().String()[:20]
	return dt, instance.SecureDetails, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*keyvaultInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *keyvaultInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*keyvaultSecureInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails as *keyvaultSecureInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*keyvault.ProvisioningParameters",
		)
	}

	outputs, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
		nil, // Go template params
		map[string]interface{}{ // ARM template params
			"keyVaultName": dt.KeyVaultName,
			"vaultSku":     instance.Plan.GetProperties().Extended["vaultSku"],
			"tenantId":     s.tenantID,
			"objectId":     pp.ObjectID,
		},
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	vaultURI, ok := outputs["vaultUri"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving vaultUri from deployment: %s",
			err,
		)
	}
	dt.VaultURI = vaultURI
	dt.ClientID = pp.ClientID
	sdt.ClientSecret = pp.ClientSecret

	return dt, sdt, nil
}
