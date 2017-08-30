// +build !unit

package lifecycle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModules(t *testing.T) {
	testCases, err := getTestCases()
	assert.Nil(t, err)

	for _, testCase := range testCases {
		// Important: Assign the value of testCase to a variable scoped within this
		// for loop-- if we don't, and simply have the function passed to t.Run()
		// below close over testCase instead, it would be closing over a variable
		// whose value will change as we continue to iterate over all the testCases.
		tc := testCase
		t.Run(tc.module.GetName(), func(t *testing.T) {
			// Run subtests in parallel!
			t.Parallel()
			err := tc.execute()
			assert.Nil(t, err)
		})
	}
}
