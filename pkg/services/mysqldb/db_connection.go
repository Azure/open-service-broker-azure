package mysqldb

import (
	"crypto/tls"
	"database/sql"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/go-sql-driver/mysql"
)

func (s *serviceManager) getDBConnection(
	dt *mysqlInstanceDetails,
) (*sql.DB, error) {
	var connectionStrTemplate string
	if dt.EnforceSSL {
		serverName := fmt.Sprintf("*.%s", s.sqlDatabaseDNSSuffix)

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
		dt.ServerName,
		dt.AdministratorLoginPassword,
		dt.FullyQualifiedDomainName,
		dt.DatabaseName,
	))
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %s", err)
	}
	return db, err
}
