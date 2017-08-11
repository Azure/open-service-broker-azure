package fake

type deployFunction func(
	deploymentName string,
	resourceGroupName string,
	location string,
	template []byte,
	params map[string]interface{},
) (map[string]interface{}, error)

type deleteFunction func(deploymentName string, resourceGroupName string) error

// Deployer is a fake implementaton of arm.Deployer used for testing
type Deployer struct {
	DeployBehavior deployFunction
	DeleteBehavior deleteFunction
}

// NewDeployer returns a new, fake implementation of arm.Deployer used for
// testing
func NewDeployer() *Deployer {
	return &Deployer{
		DeployBehavior: defaultDeployBehavior,
		DeleteBehavior: defaultDeleteBehavior,
	}
}

// Deploy simulates deployment of an ARM template
func (d *Deployer) Deploy(
	deploymentName string,
	resourceGroupName string,
	location string,
	template []byte,
	params map[string]interface{},
) (map[string]interface{}, error) {
	return d.DeployBehavior(
		deploymentName,
		resourceGroupName,
		location,
		template,
		params,
	)
}

// Delete simulates deletion of an ARM deployment
func (d *Deployer) Delete(
	deploymentName string,
	resourceGroupName string,
) error {
	return d.DeleteBehavior(deploymentName, resourceGroupName)
}

func defaultDeployBehavior(
	deploymentName string,
	resourceGroupName string,
	location string,
	template []byte,
	params map[string]interface{},
) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func defaultDeleteBehavior(
	deploymentName string,
	resourceGroupName string,
) error {
	return nil
}
