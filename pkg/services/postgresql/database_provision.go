package postgresql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *databaseManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("deployARMTemplate", d.deployARMTemplate),
		service.NewProvisioningStep("setupDatabase", d.setupDatabase),
		service.NewProvisioningStep("createExtensions", d.createExtensions),
	)
}

func (d *databaseManager) preProvision(
	context.Context,
	service.Instance,
) (service.InstanceDetails, error) {
	return &databaseInstanceDetails{
		ARMDeploymentName: uuid.NewV4().String(),
		DatabaseName:      generate.NewIdentifier(),
	}, nil
}

func (d *databaseManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	dt := instance.Details.(*databaseInstanceDetails)
	armTemplateParameters := map[string]interface{}{
		"serverName":   pdt.ServerName,
		"databaseName": dt.DatabaseName,
	}
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	_, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.Parent.ProvisioningParameters.GetString("resourceGroup"),
		instance.Parent.ProvisioningParameters.GetString("location"),
		databaseARMTemplateBytes,
		nil, // Go template params
		armTemplateParameters,
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

func (d *databaseManager) setupDatabase(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	dt := instance.Details.(*databaseInstanceDetails)
	err := setupDatabase(
		isSSLRequired(*instance.Parent.ProvisioningParameters),
		pdt.AdministratorLogin,
		pdt.ServerName,
		string(pdt.AdministratorLoginPassword),
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
	if err != nil {
		return nil, err
	}
	return instance.Details, nil
}

func (d *databaseManager) createExtensions(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	dt := instance.Details.(*databaseInstanceDetails)
	extensions := instance.ProvisioningParameters.GetStringArray("extensions")
	if len(extensions) > 0 {
		err := createExtensions(
			isSSLRequired(*instance.Parent.ProvisioningParameters),
			pdt.AdministratorLogin,
			pdt.ServerName,
			string(pdt.AdministratorLoginPassword),
			pdt.FullyQualifiedDomainName,
			dt.DatabaseName,
			extensions,
		)
		if err != nil {
			return nil, err
		}
	}
	return instance.Details, nil
}
