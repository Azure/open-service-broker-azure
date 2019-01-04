package mssqldr

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	sqlSDK "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func validateServer(
	ctx context.Context,
	serversClient *sqlSDK.ServersClient,
	resourceGroup string,
	serverName string,
	expectedVersion string,
	expectedLocation string,
) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	result, err := serversClient.Get(
		ctx,
		resourceGroup,
		serverName,
	)
	if err != nil {
		return "", fmt.Errorf("error getting the sql server: %s", err)
	}
	if result.Name == nil {
		return "", fmt.Errorf(
			"can't find sql server %s in the resource group %s",
			serverName,
			resourceGroup,
		)
	}
	if *result.Version != expectedVersion {
		return "", fmt.Errorf(
			"sql server version validation failed, "+
				"expected version: %s, actual version: %s",
			expectedVersion,
			*result.Version,
		)
	}
	expectedLocation = strings.Replace(
		strings.ToLower(expectedLocation),
		" ",
		"",
		-1,
	)
	actualLocation := strings.Replace(
		strings.ToLower(*result.Location),
		" ",
		"",
		-1,
	)
	if expectedLocation != actualLocation {
		return "", fmt.Errorf(
			"sql server location validation failed, "+
				"expected location: %s, actual location: %s",
			expectedLocation,
			actualLocation,
		)
	}
	if *result.FullyQualifiedDomainName == "" {
		return "", fmt.Errorf(
			"sql server details doesn't contain FQDN",
		)
	}
	return *result.FullyQualifiedDomainName, nil
}

func validateDatabase(
	ctx context.Context,
	databasesClient *sqlSDK.DatabasesClient,
	resourceGroup string,
	serverName string,
	databaseName string,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	_, err := databasesClient.Get(
		ctx,
		resourceGroup,
		serverName,
		databaseName,
		"",
	)
	if err != nil {
		return fmt.Errorf("error getting the sql database: %s", err)
	}
	// TODO: add the plan as param and validate?
	return nil
}

func validateFailoverGroup(
	ctx context.Context,
	failoverGroupsClient *sqlSDK.FailoverGroupsClient,
	resourceGroup string,
	priServerName string,
	secServerName string,
	databaseName string,
	failoverGroupName string,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	result, err := failoverGroupsClient.Get(
		ctx,
		resourceGroup,
		priServerName,
		failoverGroupName,
	)
	if err != nil {
		return fmt.Errorf("error getting the failover group: %s", err)
	}
	if len(*result.PartnerServers) > 1 {
		return fmt.Errorf("error unexpected more than one partner server")
	}
	partnerServerRole := (*result.PartnerServers)[0].ReplicationRole
	if partnerServerRole != sqlSDK.Secondary {
		return fmt.Errorf("error partner server is not " +
			"the secondary role in the failover group")
	}
	partnerServerID := *(*result.PartnerServers)[0].ID
	serverNameReStr := regexp.MustCompile("^.*servers/([-a-zA-Z0-9]+)$")
	partnerServerName := serverNameReStr.ReplaceAllString(
		partnerServerID,
		"$1",
	)
	if partnerServerName != secServerName {
		return fmt.Errorf("error partner server is not the one specified")
	}

	if len(*result.Databases) > 1 {
		return fmt.Errorf("error unexpected more than one database " +
			"in the failover group")
	}
	actualDatabaseID := (*result.Databases)[0]
	databaseNameReStr := regexp.MustCompile("^.*databases/([-a-zA-Z0-9]+)$")
	actualDatabaseName := databaseNameReStr.ReplaceAllString(
		actualDatabaseID,
		"$1",
	)
	if actualDatabaseName != databaseName {
		return fmt.Errorf("error unexpected database in the failover group")
	}
	return nil
}

func buildDatabaseGoTemplateParameters(
	databaseName string,
	pp service.ProvisioningParameters,
	pd planDetails,
) (map[string]interface{}, error) {
	td, err := pd.getTierProvisionParameters(pp)
	if err != nil {
		return nil, err
	}
	p := map[string]interface{}{}
	p["databaseName"] = databaseName
	for key, value := range td {
		p[key] = value
	}
	return p, nil
}

func deployDatabaseARMTemplate(
	armDeployer *arm.Deployer,
	armDeploymentName string,
	resourceGroup string,
	location string,
	serverName string,
	databaseName string,
	pp service.ProvisioningParameters,
	pd planDetails,
	tags map[string]string,
) error {
	goTemplateParams, err := buildDatabaseGoTemplateParameters(
		databaseName,
		pp,
		pd,
	)
	if err != nil {
		return err
	}
	goTemplateParams["location"] = location
	goTemplateParams["serverName"] = serverName
	_, err = (*armDeployer).Deploy(
		armDeploymentName,
		resourceGroup,
		location,
		databaseARMTemplateBytes,
		goTemplateParams,
		map[string]interface{}{}, // empty arm params
		tags,
	)
	return err
}

func deployFailoverGroupARMTemplate(
	armDeployer *arm.Deployer,
	instance service.Instance,
) error {
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	dt := instance.Details.(*databasePairInstanceDetails)
	pp := instance.ProvisioningParameters
	ppp := instance.Parent.ProvisioningParameters
	goTemplateParams := map[string]interface{}{}
	goTemplateParams["priServerName"] = pdt.PriServerName
	goTemplateParams["secServerName"] = pdt.SecServerName
	goTemplateParams["failoverGroupName"] = pp.GetString("failoverGroup")
	goTemplateParams["databaseName"] = pp.GetString("database")
	tagsObj := pp.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	_, err := (*armDeployer).Deploy(
		dt.FailoverGroupARMDeploymentName,
		ppp.GetString("primaryResourceGroup"),
		ppp.GetString("primaryLocation"),
		failoverGroupARMTemplateBytes,
		goTemplateParams,
		map[string]interface{}{}, // empty arm params
		tags,
	)
	return err
}

func deployDatabaseARMTemplateForExistingInstance(
	armDeployer *arm.Deployer,
	armDeploymentName string,
	resourceGroup string,
	location string,
	serverName string,
	databaseName string,
	tags map[string]string,
) error {
	goTemplateParams := map[string]interface{}{}
	goTemplateParams["location"] = location
	goTemplateParams["serverName"] = serverName
	goTemplateParams["databaseName"] = databaseName
	_, err := (*armDeployer).Deploy(
		armDeploymentName,
		resourceGroup,
		location,
		databaseARMTemplateBytesForExistingInstance,
		goTemplateParams,
		map[string]interface{}{}, // empty arm params
		tags,
	)
	return err
}
