package sqldb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"

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

	startIP := net.ParseIP(pp.FirewallIPStart)
	if pp.FirewallIPStart != "" && net.ParseIP(pp.FirewallIPStart) == nil {
		return service.NewValidationError(
			"firewallStartIpAddress",
			fmt.Sprintf(`invalid firewallStartIpAddress option: "%s"`, pp.FirewallIPStart),
		)
	}

	endIP := net.ParseIP(pp.FirewallIPEnd)
	if pp.FirewallIPEnd != "" && net.ParseIP(pp.FirewallIPEnd) == nil {
		return service.NewValidationError(
			"firewallEndIpAddress",
			fmt.Sprintf(`invalid firewallEndIpAddress option: "%s"`, pp.FirewallIPEnd),
		)
	}

	if startIP != nil && endIP == nil {
		return service.NewValidationError(
			"firewallEndIpAddress",
			fmt.Sprintf(`invalid firewallEndIpAddress option: "%s"`, pp.FirewallIPEnd),
		)
	}

	startBytes := startIP.To4()
	endBytes := endIP.To4()
	if bytes.Compare(startBytes, endBytes) > 0 {
		return service.NewValidationError(
			"firewallEndIpAddress",
			fmt.Sprintf(`invalid firewallEndIpAddress option: "%s"`, pp.FirewallIPEnd),
		)
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
	_ string, // instanceID
	_ service.Plan, // planID
	_ service.StandardProvisioningContext,
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

func buildARMTemplateParameters(
	plan service.Plan,
	provisioningContext *mssqlProvisioningContext,
	provisioningParameters *ProvisioningParameters,
) map[string]interface{} {
	p := map[string]interface{}{ // ARM template params
		"serverName":                 provisioningContext.ServerName,
		"administratorLogin":         provisioningContext.AdministratorLogin,
		"administratorLoginPassword": provisioningContext.AdministratorLoginPassword,
		"databaseName":               provisioningContext.DatabaseName,
		"edition":                    plan.GetProperties().Extended["edition"],
		"requestedServiceObjectiveName": plan.GetProperties().
			Extended["requestedServiceObjectiveName"],
		"maxSizeBytes": plan.GetProperties().
			Extended["maxSizeBytes"],
	}
	//Only include these if they are not empty. ARM Deployer will fail
	//if the values included are not valid IPV4 addresses (i.e. empty string wil fail)
	if provisioningParameters.FirewallIPStart != "" {
		p["firewallStartIpAddress"] = provisioningParameters.FirewallIPStart
	}
	if provisioningParameters.FirewallIPEnd != "" {
		p["firewallEndIpAddress"] = provisioningParameters.FirewallIPEnd
	}
	return p
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	_ string, // instanceID
	plan service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
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
	if pc.IsNewServer {
		armTemplateParameters := buildARMTemplateParameters(plan, pc, pp)
		// new server scenario
		outputs, err := s.armDeployer.Deploy(
			pc.ARMDeploymentName,
			standardProvisioningContext.ResourceGroup,
			standardProvisioningContext.Location,
			armTemplateNewServerBytes,
			nil, // Go template params
			armTemplateParameters,
			standardProvisioningContext.Tags,
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
			standardProvisioningContext.Tags,
		)
		if err != nil {
			return nil, fmt.Errorf("error deploying ARM template: %s", err)
		}
	}

	return pc, nil
}
