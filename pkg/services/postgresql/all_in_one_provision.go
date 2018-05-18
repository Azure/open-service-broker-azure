package postgresql

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
		service.NewProvisioningStep("setupDatabase", a.setupDatabase),
		service.NewProvisioningStep("createExtensions", a.createExtensions),
	)
}

func (a *allInOneManager) preProvision(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	serverName, err := getAvailableServerName(
		ctx,
		a.checkNameAvailabilityClient,
	)
	if err != nil {
		return nil, nil, err
	}
	dt := allInOneInstanceDetails{
		dbmsInstanceDetails: dbmsInstanceDetails{
			ARMDeploymentName: uuid.NewV4().String(),
			ServerName:        serverName,
		},
		DatabaseName: generate.NewIdentifier(),
	}

	sdt := secureAllInOneInstanceDetails{
		secureDBMSInstanceDetails: secureDBMSInstanceDetails{
			AdministratorLoginPassword: generate.NewPassword(),
		},
	}

	dtMap, err := service.GetMapFromStruct(dt)
	if err != nil {
		return nil, nil, err
	}
	sdtMap, err := service.GetMapFromStruct(sdt)
	return dtMap, sdtMap, err
}

func (a *allInOneManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := allInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	sdt := secureAllInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, nil, err
	}

	goTemplateParameters, err := buildGoTemplateParameters(
		instance.Plan,
		dt.dbmsInstanceDetails,
		sdt.secureDBMSInstanceDetails,
		instance.ProvisioningParameters,
	)

	if err != nil {
		return nil, nil, fmt.Errorf(
			"error building go template parameters :%s",
			err,
		)
	}
	goTemplateParameters["databaseName"] = dt.DatabaseName
	outputs, err := a.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		allInOneARMTemplateBytes,
		goTemplateParameters,
		map[string]interface{}{},
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	var ok bool
	dt.FullyQualifiedDomainName, ok = outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving fully qualified domain name from deployment: %s",
			err,
		)
	}

	dtMap, err := service.GetMapFromStruct(dt)
	return dtMap, instance.SecureDetails, err
}

func (a *allInOneManager) setupDatabase(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := allInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	sdt := secureAllInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, nil, err
	}
	err := setupDatabase(
		instance.ProvisioningParameters.GetString("sslEnforcement") == "enabled",
		dt.ServerName,
		sdt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
	if err != nil {
		return nil, nil, err
	}
	return instance.Details, instance.SecureDetails, nil
}

func (a *allInOneManager) createExtensions(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := allInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	sdt := secureAllInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, nil, err
	}
	extensions := instance.ProvisioningParameters.GetArray("extensions")
	if len(extensions) > 0 {
		extensionStrs := make([]string, len(extensions))
		for i, extension := range extensions {
			extensionStrs[i], _ = extension.(string)
		}
		err := createExtensions(
			instance.ProvisioningParameters.GetString("sslEnforcement") == "enabled",
			dt.ServerName,
			sdt.AdministratorLoginPassword,
			dt.FullyQualifiedDomainName,
			dt.DatabaseName,
			extensionStrs,
		)
		if err != nil {
			return nil, nil, err
		}
	}
	return instance.Details, instance.SecureDetails, nil
}
