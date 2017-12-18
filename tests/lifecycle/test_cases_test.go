// +build !unit

package lifecycle

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
)

func getTestCases(resourceGroup string) ([]moduleLifecycleTestCase, error) {
	armDeployer, err := arm.NewDeployer()
	if err != nil {
		return nil, err
	}

	testCases := []moduleLifecycleTestCase{}

	getTestCaseFuncs := []func(
		armDeployer arm.Deployer,
		resourceGroup string,
	) ([]moduleLifecycleTestCase, error){
		getRediscacheCases,
		getACICases,
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

	for _, getTestCaseFunc := range getTestCaseFuncs {
		if tcs, err := getTestCaseFunc(armDeployer, resourceGroup); err == nil {
			testCases = append(testCases, tcs...)
		} else {
			return nil, err
		}
	}

	return testCases, nil
}
