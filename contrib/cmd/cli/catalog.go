package main

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/client"
	"github.com/urfave/cli"
)

func catalog(c *cli.Context) error {
	host := c.GlobalString(flagHost)
	port := c.GlobalInt(flagPort)
	username := c.GlobalString(flagUsername)
	password := c.GlobalString(flagPassword)
	catalog, err := client.GetCatalog(host, port, username, password)
	if err != nil {
		return fmt.Errorf("error retrieving catalog: %s", err)
	}
	fmt.Println()
	for _, svc := range catalog.Services {
		fmt.Printf("service: %-20s id: %s\n", svc.Name, svc.ID)
		for _, plan := range svc.Plans {
			fmt.Printf("   plan: %-20s id: %s\n", plan.Name, plan.ID)
		}
		fmt.Println()
	}
	return nil
}
