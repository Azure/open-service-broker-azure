package main

import (
	"fmt"
	"strconv"

	parms "github.com/Azure/azure-service-broker/contrib/pkg/params"
	"github.com/urfave/cli"
)

// parseParams iterates, in turn, over string, int, float, and bool params as
// specified by the user and parses them into a map[string]interface{}.
func parseParams(c *cli.Context) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	rawParamStrs := c.StringSlice(flagParameter)
	for _, rawParamStr := range rawParamStrs {
		key, val, err := parms.Parse(rawParamStr)
		if err != nil {
			return nil, err
		}
		if err := parms.Add(params, key, val); err != nil {
			return nil, err
		}
	}
	rawParamStrs = c.StringSlice(flagIntParameter)
	for _, rawParamStr := range rawParamStrs {
		key, valStr, err := parms.Parse(rawParamStr)
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
		if err := parms.Add(params, key, val); err != nil {
			return nil, err
		}
	}
	rawParamStrs = c.StringSlice(flagFloatParameter)
	for _, rawParamStr := range rawParamStrs {
		key, valStr, err := parms.Parse(rawParamStr)
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
		if err := parms.Add(params, key, val); err != nil {
			return nil, err
		}
	}
	rawParamStrs = c.StringSlice(flagBoolParameter)
	for _, rawParamStr := range rawParamStrs {
		key, valStr, err := parms.Parse(rawParamStr)
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
		if err := parms.Add(params, key, val); err != nil {
			return nil, err
		}
	}
	return params, nil
}
