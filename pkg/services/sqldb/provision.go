package sqldb

import (
	"context"
	"errors"
	"fmt"

	az "github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (s *serviceManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"*mssql.ProvisioningParameters",
		)
	}
	if pp.ServerName != "" {
		if _, ok := s.mssqlConfig.Servers[pp.ServerName]; !ok {
			return service.NewValidationError(
				"serverName",
				fmt.Sprintf(
					`can't find serverName "%s" in Azure SQL Server configuration`,
					pp.ServerName,
				),
			)
		}
	}
	return nil
}

func (s *serviceManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", s.preProvision),
		service.NewProvisioningStep("deployARMTemplate", s.deployARMTemplate),
	)
}

func (s *serviceManager) preProvision(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*mssqlProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningContext as *mssqlProvisioningContext",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*mssql.ProvisioningParameters",
		)
	}

	if pp.ServerName == "" {
		// new server scenario
		pc.ARMDeploymentName = uuid.NewV4().String()
		pc.ServerName = uuid.NewV4().String()
		pc.IsNewServer = true
		pc.AdministratorLogin = generate.NewIdentifier()
		pc.AdministratorLoginPassword = generate.NewPassword()
		pc.DatabaseName = generate.NewIdentifier()
	} else {
		// exisiting server scenario
		servers := s.mssqlConfig.Servers
		server, ok := servers[pp.ServerName]
		if !ok {
			return nil, fmt.Errorf(
				`can't find serverName "%s" in Azure SQL Server configuration`,
				pp.ServerName,
			)
		}

		pc.ARMDeploymentName = uuid.NewV4().String()
		pc.ServerName = server.ServerName
		pc.IsNewServer = false
		pc.AdministratorLogin = server.AdministratorLogin
		pc.AdministratorLoginPassword = server.AdministratorLoginPassword
		pc.DatabaseName = generate.NewIdentifier()

		// Ensure the server configuration works
		azureConfig, err := azure.GetConfig()
		if err != nil {
			return nil, err
		}
		azureEnvironment, err := az.EnvironmentFromName(azureConfig.Environment)
		if err != nil {
			return nil, err
		}
		sqlDatabaseDNSSuffix := azureEnvironment.SQLDatabaseDNSSuffix
		pc.FullyQualifiedDomainName = fmt.Sprintf(
			"%s.%s",
			server.ServerName,
			sqlDatabaseDNSSuffix,
		)
	}
	return pc, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
	plan service.Plan,
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*mssqlProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningContext as *mssqlProvisioningContext",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*mssql.ProvisioningParameters",
		)
	}
	if pc.IsNewServer {
		// new server scenario
		outputs, err := s.armDeployer.Deploy(
			pc.ARMDeploymentName,
			instance.StandardProvisioningContext.ResourceGroup,
			instance.StandardProvisioningContext.Location,
			armTemplateNewServerBytes,
			nil, // Go template params
			map[string]interface{}{ // ARM template params
				"serverName":                 pc.ServerName,
				"administratorLogin":         pc.AdministratorLogin,
				"administratorLoginPassword": pc.AdministratorLoginPassword,
				"databaseName":               pc.DatabaseName,
				"edition":                    plan.GetProperties().Extended["edition"],
				"requestedServiceObjectiveName": plan.GetProperties().
					Extended["requestedServiceObjectiveName"],
				"maxSizeBytes": plan.GetProperties().
					Extended["maxSizeBytes"],
			},
			instance.StandardProvisioningContext.Tags,
		)
		if err != nil {
			return nil, fmt.Errorf("error deploying ARM template: %s", err)
		}
		fullyQualifiedDomainName, ok := outputs["fullyQualifiedDomainName"].(string)
		if !ok {
			return nil, fmt.Errorf(
				"error retrieving fully qualified domain name from deployment: %s",
				err,
			)
		}
		pc.FullyQualifiedDomainName = fullyQualifiedDomainName
	} else {
		// existing server scenario
		servers := s.mssqlConfig.Servers
		server, ok := servers[pp.ServerName]
		if !ok {
			return nil, fmt.Errorf(
				`can't find serverName "%s" in Azure SQL Server configuration`,
				pp.ServerName,
			)
		}

		_, err := s.armDeployer.Deploy(
			pc.ARMDeploymentName,
			server.ResourceGroupName,
			server.Location,
			armTemplateExistingServerBytes,
			nil, // Go template params
			map[string]interface{}{ // ARM template params
				"serverName":   pc.ServerName,
				"databaseName": pc.DatabaseName,
				"edition":      plan.GetProperties().Extended["edition"],
				"requestedServiceObjectiveName": plan.GetProperties().
					Extended["requestedServiceObjectiveName"],
				"maxSizeBytes": plan.GetProperties().
					Extended["maxSizeBytes"],
			},
			instance.StandardProvisioningContext.Tags,
		)
		if err != nil {
			return nil, fmt.Errorf("error deploying ARM template: %s", err)
		}
	}

	return pc, nil
}
