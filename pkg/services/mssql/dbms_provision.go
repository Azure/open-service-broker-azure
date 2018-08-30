package mssql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *dbmsManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("deployARMTemplate", d.deployARMTemplate),
	)
}

func (d *dbmsManager) preProvision(
	context.Context,
	service.Instance,
) (service.InstanceDetails, error) {
	return &dbmsInstanceDetails{
		ARMDeploymentName:          uuid.NewV4().String(),
		ServerName:                 uuid.NewV4().String(),
		AdministratorLogin:         generate.NewIdentifier(),
		AdministratorLoginPassword: service.SecureString(generate.NewPassword()),
	}, nil
}

func (d *dbmsManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsInstanceDetails)
	version := instance.Service.GetProperties().Extended["version"].(string)
	goTemplateParams := buildDBMSGoTemplateParameters(
		dt,
		*instance.ProvisioningParameters,
		version,
	)
	goTemplateParams["location"] =
		instance.ProvisioningParameters.GetString("location")
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	outputs, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		instance.ProvisioningParameters.GetString("location"),
		dbmsARMTemplateBytes,
		goTemplateParams,
		map[string]interface{}{},
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	var ok bool
	dt.FullyQualifiedDomainName, ok = outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving fully qualified domain name from deployment: %s",
			err,
		)
	}
	return dt, err
}
