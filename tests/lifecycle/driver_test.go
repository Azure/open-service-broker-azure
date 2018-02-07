// +build !unit

package lifecycle

import (
	"log"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestServices(t *testing.T) {
	resourceGroup := "test-" + uuid.NewV4().String()

	log.Printf("----> creating resource group \"%s\"\n", resourceGroup)
	err := ensureResourceGroup(resourceGroup)
	assert.Nil(t, err)
	log.Printf("----> created resource group \"%s\"\n", resourceGroup)

	log.Printf("----> testing in resource group \"%s\"\n", resourceGroup)

	// Make sure we clean up after ourselves
	defer func() {
		log.Printf("----> deleting resource group \"%s\"\n", resourceGroup)
		if err = deleteResourceGroup(resourceGroup); err != nil {
			log.Printf("----> error deleting resource group: %s", err)
		} else {
			log.Printf(
				"----> deleted resource group \"%s\"\n",
				resourceGroup,
			)
		}
	}()

	testCases, err := getTestCases(resourceGroup)
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
				err := tc.execute(t, resourceGroup)
				assert.Nil(t, err)
			})
		}
	})

}
