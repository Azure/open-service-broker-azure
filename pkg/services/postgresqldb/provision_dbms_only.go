package postgresqldb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *dbmsOnlyManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
	_ service.SecureProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ServerProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"*postgresql.ServerProvisioningParameters",
		)
	}
	return validateServerParameters(
		pp.SSLEnforcement,
		pp.FirewallIPStart,
		pp.FirewallIPEnd,
	)
}

func (d *dbmsOnlyManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("deployARMTemplate", d.deployARMTemplate),
	)
}

func (d *dbmsOnlyManager) preProvision(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt, ok := instance.Details.(*dbmsOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *dbmsOnlyPostgresqlInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*dbmsOnlyPostgresqlSecureInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails as " +
				"*dbmsOnlyPostgresqlSecureInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ServerProvisioningParameters)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*postgresql.ServerProvisioningParameters",
		)
	}

	dt.ARMDeploymentName = uuid.NewV4().String()

	var err error
	if dt.ServerName, err = getAvailableServerName(
		ctx,
		d.checkNameAvailabilityClient,
	); err != nil {
		return nil, nil, err
	}

	sdt.AdministratorLoginPassword = generate.NewPassword()

	sslEnforcement := strings.ToLower(pp.SSLEnforcement)
	switch sslEnforcement {
	case "", enabled:
		dt.EnforceSSL = true
	case disabled:
		dt.EnforceSSL = false
	}

	return dt, instance.SecureDetails, nil
}

func (d *dbmsOnlyManager) buildARMTemplateParameters(
	plan service.Plan,
	details *dbmsOnlyPostgresqlInstanceDetails,
	secureDetails *dbmsOnlyPostgresqlSecureInstanceDetails,
	provisioningParameters *ServerProvisioningParameters,
) map[string]interface{} {
	var sslEnforcement string
	if details.EnforceSSL {
		sslEnforcement = "Enabled"
	} else {
		sslEnforcement = "Disabled"
	}
	p := map[string]interface{}{ // ARM template params
		"administratorLoginPassword": secureDetails.AdministratorLoginPassword,
		"serverName":                 details.ServerName,
		"skuName":                    plan.GetProperties().Extended["skuName"],
		"skuTier":                    plan.GetProperties().Extended["skuTier"],
		"skuCapacityDTU": plan.GetProperties().
			Extended["skuCapacityDTU"],
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

func (d *dbmsOnlyManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*dbmsOnlyPostgresqlInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *dbmsOnlyPostgresqlInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*dbmsOnlyPostgresqlSecureInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails as " +
				"*dbmsOnlyPostgresqlSecureInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ServerProvisioningParameters)
	if !ok {
		return nil, nil, errors.New(
			"error casting provisioningParameters as " +
				"*postgresql.ServerProvisioningParameters",
		)
	}
	armTemplateParameters := d.buildARMTemplateParameters(
		instance.Plan,
		dt,
		sdt,
		pp,
	)
	outputs, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateDBMSOnlyBytes,
		nil, // Go template params
		armTemplateParameters,
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
