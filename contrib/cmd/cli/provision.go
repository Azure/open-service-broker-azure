package main

import (
	"fmt"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/client"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func provision(c *cli.Context) error {
	useSSL := c.GlobalBool(flagSSL)
	skipCertValidation := c.GlobalBool(flagInsecure)
	host := c.GlobalString(flagHost)
	port := c.GlobalInt(flagPort)
	username := c.GlobalString(flagUsername)
	password := c.GlobalString(flagPassword)
	serviceID := c.String(flagServiceID)
	if serviceID == "" {
		return fmt.Errorf("--%s is a required flag", flagServiceID)
	}
	planID := c.String(flagPlanID)
	if planID == "" {
		return fmt.Errorf("--%s is a required flag", flagPlanID)
	}
	params, err := parseParams(c)
	if err != nil {
		return err
	}
	instanceID, err := client.Provision(
		useSSL,
		skipCertValidation,
		host,
		port,
		username,
		password,
		serviceID,
		planID,
		params,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nProvisioning service instance %s\n\n", instanceID)
	if c.Bool(flagPoll) {
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()
		for range ticker.C {
			result, err := client.Poll(
				useSSL,
				skipCertValidation,
				host,
				port,
				username,
				password,
				instanceID,
				client.OperationProvisioning,
			)
			if err != nil {
				return fmt.Errorf("error polling for provisioning status: %s", err)
			}
			switch result {
			case client.OperationStateInProgress:
				fmt.Print(".")
			case client.OperationStateSucceeded:
				fmt.Printf(
					"\n\nService instance %s has been successfully provisioned\n\n",
					instanceID,
				)
				return nil
			case client.OperationStateFailed:
				return fmt.Errorf(
					"Provisioning service instance %s has failed",
					instanceID,
				)
			default:
				return fmt.Errorf("Unrecognized operation status: %s", result)
			}
		}
	}
	return nil
}
