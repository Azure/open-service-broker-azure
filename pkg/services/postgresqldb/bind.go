package postgresqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
)

func (s *serviceManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to PostgreSQL, so there is nothing
	// to validate
	return nil
}

func (s *serviceManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	dt, ok := instance.Details.(*postgresqlInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *postgresqlInstanceDetails",
		)
	}

	roleName := generate.NewIdentifier()
	password := generate.NewPassword()

	db, err := getDBConnection(dt, primaryDB)
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
		fmt.Sprintf("grant %s to %s", dt.DatabaseName, roleName),
	); err != nil {
		return nil, nil, fmt.Errorf(
			`error adding role "%s" to role "%s": %s`,
			dt.DatabaseName,
			roleName,
			err,
		)
	}
	if _, err = tx.Exec(
		fmt.Sprintf("alter role %s set role %s", roleName, dt.DatabaseName),
	); err != nil {
		return nil, nil, fmt.Errorf(
			`error making "%s" the default role for "%s" sessions: %s`,
			dt.DatabaseName,
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
			Host:     dt.FullyQualifiedDomainName,
			Port:     5432,
			Database: dt.DatabaseName,
			Username: fmt.Sprintf("%s@%s", roleName, dt.ServerName),
			Password: password,
		},
		nil
}
