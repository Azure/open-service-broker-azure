package mysqldb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	uuid "github.com/satori/go.uuid"
)

const (
	enabled  = "enabled"
	disabled = "disabled"
)

func validateServerProvisionParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ServerProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters " +
				"as *mysql.ServerProvisioningParameters",
		)
	}
	sslEnforcement := strings.ToLower(pp.SSLEnforcement)
	if sslEnforcement != "" && sslEnforcement != enabled &&
		sslEnforcement != disabled {
		return service.NewValidationError(
			"sslEnforcement",
			fmt.Sprintf(`invalid option: "%s"`, pp.SSLEnforcement),
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

func (a *allInOneManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	return validateServerProvisionParameters(provisioningParameters)
}

func (v *vmOnlyManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	return validateServerProvisionParameters(provisioningParameters)
}

func (d *dbOnlyManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
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
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*allInOneMysqlInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *allInOneMysqlInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ServerProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*mysql.ServerProvisioningParameters",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.ServerName = uuid.NewV4().String()
	dt.AdministratorLoginPassword = generate.NewPassword()
	dt.DatabaseName = generate.NewIdentifier()

	sslEnforcement := strings.ToLower(pp.SSLEnforcement)
	switch sslEnforcement {
	case "", enabled:
		dt.EnforceSSL = true
	case disabled:
		dt.EnforceSSL = false
	}

	return dt, nil
}

func (v *vmOnlyManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*vmOnlyMysqlInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *vmOnlyMysqlInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ServerProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*mysql.ServerProvisioningParameters",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.ServerName = uuid.NewV4().String()
	dt.AdministratorLoginPassword = generate.NewPassword()

	sslEnforcement := strings.ToLower(pp.SSLEnforcement)
	switch sslEnforcement {
	case "", enabled:
		dt.EnforceSSL = true
	case disabled:
		dt.EnforceSSL = false
	}

	return dt, nil
}

func (d *dbOnlyManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*dbOnlyMysqlInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *dbOnlyMysqlInstanceDetails",
		)
	}
	//We aren't using any of these, but validate it can be type cast
	_, ok = instance.ProvisioningParameters.(*DatabaseProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*mysqldb.DatabaseProvisioningParameters",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.DatabaseName = generate.NewIdentifier()

	return dt, nil
}

func (a *allInOneManager) buildARMTemplateParameters(
	plan service.Plan,
	details *allInOneMysqlInstanceDetails,
	provisioningParameters *ServerProvisioningParameters,
) map[string]interface{} {
	var sslEnforcement string
	if details.EnforceSSL {
		sslEnforcement = "Enabled"
	} else {
		sslEnforcement = "Disabled"
	}
	p := map[string]interface{}{ // ARM template params
		"administratorLoginPassword": details.AdministratorLoginPassword,
		"serverName":                 details.ServerName,
		"databaseName":               details.DatabaseName,
		"skuName":                    plan.GetProperties().Extended["skuName"],
		"skuTier":                    plan.GetProperties().Extended["skuTier"],
		"skuCapacityDTU": plan.GetProperties().
			Extended["skuCapacityDTU"],
		"skuSizeMB":      plan.GetProperties().Extended["skuSizeMB"],
		"sslEnforcement": sslEnforcement,
	}
	//Only include these if they are not empty.
	//ARM Deployer will fail if the values included are not
	//valid IPV4 addresses (i.e. empty string will fail)
	if provisioningParameters.FirewallIPStart != "" {
		p["firewallStartIpAddress"] = provisioningParameters.FirewallIPStart
	}
	if provisioningParameters.FirewallIPEnd != "" {
		p["firewallEndIpAddress"] = provisioningParameters.FirewallIPEnd
	}
	return p
}

func (v *vmOnlyManager) buildARMTemplateParameters(
	plan service.Plan,
	details *vmOnlyMysqlInstanceDetails,
	provisioningParameters *ServerProvisioningParameters,
) map[string]interface{} {
	var sslEnforcement string
	if details.EnforceSSL {
		sslEnforcement = "Enabled"
	} else {
		sslEnforcement = "Disabled"
	}
	p := map[string]interface{}{ // ARM template params
		"administratorLoginPassword": details.AdministratorLoginPassword,
		"serverName":                 details.ServerName,
		"skuName":                    plan.GetProperties().Extended["skuName"],
		"skuTier":                    plan.GetProperties().Extended["skuTier"],
		"skuCapacityDTU": plan.GetProperties().
			Extended["skuCapacityDTU"],
		"skuSizeMB":      plan.GetProperties().Extended["skuSizeMB"],
		"sslEnforcement": sslEnforcement,
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

func (a *allInOneManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*allInOneMysqlInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *allInOneMysqlInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ServerProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters " +
				"as *mysql.ServerProvisioningParameters",
		)
	}
	armTemplateParameters := a.buildARMTemplateParameters(instance.Plan, dt, pp)
	outputs, err := a.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		allInOneArmTemplateBytes,
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

	return dt, nil
}

func (v *vmOnlyManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*vmOnlyMysqlInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *vmOnlyMysqlInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ServerProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters as " +
				"*mysql.ServerProvisioningParameters",
		)
	}
	armTemplateParameters := v.buildARMTemplateParameters(instance.Plan, dt, pp)
	outputs, err := v.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		vmOnlyArmTemplateBytes,
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

	return dt, nil
}

func (d *dbOnlyManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pdt, ok := instance.Parent.Details.(*vmOnlyMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Parent.Details as *vmOnlyMysqlInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*dbOnlyMysqlInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *dbOnlyMysqlInstanceDetails",
		)
	}
	_, ok = instance.ProvisioningParameters.(*DatabaseProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters " +
				"as *mysql.DatabaseProvisioningParameters",
		)
	}

	armTemplateParameters := map[string]interface{}{ // ARM template params
		"serverName":   pdt.ServerName,
		"databaseName": dt.DatabaseName,
	}
	_, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.Parent.ResourceGroup,
		instance.Parent.Location,
		dbOnlyArmTemplateBytes,
		nil, // Go template params
		armTemplateParameters,
		instance.Tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	return dt, nil
}
