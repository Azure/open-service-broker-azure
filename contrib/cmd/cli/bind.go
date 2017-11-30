package main

import (
	"fmt"
	"log"

	"github.com/Azure/open-service-broker-azure/contrib/pkg/client"
	"github.com/urfave/cli"
)

func bind(c *cli.Context) error {
	host := c.GlobalString(flagHost)
	port := c.GlobalInt(flagPort)
	username := c.GlobalString(flagUsername)
	password := c.GlobalString(flagPassword)
	instanceID := c.String(flagInstanceID)
	if instanceID == "" {
		return fmt.Errorf("--%s is a required flag", flagInstanceID)
	}
	params, err := parseParams(c)
	if err != nil {
		return err
	}
	bindingID, credentialMap, err := client.Bind(
		host,
		port,
		username,
		password,
		instanceID,
		params,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"\nBinding %s created for service instance %s\n",
		bindingID,
		instanceID,
	)
	fmt.Println("Credentials:")
	for k, v := range credentialMap {
		fmt.Printf(
			"   %-20s %v\n",
			fmt.Sprintf("%s:", k),
			v,
		)
	}
	fmt.Println()
	return nil
}
