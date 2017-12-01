package main

import (
	"fmt"
	"log"

	"github.com/Azure/open-service-broker-azure/contrib/pkg/client"
	"github.com/urfave/cli"
)

func unbind(c *cli.Context) error {
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
	err := client.Unbind(host, port, username, password, instanceID, bindingID)
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
