package postgresql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *dbmsManager) ValidateProvisioningParameters(
	plan service.Plan,
	provisioningParameters service.ProvisioningParameters,
	_ service.SecureProvisioningParameters,
) error {
	pp := dbmsProvisioningParameters{}
	if err := service.GetStructFromMap(provisioningParameters, &pp); err != nil {
		return err
	}
	return validateDBMSProvisionParameters(plan, pp)
}

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
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	serverName, err := getAvailableServerName(
		ctx,
		d.checkNameAvailabilityClient,
	)
	if err != nil {
		return nil, nil, err
	}
	dt := dbmsInstanceDetails{
		ARMDeploymentName: uuid.NewV4().String(),
		ServerName:        serverName,
	}

	sdt := secureDBMSInstanceDetails{
		AdministratorLoginPassword: generate.NewPassword(),
	}

	dtMap, err := service.GetMapFromStruct(dt)
	if err != nil {
		return nil, nil, err
	}
	sdtMap, err := service.GetMapFromStruct(sdt)
	return dtMap, sdtMap, err
}

func (d *dbmsManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := dbmsInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	sdt := secureDBMSInstanceDetails{}
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, nil, err
	}
	pp := dbmsProvisioningParameters{}
	if err :=
		service.GetStructFromMap(instance.ProvisioningParameters, &pp); err != nil {
		return nil, nil, err
	}
	goTemplateParameters, err := buildGoTemplateParameters(instance)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to build go template parameters: %s", err)
	}
	outputs, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		dbmsARMTemplateBytes,
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
