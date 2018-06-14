package postgresql

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
	ctx context.Context,
	_ service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	serverName, err := getAvailableServerName(
		ctx,
		d.checkNameAvailabilityClient,
	)
	if err != nil {
		return nil, err
	}
	return &dbmsInstanceDetails{
		ARMDeploymentName:          uuid.NewV4().String(),
		ServerName:                 serverName,
		AdministratorLoginPassword: service.SecureString(generate.NewPassword()),
	}, nil
}

func (d *dbmsManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsInstanceDetails)
	version := instance.Service.GetProperties().Extended["version"].(string)
	goTemplateParameters, err := buildGoTemplateParameters(
		instance.Plan,
		version,
		dt,
		*instance.ProvisioningParameters,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to build go template parameters: %s", err)
	}
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
		goTemplateParameters,
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
	return dt, nil
}
