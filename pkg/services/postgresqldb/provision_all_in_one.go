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

func (a *allInOneManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*AllInOneProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"*postgresql.AllInOneProvisioningParameters",
		)
	}
	return validateServerParameters(
		pp.SSLEnforcement,
		pp.FirewallIPStart,
		pp.FirewallIPEnd,
	)
}

func (a *allInOneManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", a.preProvision),
		service.NewProvisioningStep("deployARMTemplate", a.deployARMTemplate),
		service.NewProvisioningStep("setupDatabase", a.setupDatabase),
		service.NewProvisioningStep("createExtensions", a.createExtensions),
	)
}

func (a *allInOneManager) preProvision(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt, ok := instance.Details.(*allInOnePostgresqlInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details" +
				"as *allInOnePostgresqlInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*allInOnePostgresqlSecureInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails" +
				"as *allInOnePostgresqlSecureInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*AllInOneProvisioningParameters)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*postgresql.AllInOneProvisioningParameters",
		)
	}

	dt.ARMDeploymentName = uuid.NewV4().String()

	var err error
	if dt.ServerName, err = getAvailableServerName(
		ctx,
		a.checkNameAvailabilityClient,
	); err != nil {
		return nil, nil, err
	}

	sdt.AdministratorLoginPassword = generate.NewPassword()
	dt.DatabaseName = generate.NewIdentifier()

	sslEnforcement := strings.ToLower(pp.SSLEnforcement)
	switch sslEnforcement {
	case "", enabled:
		dt.EnforceSSL = true
	case disabled:
		dt.EnforceSSL = false
	}

	return dt, sdt, nil
}

func (a *allInOneManager) buildARMTemplateParameters(
	plan service.Plan,
	details *allInOnePostgresqlInstanceDetails,
	secureDetails *allInOnePostgresqlSecureInstanceDetails,
	provisioningParameters *AllInOneProvisioningParameters,
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
		"databaseName":               details.DatabaseName,
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

func (a *allInOneManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*allInOnePostgresqlInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details " +
				"as *allInOnePostgresqlInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*allInOnePostgresqlSecureInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails " +
				"as *allInOnePostgresqlSecureInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*AllInOneProvisioningParameters)
	if !ok {
		return nil, nil, errors.New(
			"error casting provisioningParameters as " +
				"*postgresql.AllInOneProvisioningParameters",
		)
	}
	armTemplateParameters := a.buildARMTemplateParameters(
		instance.Plan,
		dt,
		sdt,
		pp,
	)
	outputs, err := a.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
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

func (a *allInOneManager) setupDatabase(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*allInOnePostgresqlInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details " +
				"as *allInOnePostgresqlInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*allInOnePostgresqlSecureInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details " +
				"as *allInOnePostgresqlInstanceDetails",
		)
	}
	err := setupDatabase(
		dt.EnforceSSL,
		dt.ServerName,
		sdt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
	if err != nil {
		return nil, nil, err
	}
	return dt, sdt, nil
}

func (a *allInOneManager) createExtensions(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*allInOnePostgresqlInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as " +
				"*allInOnePostgresqlInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*allInOnePostgresqlSecureInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails as " +
				"*allInOnePostgresqlSecureInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*AllInOneProvisioningParameters)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*postgresql.AllInOneProvisioningParameters",
		)
	}

	if len(pp.Extensions) > 0 {
		err := createExtensions(
			dt.EnforceSSL,
			dt.ServerName,
			sdt.AdministratorLoginPassword,
			dt.FullyQualifiedDomainName,
			dt.DatabaseName,
			pp.Extensions,
		)
		if err != nil {
			return nil, nil, err
		}
	}
	return dt, instance.SecureDetails, nil
}
