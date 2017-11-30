package main

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/contrib/pkg/client"
	"github.com/Azure/open-service-broker-azure/pkg/api"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func poll(c *cli.Context) error {
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
	if operation != api.OperationProvisioning &&
		operation != api.OperationDeprovisioning &&
		operation != api.OperationUpdating {
		return fmt.Errorf("invalid value for flag --%s", flagOperation)
	}
	status, err := client.Poll(
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
