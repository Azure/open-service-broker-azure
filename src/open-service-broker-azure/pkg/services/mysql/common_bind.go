package mysql

import (
	"fmt"
	"net/url"

	"open-service-broker-azure/pkg/generate"
	"open-service-broker-azure/pkg/service"
)

func createBinding(
	enforceSSL bool,
	dnsSuffix string,
	serverName string,
	adminPassword string,
	fqdn string,
	databaseName string,
) (service.BindingDetails, service.SecureBindingDetails, error) {

	userName := generate.NewIdentifier()
	password := generate.NewPassword()

	db, err := createDBConnection(
		enforceSSL,
		dnsSuffix,
		serverName,
		adminPassword,
		fqdn,
		databaseName,
	)
	if err != nil {
		return nil, nil, err
	}
	defer db.Close() // nolint: errcheck

	// Open doesn't open a connection. Validate DSN data:
	if err = db.Ping(); err != nil {
		return nil, nil, err
	}

	if _, err = db.Exec(
		fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED BY '%s'", userName, password),
	); err != nil {
		return nil, nil, fmt.Errorf(
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
			databaseName, userName)); err != nil {
		return nil, nil, fmt.Errorf(
			`error granting permission to "%s": %s`,
			userName,
			err,
		)
	}

	bd := bindingDetails{
		LoginName: userName,
	}
	sbd := secureBindingDetails{
		Password: password,
	}

	bdMap, err := service.GetMapFromStruct(bd)
	if err != nil {
		return nil, nil, err
	}
	sbdMap, err := service.GetMapFromStruct(sbd)
	return bdMap, sbdMap, err
}

// Create a credential to be returned for binding purposes. This includes a CF
// compatible uri string and a flag to indicate if this connection should
// use ssl. URI is built with the username passed to url.QueryEscape to escape
// the @ in the username
func createCredential(
	fqdn string,
	sslRequired bool,
	serverName string,
	databaseName string,
	bindingDetails bindingDetails,
	secureBidningDetails secureBindingDetails,
) credentials {
	username := fmt.Sprintf("%s@%s", bindingDetails.LoginName, serverName)
	connectionTemplate := "mysql://%s:%s@%s:3306/%s?useSSL=true&requireSSL=true"
	connectionString := fmt.Sprintf(
		connectionTemplate,
		url.QueryEscape(username),
		secureBidningDetails.Password,
		fqdn,
		databaseName,
	)
	return credentials{
		Host:        fqdn,
		Port:        3306,
		Database:    databaseName,
		Username:    username,
		Password:    secureBidningDetails.Password,
		SSLRequired: sslRequired,
		URI:         connectionString,
		Tags:        []string{"mysql"},
	}
}
