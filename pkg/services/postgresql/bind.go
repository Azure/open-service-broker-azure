package postgresql

import (
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/generate"
	"github.com/Azure/azure-service-broker/pkg/service"
	log "github.com/Sirupsen/logrus"
)

func (m *module) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to PostgreSQL, so there is nothing
	// to validate
	return nil
}

func (m *module) Bind(
	provisioningContext service.ProvisioningContext,
	bindingParameters service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	pc, ok := provisioningContext.(*postgresqlProvisioningContext)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting provisioningContext as *postgresqlProvisioningContext",
		)
	}

	roleName := generate.NewIdentifier()
	password := generate.NewPassword()

	db, err := getDBConnection(pc)
	if err != nil {
		return nil, nil, err
	}
	defer db.Close() // nolint: errcheck

	tx, err := db.Begin()
	if err != nil {
		return nil, nil, fmt.Errorf("error starting transaction: %s", err)
	}
	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.WithField("error", err).Error("error rolling back transaction")
			}
		}
	}()
	if _, err = tx.Exec(
		fmt.Sprintf("create role %s with password '%s' login", roleName, password),
	); err != nil {
		return nil, nil, fmt.Errorf(
			`error creating role "%s": %s`,
			roleName,
			err,
		)
	}
	if _, err = tx.Exec(
		fmt.Sprintf("grant %s to %s", pc.DatabaseName, roleName),
	); err != nil {
		return nil, nil, fmt.Errorf(
			`error adding role "%s" to role "%s": %s`,
			pc.DatabaseName,
			roleName,
			err,
		)
	}
	if _, err = tx.Exec(
		fmt.Sprintf("alter role %s set role %s", roleName, pc.DatabaseName),
	); err != nil {
		return nil, nil, fmt.Errorf(
			`error making "%s" the default role for "%s" sessions: %s`,
			pc.DatabaseName,
			roleName,
			err,
		)
	}
	if err = tx.Commit(); err != nil {
		return nil, nil, fmt.Errorf("error committing transaction: %s", err)
	}

	return &postgresqlBindingContext{
			LoginName: roleName,
		},
		&Credentials{
			Host:     pc.FullyQualifiedDomainName,
			Port:     5432,
			Database: pc.DatabaseName,
			Username: fmt.Sprintf("%s@%s", roleName, pc.ServerName),
			Password: password,
		},
		nil
}
