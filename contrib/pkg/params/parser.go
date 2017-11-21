package params

import (
	"fmt"
	"regexp"
	"strings"
)

var arrayParamKeyRegex = regexp.MustCompile(`^([\-\w]+)\[(\d+)\]$`)

// Parse splits a single param into key and value components, using "=" as a
// delimiter.
func Parse(rawParamStr string) (string, string, error) {
	rawParamStr = strings.TrimSpace(rawParamStr)
	tokens := strings.SplitN(rawParamStr, "=", 2)
	if len(tokens) != 2 {
		return "", "", fmt.Errorf(
			`parameter string "%s" is incorrectly formatted`,
			rawParamStr,
		)
	}
	return strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1]), nil
}
