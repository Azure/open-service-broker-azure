package main

import (
	"fmt"
	"time"

	"github.com/Azure/azure-service-broker/contrib/pkg/client"
	"github.com/Azure/azure-service-broker/pkg/api"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func deprovision(c *cli.Context) error {
	host := c.GlobalString(flagHost)
	port := c.GlobalInt(flagPort)
	username := c.GlobalString(flagUsername)
	password := c.GlobalString(flagPassword)
	instanceID := c.String(flagInstanceID)
	if instanceID == "" {
		return fmt.Errorf("--%s is a required flag", flagInstanceID)
	}
	err := client.Deprovision(host, port, username, password, instanceID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nDeprovisioning service instance %s\n\n", instanceID)
	if c.Bool(flagPoll) {
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()
		for range ticker.C {
			result, err := client.Poll(
				host,
				port,
				username,
				password,
				instanceID,
				api.OperationDeprovisioning,
			)
			if err != nil {
				return fmt.Errorf("error polling for deprovisioning status: %s", err)
			}
			switch result {
			case api.OperationStateInProgress:
				fmt.Print(".")
			case api.OperationStateGone:
				fmt.Printf(
					"\n\nService instance %s has been successfully deprovisioned\n\n",
					instanceID,
				)
				return nil
			case api.OperationStateFailed:
				return fmt.Errorf(
					"Deprovisioning service instance %s has failed",
					instanceID,
				)
			default:
				return fmt.Errorf("Unrecognized operation status: %s", result)
			}
		}
	}
	return nil
}
