// +build !unit

package lifecycle

import (
	"log"
	"os"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
)

func getTestCases(resourceGroup string) ([]serviceLifecycleTestCase, error) {
	armDeployer, err := arm.NewDeployer()
	if err != nil {
		return nil, err
	}

	testCases := []serviceLifecycleTestCase{}

	getTestCaseFuncs := []func(
		armDeployer arm.Deployer,
		resourceGroup string,
	) ([]serviceLifecycleTestCase, error){
		getRediscacheCases,
		getACICases,
		getAcrCases,
		getCosmosdbCases,
		getEventhubCases,
		getKeyvaultCases,
		getMssqlCases,
		getMysqlCases,
		getPostgresqlCases,
		getSearchCases,
		getServicebusCases,
		getStorageCases,
	}

	testFilters := getTestFilters()

	for _, getTestCaseFunc := range getTestCaseFuncs {
		if tcs, err := getTestCaseFunc(armDeployer, resourceGroup); err == nil {
			testCases = filter(append(testCases, tcs...), testFilters)
		} else {
			return nil, err
		}
	}
	if len(testCases) == 0 {
		log.Print("No test cases selected. Please check TEST_MODULES variable.")
	}
	return testCases, nil
}

func filter(
	testCases []serviceLifecycleTestCase,
	filters map[string]struct{},
) []serviceLifecycleTestCase {
	//If filters is empty, we are not filtering so include all the testcases
	if len(filters) == 0 {
		return testCases
	}

	//If filters is not empty, see if the testcase's module name is in the filter
	//map
	filtered := testCases[:0]
	for _, testCase := range testCases {
		_, ok := filters[testCase.module.GetName()]
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
