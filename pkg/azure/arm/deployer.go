package arm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/open-service-broker-azure/pkg/template"
	log "github.com/Sirupsen/logrus"
)

type deploymentStatus string

const (
	deploymentStatusNotFound  deploymentStatus = "NOT_FOUND"
	deploymentStatusRunning   deploymentStatus = "RUNNING"
	deploymentStatusSucceeded deploymentStatus = "SUCCEEDED"
	deploymentStatusFailed    deploymentStatus = "FAILED"
	deploymentStatusUnknown   deploymentStatus = "UNKNOWN"
)

// Deployer is an interface to be implemented by any component capable of
// deploying resource to Azure using an ARM template
type Deployer interface {
	Deploy(
		deploymentName string,
		resourceGroupName string,
		location string,
		template []byte,
		goParams interface{},
		armParams map[string]interface{},
		tags map[string]string,
	) (map[string]interface{}, error)
	Delete(deploymentName string, resourceGroupName string) error
}

// deployer is an ARM-based implementation of the Deployer interface
type deployer struct {
	groupsClient      resources.GroupsClient
	deploymentsClient resources.DeploymentsClient
}

// NewDeployer returns a new ARM-based implementation of the Deployer interface
func NewDeployer(
	groupsClient resources.GroupsClient,
	deploymentsClient resources.DeploymentsClient,
) Deployer {
	return &deployer{
		groupsClient:      groupsClient,
		deploymentsClient: deploymentsClient,
	}
}

// Deploy idempotently handles ARM deployments. To do this, it checks for the
// existence and status of a deployment before choosing to create a new one,
// poll until success or failure, or return an error.
func (d *deployer) Deploy(
	deploymentName string,
	resourceGroupName string,
	location string,
	template []byte,
	goParams interface{},
	armParams map[string]interface{},
	tags map[string]string,
) (map[string]interface{}, error) {
	logFields := log.Fields{
		"resourceGroup": resourceGroupName,
		"deployment":    deploymentName,
	}

	// Get the deployment and its current status
	deployment, ds, err := d.getDeploymentAndStatus(
		deploymentName,
		resourceGroupName,
	)
	if err != nil {
		return nil, fmt.Errorf(
			`error deploying "%s" in resource group "%s": error getting `+
				`deployment: %s`,
			deploymentName,
			resourceGroupName,
			err,
		)
	}

	// Handle according to status...
	switch ds {
	case deploymentStatusNotFound:
		// The deployment wasn't found, which means we are free to proceed with
		// initiating a new deployment
		log.WithFields(logFields).Debug(
			"deployment does not already exist; beginning new deployment",
		)
		if deployment, err = d.doNewDeployment(
			deploymentName,
			resourceGroupName,
			location,
			template,
			goParams,
			armParams,
			tags,
		); err != nil {
			return nil, fmt.Errorf(
				`error deploying "%s" in resource group "%s": %s`,
				deploymentName,
				resourceGroupName,
				err,
			)
		}
	case deploymentStatusRunning:
		// The deployment exists and is currently running, which means we'll poll
		// until it completes. The return at the end of the function will return the
		// deployment's outputs.
		log.WithFields(logFields).Debug(
			"deployment exists and is in-progress; polling until complete",
		)
		if deployment, err = d.pollUntilComplete(
			deploymentName,
			resourceGroupName,
		); err != nil {
			return nil, fmt.Errorf(
				`error deploying "%s" in resource group "%s": %s`,
				deploymentName,
				resourceGroupName,
				err,
			)
		}
	case deploymentStatusSucceeded:
		// The deployment exists and has succeeded already. There's nothing to do.
		// The return at the end of the function will return the deployment's
		// outputs.
		log.WithFields(logFields).Debug(
			"deployment exists and has already succeeded",
		)
	case deploymentStatusFailed:
		// The deployment exists and has failed already.
		return nil, fmt.Errorf(
			`error deploying "%s" in resource group "%s": deployment is in failed `+
				`state`,
			deploymentName,
			resourceGroupName,
		)
	case deploymentStatusUnknown:
		fallthrough
	default:
		// Unrecognized state
		return nil, fmt.Errorf(
			`error deploying "%s" in resource group "%s": deployment is in an `+
				`unrecognized state`,
			deploymentName,
			resourceGroupName,
		)
	}

	return getOutputs(deployment), nil
}

func (d *deployer) Delete(
	deploymentName string,
	resourceGroupName string,
) error {
	cancelCh := make(chan struct{})
	defer close(cancelCh)
	_, errChan := d.deploymentsClient.Delete(
		resourceGroupName,
		deploymentName,
		cancelCh,
	)
	timer := time.NewTimer(time.Minute * 20)
	defer timer.Stop()
	select {
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf(
				`error deleting deployment "%s" from resource group "%s": %s`,
				deploymentName,
				resourceGroupName,
				err,
			)
		}
	case <-timer.C:
		return fmt.Errorf(
			`timed out deleting deployment "%s" from resource group "%s"`,
			deploymentName,
			resourceGroupName,
		)
	}

	return nil
}

// getDeploymentAndStatus attempts to retrieve and return a deployment. Whether
// it's found or not, a status is returned. (It's not enough to just return the
// deployment and let the caller check status itself, because in the case a
// given deployment doesn't exist, there isn't one to return. Returning a
// separate status indicator resolves that problem.)
func (d *deployer) getDeploymentAndStatus(
	deploymentName string,
	resourceGroupName string,
) (*resources.DeploymentExtended, deploymentStatus, error) {
	deployment, err := d.deploymentsClient.Get(resourceGroupName, deploymentName)
	if err != nil {
		detailedErr, ok := err.(autorest.DetailedError)
		if !ok || detailedErr.StatusCode != http.StatusNotFound {
			return nil, "", err
		}
		return nil, deploymentStatusNotFound, nil
	}
	switch *deployment.Properties.ProvisioningState {
	case "Running":
		return &deployment, deploymentStatusRunning, nil
	case "Succeeded":
		return &deployment, deploymentStatusSucceeded, nil
	case "Failed":
		return &deployment, deploymentStatusFailed, nil
	default:
		return &deployment, deploymentStatusUnknown, nil
	}
}

func (d *deployer) doNewDeployment(
	deploymentName string,
	resourceGroupName string,
	location string,
	armTemplate []byte,
	goParams interface{},
	armParams map[string]interface{},
	tags map[string]string,
) (*resources.DeploymentExtended, error) {
	res, err := d.groupsClient.CheckExistence(resourceGroupName)
	if err != nil {
		return nil, fmt.Errorf(
			"error checking existence of resource group: %s",
			err,
		)
	}
	if res.StatusCode == http.StatusNotFound {
		if _, err = d.groupsClient.CreateOrUpdate(
			resourceGroupName,
			resources.Group{
				Name:     &resourceGroupName,
				Location: &location,
			},
		); err != nil {
			return nil, fmt.Errorf(
				"error creating resource group: %s",
				err,
			)
		}
	}

	finalArmTemplate := armTemplate

	// The template could be a Go text template that renders down to an ARM
	// template, so deal with that possibility first.
	if goParams != nil {
		finalArmTemplate, err = template.Render(armTemplate, goParams)
		if err != nil {
			return nil, err
		}
	}

	// Unmarshal the template into a map
	var armTemplateMap map[string]interface{}
	err = json.Unmarshal(finalArmTemplate, &armTemplateMap)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling ARM template: %s", err)
	}

	// Augment the params with location
	armParams["location"] = location

	// Deal with the possibility that tags == nil
	if tags == nil {
		tags = make(map[string]string)
	}

	// Augment the provided tags with heritage information
	tags["heritage"] = "open-service-broker-azure"

	// Deal with the possiiblity that params == nil
	if armParams == nil {
		armParams = make(map[string]interface{})
	}

	// Augment the params with tags
	armParams["tags"] = tags

	// Convert a simple map[string]interface{} to the more complex
	// map[string]map[string]interface{} required by the deployments client
	armParamsMap := map[string]interface{}{}
	for key, val := range armParams {
		armParamsMap[key] = map[string]interface{}{
			"value": val,
		}
	}

	// Deploy the template
	cancelCh := make(chan struct{})
	defer close(cancelCh)
	_, errChan := d.deploymentsClient.CreateOrUpdate(
		resourceGroupName,
		deploymentName,
		resources.Deployment{
			Properties: &resources.DeploymentProperties{
				Template:   &armTemplateMap,
				Parameters: &armParamsMap,
				Mode:       resources.Incremental,
			},
		},
		cancelCh,
	)
	timer := time.NewTimer(time.Minute * 30)
	defer timer.Stop()
	select {
	case err = <-errChan:
		if err != nil {
			return nil, fmt.Errorf("error submitting ARM template: %s", err)
		}
	case <-timer.C:
		return nil, errors.New("timed out waiting for deployment to complete")
	}

	// Deployment object found on the result channel doesn't include properties,
	// so we need to make a separate call to retrieve the deployment
	deployment, err := d.deploymentsClient.Get(resourceGroupName, deploymentName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving completed deployment: %s", err)
	}

	return &deployment, nil
}

// pollUntilComplete polls the status of a deployment periodically until the
// deployment succeeds or fails, polling fails, or a timeout is reached
func (d *deployer) pollUntilComplete(
	deploymentName string,
	resourceGroupName string,
) (*resources.DeploymentExtended, error) {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	timer := time.NewTimer(time.Minute * 30)
	defer timer.Stop()
	var deployment *resources.DeploymentExtended
	var ds deploymentStatus
	var err error
	for {
		select {
		case <-ticker.C:
			if deployment, ds, err = d.getDeploymentAndStatus(
				deploymentName,
				resourceGroupName,
			); err != nil {
				return nil, err
			}
			switch ds {
			case deploymentStatusNotFound:
				// This is an error. We'd only be polling for status on a deployment
				// that exists. If it no longer exists, something is very wrong.
				return nil, errors.New(
					"error polling deployment status; deployment should exist, but " +
						"does not",
				)
			case deploymentStatusRunning:
				// Do nothing == continue the loop
			case deploymentStatusSucceeded:
				// We're done
				return deployment, nil
			case deploymentStatusFailed:
				// The deployment has failed
				return nil, errors.New("deployment has failed")
			case deploymentStatusUnknown:
				fallthrough
			default:
				// The deployment has entered an unknown state
				return nil, errors.New("deployment is in an unrecognized state")
			}
		case <-timer.C:
			// We've reached a timeout
			return nil, errors.New("timed out waiting for deployment to complete")
		}
	}
}

func getOutputs(
	deployment *resources.DeploymentExtended,
) map[string]interface{} {
	outputs := make(map[string]interface{})
	for k, v := range *deployment.Properties.Outputs {
		outputs[k] = v.(map[string]interface{})["value"]
	}
	return outputs
}
