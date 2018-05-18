package postgresql

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
		service.NewProvisioningStep("setupDatabase", d.setupDatabase),
		service.NewProvisioningStep("createExtensions", d.createExtensions),
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
	pdt := dbmsInstanceDetails{}
	if err :=
		service.GetStructFromMap(instance.Parent.Details, &pdt); err != nil {
		return nil, nil, err
	}
	dt := databaseInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	goTemplateParameters := map[string]interface{}{
		"serverName":   pdt.ServerName,
		"databaseName": dt.DatabaseName,
	}
	_, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.Parent.ResourceGroup,
		instance.Parent.Location,
		databaseARMTemplateBytes,
		goTemplateParameters,
		map[string]interface{}{}, // ARM template params
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, instance.SecureDetails, nil
}

func (d *databaseManager) setupDatabase(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	pdt := dbmsInstanceDetails{}
	if err :=
		service.GetStructFromMap(instance.Parent.Details, &pdt); err != nil {
		return nil, nil, err
	}
	spdt := secureDBMSInstanceDetails{}
	if err :=
		service.GetStructFromMap(instance.Parent.SecureDetails, &spdt); err != nil {
		return nil, nil, err
	}
	dt := databaseInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	err := setupDatabase(
		instance.Parent.ProvisioningParameters.GetString("sslEnforcement") ==
			"enabled",
		pdt.ServerName,
		spdt.AdministratorLoginPassword,
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
	if err != nil {
		return nil, nil, err
	}
	return instance.Details, instance.SecureDetails, nil
}

func (d *databaseManager) createExtensions(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	pdt := dbmsInstanceDetails{}
	if err :=
		service.GetStructFromMap(instance.Parent.Details, &pdt); err != nil {
		return nil, nil, err
	}
	spdt := secureDBMSInstanceDetails{}
	if err :=
		service.GetStructFromMap(instance.Parent.SecureDetails, &spdt); err != nil {
		return nil, nil, err
	}
	dt := databaseInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	extensions := instance.ProvisioningParameters.GetArray("extensions")
	if len(extensions) > 0 {
		extensionStrs := make([]string, len(extensions))
		for i, extension := range extensions {
			extensionStrs[i], _ = extension.(string)
		}
		err := createExtensions(
			instance.Parent.ProvisioningParameters.GetString("sslEnforcement") ==
				"enabled",
			pdt.ServerName,
			spdt.AdministratorLoginPassword,
			pdt.FullyQualifiedDomainName,
			dt.DatabaseName,
			extensionStrs,
		)
		if err != nil {
			return nil, nil, err
		}
	}
	return instance.Details, instance.SecureDetails, nil
}
