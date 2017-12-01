package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-service-broker/pkg/generate"
	"github.com/Azure/azure-service-broker/pkg/service"
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
			fmt.Sprintf(`invalid sslEnforcement option: "%s"`, pp.SSLEnforcement),
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
	_ string, // instanceID
	_ service.Plan,
	_ service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	provisioningParameters service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*postgresqlProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *postgresqlProvisioningContext",
		)
	}
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters as " +
				"*postgresql.ProvisioningParameters",
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

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	_ string, // instanceID
	plan service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	_ service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*postgresqlProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *postgresqlProvisioningContext",
		)
	}
	var sslEnforcement string
	if pc.EnforceSSL {
		sslEnforcement = "Enabled"
	} else {
		sslEnforcement = "Disabled"
	}
	outputs, err := s.armDeployer.Deploy(
		pc.ARMDeploymentName,
		standardProvisioningContext.ResourceGroup,
		standardProvisioningContext.Location,
		armTemplateBytes,
		nil, // Go template params
		map[string]interface{}{ // ARM template params
			"administratorLoginPassword": pc.AdministratorLoginPassword,
			"serverName":                 pc.ServerName,
			"databaseName":               pc.DatabaseName,
			"skuName":                    plan.GetProperties().Extended["skuName"],
			"skuTier":                    plan.GetProperties().Extended["skuTier"],
			"skuCapacityDTU": plan.GetProperties().
				Extended["skuCapacityDTU"],
			"sslEnforcement": sslEnforcement,
		},
		standardProvisioningContext.Tags,
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

func (s *serviceManager) setupDatabase(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
	_ service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	_ service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*postgresqlProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *postgresqlProvisioningContext",
		)
	}

	db, err := getDBConnection(pc, primaryDB)
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
		fmt.Sprintf("create role %s", pc.DatabaseName),
	); err != nil {
		return nil, fmt.Errorf(`error creating role "%s": %s`, pc.DatabaseName, err)
	}
	if _, err = tx.Exec(
		fmt.Sprintf("grant %s to postgres", pc.DatabaseName),
	); err != nil {
		return nil, fmt.Errorf(
			`error adding role "%s" to role "postgres": %s`,
			pc.DatabaseName,
			err,
		)
	}
	if _, err = tx.Exec(
		fmt.Sprintf(
			"alter database %s owner to %s",
			pc.DatabaseName,
			pc.DatabaseName,
		),
	); err != nil {
		return nil, fmt.Errorf(
			`error updating database owner"%s": %s`,
			pc.DatabaseName,
			err,
		)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %s", err)
	}

	return pc, nil
}

func (s *serviceManager) createExtensions(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
	_ service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	provisioningParameters service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*postgresqlProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *postgresqlProvisioningContext",
		)
	}
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters as " +
				"*postgresql.ProvisioningParameters",
		)
	}

	if len(pp.Extensions) > 0 {
		db, err := getDBConnection(pc, pc.DatabaseName)
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
	return pc, nil
}
