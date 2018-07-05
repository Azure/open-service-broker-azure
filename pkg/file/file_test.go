package file

import (
	"fmt"
	"go/build"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	file := fmt.Sprintf(
		"%s/src/github.com/Azure/open-service-broker-azure/pkg/file/file_test.go",
		gopath,
	)
	assert.True(t, Exists(file))
	file = fmt.Sprintf(
		"%s/src/github.com/Azure/open-service-broker-azure/pkg/file/bogus.go",
		gopath,
	)
	assert.False(t, Exists(file))
}
