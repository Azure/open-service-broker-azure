package postgresql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (a *allInOneManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", a.preProvision),
		service.NewProvisioningStep("deployARMTemplate", a.deployARMTemplate),
		service.NewProvisioningStep("setupDatabase", a.setupDatabase),
		service.NewProvisioningStep("createExtensions", a.createExtensions),
	)
}

func (a *allInOneManager) preProvision(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	dbmsInstanceDetails, err := generateDBMSInstanceDetails(
		ctx,
		instance,
		a.checkNameAvailabilityClient,
	)
	if err != nil {
		return nil, fmt.Errorf("error generating instance detail: %v", err)
	}

	return &allInOneInstanceDetails{
		dbmsInstanceDetails: *dbmsInstanceDetails,
		DatabaseName:        generate.NewIdentifier(),
	}, nil
}

func (a *allInOneManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*allInOneInstanceDetails)
	version := instance.Service.GetProperties().Extended["version"].(string)
	goTemplateParameters, err := buildGoTemplateParameters(
		instance.Plan,
		version,
		&dt.dbmsInstanceDetails,
		*instance.ProvisioningParameters,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error building go template parameters :%s",
			err,
		)
	}
	goTemplateParameters["databaseName"] = dt.DatabaseName
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
	return dt, err
}

func (a *allInOneManager) setupDatabase(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*allInOneInstanceDetails)
	err := setupDatabase(
		isSSLRequired(*instance.ProvisioningParameters),
		dt.AdministratorLogin,
		dt.ServerName,
		string(dt.AdministratorLoginPassword),
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
	if err != nil {
		return nil, err
	}
	return instance.Details, nil
}

func (a *allInOneManager) createExtensions(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*allInOneInstanceDetails)
	extensions := instance.ProvisioningParameters.GetStringArray("extensions")
	if len(extensions) > 0 {
		err := createExtensions(
			isSSLRequired(*instance.ProvisioningParameters),
			dt.AdministratorLogin,
			dt.ServerName,
			string(dt.AdministratorLoginPassword),
			dt.FullyQualifiedDomainName,
			dt.DatabaseName,
			extensions,
		)
		if err != nil {
			return nil, err
		}
	}
	return instance.Details, nil
}
