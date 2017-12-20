package arm

import (
	"fmt"
)

func userAgent(version string) string {
	return fmt.Sprintf("open-service-broker/%s", version)
}
