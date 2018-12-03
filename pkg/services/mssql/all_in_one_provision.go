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
		service.NewProvisioningStep("setConnectionPolicy", a.setConnectionPolicy),
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
	pp := instance.ProvisioningParameters
	version := instance.Service.GetProperties().Extended["version"].(string)
	goTemplateParams, err := buildDBMSGoTemplateParameters(
		&dt.dbmsInstanceDetails,
		*pp,
		version,
	)
	if err != nil {
		return nil, err
	}
	pd := instance.Plan.GetProperties().Extended["tierDetails"].(planDetails)
	dbParams, err := buildDatabaseGoTemplateParameters(
		dt.DatabaseName,
		*pp,
		pd,
	)
	if err != nil {
		return nil, err
	}
	for key, value := range dbParams {
		goTemplateParams[key] = value
	}
	goTemplateParams["location"] = pp.GetString("location")
	tagsObj := pp.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	outputs, err := a.armDeployer.Deploy(
		dt.ARMDeploymentName,
		pp.GetString("resourceGroup"),
		pp.GetString("location"),
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

func (a *allInOneManager) setConnectionPolicy(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*allInOneInstanceDetails)
	pp := instance.ProvisioningParameters
	connectionPolicy := pp.GetString("connectionPolicy")
	var err error
	if connectionPolicy != "" {
		err = setConnectionPolicy(
			ctx,
			&a.serverConnectionPoliciesClient,
			pp.GetString("resourceGroup"),
			dt.ServerName,
			pp.GetString("location"),
			connectionPolicy,
		)
	}
	return instance.Details, err
}
