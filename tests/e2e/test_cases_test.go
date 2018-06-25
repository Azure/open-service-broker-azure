// +build !unit

package e2e

import (
	"log"
	"os"
	"strings"
)

func getTestCases() ([]e2eTestCase, error) {
	testCases := getPostgreSQLTestCases()
	// testCases = append(testCases, rediscacheTestCases...)
	// testCases = append(testCases, aciTestCases...)
	// testCases = append(testCases, cosmosdbTestCases...)
	// testCases = append(testCases, eventhubsTestCases...)
	// testCases = append(testCases, keyvaultTestCases...)
	testCases = append(testCases, getMSSQLTestCases()...)
	testCases = append(testCases, getMySQLTestCases()...)
	// testCases = append(testCases, searchTestCases...)
	testCases = append(testCases, servicebusTestCases...)
	// testCases = append(testCases, storageTestCases...)

	testCases = filterTestCases(testCases, getTestFilters())

	if len(testCases) == 0 {
		log.Print("No test cases selected. Please check TEST_MODULES variable.")
	}

	return testCases, nil
}

func filterTestCases(
	testCases []e2eTestCase,
	filters map[string]struct{},
) []e2eTestCase {
	// If filters is empty, we are not filtering so include all the testcases
	if len(filters) == 0 {
		return testCases
	}

	// If filters is not empty, see if the testcase's module name is in the filter
	// map
	filtered := testCases[:0]
	for _, testCase := range testCases {
		_, ok := filters[testCase.group]
		if ok {
			filtered = append(filtered, testCase)
		}
	}
	return filtered
}

func getTestFilters() map[string]struct{} {
	f := make(map[string]struct{})
	if env := os.Getenv("TEST_MODULES"); env != "" {
		filters := strings.Split(env, ",")
		log.Printf("Running tests for modules: %v", filters)
		for _, filter := range filters {
			f[filter] = struct{}{}
		}
	}
	return f
}
