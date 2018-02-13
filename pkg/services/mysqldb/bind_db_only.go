package mysqldb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbOnlyManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to MySQL, so there is nothing
	// to validate
	return nil
}

func (d *dbOnlyManager) Bind(
	instance service.Instance,
	_ service.BindingParameters,
) (service.BindingDetails, error) {
	pdt, ok := instance.Parent.Details.(*dbmsOnlyMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Parent.Details " +
				"as *dbmsOnlyMysqlInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*dbOnlyMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *dbOnlyMysqlInstanceDetails",
		)
	}

	userName := generate.NewIdentifier()
	password := generate.NewPassword()

	db, err := d.getDBConnection(pdt, dt)
	if err != nil {
		return nil, err
	}
	defer db.Close() // nolint: errcheck

	// Open doesn't open a connection. Validate DSN data:
	if err = db.Ping(); err != nil {
		return nil, err
	}

	if _, err = db.Exec(
		fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED BY '%s'", userName, password),
	); err != nil {
		return nil, fmt.Errorf(
			`error creating user "%s": %s`,
			userName,
			err,
		)
	}

	if _, err = db.Exec(
		fmt.Sprintf("GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, "+
			"INDEX, ALTER, CREATE TEMPORARY TABLES, LOCK TABLES, "+
			"CREATE VIEW, SHOW VIEW, CREATE ROUTINE, ALTER ROUTINE, "+
			"EXECUTE, REFERENCES, EVENT, "+
			"TRIGGER ON %s.* TO '%s'@'%%'",
			dt.DatabaseName, userName)); err != nil {
		return nil, fmt.Errorf(
			`error granting permission to "%s": %s`,
			userName,
			err,
		)
	}

	return &mysqlBindingDetails{
		LoginName: userName,
		Password:  password,
	}, nil
}

func (d *dbOnlyManager) GetCredentials(
	instance service.Instance,
	binding service.Binding,
) (service.Credentials, error) {
	pdt, ok := instance.Parent.Details.(*dbmsOnlyMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Parent.Details " +
				"as *dbmsOnlyMysqlInstanceDetails",
		)
	}
	dt, ok := instance.Details.(*dbOnlyMysqlInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *dbOnlyMysqlInstanceDetails",
		)
	}
	bd, ok := binding.Details.(*mysqlBindingDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting binding.Details as *mysqlBindingDetails",
		)
	}

	return &Credentials{
		Host:     pdt.FullyQualifiedDomainName,
		Port:     3306,
		Database: dt.DatabaseName,
		Username: fmt.Sprintf("%s@%s", bd.LoginName, pdt.ServerName),
		Password: bd.Password,
	}, nil
}
