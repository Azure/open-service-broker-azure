package mssql

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
	)
}

func (d *databaseManager) preProvision(
	context.Context,
	service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := databaseInstanceDetails{
		ARMDeploymentName: uuid.NewV4().String(),
		DatabaseName:      generate.NewIdentifier(),
	}
	dtMap, err := service.GetMapFromStruct(dt)
	return dtMap, nil, err
}

func (d *databaseManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := databaseInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	pdt := dbmsInstanceDetails{}
	if err :=
		service.GetStructFromMap(instance.Parent.Details, &pdt); err != nil {
		return nil, nil, err
	}
	goTemplateParams, err := buildDatabaseGoTemplateParameters(instance)
	if err != nil {
		return nil, nil, err
	}
	goTemplateParams["location"] =
		instance.Parent.ProvisioningParameters.GetString("location")
	goTemplateParams["serverName"] = pdt.ServerName
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	// No output, so ignore the output
	_, err = d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.Parent.ProvisioningParameters.GetString("resourceGroup"),
		instance.Parent.ProvisioningParameters.GetString("location"),
		databaseARMTemplateBytes,
		goTemplateParams,
		map[string]interface{}{}, // empty arm params
		tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, instance.SecureDetails, nil
}
