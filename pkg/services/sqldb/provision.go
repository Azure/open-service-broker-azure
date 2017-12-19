package sqldb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (a *allServiceManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ServerProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"*mssql.ServerProvisioningParameters",
		)
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

//TODO implement DB Only valiation
func (d *dbServiceManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	return nil
}

func (s *vmServiceManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ServerProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"*mssql.ProvisioningParameters",
		)
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

func (a *allServiceManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", a.preProvision),
		service.NewProvisioningStep("deployARMTemplate", a.deployARMTemplate),
	)
}

func (s *vmServiceManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", s.preProvision),
		service.NewProvisioningStep("deployARMTemplate", s.deployARMTemplate),
	)
}

func (d *dbServiceManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("deployARMTemplate", d.deployARMTemplate),
	)
}

func (a *allServiceManager) preProvision(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlAllInOneInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *mssqlAllInOneInstanceDetails",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.ServerName = uuid.NewV4().String()
	dt.AdministratorLogin = generate.NewIdentifier()
	dt.AdministratorLoginPassword = generate.NewPassword()
	dt.DatabaseName = generate.NewIdentifier()
	return dt, nil
}

func (s *vmServiceManager) preProvision(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlServerOnlyInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *mssqlServerOnlyInstanceDetails",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.ServerName = uuid.NewV4().String()
	dt.AdministratorLogin = generate.NewIdentifier()
	dt.AdministratorLoginPassword = generate.NewPassword()

	return dt, nil
}

//TODO implement DB Only preProvision
func (d *dbServiceManager) preProvision(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	return instance.Details, nil
}

func (a *allServiceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
	plan service.Plan,
	_ service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlAllInOneInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *mssqlAllInOneInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ServerProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*mssql.ServerProvisioningParameters",
		)
	}
	armTemplateParams := map[string]interface{}{ // ARM template params
		"serverName":                 dt.ServerName,
		"administratorLogin":         dt.AdministratorLogin,
		"administratorLoginPassword": dt.AdministratorLoginPassword,
	}
	//Only include these if they are not empty.
	//ARM Deployer will fail if the values included are not
	//valid IPV4 addresses (i.e. empty string wil fail)
	if pp.FirewallIPStart != "" {
		armTemplateParams["firewallStartIpAddress"] = pp.FirewallIPStart
	}
	if pp.FirewallIPEnd != "" {
		armTemplateParams["firewallEndIpAddress"] = pp.FirewallIPEnd
	}
	//These are the database related properties that are needed
	//for this deployment
	armTemplateParams["databaseName"] = dt.DatabaseName
	armTemplateParams["edition"] = plan.GetProperties().Extended["edition"]
	armTemplateParams["requestedServiceObjectiveName"] = plan.GetProperties().
		Extended["requestedServiceObjectiveName"]
	armTemplateParams["maxSizeBytes"] = plan.GetProperties().
		Extended["maxSizeBytes"]
	outputs, err := a.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateNewServerBytes,
		nil, // Go template params
		armTemplateParams,
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

	return dt, nil
}

func (s *vmServiceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
	plan service.Plan,
	_ service.Instance, //referenace instance
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlServerOnlyInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *mssqlServerOnlyInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ServerProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*mssql.ServerProvisioningParameters",
		)
	}
	armTemplateParams := map[string]interface{}{ // ARM template params
		"serverName":                 dt.ServerName,
		"administratorLogin":         dt.AdministratorLogin,
		"administratorLoginPassword": dt.AdministratorLoginPassword,
	}
	if pp.FirewallIPStart != "" {
		armTemplateParams["firewallStartIpAddress"] = pp.FirewallIPStart
	}
	if pp.FirewallIPEnd != "" {
		armTemplateParams["firewallEndIpAddress"] = pp.FirewallIPEnd
	}
	// new server scenario
	outputs, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateServerOnlyBytes,
		nil, // Go template params
		armTemplateParams,
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
	return dt, nil
}

//TODO implement DB only scenario
func (d *dbServiceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
	_ service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	return instance.Details, nil
}
