package postgresqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"

	log "github.com/Sirupsen/logrus"
)

func createBinding(
	enforceSSL bool,
	serverName string,
	administratorLoginPassword string,
	fullyQualifiedDomainName string,
	databaseName string,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	roleName := generate.NewIdentifier()
	password := generate.NewPassword()

	db, err := getDBConnection(
		enforceSSL,
		serverName,
		administratorLoginPassword,
		fullyQualifiedDomainName,
		primaryDB,
	)
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
		fmt.Sprintf("grant %s to %s", databaseName, roleName),
	); err != nil {
		return nil, nil, fmt.Errorf(
			`error adding role "%s" to role "%s": %s`,
			databaseName,
			roleName,
			err,
		)
	}
	if _, err = tx.Exec(
		fmt.Sprintf("alter role %s set role %s", roleName, databaseName),
	); err != nil {
		return nil, nil, fmt.Errorf(
			`error making "%s" the default role for "%s" sessions: %s`,
			databaseName,
			roleName,
			err,
		)
	}
	if err = tx.Commit(); err != nil {
		return nil, nil, fmt.Errorf("error committing transaction: %s", err)
	}

	return &postgresqlBindingDetails{
			LoginName: roleName,
		},
		&postgresqlSecureBindingDetails{
			Password: password,
		},
		nil
}
