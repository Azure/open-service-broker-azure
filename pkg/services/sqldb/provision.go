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

func (a *allInOneManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ServerProvisioningParams)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"*mssql.ServerProvisioningParams",
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

func (v *vmOnlyManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ServerProvisioningParams)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"*mssql.ServerProvisioningParams",
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

//TODO: implement db only validation
func (d *dbOnlyManager) ValidateProvisioningParameters(
	_ service.ProvisioningParameters,
) error {
	return nil
}

func (a *allInOneManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", a.preProvision),
		service.NewProvisioningStep("deployARMTemplate", a.deployARMTemplate),
	)
}

func (v *vmOnlyManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", v.preProvision),
		service.NewProvisioningStep("deployARMTemplate", v.deployARMTemplate),
	)
}

func (d *dbOnlyManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("deployARMTemplate", d.deployARMTemplate),
	)
}

func (a *allInOneManager) preProvision(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
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

func (v *vmOnlyManager) preProvision(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlVMOnlyInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *mssqlVMOnlyInstanceDetails",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.ServerName = uuid.NewV4().String()
	dt.AdministratorLogin = generate.NewIdentifier()
	dt.AdministratorLoginPassword = generate.NewPassword()
	return dt, nil
}

func (d *dbOnlyManager) preProvision(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlDBOnlyInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *mssqlDBOnlyInstanceDetails",
		)
	}
	//Assume refererence instance is a vm only instance. Fail if not
	pdt, ok := instance.Parent.Details.(*mssqlVMOnlyInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Parent.Details as " +
				"*mssqlVMOnlyInstanceDetails",
		)
	}

	azureConfig, err := azure.GetConfig()
	if err != nil {
		return nil, err
	}
	azureEnvironment, err := az.EnvironmentFromName(azureConfig.Environment)
	if err != nil {
		return nil, err
	}

	dt.DatabaseName = generate.NewIdentifier()

	sqlDatabaseDNSSuffix := azureEnvironment.SQLDatabaseDNSSuffix
	dt.FullyQualifiedDomainName = fmt.Sprintf(
		"%s.%s",
		pdt.ServerName,
		sqlDatabaseDNSSuffix,
	)

	return dt, nil
}

func (a *allInOneManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
	plan service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlAllInOneInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *mssqlAllInOneInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ServerProvisioningParams)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters as " +
				"*mssql.ServerProvisioningParams",
		)
	}
	p := map[string]interface{}{ // ARM template params
		"serverName":                 dt.ServerName,
		"administratorLogin":         dt.AdministratorLogin,
		"administratorLoginPassword": dt.AdministratorLoginPassword,
		"databaseName":               dt.DatabaseName,
		"edition":                    plan.GetProperties().Extended["edition"],
		"requestedServiceObjectiveName": plan.GetProperties().
			Extended["requestedServiceObjectiveName"],
		"maxSizeBytes": plan.GetProperties().
			Extended["maxSizeBytes"],
	}
	//Only include these if they are not empty.
	//ARM Deployer will fail if the values included are not
	//valid IPV4 addresses (i.e. empty string wil fail)
	if pp.FirewallIPStart != "" {
		p["firewallStartIpAddress"] = pp.FirewallIPStart
	}
	if pp.FirewallIPEnd != "" {
		p["firewallEndIpAddress"] = pp.FirewallIPEnd
	}
	// new server scenario
	outputs, err := a.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateNewServerBytes,
		nil, // Go template params
		p,
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

func (v *vmOnlyManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlVMOnlyInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *mssqlVMOnlyInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ServerProvisioningParams)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters as " +
				"*mssql.ServerProvisioningParams",
		)
	}
	p := map[string]interface{}{ // ARM template params
		"serverName":                 dt.ServerName,
		"administratorLogin":         dt.AdministratorLogin,
		"administratorLoginPassword": dt.AdministratorLoginPassword,
	}
	//Only include these if they are not empty.
	//ARM Deployer will fail if the values included are not
	//valid IPV4 addresses (i.e. empty string wil fail)
	if pp.FirewallIPStart != "" {
		p["firewallStartIpAddress"] = pp.FirewallIPStart
	}
	if pp.FirewallIPEnd != "" {
		p["firewallEndIpAddress"] = pp.FirewallIPEnd
	}
	// new server scenario
	outputs, err := v.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateServerOnlyBytes,
		nil, // Go template params
		p,
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

func (d *dbOnlyManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
	plan service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlDBOnlyInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *mssqlDBOnlyInstanceDetails",
		)
	}
	pdt, ok := instance.Details.(*mssqlVMOnlyInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Parent.Details as " +
				"*mssqlVMOnlyInstanceDetails",
		)
	}
	p := map[string]interface{}{ // ARM template params
		"serverName":   pdt.ServerName,
		"databaseName": dt.DatabaseName,
		"edition":      plan.GetProperties().Extended["edition"],
		"requestedServiceObjectiveName": plan.GetProperties().
			Extended["requestedServiceObjectiveName"],
		"maxSizeBytes": plan.GetProperties().
			Extended["maxSizeBytes"],
	}
	//No output, so ignore the output
	_, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.Parent.ResourceGroup,
		instance.Parent.Location,
		armTemplateDBOnlyBytes,
		nil, // Go template params
		p,
		instance.Tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return dt, nil
}
