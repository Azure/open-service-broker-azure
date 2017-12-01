package mysql

import (
	"crypto/tls"
	"database/sql"
	"fmt"

	"github.com/Azure/go-autorest/autorest/azure"
	az "github.com/Azure/open-service-broker-azure/pkg/azure"
	log "github.com/Sirupsen/logrus"
	"github.com/go-sql-driver/mysql"
)

func getDBConnection(pc *mysqlProvisioningContext) (*sql.DB, error) {
	var connectionStrTemplate string
	if pc.EnforceSSL {
		azureConfig, err := az.GetConfig()
		if err != nil {
			return nil, err
		}
		azureEnvironment, err := azure.EnvironmentFromName(azureConfig.Environment)
		if err != nil {
			return nil, err
		}
		sqlDatabaseDNSSuffix := azureEnvironment.SQLDatabaseDNSSuffix
		serverName := fmt.Sprintf("*.%s", sqlDatabaseDNSSuffix)

		log.WithField(
			"serverName", serverName,
		).Debug("Azure ENV SQLDatabaseDNSSuffix")

		err = mysql.RegisterTLSConfig("custom", &tls.Config{
			ServerName: serverName,
		})
		if err != nil {
			return nil, fmt.Errorf("error registering tlsconfig"+
				" for the database: %s", err)
		}
		connectionStrTemplate =
			"azureuser@%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true&tls=custom"
	} else {
		connectionStrTemplate =
			"azureuser@%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true"
	}

	db, err := sql.Open("mysql", fmt.Sprintf(
		connectionStrTemplate,
		pc.ServerName,
		pc.AdministratorLoginPassword,
		pc.FullyQualifiedDomainName,
		pc.DatabaseName,
	))
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %s", err)
	}
	return db, err
}
