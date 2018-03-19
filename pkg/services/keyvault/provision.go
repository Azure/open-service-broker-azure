package keyvault

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (s *serviceManager) ValidateProvisioningParameters(
	pp service.ProvisioningParameters,
	spp service.SecureProvisioningParameters,
) error {
	kvPP := provisioningParameters{}
	if err := service.GetStructFromMap(pp, &kvPP); err != nil {
		return err
	}
	kvSPP := secureProvisioningParameters{}
	if err := service.GetStructFromMap(spp, &kvSPP); err != nil {
		return err
	}
	if kvPP.ObjectID == "" {
		return service.NewValidationError(
			"objectid",
			fmt.Sprintf(`invalid service principal objectid: "%s"`, kvPP.ObjectID),
		)
	}
	if kvPP.ClientID == "" {
		return service.NewValidationError(
			"clientId",
			fmt.Sprintf(`invalid service principal clientId: "%s"`, kvPP.ClientID),
		)
	}
	if kvSPP.ClientSecret == "" {
		return service.NewValidationError(
			"clientSecret",
			fmt.Sprintf(
				`invalid service principal clientSecret: "%s"`,
				kvSPP.ClientSecret,
			),
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
	context.Context,
	service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := instanceDetails{
		ARMDeploymentName: uuid.NewV4().String(),
		KeyVaultName:      "sb" + uuid.NewV4().String()[:20],
	}
	dtMap, err := service.GetMapFromStruct(dt)
	return dtMap, nil, err
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := instanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	pp := provisioningParameters{}
	if err :=
		service.GetStructFromMap(instance.ProvisioningParameters, &pp); err != nil {
		return nil, nil, err
	}
	spp := secureProvisioningParameters{}
	if err := service.GetStructFromMap(
		instance.SecureProvisioningParameters,
		&spp,
	); err != nil {
		return nil, nil, err
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

	var ok bool
	dt.VaultURI, ok = outputs["vaultUri"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving vaultUri from deployment: %s",
			err,
		)
	}
	dt.ClientID = pp.ClientID

	dtMap, err := service.GetMapFromStruct(dt)
	return dtMap, instance.SecureDetails, err
}
