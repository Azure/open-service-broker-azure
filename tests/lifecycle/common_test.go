// +build !unit

package lifecycle

import uuid "github.com/satori/go.uuid"

func newTestResourceGroupName() string {
	return "test-" + uuid.NewV4().String()
}
