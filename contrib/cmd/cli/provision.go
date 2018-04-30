package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/client"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func provision(c *cli.Context) error {
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
	tags := make(map[string]string)
	rawTagStrs := c.StringSlice(flagTag)
	for _, rawTagStr := range rawTagStrs {
		rawTagStr = strings.TrimSpace(rawTagStr)
		tokens := strings.Split(rawTagStr, "=")
		if len(tokens) != 2 {
			return errors.New("tag string is incorrectly formatted")
		}
		key := strings.TrimSpace(tokens[0])
		value := strings.TrimSpace(tokens[1])
		tags[key] = value
	}
	instanceID, err := client.Provision(
		host,
		port,
		username,
		password,
		serviceID,
		planID,
		params,
		tags,
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
