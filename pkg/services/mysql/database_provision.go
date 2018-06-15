package mysql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	uuid "github.com/satori/go.uuid"
)

func (d *databaseManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("deployARMTemplate", d.deployARMTemplate),
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
	armTemplateParameters := map[string]interface{}{ // ARM template params
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
