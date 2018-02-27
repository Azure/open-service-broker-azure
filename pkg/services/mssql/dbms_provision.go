package mssql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *dbmsManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
	_ service.SecureProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*DBMSProvisioningParams)

	if !ok {
		return errors.New(
			"error casting provisioningParameters as *mssql.DBMSProvisioningParams",
		)
	}
	return validateDBMSProvisionParameters(pp)
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
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*dbmsInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *mssql.dbmsInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*secureDBMSInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails as " +
				"*mssql.secureDBMSInstanceDetails",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.ServerName = uuid.NewV4().String()
	dt.AdministratorLogin = generate.NewIdentifier()
	sdt.AdministratorLoginPassword = generate.NewPassword()
	return dt, sdt, nil
}

func (d *dbmsManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*dbmsInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *mssql.dbmsInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*secureDBMSInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails as " +
				"*mssql.secureDBMSInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*DBMSProvisioningParams)
	if !ok {
		return nil, nil, errors.New(
			"error casting provisioningParameters as *mssql.DBMSProvisioningParams",
		)
	}
	p := map[string]interface{}{ // ARM template params
		"serverName":                 dt.ServerName,
		"administratorLogin":         dt.AdministratorLogin,
		"administratorLoginPassword": sdt.AdministratorLoginPassword,
	}
	goTemplateParams := buildGoTemplateParameters(pp)
	outputs, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		dbmsARMTemplateBytes,
		goTemplateParams,
		p,
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	fullyQualifiedDomainName, ok := outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving fully qualified domain name from deployment: %s",
			err,
		)
	}
	dt.FullyQualifiedDomainName = fmt.Sprintf(
		"%s.%s",
		fullyQualifiedDomainName,
		d.sqlDatabaseDNSSuffix,
	)
	return dt, sdt, nil
}
