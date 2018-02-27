package main

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli"
)

// parseParams iterates, in turn, over string, int, float, and bool params as
// specified by the user and parses them into a map[string]interface{}.
func parseParams(c *cli.Context) (map[string]interface{}, error) {
	params := map[string]interface{}{}
	paramsStr := c.String(flagParameters)
	if paramsStr == "" {
		return params, nil
	}
	jsonBytes := []byte(paramsStr)
	err := json.Unmarshal(jsonBytes, &params)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON parameters: %s", err)
	}
	return params, nil
}
