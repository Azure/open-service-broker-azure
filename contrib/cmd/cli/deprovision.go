package main

import (
	"fmt"
	"log"

	"github.com/Azure/azure-service-broker/contrib/pkg/client"
	"github.com/urfave/cli"
)

func deprovision(c *cli.Context) error {
	host := c.GlobalString(flagHost)
	port := c.GlobalInt(flagPort)
	instanceID := c.String(flagInstanceID)
	if instanceID == "" {
		return fmt.Errorf("--%s is a required flag", flagInstanceID)
	}
	err := client.Deprovision(host, port, instanceID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nDeprovisioning service instance %s\n\n", instanceID)
	return nil
}
