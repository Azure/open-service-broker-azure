package main

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/client"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func poll(c *cli.Context) error {
	useSSL := c.GlobalBool(flagSSL)
	skipCertValidation := c.GlobalBool(flagInsecure)
	host := c.GlobalString(flagHost)
	port := c.GlobalInt(flagPort)
	username := c.GlobalString(flagUsername)
	password := c.GlobalString(flagPassword)
	instanceID := c.String(flagInstanceID)
	if instanceID == "" {
		return fmt.Errorf("--%s is a required flag", flagInstanceID)
	}
	operation := c.String(flagOperation)
	if operation == "" {
		return fmt.Errorf("--%s is a required flag", flagOperation)
	}
	if operation != client.OperationProvisioning &&
		operation != client.OperationDeprovisioning &&
		operation != client.OperationUpdating {
		return fmt.Errorf("invalid value for flag --%s", flagOperation)
	}
	status, err := client.Poll(
		useSSL,
		skipCertValidation,
		host,
		port,
		username,
		password,
		instanceID,
		operation,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nInstance %s %s state: %s\n\n", instanceID, operation, status)
	return nil
}
