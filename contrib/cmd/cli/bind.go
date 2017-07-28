package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-service-broker/contrib/pkg/client"
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
