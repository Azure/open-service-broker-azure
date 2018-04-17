package cosmosdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (s *sqlAllInOneManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", s.preProvision),
		service.NewProvisioningStep("deployARMTemplate", s.deployARMTemplate),
		service.NewProvisioningStep("createDatabase", s.createDatabase),
	)
}

func (s *sqlAllInOneManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := sqlAllInOneInstanceDetails{
		ARMDeploymentName:   uuid.NewV4().String(),
		DatabaseAccountName: generateAccountName(instance.Location),
		DatabaseName:        uuid.NewV4().String(),
	}
	dtMap, err := service.GetMapFromStruct(dt)
	return dtMap, nil, err
}

func (s *sqlAllInOneManager) buildGoTemplateParams(
	pp *provisioningParameters,
	dt *sqlAllInOneInstanceDetails,
	kind string,
) map[string]interface{} {
	p := map[string]interface{}{}
	p["name"] = dt.DatabaseAccountName
	p["kind"] = kind

	filters := []string{}

	if pp.IPFilterRules != nil {
		allowAzure := strings.ToLower(pp.IPFilterRules.AllowPortal)
		allowPortal := strings.ToLower(pp.IPFilterRules.AllowPortal)
		if allowAzure != "disable" {
			filters = append(filters, "0.0.0.0")
		} else if allowPortal != "disable" {
			// Azure Portal IP Addresses per:
			// https://aka.ms/Vwxndo
			//|| Region            || IP address(es) ||
			//||=====================================||
			//|| China             || 139.217.8.252  ||
			//||===================||================||
			//|| Germany           || 51.4.229.218   ||
			//||===================||================||
			//|| US Gov            || 52.244.48.71   ||
			//||===================||================||
			//|| All other regions || 104.42.195.92  ||
			//||                   || 40.76.54.131   ||
			//||                   || 52.176.6.30    ||
			//||                   || 52.169.50.45   ||
			//||                   || 52.187.184.26  ||
			//=======================================||
			// Given that we don't really have context of the cloud
			// we are provisioning with right now, use all of the above
			// addresses.
			filters = append(filters,
				"104.42.195.92",
				"40.76.54.131",
				"52.176.6.30",
				"52.169.50.45",
				"52.187.184.26",
				"51.4.229.218",
				"139.217.8.252",
				"52.244.48.71",
			)
		}
		filters = append(filters, pp.IPFilterRules.Filters...)
	} else {
		filters = append(filters, "0.0.0.0")
	}
	if len(filters) > 0 {
		p["ipFilters"] = strings.Join(filters, ",")
	}
	return p
}

func (s *sqlAllInOneManager) deployARMTemplate(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {

	pp := &provisioningParameters{}
	if err :=
		service.GetStructFromMap(instance.ProvisioningParameters, pp); err != nil {
		return nil, nil, err
	}

	dt := &sqlAllInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}

	p := s.buildGoTemplateParams(pp, dt, "GlobalDocumentDB")
	if instance.Tags == nil {
		instance.Tags = make(map[string]string)
	}

	outputs, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
		p, // Go template params
		map[string]interface{}{},
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	var ok bool
	dt.FullyQualifiedDomainName, ok = outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving fully qualified domain name from deployment: %s",
			err,
		)
	}

	primaryKey, ok := outputs["primaryKey"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving primary key from deployment: %s",
			err,
		)
	}

	sdt := cosmosdbSecureInstanceDetails{
		PrimaryKey: primaryKey,
	}

	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	sdt.ConnectionString = fmt.Sprintf("AccountEndpoint=%s;AccountKey=%s;",
		dt.FullyQualifiedDomainName,
		sdt.PrimaryKey,
	)

	dtMap, err := service.GetMapFromStruct(dt)
	if err != nil {
		return nil, nil, err
	}
	sdtMap, err := service.GetMapFromStruct(sdt)
	return dtMap, sdtMap, err
}

func (s *sqlAllInOneManager) createDatabase(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := &sqlAllInOneInstanceDetails{}
	err := service.GetStructFromMap(instance.Details, &dt)
	if err != nil {
		fmt.Printf("Failed to get DT Struct from Map %s", err)
		return nil, nil, err
	}

	sdt := &cosmosdbSecureInstanceDetails{}
	err = service.GetStructFromMap(instance.SecureDetails, &sdt)
	if err != nil {
		fmt.Printf("Failed to get SDT Struct from Map %s", err)
		return nil, nil, err
	}
	err = createDatabase(dt.DatabaseAccountName, dt.DatabaseName, sdt.PrimaryKey)
	if err != nil {
		return nil, nil, err
	}
	return instance.Details, instance.SecureDetails, nil
}
