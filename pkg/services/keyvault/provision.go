package keyvault

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

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
) (service.InstanceDetails, error) {
	return &instanceDetails{
		ARMDeploymentName: uuid.NewV4().String(),
		KeyVaultName:      "sb" + uuid.NewV4().String()[:20],
	}, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)
	pp := instance.ProvisioningParameters
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	outputs, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		pp.GetString("resourceGroup"),
		pp.GetString("location"),
		armTemplateBytes,
		map[string]interface{}{
			"location":     pp.GetString("location"),
			"keyVaultName": dt.KeyVaultName,
			"vaultSku":     instance.Plan.GetProperties().Extended["vaultSku"],
			"tenantId":     s.tenantID,
			"objectId":     pp.GetString("objectId"),
		},
		map[string]interface{}{},
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	var ok bool
	dt.VaultURI, ok = outputs["vaultUri"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving vaultUri from deployment: %s",
			err,
		)
	}
	return dt, err
}
