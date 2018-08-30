package mssql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (a *allInOneManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", a.preProvision),
		service.NewProvisioningStep("deployARMTemplate", a.deployARMTemplate),
	)
}

func (a *allInOneManager) preProvision(
	context.Context,
	service.Instance,
) (service.InstanceDetails, error) {
	return &allInOneInstanceDetails{
		dbmsInstanceDetails: dbmsInstanceDetails{
			ARMDeploymentName:          uuid.NewV4().String(),
			ServerName:                 uuid.NewV4().String(),
			AdministratorLogin:         generate.NewIdentifier(),
			AdministratorLoginPassword: service.SecureString(generate.NewPassword()),
		},
		DatabaseName: generate.NewIdentifier(),
	}, nil
}

func (a *allInOneManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*allInOneInstanceDetails)
	version := instance.Service.GetProperties().Extended["version"].(string)
	goTemplateParams := buildDBMSGoTemplateParameters(
		&dt.dbmsInstanceDetails,
		*instance.ProvisioningParameters,
		version,
	)
	pd := instance.Plan.GetProperties().Extended["tierDetails"].(planDetails)
	dbParams, err := buildDatabaseGoTemplateParameters(
		dt.DatabaseName,
		*instance.ProvisioningParameters,
		pd,
	)
	if err != nil {
		return nil, err
	}
	for key, value := range dbParams {
		goTemplateParams[key] = value
	}
	goTemplateParams["location"] =
		instance.ProvisioningParameters.GetString("location")
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	outputs, err := a.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		instance.ProvisioningParameters.GetString("location"),
		allInOneARMTemplateBytes,
		goTemplateParams,
		map[string]interface{}{}, // empty arm template params
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
