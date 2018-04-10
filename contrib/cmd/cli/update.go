package main

import (
	"fmt"
	"time"

	"github.com/Azure/open-service-broker-azure/contrib/pkg/client"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func update(c *cli.Context) error {
	host := c.GlobalString(flagHost)
	port := c.GlobalInt(flagPort)
	username := c.GlobalString(flagUsername)
	password := c.GlobalString(flagPassword)
	instanceID := c.String(flagInstanceID)
	if instanceID == "" {
		return fmt.Errorf("--%s is a required flag", flagInstanceID)
	}
	serviceID := c.String(flagServiceID)
	if serviceID == "" {
		return fmt.Errorf("--%s is a required flag", flagServiceID)
	}
	planID := c.String(flagPlanID)
	params, err := parseParams(c)
	if err != nil {
		return err
	}
	if err := client.Update(
		host,
		port,
		username,
		password,
		instanceID,
		serviceID,
		planID,
		params,
	); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nUpdating service instance %s\n\n", instanceID)
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
				client.OperationUpdating,
			)
			if err != nil {
				return fmt.Errorf("error polling for updating status: %s", err)
			}
			switch result {
			case client.OperationStateInProgress:
				fmt.Print(".")
			case client.OperationStateSucceeded:
				fmt.Printf(
					"\n\nService instance %s has been successfully updated\n\n",
					instanceID,
				)
				return nil
			case client.OperationStateFailed:
				return fmt.Errorf(
					"Updating service instance %s has failed",
					instanceID,
				)
			default:
				return fmt.Errorf("Unrecognized operation status: %s", result)
			}
		}
	}
	return nil
}
