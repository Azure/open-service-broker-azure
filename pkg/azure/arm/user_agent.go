package arm

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/version"
)

func userAgent() string {
	return fmt.Sprintf("open-service-broker/%s", version.GetVersion())
}
