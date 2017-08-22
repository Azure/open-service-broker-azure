package mysql

import (
	"crypto/tls"
	"database/sql"
	"fmt"

	az "github.com/Azure/azure-service-broker/pkg/azure"
	"github.com/Azure/go-autorest/autorest/azure"
	log "github.com/Sirupsen/logrus"
	"github.com/go-sql-driver/mysql"
)

func getDBConnection(pc *mysqlProvisioningContext) (*sql.DB, error) {

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
	).Info("Azure ENV SQLDatabaseDNSSuffix")

	err = mysql.RegisterTLSConfig("custom", &tls.Config{
		ServerName: serverName,
	})
	if err != nil {
		return nil, fmt.Errorf("error registering tlsconfig"+
			" for the database: %s", err)
	}

	db, err := sql.Open("mysql", fmt.Sprintf(
		"azureuser@%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true&tls=custom",
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
