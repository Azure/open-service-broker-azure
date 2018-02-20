package sqldb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

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
	return validateServerProvisionParameters(pp)
}

func (v *vmOnlyManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", v.preProvision),
		service.NewProvisioningStep("deployARMTemplate", v.deployARMTemplate),
	)
}

func (v *vmOnlyManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlVMOnlyInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *mssqlVMOnlyInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*mssqlVMOnlySecureInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails as *mssqlVMOnlySecureInstanceDetails",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.ServerName = uuid.NewV4().String()
	dt.AdministratorLogin = generate.NewIdentifier()
	sdt.AdministratorLoginPassword = generate.NewPassword()
	return dt, sdt, nil
}
func (v *vmOnlyManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*mssqlVMOnlyInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *mssqlVMOnlyInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*mssqlVMOnlySecureInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails as " +
				"*mssqlVMOnlySecureInstanceDetails",
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
