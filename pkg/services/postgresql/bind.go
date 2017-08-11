package postgresql

import (
	"fmt"
)

const (
	passwordLength = 16
	passwordChars  = lowerAlphaChars + upperAlphaChars + numberChars
)

func (m *module) ValidateBindingParameters(
	bindingParameters interface{},
) error {
	// There are no parameters for binding to PostgreSQL, so there is nothing
	// to validate
	return nil
}

func (m *module) Bind(
	provisioningContext interface{},
	bindingParameters interface{},
) (interface{}, interface{}, error) {
	pc, ok := provisioningContext.(*postgresqlProvisioningContext)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting provisioningContext as postgresqlProvisioningContext",
		)
	}

	roleName := generateIdentifier()
	password := generatePassword()

	db, err := getDBConnection(pc)
	if err != nil {
		return nil, nil, err
	}
	defer db.Close() // nolint: errcheck

	tx, err := db.Begin()
	if err != nil {
		return nil, nil, fmt.Errorf("error starting transaction: %s", err)
	}
	_, err = tx.Exec(
		fmt.Sprintf("create role %s with password '%s' login", roleName, password),
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			`error creating role "%s": %s`,
			roleName,
			err,
		)
	}
	_, err = tx.Exec(
		fmt.Sprintf("grant %s to %s", pc.DatabaseName, roleName),
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			`error adding role "%s" to role "%s": %s`,
			pc.DatabaseName,
			roleName,
			err,
		)
	}
	_, err = tx.Exec(
		fmt.Sprintf("alter role %s set role %s", roleName, pc.DatabaseName),
	)
	if err != nil {
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
		&postgresqlCredentials{
			Host:     pc.FullyQualifiedDomainName,
			Port:     5432,
			Database: pc.DatabaseName,
			Username: roleName,
			Password: password,
		},
		nil
}
