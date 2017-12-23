package postgresqldb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq" // Postgres SQL driver
	uuid "github.com/satori/go.uuid"
)

func (s *serviceManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"*postgresql.ProvisioningParameters",
		)
	}
	sslEnforcement := strings.ToLower(pp.SSLEnforcement)
	if sslEnforcement != "" && sslEnforcement != "enabled" &&
		sslEnforcement != "disabled" {
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

func (s *serviceManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", s.preProvision),
		service.NewProvisioningStep("deployARMTemplate", s.deployARMTemplate),
		service.NewProvisioningStep("setupDatabase", s.setupDatabase),
		service.NewProvisioningStep("createExtensions", s.createExtensions),
	)
}

func (s *serviceManager) preProvision(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*postgresqlInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *postgresqlInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*postgresql.ProvisioningParameters",
		)
	}

	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.ServerName = uuid.NewV4().String()
	dt.AdministratorLoginPassword = generate.NewPassword()
	dt.DatabaseName = generate.NewIdentifier()

	sslEnforcement := strings.ToLower(pp.SSLEnforcement)
	switch sslEnforcement {
	case "", "enabled":
		dt.EnforceSSL = true
	case "disabled":
		dt.EnforceSSL = false
	}

	return dt, nil
}

func buildARMTemplateParameters(
	plan service.Plan,
	details *postgresqlInstanceDetails,
	provisioningParameters *ProvisioningParameters,
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

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
	plan service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*postgresqlInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *postgresqlInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters as " +
				"*postgresql.ProvisioningParameters",
		)
	}
	armTemplateParameters := buildARMTemplateParameters(plan, dt, pp)
	outputs, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
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

func (s *serviceManager) setupDatabase(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*postgresqlInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *postgresqlInstanceDetails",
		)
	}

	db, err := getDBConnection(dt, primaryDB)
	if err != nil {
		return nil, err
	}
	defer db.Close() // nolint: errcheck

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %s", err)
	}
	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.WithField("error", err).Error("error rolling back transaction")
			}
		}
	}()
	if _, err = tx.Exec(
		fmt.Sprintf("create role %s", dt.DatabaseName),
	); err != nil {
		return nil, fmt.Errorf(`error creating role "%s": %s`, dt.DatabaseName, err)
	}
	if _, err = tx.Exec(
		fmt.Sprintf("grant %s to postgres", dt.DatabaseName),
	); err != nil {
		return nil, fmt.Errorf(
			`error adding role "%s" to role "postgres": %s`,
			dt.DatabaseName,
			err,
		)
	}
	if _, err = tx.Exec(
		fmt.Sprintf(
			"alter database %s owner to %s",
			dt.DatabaseName,
			dt.DatabaseName,
		),
	); err != nil {
		return nil, fmt.Errorf(
			`error updating database owner"%s": %s`,
			dt.DatabaseName,
			err,
		)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %s", err)
	}

	return dt, nil
}

func (s *serviceManager) createExtensions(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*postgresqlInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *postgresqlInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*postgresql.ProvisioningParameters",
		)
	}

	if len(pp.Extensions) > 0 {
		db, err := getDBConnection(dt, dt.DatabaseName)
		if err != nil {
			return nil, err
		}
		defer db.Close() // nolint: errcheck

		tx, err := db.Begin()
		if err != nil {
			return nil, fmt.Errorf("error starting transaction: %s", err)
		}
		defer func() {
			if err != nil {
				if err = tx.Rollback(); err != nil {
					log.WithField("error", err).Error("error rolling back transaction")
				}
			}
		}()
		for _, extension := range pp.Extensions {
			if _, err = tx.Exec(
				fmt.Sprintf(`create extension "%s"`, extension),
			); err != nil {
				return nil, fmt.Errorf(
					`error creating extension "%s": %s`,
					extension,
					err,
				)
			}
		}
		if err = tx.Commit(); err != nil {
			return nil, fmt.Errorf("error committing transaction: %s", err)
		}
	}
	return dt, nil
}
