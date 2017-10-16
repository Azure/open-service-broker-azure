// +build !unit

package lifecycle

import (
	"github.com/Azure/azure-service-broker/pkg/azure/arm"
)

func getTestCases() ([]moduleLifecycleTestCase, error) {
	armDeployer, err := arm.NewDeployer()
	if err != nil {
		return nil, err
	}

	testCases := []moduleLifecycleTestCase{}

	getTestCaseFuncs := []func(arm.Deployer) ([]moduleLifecycleTestCase, error){
		getACICases,
		getCosmosdbCases,
		getEventhubCases,
		getKeyvaultCases,
		getMssqlCases,
		getMysqlCases,
		getPostgresqlCases,
		getRediscacheCases,
		getSearchCases,
		getServicebusCases,
		getStorageCases,
	}

	for _, getTestCaseFunc := range getTestCaseFuncs {
		if tcs, err := getTestCaseFunc(armDeployer); err == nil {
			testCases = append(testCases, tcs...)
		} else {
			return nil, err
		}
	}

	return testCases, nil
}
