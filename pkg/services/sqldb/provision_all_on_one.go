package sqldb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (a *allInOneManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
	_ service.SecureProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ServerProvisioningParams)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"*mssql.ServerProvisioningParams",
		)
	}
	return validateServerProvisionParameters(pp)
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
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlAllInOneInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *mssqlAllInOneInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*mssqlAllInOneSecureInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails as " +
				"*mssqlAllInOneSecureInstanceDetails",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.ServerName = uuid.NewV4().String()
	dt.AdministratorLogin = generate.NewIdentifier()
	sdt.AdministratorLoginPassword = generate.NewPassword()
	dt.DatabaseName = generate.NewIdentifier()
	return dt, sdt, nil
}

func (a *allInOneManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlAllInOneInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *mssqlAllInOneInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*mssqlAllInOneSecureInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails as " +
				"*mssqlAllInOneSecureInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ServerProvisioningParams)
	if !ok {
		return nil, nil, errors.New(
			"error casting provisioningParameters as " +
				"*mssql.ServerProvisioningParams",
		)
	}
	p := map[string]interface{}{ // ARM template params
		"serverName":                 dt.ServerName,
		"administratorLogin":         dt.AdministratorLogin,
		"administratorLoginPassword": sdt.AdministratorLoginPassword,
		"databaseName":               dt.DatabaseName,
		"edition": instance.Plan.GetProperties().
			Extended["edition"],
		"requestedServiceObjectiveName": instance.Plan.GetProperties().
			Extended["requestedServiceObjectiveName"],
		"maxSizeBytes": instance.Plan.GetProperties().Extended["maxSizeBytes"],
	}
	goTemplateParams := buildGoTemplateParameters(pp)
	// new server scenario
	outputs, err := a.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateDBMSOnlyBytes,
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
	dt.FullyQualifiedDomainName = fullyQualifiedDomainName
	return dt, sdt, nil
}
