package mssql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/azure"
	"github.com/Azure/azure-service-broker/pkg/generate"
	"github.com/Azure/azure-service-broker/pkg/service"
	az "github.com/Azure/go-autorest/autorest/azure"
	uuid "github.com/satori/go.uuid"
)

func (m *module) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"*mssql.ProvisioningParameters",
		)
	}
	if pp.ServerName == "" {
		if !azure.IsValidLocation(pp.Location) {
			return service.NewValidationError(
				"location",
				fmt.Sprintf(`invalid location: "%s"`, pp.Location),
			)
		}
	} else {
		if _, ok := m.mssqlConfig.Servers[pp.ServerName]; !ok {
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

func (m *module) GetProvisioner(string, string) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", m.preProvision),
		service.NewProvisioningStep("deployARMTemplate", m.deployARMTemplate),
	)
}

func (m *module) preProvision(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string, // nolint: unparam
	planID string, // nolint: unparam
	provisioningContext service.ProvisioningContext,
	provisioningParameters service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*mssqlProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *mssqlProvisioningContext",
		)
	}
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters as " +
				"*mssql.ProvisioningParameters",
		)
	}

	if pp.ServerName == "" {
		// new server scenario
		if pp.ResourceGroup != "" {
			pc.ResourceGroupName = pp.ResourceGroup
		} else {
			pc.ResourceGroupName = uuid.NewV4().String()
		}
		pc.ARMDeploymentName = uuid.NewV4().String()
		pc.ServerName = uuid.NewV4().String()
		pc.IsNewServer = true
		pc.Location = pp.Location
		pc.AdministratorLogin = generate.NewIdentifier()
		pc.AdministratorLoginPassword = generate.NewPassword()
		pc.DatabaseName = generate.NewIdentifier()
	} else {
		// exisiting server scenario
		servers := m.mssqlConfig.Servers
		server, ok := servers[pp.ServerName]
		if !ok {
			return nil, fmt.Errorf(
				`can't find serverName "%s" in Azure SQL Server configuration`,
				pp.ServerName,
			)
		}

		pc.ResourceGroupName = server.ResourceGroupName
		pc.ARMDeploymentName = uuid.NewV4().String()
		pc.ServerName = server.ServerName
		pc.IsNewServer = false
		pc.Location = server.Location
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

func (m *module) deployARMTemplate(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string,
	planID string,
	provisioningContext service.ProvisioningContext,
	provisioningParameters service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*mssqlProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *mssqlProvisioningContext",
		)
	}
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters as " +
				"*mssql.ProvisioningParameters",
		)
	}
	catalog, err := m.GetCatalog()
	if err != nil {
		return nil, fmt.Errorf("error retrieving catalog: %s", err)
	}
	service, ok := catalog.GetService(serviceID)
	if !ok {
		return nil, fmt.Errorf(
			`service "%s" not found in the "%s" module catalog`,
			serviceID,
			m.GetName(),
		)
	}
	plan, ok := service.GetPlan(planID)
	if !ok {
		return nil, fmt.Errorf(
			`plan "%s" not found for service "%s"`,
			planID,
			serviceID,
		)
	}

	if pp.ServerName == "" {
		// new server scenario
		outputs, err := m.armDeployer.Deploy(
			pc.ARMDeploymentName,
			pc.ResourceGroupName,
			pc.Location,
			armTemplateNewServerBytes,
			map[string]interface{}{
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
			pp.Tags,
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
		_, err := m.armDeployer.Deploy(
			pc.ARMDeploymentName,
			pc.ResourceGroupName,
			pc.Location,
			armTemplateExistingServerBytes,
			map[string]interface{}{
				"serverName":   pc.ServerName,
				"databaseName": pc.DatabaseName,
				"edition":      plan.GetProperties().Extended["edition"],
				"requestedServiceObjectiveName": plan.GetProperties().
					Extended["requestedServiceObjectiveName"],
				"maxSizeBytes": plan.GetProperties().
					Extended["maxSizeBytes"],
			},
			pp.Tags,
		)
		if err != nil {
			return nil, fmt.Errorf("error deploying ARM template: %s", err)
		}
	}

	return pc, nil
}
