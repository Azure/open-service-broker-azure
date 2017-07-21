package main

import (
	"fmt"

	"github.com/Azure/azure-service-broker/contrib/pkg/client"
	"github.com/urfave/cli"
)

func catalog(c *cli.Context) error {
	host := c.GlobalString(flagHost)
	port := c.GlobalInt(flagPort)
	catalog, err := client.GetCatalog(host, port)
	if err != nil {
		return fmt.Errorf("error retrieving catalog: %s", err)
	}
	fmt.Println()
	for _, svc := range catalog.GetServices() {
		fmt.Printf("service: %-20s id: %s\n", svc.GetName(), svc.GetID())
		for _, plan := range svc.GetPlans() {
			fmt.Printf("   plan: %-20s id: %s\n", plan.GetName(), plan.GetID())
		}
		fmt.Println()
	}
	return nil
}
