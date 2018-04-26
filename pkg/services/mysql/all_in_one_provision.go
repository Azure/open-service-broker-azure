package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	uuid "github.com/satori/go.uuid"
)

func (a *allInOneManager) ValidateProvisioningParameters(
	plan service.Plan,
	provisioningParameters service.ProvisioningParameters,
	_ service.SecureProvisioningParameters,
) error {
	pp := allInOneProvisioningParameters{}
	if err := service.GetStructFromMap(
		provisioningParameters,
		&pp,
	); err != nil {
		return err
	}
	return validateDBMSProvisionParameters(plan, pp.dbmsProvisioningParameters)
}

func (a *allInOneManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", a.preProvision),
		service.NewProvisioningStep("deployARMTemplate", a.deployARMTemplate),
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
	pp := allInOneProvisioningParameters{}
	if err :=
		service.GetStructFromMap(
			instance.ProvisioningParameters,
			&pp,
		); err != nil {
		return nil, nil, err
	}
	sslEnforcement := strings.ToLower(pp.SSLEnforcement)
	var enforceSSL bool
	if sslEnforcement == "" || sslEnforcement == enabled {
		enforceSSL = true
	}
	dt := allInOneInstanceDetails{
		dbmsInstanceDetails: dbmsInstanceDetails{
			ARMDeploymentName: uuid.NewV4().String(),
			ServerName:        serverName,
			EnforceSSL:        enforceSSL,
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

func (a *allInOneManager) buildARMTemplateParameters(
	details allInOneInstanceDetails,
	secureDetails secureAllInOneInstanceDetails,
) map[string]interface{} {
	var sslEnforcement string
	if details.EnforceSSL {
		sslEnforcement = enabledARMString
	} else {
		sslEnforcement = disabledARMString
	}
	p := map[string]interface{}{ // ARM template params
		"administratorLoginPassword": secureDetails.AdministratorLoginPassword,
		"serverName":                 details.ServerName,
		"databaseName":               details.DatabaseName,
		"sslEnforcement":             sslEnforcement,
	}
	return p
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
	if err := service.GetStructFromMap(
		instance.SecureDetails,
		&sdt,
	); err != nil {
		return nil, nil, err
	}
	pp := allInOneProvisioningParameters{}
	if err :=
		service.GetStructFromMap(
			instance.ProvisioningParameters,
			&pp,
		); err != nil {
		return nil, nil, err
	}
	armTemplateParameters := a.buildARMTemplateParameters(
		dt,
		sdt,
	)
	goTemplateParameters, err := buildGoTemplateParameters(instance)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"error building ARM template parameters %s",
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
		armTemplateParameters,
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	var ok bool
	dt.FullyQualifiedDomainName, ok =
		outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving fully qualified domain name from deployment: %s",
			err,
		)
	}

	dtMap, err := service.GetMapFromStruct(dt)
	return dtMap, instance.SecureDetails, err
}
