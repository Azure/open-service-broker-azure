package mssql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *databaseFeManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("getDatabase", d.getDatabase),
		service.NewProvisioningStep("deployARMTemplate", d.deployARMTemplate),
	)
}

func (d *databaseFeManager) preProvision(
	_ context.Context,
	_ service.Instance,
) (service.InstanceDetails, error) {
	return &databaseInstanceDetails{
		ARMDeploymentName: uuid.NewV4().String(),
	}, nil
}

func (d *databaseFeManager) getDatabase(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	resourceGroup :=
		instance.Parent.ProvisioningParameters.GetString("resourceGroup")
	databaseName :=
		instance.ProvisioningParameters.GetString("database")
	_, err := d.databasesClient.Get(
		ctx,
		resourceGroup,
		pdt.ServerName,
		databaseName,
		"",
	)
	if err != nil {
		return nil, fmt.Errorf("error getting sql database: %s", err)
	}
	return instance.Details, nil
}

func (d *databaseFeManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databaseInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	goTemplateParams := map[string]interface{}{}
	location := instance.Parent.ProvisioningParameters.GetString("location")
	databaseName := instance.ProvisioningParameters.GetString("database")
	goTemplateParams["serverName"] = pdt.ServerName
	goTemplateParams["location"] = location
	goTemplateParams["databaseName"] = databaseName
	instance.ProvisioningParameters.GetString("database")
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	// No output, so ignore the output
	_, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.Parent.ProvisioningParameters.GetString("resourceGroup"),
		location,
		databaseFeARMTemplateBytes,
		goTemplateParams,
		map[string]interface{}{}, // empty arm params
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	dt.DatabaseName = databaseName
	return instance.Details, nil
}
