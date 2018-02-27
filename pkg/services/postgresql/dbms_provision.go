package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *dbmsManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
	_ service.SecureProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*DBMSProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"*postgresql.DBMSProvisioningParameters",
		)
	}
	return validateDBMSProvisionParameters(pp)
}

func (d *dbmsManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("deployARMTemplate", d.deployARMTemplate),
	)
}

func (d *dbmsManager) preProvision(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt, ok := instance.Details.(*dbmsInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *postgresql.dbmsInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*secureDBMSInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails as " +
				"*postgresql.secureDBMSInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*DBMSProvisioningParameters)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*postgresql.DBMSProvisioningParameters",
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

func (d *dbmsManager) buildARMTemplateParameters(
	plan service.Plan,
	details *dbmsInstanceDetails,
	secureDetails *secureDBMSInstanceDetails,
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
	return p
}

func (d *dbmsManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*dbmsInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *postgresql.dbmsInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*secureDBMSInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails as " +
				"*postgresql.secureDBMSInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*DBMSProvisioningParameters)
	if !ok {
		return nil, nil, errors.New(
			"error casting provisioningParameters as " +
				"*postgresql.DBMSProvisioningParameters",
		)
	}
	armTemplateParameters := d.buildARMTemplateParameters(
		instance.Plan,
		dt,
		sdt,
	)
	goTemplateParameters := buildGoTemplateParameters(pp)
	outputs, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		dbmsARMTemplateBytes,
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

	return dt, sdt, nil
}
