package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-service-broker/pkg/azure"
	"github.com/Azure/azure-service-broker/pkg/generate"
	"github.com/Azure/azure-service-broker/pkg/service"
	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq" // Postgres SQL driver
	uuid "github.com/satori/go.uuid"
)

func (m *module) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"*postgresql.ProvisioningParameters",
		)
	}
	if !azure.IsValidLocation(pp.Location) {
		return service.NewValidationError(
			"location",
			fmt.Sprintf(`invalid location: "%s"`, pp.Location),
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

func (m *module) GetProvisioner(string, string) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", m.preProvision),
		service.NewProvisioningStep("deployARMTemplate", m.deployARMTemplate),
		service.NewProvisioningStep("setupDatabase", m.setupDatabase),
		service.NewProvisioningStep("createExtensions", m.createExtensions),
	)
}

func (m *module) preProvision(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string, // nolint: unparam
	planID string, // nolint: unparam
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
	if pp.ResourceGroup != "" {
		pc.ResourceGroupName = pp.ResourceGroup
	} else {
		pc.ResourceGroupName = uuid.NewV4().String()
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

func (m *module) deployARMTemplate(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string,
	planID string,
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
	catalog, err := m.GetCatalog()
	if err != nil {
		return nil, fmt.Errorf("error retrieving catalog: %s", err)
	}
	svc, ok := catalog.GetService(serviceID)
	if !ok {
		return nil, fmt.Errorf(
			`service "%s" not found in the "%s" module catalog`,
			serviceID,
			m.GetName(),
		)
	}
	plan, ok := svc.GetPlan(planID)
	if !ok {
		return nil, fmt.Errorf(
			`plan "%s" not found for service "%s"`,
			planID,
			serviceID,
		)
	}
	var sslEnforcement string
	if pc.EnforceSSL {
		sslEnforcement = "Enabled"
	} else {
		sslEnforcement = "Disabled"
	}
	outputs, err := m.armDeployer.Deploy(
		pc.ARMDeploymentName,
		pc.ResourceGroupName,
		pp.Location,
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
		pp.Tags,
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

func (m *module) setupDatabase(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string, // nolint: unparam
	planID string, // nolint: unparam
	provisioningContext service.ProvisioningContext,
	provisioningParameters service.ProvisioningParameters, // nolint: unparam
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

func (m *module) createExtensions(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string, // nolint: unparam
	planID string, // nolint: unparam
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
