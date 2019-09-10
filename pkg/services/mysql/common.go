package mysql

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
	username string,
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
			"%s@%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true&tls=custom"
	} else {
		connectionStrTemplate =
			"%s@%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true"
	}

	db, err := sql.Open("mysql", fmt.Sprintf(
		connectionStrTemplate,
		username,
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
