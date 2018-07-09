package main

import (
	"fmt"
	"log"

	"github.com/Azure/open-service-broker-azure/pkg/client"
	"github.com/urfave/cli"
)

func unbind(c *cli.Context) error {
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
	bindingID := c.String(flagBindingID)
	if bindingID == "" {
		return fmt.Errorf("--%s is a required flag", flagBindingID)
	}
	err := client.Unbind(
		useSSL,
		skipCertValidation,
		host,
		port,
		username,
		password,
		instanceID,
		bindingID,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"\nUnbound binding %s to service instance %s\n\n",
		bindingID,
		instanceID,
	)
	return nil
}
