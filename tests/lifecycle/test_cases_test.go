// +build !unit

package lifecycle

import (
	"log"
	"os"
	"strings"
	"time"

	resourcesSDK "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources" // nolint: lll
	"github.com/Azure/go-autorest/autorest"
	azureSDK "github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
)

func getTestCases() ([]serviceLifecycleTestCase, error) {
	azureConfig, err := azure.GetConfigFromEnvironment()
	if err != nil {
		return nil, err
	}

	authorizer, err := azure.GetBearerTokenAuthorizer(
		azureConfig.Environment,
		azureConfig.TenantID,
		azureConfig.ClientID,
		azureConfig.ClientSecret,
	)
	if err != nil {
		return nil, err
	}

	resourceGroupsClient := resourcesSDK.NewGroupsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	resourceGroupsClient.Authorizer = authorizer
	resourceDeploymentsClient := resourcesSDK.NewDeploymentsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	resourceDeploymentsClient.Authorizer = authorizer
	resourceDeploymentsClient.PollingDuration = time.Minute * 30
	armDeployer := arm.NewDeployer(
		resourceGroupsClient,
		resourceDeploymentsClient,
	)

	testCases := []serviceLifecycleTestCase{}

	getTestCaseFuncs := []func(
		azureEnvironment azureSDK.Environment,
		subscriptionID string,
		authorizer autorest.Authorizer,
		armDeployer arm.Deployer,
	) ([]serviceLifecycleTestCase, error){
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

	testFilters := getTestFilters()

	for _, getTestCaseFunc := range getTestCaseFuncs {
		if tcs, err := getTestCaseFunc(
			azureConfig.Environment,
			azureConfig.SubscriptionID,
			authorizer,
			armDeployer,
		); err == nil {
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
	// If filters is empty, we are not filtering so include all the testcases
	if len(filters) == 0 {
		return testCases
	}

	// If filters is not empty, see if the testcase's module name is in the filter
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
