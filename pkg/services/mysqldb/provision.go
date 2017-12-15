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

func (s *serviceManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"*mysql.ProvisioningParameters",
		)
	}
	sslEnforcement := strings.ToLower(pp.SSLEnforcement)
	if sslEnforcement != "" && sslEnforcement != "enabled" &&
		sslEnforcement != "disabled" {
		return service.NewValidationError(
			"sslEnforcement",
			fmt.Sprintf(`invalid sslEnforcement option: "%s"`, pp.SSLEnforcement),
		)
	}

	startIP := net.ParseIP(pp.FirewallIPStart)
	if pp.FirewallIPStart != "" && net.ParseIP(pp.FirewallIPStart) == nil {
		return service.NewValidationError(
			"firewallStartIpAddress",
			fmt.Sprintf(`invalid firewallStartIpAddress option: "%s"`, pp.SSLEnforcement),
		)
	}

	endIP := net.ParseIP(pp.FirewallIPEnd)
	if pp.FirewallIPEnd != "" && net.ParseIP(pp.FirewallIPEnd) == nil {
		return service.NewValidationError(
			"firewallEndIpAddress",
			fmt.Sprintf(`invalid firewallEndIpAddress option: "%s"`, pp.SSLEnforcement),
		)
	}

	if startIP != nil && endIP == nil {
		return service.NewValidationError(
			"firewallEndIpAddress",
			fmt.Sprintf(`invalid firewallEndIpAddress option: "%s"`, pp.SSLEnforcement),
		)
	}

	startBytes := startIP.To4()
	endBytes := endIP.To4()
	if bytes.Compare(startBytes, endBytes) > 0 {
		return service.NewValidationError(
			"firewallEndIpAddress",
			fmt.Sprintf(`invalid firewallEndIpAddress option: "%s"`, pp.SSLEnforcement),
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
	_ service.Plan,
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*mysqlProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningContext as " +
				"*mysqlProvisioningContext",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*mysql.ProvisioningParameters",
		)
	}
	pc.ARMDeploymentName = uuid.NewV4().String()
	pc.ServerName = uuid.NewV4().String()
	pc.AdministratorLoginPassword = generate.NewPassword()
	pc.DatabaseName = generate.NewIdentifier()

	sslEnforcement := strings.ToLower(pp.SSLEnforcement)
	switch sslEnforcement {
	case "", "enabled":
		pc.EnforceSSL = true
	case "disabled":
		pc.EnforceSSL = false
	}

	return pc, nil
}

func buildARMTemplateParameters(
	plan service.Plan,
	provisioningContext *mysqlProvisioningContext,
	provisioningParameters *ProvisioningParameters,
) map[string]interface{} {
	var sslEnforcement string
	if provisioningContext.EnforceSSL {
		sslEnforcement = "Enabled"
	} else {
		sslEnforcement = "Disabled"
	}
	p := map[string]interface{}{ // ARM template params
		"administratorLoginPassword": provisioningContext.AdministratorLoginPassword,
		"serverName":                 provisioningContext.ServerName,
		"databaseName":               provisioningContext.DatabaseName,
		"skuName":                    plan.GetProperties().Extended["skuName"],
		"skuTier":                    plan.GetProperties().Extended["skuTier"],
		"skuCapacityDTU": plan.GetProperties().
			Extended["skuCapacityDTU"],
		"skuSizeMB":      plan.GetProperties().Extended["skuSizeMB"],
		"sslEnforcement": sslEnforcement,
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
	instance service.Instance,
	plan service.Plan,
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*mysqlProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningContext " +
				"as *mysqlProvisioningContext",
		)
	}
	armTemplateParameters := buildARMTemplateParameters(plan, pc, pp)
	outputs, err := s.armDeployer.Deploy(
		pc.ARMDeploymentName,
		instance.StandardProvisioningContext.ResourceGroup,
		instance.StandardProvisioningContext.Location,
		armTemplateBytes,
		nil, // Go template params
		armTemplateParameters,
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

	return pc, nil
}
