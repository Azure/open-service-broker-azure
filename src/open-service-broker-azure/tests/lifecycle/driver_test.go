// +build !unit

package lifecycle

import (
	"log"
	"os"
	"testing"
	"time"

	"open-service-broker-azure/pkg/azure"
	"open-service-broker-azure/pkg/boot"
	"open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var resourceGroup string

func TestServices(t *testing.T) {
	log.Printf("----> testing in resource group \"%s\"\n", resourceGroup)

	azureConfig, err := azure.GetConfigFromEnvironment()
	assert.Nil(t, err)

	catalogConfig := service.NewCatalogConfigWithDefaults()
	catalogConfig.MinStability = service.StabilityExperimental

	catalog, err := boot.GetCatalog(catalogConfig, azureConfig)
	assert.Nil(t, err)

	testCases, err := getTestCases()
	assert.Nil(t, err)

	t.Run("lifecycle", func(t *testing.T) {
		for _, testCase := range testCases {
			// Important: Assign the value of testCase to a variable scoped within this
			// for loop-- if we don't, and simply have the function passed to t.Run()
			// below close over testCase instead, it would be closing over a variable
			// whose value will change as we continue to iterate over all the testCases.
			tc := testCase
			t.Run(tc.getName(), func(t *testing.T) {
				// Run subtests in parallel!
				t.Parallel()
				err := tc.execute(t, catalog, resourceGroup)
				assert.Nil(t, err)
			})
		}
	})
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		os.Exit(-1)
	}
	exitCode := m.Run()
	if err := tearDown(); err != nil {
		os.Exit(-1)
	}
	os.Exit(exitCode)
}

func setup() error {
	resourceGroup = "test-" + uuid.NewV4().String()

	log.Printf("----> creating resource group \"%s\"\n", resourceGroup)
	if err := ensureResourceGroup(resourceGroup); err != nil {
		return err
	}
	log.Printf("----> created resource group \"%s\"\n", resourceGroup)

	return nil
}

func tearDown() error {
	log.Printf("----> deleting resource group \"%s\"\n", resourceGroup)
	done := make(chan struct{})
	failed := make(chan error)
	t := time.NewTicker(time.Minute * 5).C
	timeout := time.NewTimer(time.Minute * 30).C
	go func() {
		if err := deleteResourceGroup(resourceGroup); err != nil {
			failed <- err
		} else {
			done <- struct{}{}
		}
	}()
	for {
		select {
		case err := <-failed:
			log.Printf("----> error deleting resource group: %s", err)
			return err
		case <-done:
			log.Printf(
				"----> deleted resource group \"%s\"\n",
				resourceGroup,
			)
			return nil
		case <-t:
			//Periodically emit a message
			log.Printf(
				"----> delete resource group \"%s\": in progress\n",
				resourceGroup,
			)
		case <-timeout:
			//Also use a timeout channel to enforce some (un)reasonable
			//lenght for the delete RG to get killed in
			log.Printf("----> error deleting resource group: timeout")
			return nil
		}
	}
}
