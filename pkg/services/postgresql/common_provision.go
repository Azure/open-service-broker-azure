package postgresql

import (
	"context"
	"fmt"

	postgresSDK "github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-04-30-preview/postgresql" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

func getAvailableServerName(
	ctx context.Context,
	checkNameAvailabilityClient postgresSDK.CheckNameAvailabilityClient,
) (string, error) {
	for {
		serverName := uuid.NewV4().String()
		nameAvailability, err := checkNameAvailabilityClient.Execute(
			ctx,
			postgresSDK.NameAvailabilityRequest{
				Name: &serverName,
			},
		)
		if err != nil {
			return "", fmt.Errorf(
				"error determining server name availability: %s",
				err,
			)
		}
		if *nameAvailability.NameAvailable {
			return serverName, nil
		}
	}
}

func setupDatabase(
	enforceSSL bool,
	serverName string,
	administratorLoginPassword string,
	fullyQualifiedDomainName string,
	dbName string,
) error {
	db, err := getDBConnection(
		enforceSSL,
		serverName,
		administratorLoginPassword,
		fullyQualifiedDomainName,
		primaryDB,
	)
	if err != nil {
		return err
	}
	defer db.Close() // nolint: errcheck

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %s", err)
	}
	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.WithField("error", err).Error("error rolling back transaction")
			}
		}
	}()
	if _, err = tx.Exec(
		fmt.Sprintf("create role %s", dbName),
	); err != nil {
		return fmt.Errorf(`error creating role "%s": %s`, dbName, err)
	}
	if _, err = tx.Exec(
		fmt.Sprintf("grant %s to postgres", dbName),
	); err != nil {
		return fmt.Errorf(
			`error adding role "%s" to role "postgres": %s`,
			dbName,
			err,
		)
	}
	if _, err = tx.Exec(
		fmt.Sprintf(
			"alter database %s owner to %s",
			dbName,
			dbName,
		),
	); err != nil {
		return fmt.Errorf(
			`error updating database owner"%s": %s`,
			dbName,
			err,
		)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %s", err)
	}

	return nil
}

func createExtensions(
	enforceSSL bool,
	serverName string,
	administratorLoginPassword string,
	fullyQualifiedDomainName string,
	dbName string,
	extensions []string,
) error {
	db, err := getDBConnection(
		enforceSSL,
		serverName,
		administratorLoginPassword,
		fullyQualifiedDomainName,
		dbName,
	)
	if err != nil {
		return err
	}
	defer db.Close() // nolint: errcheck

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %s", err)
	}
	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.WithField("error", err).Error("error rolling back transaction")
			}
		}
	}()
	for _, extension := range extensions {
		if _, err = tx.Exec(
			fmt.Sprintf(`create extension "%s"`, extension),
		); err != nil {
			return fmt.Errorf(
				`error creating extension "%s": %s`,
				extension,
				err,
			)
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %s", err)
	}
	return nil
}

func buildGoTemplateParameters(
	plan service.Plan,
	dt dbmsInstanceDetails,
	sdt secureDBMSInstanceDetails,
	pp service.Parameters,
) (map[string]interface{}, error) {
	p := map[string]interface{}{}
	p["sku"] = getSku(pp)
	p["tier"] = pp.GetString("tier")
	p["cores"] = pp.GetInt64("cores")
	p["storage"] = pp.GetInt64("storage") * 1024 // storage is in MB to arm :/
	p["backupRetention"] = pp.GetInt64("backupRetention")
	p["hardwareFamily"] = pp.GetString("hardwareFamily")
	if pp.GetString("geoRedundantBackup") == "enabled" {
		p["geoRedundantBackup"] = "Enabled"
	}
	p["version"] = pp.GetString("version")
	p["serverName"] = dt.ServerName
	p["administratorLoginPassword"] = sdt.AdministratorLoginPassword
	if pp.GetString("sslEnforcement") == "enabled" {
		p["sslEnforcement"] = "Enabled"
	} else {
		p["sslEnforcement"] = "Disabled"
	}
	p["firewallRules"] = pp.GetArray("firewallRules")
	return p, nil
}

func getSku(pp service.Parameters) string {
	// The name of the sku, typically:
	// tier + family + cores, e.g. B_Gen4_1, GP_Gen5_8.
	var tier string
	switch pp.GetString("tier") {
	case "Basic":
		tier = "B"
	case "GeneralPurpose":
		tier = "GP"
	case "MemoryOptimized":
		tier = "MO"
	}
	sku := fmt.Sprintf(
		"%s_%s_%d",
		tier,
		pp.GetString("hardwareFamily"),
		pp.GetInt64("cores"),
	)
	return sku
}
