package mysqldb

import (
	"crypto/tls"
	"database/sql"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/go-sql-driver/mysql"
)

func createDBConnection(
	enforceSSL bool,
	sqlDatabaseDNSSuffix string,
	server string,
	password string,
	fqdn string,
	dbname string,
) (*sql.DB, error) {
	var connectionStrTemplate string
	if enforceSSL {
		serverName := fmt.Sprintf("*.%s", sqlDatabaseDNSSuffix)

		log.WithField(
			"serverName", serverName,
		).Debug("Azure ENV SQLDatabaseDNSSuffix")

		err := mysql.RegisterTLSConfig("custom", &tls.Config{
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
		server,
		password,
		fqdn,
		dbname,
	))
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %s", err)
	}
	return db, err
}
func (s *allInOneManager) getDBConnection(
	dt *allInOneMysqlInstanceDetails,
) (*sql.DB, error) {
	return createDBConnection(
		dt.EnforceSSL,
		s.sqlDatabaseDNSSuffix,
		dt.ServerName,
		dt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
}

func (d *dbOnlyManager) getDBConnection(
	pdt *dbmsOnlyMysqlInstanceDetails,
	dt *dbOnlyMysqlInstanceDetails,
) (*sql.DB, error) {
	return createDBConnection(
		pdt.EnforceSSL,
		d.sqlDatabaseDNSSuffix,
		pdt.ServerName,
		pdt.AdministratorLoginPassword,
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
}
