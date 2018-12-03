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
		service.NewProvisioningStep("setConnectionPolicy", d.setConnectionPolicy),
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
	pp := instance.ProvisioningParameters
	version := instance.Service.GetProperties().Extended["version"].(string)
	goTemplateParams, err := buildDBMSGoTemplateParameters(
		dt,
		*pp,
		version,
	)
	if err != nil {
		return nil, err
	}
	goTemplateParams["location"] =
		pp.GetString("location")
	tagsObj := pp.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	outputs, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		pp.GetString("resourceGroup"),
		pp.GetString("location"),
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

func (d *dbmsManager) setConnectionPolicy(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	connectionPolicy := pp.GetString("connectionPolicy")
	var err error
	if connectionPolicy != "" {
		err = setConnectionPolicy(
			ctx,
			&d.serverConnectionPoliciesClient,
			pp.GetString("resourceGroup"),
			pp.GetString("server"),
			connectionPolicy,
		)
	}
	return instance.Details, err
}
