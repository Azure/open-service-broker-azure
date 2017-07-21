package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-service-broker/contrib/pkg/client"
	"github.com/urfave/cli"
)

func provision(c *cli.Context) error {
	host := c.GlobalString(flagHost)
	port := c.GlobalInt(flagPort)
	serviceID := c.String(flagServiceID)
	if serviceID == "" {
		return fmt.Errorf("--%s is a required flag", flagServiceID)
	}
	planID := c.String(flagPlanID)
	if planID == "" {
		return fmt.Errorf("--%s is a required flag", flagPlanID)
	}
	params := make(map[string]string)
	rawParamStrs := c.StringSlice(flagParameter)
	for _, rawParamStr := range rawParamStrs {
		rawParamStr = strings.TrimSpace(rawParamStr)
		tokens := strings.Split(rawParamStr, "=")
		if len(tokens) != 2 {
			return errors.New("parameter string is incorrectly formatted")
		}
		key := strings.TrimSpace(tokens[0])
		value := strings.TrimSpace(tokens[1])
		params[key] = value
	}
	instanceID, err := client.Provision(host, port, serviceID, planID, params)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nProvisioning service instance %s\n\n", instanceID)
	return nil
}
