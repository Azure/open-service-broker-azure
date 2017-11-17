package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-service-broker/contrib/pkg/client"
	"github.com/urfave/cli"
)

func parseParams(c *cli.Context) (map[string]interface{}, error) {
	params := client.ProvisioningParameters{}
	rawParamStrs := c.StringSlice(flagParameter)
	for _, rawParamStr := range rawParamStrs {
		key, val, err := parseParam(rawParamStr)
		if err != nil {
			return nil, err
		}
		params[key] = val
	}
	rawParamStrs = c.StringSlice(flagIntParameter)
	for _, rawParamStr := range rawParamStrs {
		key, valStr, err := parseParam(rawParamStr)
		if err != nil {
			return nil, err
		}
		val, err := strconv.Atoi(valStr)
		if err != nil {
			return nil, fmt.Errorf(
				`error parsing int value from parameter string "%s"`,
				rawParamStr,
			)
		}
		params[key] = val
	}
	rawParamStrs = c.StringSlice(flagFloatParameter)
	for _, rawParamStr := range rawParamStrs {
		key, valStr, err := parseParam(rawParamStr)
		if err != nil {
			return nil, err
		}
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return nil, fmt.Errorf(
				`error parsing float value from parameter string "%s"`,
				rawParamStr,
			)
		}
		params[key] = val
	}
	rawParamStrs = c.StringSlice(flagBoolParameter)
	for _, rawParamStr := range rawParamStrs {
		key, valStr, err := parseParam(rawParamStr)
		if err != nil {
			return nil, err
		}
		val, err := strconv.ParseBool(valStr)
		if err != nil {
			return nil, fmt.Errorf(
				`error parsing bool value from parameter string "%s"`,
				rawParamStr,
			)
		}
		params[key] = val
	}
	return params, nil
}

func parseParam(rawParamStr string) (string, string, error) {
	rawParamStr = strings.TrimSpace(rawParamStr)
	tokens := strings.Split(rawParamStr, "=")
	if len(tokens) != 2 {
		return "", "", fmt.Errorf(
			`parameter string "%s" is incorrectly formatted`,
			rawParamStr,
		)
	}
	return strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1]), nil
}
