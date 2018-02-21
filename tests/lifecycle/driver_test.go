// +build !unit

package lifecycle

import (
	"log"
	"testing"
	"time"

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
		done := make(chan struct{})
		failed := make(chan error)
		t := time.NewTicker(time.Minute * 5).C
		timeout := time.NewTimer(time.Minute * 30).C
		go func() {
			if err = deleteResourceGroup(resourceGroup); err != nil {
				failed <- err
			} else {
				done <- struct{}{}
			}
		}()
		for {
			select {
			case <-failed:
				log.Printf("----> error deleting resource group: %s", err)
				return
			case <-done:
				log.Printf(
					"----> deleted resource group \"%s\"\n",
					resourceGroup,
				)
				return
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
				return
			}
		}

	}()

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
				err := tc.execute(t, resourceGroup)
				assert.Nil(t, err)
			})
		}
	})

}
