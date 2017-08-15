package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/azure"
	"github.com/Azure/azure-service-broker/pkg/service"
	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq" // Postgres SQL driver
	uuid "github.com/satori/go.uuid"
)

const (
	identifierLength = 10
	identifierChars  = lowerAlphaChars + numberChars
)

func (m *module) ValidateProvisioningParameters(
	provisioningParameters interface{},
) error {
	pp, ok := provisioningParameters.(*postgresqlProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"postgresqlProvisioningParameters",
		)
	}
	if !azure.IsValidLocation(pp.Location) {
		return service.NewValidationError(
			"location",
			fmt.Sprintf(`invalid location: "%s"`, pp.Location),
		)
	}
	return nil
}

func (m *module) GetProvisioner(string, string) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", m.preProvision),
		service.NewProvisioningStep("deployARMTemplate", m.deployARMTemplate),
		service.NewProvisioningStep("setupDatabase", m.setupDatabase),
	)
}

func (m *module) preProvision(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string, // nolint: unparam
	planID string, // nolint: unparam
	provisioningContext interface{},
	provisioningParameters interface{}, // nolint: unparam
) (interface{}, error) {
	pc, ok := provisioningContext.(*postgresqlProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as postgresqlProvisioningContext",
		)
	}
	pc.ResourceGroupName = uuid.NewV4().String()
	pc.ARMDeploymentName = uuid.NewV4().String()
	pc.ServerName = uuid.NewV4().String()
	pc.AdministratorLoginPassword = generatePassword()
	pc.DatabaseName = generateIdentifier()
	return pc, nil
}

func (m *module) deployARMTemplate(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string,
	planID string,
	provisioningContext interface{},
	provisioningParameters interface{},
) (interface{}, error) {
	pc, ok := provisioningContext.(*postgresqlProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as postgresqlProvisioningContext",
		)
	}
	pp, ok := provisioningParameters.(*postgresqlProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters as " +
				"postgresqlProvisioningParameters",
		)
	}
	catalog, err := m.GetCatalog()
	if err != nil {
		return nil, fmt.Errorf("error retrieving catalog: %s", err)
	}
	service, ok := catalog.GetService(serviceID)
	if !ok {
		return nil, fmt.Errorf(
			`service "%s" not found in the "%s" module catalog`,
			serviceID,
			m.GetName(),
		)
	}
	plan, ok := service.GetPlan(planID)
	if !ok {
		return nil, fmt.Errorf(
			`plan "%s" not found for service "%s"`,
			planID,
			serviceID,
		)
	}
	outputs, err := m.armDeployer.Deploy(
		pc.ARMDeploymentName,
		pc.ResourceGroupName,
		pp.Location,
		armTemplateBytes,
		map[string]interface{}{
			"administratorLoginPassword": pc.AdministratorLoginPassword,
			"serverName":                 pc.ServerName,
			"databaseName":               pc.DatabaseName,
			"skuName":                    plan.GetProperties().Extended["skuName"],
			"skuTier":                    plan.GetProperties().Extended["skuTier"],
			"skuCapacityDTU": plan.GetProperties().
				Extended["skuCapacityDTU"],
		},
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
	provisioningContext interface{},
	provisioningParameters interface{}, // nolint: unparam
) (interface{}, error) {
	pc, ok := provisioningContext.(*postgresqlProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as postgresqlProvisioningContext",
		)
	}

	db, err := getDBConnection(pc)
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
