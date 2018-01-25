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
	if pp.FirewallIPStart != "" || pp.FirewallIPEnd != "" {
		if pp.FirewallIPStart == "" {
			return service.NewValidationError(
				"firewallStartIPAddress",
				"must be set when firewallEndIPAddress is set",
			)
		}
		if pp.FirewallIPEnd == "" {
			return service.NewValidationError(
				"firewallEndIPAddress",
				"must be set when firewallStartIPAddress is set",
			)
		}
	}
	startIP := net.ParseIP(pp.FirewallIPStart)
	if pp.FirewallIPStart != "" && startIP == nil {
		return service.NewValidationError(
			"firewallStartIPAddress",
			fmt.Sprintf(`invalid value: "%s"`, pp.FirewallIPStart),
		)
	}
	endIP := net.ParseIP(pp.FirewallIPEnd)
	if pp.FirewallIPEnd != "" && endIP == nil {
		return service.NewValidationError(
			"firewallEndIPAddress",
			fmt.Sprintf(`invalid value: "%s"`, pp.FirewallIPEnd),
		)
	}
	//The net.IP.To4 method returns a 4 byte representation of an IPv4 address.
	//Once converted,comparing two IP addresses can be done by using the
	//bytes. Compare function. Per the ARM template documentation,
	//startIP must be <= endIP.
	startBytes := startIP.To4()
	endBytes := endIP.To4()
	if bytes.Compare(startBytes, endBytes) > 0 {
		return service.NewValidationError(
			"firewallEndIPAddress",
			fmt.Sprintf(`invalid value: "%s". must be 
				greater than or equal to firewallStartIPAddress`, pp.FirewallIPEnd),
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
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *mssqlInstanceDetails",
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
		dt.ARMDeploymentName = uuid.NewV4().String()
		dt.ServerName = uuid.NewV4().String()
		dt.IsNewServer = true
		dt.AdministratorLogin = generate.NewIdentifier()
		dt.AdministratorLoginPassword = generate.NewPassword()
		dt.DatabaseName = generate.NewIdentifier()
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

		dt.ARMDeploymentName = uuid.NewV4().String()
		dt.ServerName = server.ServerName
		dt.IsNewServer = false
		dt.AdministratorLogin = server.AdministratorLogin
		dt.AdministratorLoginPassword = server.AdministratorLoginPassword
		dt.DatabaseName = generate.NewIdentifier()

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
		dt.FullyQualifiedDomainName = fmt.Sprintf(
			"%s.%s",
			server.ServerName,
			sqlDatabaseDNSSuffix,
		)
	}
	return dt, nil
}

func buildARMTemplateParameters(
	plan service.Plan,
	details *mssqlInstanceDetails,
	provisioningParameters *ProvisioningParameters,
) map[string]interface{} {
	p := map[string]interface{}{ // ARM template params
		"serverName":                 details.ServerName,
		"administratorLogin":         details.AdministratorLogin,
		"administratorLoginPassword": details.AdministratorLoginPassword,
		"databaseName":               details.DatabaseName,
		"edition":                    plan.GetProperties().Extended["edition"],
		"requestedServiceObjectiveName": plan.GetProperties().
			Extended["requestedServiceObjectiveName"],
		"maxSizeBytes": plan.GetProperties().
			Extended["maxSizeBytes"],
	}
	//Only include these if they are not empty.
	//ARM Deployer will fail if the values included are not
	//valid IPV4 addresses (i.e. empty string wil fail)
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
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *mssqlInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*mssql.ProvisioningParameters",
		)
	}
	plan := instance.Plan
	if dt.IsNewServer {
		armTemplateParameters := buildARMTemplateParameters(plan, dt, pp)
		// new server scenario
		outputs, err := s.armDeployer.Deploy(
			dt.ARMDeploymentName,
			instance.ResourceGroup,
			instance.Location,
			armTemplateNewServerBytes,
			nil, // Go template params
			armTemplateParameters,
			instance.Tags,
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
		dt.FullyQualifiedDomainName = fullyQualifiedDomainName
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
			dt.ARMDeploymentName,
			server.ResourceGroupName,
			server.Location,
			armTemplateExistingServerBytes,
			nil, // Go template params
			map[string]interface{}{ // ARM template params
				"serverName":   dt.ServerName,
				"databaseName": dt.DatabaseName,
				"edition":      plan.GetProperties().Extended["edition"],
				"requestedServiceObjectiveName": plan.GetProperties().
					Extended["requestedServiceObjectiveName"],
				"maxSizeBytes": plan.GetProperties().
					Extended["maxSizeBytes"],
			},
			instance.Tags,
		)
		if err != nil {
			return nil, fmt.Errorf("error deploying ARM template: %s", err)
		}
	}

	return dt, nil
}
