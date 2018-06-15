// +build !unit

package e2e

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"open-service-broker-azure/pkg/api"
	"open-service-broker-azure/pkg/client"
	"github.com/stretchr/testify/assert"
)

const (
	host = "broker"
	port = 8080
)

// e2eTestCase encapsulates all the required details for an end-to-end test
// case.
type e2eTestCase struct {
	group                  string
	name                   string
	serviceID              string
	planID                 string
	provisioningParameters map[string]interface{}
	bind                   bool
	bindingParameters      map[string]interface{}
	childTestCases         []*e2eTestCase
}

func (e e2eTestCase) getName() string {
	return fmt.Sprintf("TestE2E/%s/%s", e.group, e.name)
}

func (e e2eTestCase) execute(
	t *testing.T,
	resourceGroup string,
	basicAuthConfig api.BasicAuthConfig,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*30)
	defer cancel()

	name := e.getName()

	log.Printf("----> %s: starting\n", name)

	defer log.Printf("----> %s: completed\n", name)

	// This will periodically send status to stdout until the context is canceled.
	// THIS is what stops CI from timing out these tests!
	go e.showStatus(ctx)

	// Force the resource group to be something known to this test executor to
	// ensure good cleanup...
	if e.provisioningParameters == nil {
		e.provisioningParameters = map[string]interface{}{}
	}
	// Overwrite the resource group if there's a placeholder once already--
	// otherwise, this is a service that doesn't accept a resourceGroup
	// parameter.
	if _, ok := e.provisioningParameters["resourceGroup"]; ok {
		e.provisioningParameters["resourceGroup"] = resourceGroup
	}

	// Provision...
	instanceID, err := client.Provision(
		host,
		port,
		basicAuthConfig.GetUsername(),
		basicAuthConfig.GetPassword(),
		e.serviceID,
		e.planID,
		e.provisioningParameters,
	)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()
provisionPollLoop:
	for {
		select {
		case <-ticker.C:
			var result string
			result, err = client.Poll(
				host,
				port,
				basicAuthConfig.GetUsername(),
				basicAuthConfig.GetPassword(),
				instanceID,
				api.OperationProvisioning,
			)
			if err != nil {
				return fmt.Errorf("error polling for provisioning status: %s", err)
			}
			switch result {
			case api.OperationStateInProgress:
				continue provisionPollLoop
			case api.OperationStateSucceeded:
				break provisionPollLoop
			case api.OperationStateFailed:
				return fmt.Errorf(
					"Provisioning service instance %s has failed",
					instanceID,
				)
			default:
				return fmt.Errorf("Unrecognized operation status: %s", result)
			}
		case <-ctx.Done():
			return fmt.Errorf("context canceled with provisioning incomplete")
		}
	}

	var bindingID string
	// Bind...
	if e.bind {
		bindingID, _, err = client.Bind(
			host,
			port,
			basicAuthConfig.GetUsername(),
			basicAuthConfig.GetPassword(),
			instanceID,
			e.bindingParameters,
		)
		if err != nil {
			return err
		}
	}

	// Iterate through any child test cases
	for _, childTestCase := range e.childTestCases {
		t.Run(childTestCase.getName(), func(t *testing.T) {
			err = childTestCase.execute(t, resourceGroup, basicAuthConfig)
			//This will fail this subtest, but also the parent lifecycle test
			assert.Nil(t, err)
		})
	}

	// Unbind...
	if e.bind {
		err = client.Unbind(
			host,
			port,
			basicAuthConfig.GetUsername(),
			basicAuthConfig.GetPassword(),
			instanceID,
			bindingID,
		)
		if err != nil {
			return err
		}
	}

	// Deprovision...
	err = client.Deprovision(
		host,
		port,
		basicAuthConfig.GetUsername(),
		basicAuthConfig.GetPassword(),
		instanceID,
	)
	if err != nil {
		return err
	}
deprovisionPollLoop:
	for {
		select {
		case <-ticker.C:
			result, err := client.Poll(
				host,
				port,
				basicAuthConfig.GetUsername(),
				basicAuthConfig.GetPassword(),
				instanceID,
				api.OperationDeprovisioning,
			)
			if err != nil {
				return fmt.Errorf("error polling for deprovisioning status: %s", err)
			}
			switch result {
			case api.OperationStateInProgress:
				continue deprovisionPollLoop
			case api.OperationStateGone:
				break deprovisionPollLoop
			case api.OperationStateFailed:
				return fmt.Errorf(
					"Deprovisioning service instance %s has failed",
					instanceID,
				)
			default:
				return fmt.Errorf("Unrecognized operation status: %s", result)
			}
		case <-ctx.Done():
			return fmt.Errorf("context canceled with deprovisioning incomplete")
		}
	}

	return nil
}

func (e e2eTestCase) showStatus(ctx context.Context) {
	name := e.getName()
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Printf("----> %s: in progress\n", name)
		case <-ctx.Done():
			return
		}
	}
}
