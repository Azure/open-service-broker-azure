package postgresqldb

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strings"

	postgresSDK "github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-04-30-preview/postgresql" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq" // Postgres SQL driver
	uuid "github.com/satori/go.uuid"
)

const (
	enabled  = "enabled"
	disabled = "disabled"
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

func validateServerParameters(
	sslEnforcementParam string,
	firewallRules []FirewallRule,
) error {
	sslEnforcement := strings.ToLower(sslEnforcementParam)
	if sslEnforcement != "" && sslEnforcement != enabled &&
		sslEnforcement != disabled {
		return service.NewValidationError(
			"sslEnforcement",
			fmt.Sprintf(`invalid option: "%s"`, sslEnforcementParam),
		)
	}
	for _, firewallRule := range firewallRules {
		if firewallRule.FirewallRuleName == "" {
			return service.NewValidationError(
				"firewallRuleName",
				"must be set",
			)
		}
		if firewallRule.FirewallIPStart != "" || firewallRule.FirewallIPEnd != "" {
			if firewallRule.FirewallIPStart == "" {
				return service.NewValidationError(
					"firewallStartIPAddress",
					"must be set when firewallEndIPAddress is set",
				)
			}
			if firewallRule.FirewallIPEnd == "" {
				return service.NewValidationError(
					"firewallEndIPAddress",
					"must be set when firewallStartIPAddress is set",
				)
			}
		}
		startIP := net.ParseIP(firewallRule.FirewallIPStart)
		if firewallRule.FirewallIPStart != "" && startIP == nil {
			return service.NewValidationError(
				"firewallStartIPAddress",
				fmt.Sprintf(`invalid value: "%s"`, firewallRule.FirewallIPStart),
			)
		}
		endIP := net.ParseIP(firewallRule.FirewallIPStart)
		if firewallRule.FirewallIPEnd != "" && endIP == nil {
			return service.NewValidationError(
				"firewallEndIPAddress",
				fmt.Sprintf(
					`invalid value: "%s"`,
					firewallRule.FirewallIPEnd,
				),
			)
		}
		//The net.IP.To4 method returns a 4 byte representation of an IPv4 address.
		//Once converted,comparing two IP addresses can be done by using the
		//bytes. Compare function. Per the ARM template documentation,
		//startIP must be <= endIP.
		startBytes := startIP.To4()
		endBytes := endIP.To4()
		if bytes.Compare(startBytes, endBytes) > 0 {
			return service.NewValidationError(
				"firewallEndIPAddress",
				fmt.Sprintf(`invalid value: "%s". must be 
				greater than or equal to firewallStartIPAddress`,
					firewallRule.FirewallIPEnd),
			)
		}
	}
	return nil
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
	provisioningParameters *ServerProvisioningParameters,
) map[string]interface{} {
	p := map[string]interface{}{}
	//Only include these if they are not empty.
	//ARM Deployer will fail if the values included are not
	//valid IPV4 addresses (i.e. empty string wil fail)
	if len(provisioningParameters.FirewallRules) > 0 {
		p["firewallRules"] = provisioningParameters.FirewallRules
	} else {
		//Build the azure default
		p["firewallRules"] = []FirewallRule{
			{
				FirewallRuleName: "AllowAzure",
				FirewallIPStart:  "0.0.0.0",
				FirewallIPEnd:    "0.0.0.0",
			},
		}
	}
	return p
}
