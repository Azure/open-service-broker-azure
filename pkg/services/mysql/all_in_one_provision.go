package mysql

import (
	"context"
	"errors"
	"fmt"
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

func (a *allInOneManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
	_ service.SecureProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*AllInOneProvisioningParameters)
	if !ok {
		return fmt.Errorf(
			"error casting provisioningParameters as " +
				"*mysql.AllInOneProvisioningParameters",
		)
	}
	return validateDBMSProvisionParameters(&pp.DBMSProvisioningParameters)
}

func (a *allInOneManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", a.preProvision),
		service.NewProvisioningStep("deployARMTemplate", a.deployARMTemplate),
	)
}

func (a *allInOneManager) preProvision(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt, ok := instance.Details.(*allInOneInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *mysql.allInOneInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*secureAllInOneInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails as " +
				"*mysql.secureAllInOneInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*AllInOneProvisioningParameters)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*mysql.AllInOneProvisioningParameters",
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
	details *allInOneInstanceDetails,
	secureDetails *secureAllInOneInstanceDetails,
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
		"skuSizeMB":      plan.GetProperties().Extended["skuSizeMB"],
		"sslEnforcement": sslEnforcement,
	}
	return p
}

func (a *allInOneManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*allInOneInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *mysql.allInOneInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*secureAllInOneInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails as " +
				"*mysql.secureAllInOneInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*AllInOneProvisioningParameters)
	if !ok {
		return nil, nil, errors.New(
			"error casting provisioningParameters as " +
				"*mysql.AllInOneProvisioningParameters",
		)
	}
	armTemplateParameters := a.buildARMTemplateParameters(
		instance.Plan,
		dt,
		sdt,
	)
	goTemplateParameters :=
		buildGoTemplateParameters(&pp.DBMSProvisioningParameters)
	outputs, err := a.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		allInOneARMTemplateBytes,
		goTemplateParameters,
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

	return dt, instance.SecureDetails, nil
}
