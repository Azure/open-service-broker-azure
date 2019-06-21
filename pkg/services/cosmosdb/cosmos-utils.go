package cosmosdb

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	cosmosSDK "github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb" //nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

// This method implements the CosmosDB API authentication token generation
// scheme. For reference, please see the CosmosDB REST API at:
// https://aka.ms/Fyra7j
func generateAuthToken(verb, id, date, key string) (string, error) {
	resource := "dbs"
	var resourceID string
	if id != "" {
		resourceID = fmt.Sprintf("%s/%s", strings.ToLower(resource), id)
	} else {
		resourceID = id
	}
	payload := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n",
		strings.ToLower(verb),
		strings.ToLower(resource),
		resourceID,
		strings.ToLower(date),
		"",
	)

	decodedKey, _ := base64.StdEncoding.DecodeString(key)
	hmac := hmac.New(sha256.New, decodedKey)
	_, err := hmac.Write([]byte(payload))
	if err != nil {
		return "", err
	}
	b := hmac.Sum(nil)
	authHash := base64.StdEncoding.EncodeToString(b)
	authHeader := url.QueryEscape("type=master&ver=1.0&sig=" + authHash)
	return authHeader, nil
}

func createRequest(
	accountName string,
	method string,
	resourceID string,
	key string,
	body interface{},
) (*http.Request, error) {
	resourceType := "dbs" // If we support other types, parameterize this
	path := fmt.Sprintf("%s/%s", resourceType, resourceID)
	url := fmt.Sprintf("https://%s.documents.azure.com/%s", accountName, path)
	var buf *bytes.Buffer
	var err error
	var req *http.Request
	if body != nil {
		var b []byte
		b, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf(
				"error building comsosdb request body: %s",
				err,
			)
		}
		buf = bytes.NewBuffer(b)
		req, err = http.NewRequest(method, url, buf)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("error building comsosdb request: %s", err)
	}

	dateStr := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	authHeader, err := generateAuthToken(
		method,
		resourceID,
		dateStr,
		key,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Ms-Date", dateStr)
	req.Header.Add("X-Ms-version", "2017-02-22")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", authHeader)

	return req, nil
}

func createDatabase(
	accountName string,
	id string,
	key string,
) error {
	request := &databaseCreationRequest{
		ID: id,
	}
	databaseName := ""
	req, err := createRequest(
		accountName,
		"POST",
		databaseName,
		key,
		request,
	)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf(
			"error making create comsosdb database request: %s",
			err,
		)
	}
	if resp.StatusCode != 201 { // CosmosDB returns a 201 on success
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf(
				"error creating database %d : unable to get body",
				resp.StatusCode,
			)
		}
		return fmt.Errorf(
			"error creating database %d : %s",
			resp.StatusCode,
			string(body),
		)
	}
	return nil
}

func deleteDatabase(
	databaseAccount string,
	databaseName string,
	key string,
) error {
	req, err := createRequest(
		databaseAccount,
		"DELETE",
		databaseName,
		key,
		nil, //No Body here
	)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf(
			"error making delete comsosdb database request: %s",
			err,
		)
	}
	if resp.StatusCode != 204 { // CosmosDB returns a 204 on success
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf(
				"error deleting database %d : unable to get body",
				resp.StatusCode,
			)
		}
		return fmt.Errorf(
			"error deleting database %d : %s",
			resp.StatusCode,
			string(body),
		)
	}
	return nil
}

const succeeded = "succeeded"

// This method will return when any of following situations is satisfied:
// 1. The parent context is cancelled
// 2. The timeout expired. Currently, the timeout is calculated as :
// the_number_of_read_locations * 7 minutes
// 3. Every location in parameter `readLocations` is created successfully
func pollingUntilReadLocationsReady(
	ctx context.Context,
	resourceGroupName string,
	accountName string,
	databaseAccountClient cosmosSDK.DatabaseAccountsClient,
	location string,
	readLocations []string,
	// When updating an existing instance, data in existing database will
	// be copied and synchronized across all regions, it's hard to estimate
	// how long the process will take so we disable timeout when updating.
	enableTimeout bool, // nolint: unparam
) error {
	const timeForOneReadLocation = time.Minute * 7
	readLocations = append([]string{location}, readLocations...)

	var cancel context.CancelFunc
	if enableTimeout {
		ctx, cancel = context.WithDeadline(
			ctx,
			time.Now().Add(
				time.Duration(len(readLocations))*timeForOneReadLocation,
			),
		)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			result, err := databaseAccountClient.Get(
				ctx,
				resourceGroupName,
				accountName,
			)
			if err != nil {
				return err
			}

			currentLocations := *(result.DatabaseAccountProperties.ReadLocations)
			if isCreationSucceeded(
				readLocations,
				currentLocations,
			) {
				// This is a temporary fix and should be removed after #617 is resolved
				time.Sleep(time.Second * 20)
				return nil
			}
		}
	}
}

func validateReadLocations(
	context string,
	regions []string,
) error {
	allowedLocations := make(map[string]bool)
	for _, item := range allowedReadLocations() {
		allowedLocations[item.Value] = true
	}
	occurred := make(map[string]bool)
	for i := range regions {
		region := regions[i]
		if !allowedLocations[region] {
			return service.NewValidationError(
				fmt.Sprintf("%s.readRegions", context),
				fmt.Sprintf("given read region %s is not allowed", region),
			)
		}
		if occurred[region] {
			return service.NewValidationError(
				fmt.Sprintf("%s.readRegions", context),
				fmt.Sprintf(
					"given read region %s can only occur once",
					region,
				),
			)
		}
		occurred[region] = true
	}
	return nil
}

// Allowed CosmosDB read locations, it is different from Azure regions.
func allowedReadLocations() []service.EnumValue {
	azurePublicCloudCosmosDBLocations := []service.EnumValue{
		{Value: "westus2", Title: "West US 2"},
		{Value: "westus", Title: "West US"},
		{Value: "southcentralus", Title: "South Central US"},
		{Value: "centralus", Title: "Central US"},
		{Value: "northcentralus", Title: "North Central US"},
		{Value: "canadacentral", Title: "Canada Central"},
		{Value: "eastus", Title: "East US"},
		{Value: "eastus2", Title: "East US 2"},
		{Value: "canadaeast", Title: "Canada East"},
		{Value: "brazilsouth", Title: "Brazil South"},
		{Value: "northeurope", Title: "North Europe"},
		{Value: "ukwest", Title: "UK West"},
		{Value: "uksouth", Title: "UK South"},
		{Value: "francecentral", Title: "France Central"},
		{Value: "westeurope", Title: "West Europe"},
		{Value: "westindia", Title: "West India"},
		{Value: "centralindia", Title: "Central India"},
		{Value: "southindia", Title: "South India"},
		{Value: "southeastasia", Title: "Southeast Asia"},
		{Value: "eastasia", Title: "East Asia"},
		{Value: "koreacentral", Title: "Korea Central"},
		{Value: "koreasouth", Title: "Korea South"},
		{Value: "japaneast", Title: "Japan East"},
		{Value: "japanwest", Title: "Japan West"},
		{Value: "australiasoutheast", Title: "Australia Southeast"},
		{Value: "australiaeast", Title: "Australia East"},
	}

	azureChinaCloudCosmosDBLocations := []service.EnumValue{
		{Value: "chinanorth2", Title: "China North 2"},
		{Value: "chinaeast2", Title: "China East 2"},
	}

	environmentName := azure.GetEnvironmentName()
	switch environmentName {
	case "AzurePublicCloud":
		return azurePublicCloudCosmosDBLocations
	case "AzureChinaCloud":
		return azureChinaCloudCosmosDBLocations
	}
	return azurePublicCloudCosmosDBLocations
}

func (c *cosmosAccountManager) buildGoTemplateParamsCore(
	pp *service.ProvisioningParameters,
	dt *cosmosdbInstanceDetails,
	kind string,
	readLocations []string,
) (map[string]interface{}, error) {
	p := map[string]interface{}{}
	p["name"] = dt.DatabaseAccountName
	p["kind"] = kind
	p["location"] = pp.GetString("location")
	p["readLocations"] = buildReadLocationInformation(
		readLocations,
		dt.DatabaseAccountName,
	)
	if pp.GetString("autoFailoverEnabled") == enabled {
		p["enableAutomaticFailover"] = true
	} else {
		p["enableAutomaticFailover"] = false
	}
	if pp.GetString("multipleWriteRegionsEnabled") == enabled {
		p["enableMultipleWriteLocations"] = true
	} else {
		p["enableMultipleWriteLocations"] = false
	}

	filters := []string{}
	ipFilters := pp.GetObject("ipFilters")
	if ipFilters.GetString("allowAzure") == disabled &&
		ipFilters.GetString("allowPortal") != disabled {
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
	} else {
		filters = append(filters, "0.0.0.0")
	}
	filters = append(filters, ipFilters.GetStringArray("allowedIPRanges")...)
	if len(filters) > 0 {
		p["ipFilters"] = strings.Join(filters, ",")
	}
	p["consistencyPolicy"] = pp.GetObject("consistencyPolicy").Data
	return p, nil
}

type readLocationInformation struct {
	ID       string
	Location string
	Priority int
}

func buildReadLocationInformation(
	readLocations []string,
	databaseAccountName string,
) []readLocationInformation {
	informations := []readLocationInformation{}
	for i := range readLocations {
		information := readLocationInformation{}
		information.Location = readLocations[i]
		information.Priority = i
		information.ID = generateIDForReadLocation(
			databaseAccountName,
			readLocations[i],
		)
		informations = append(informations, information)
	}
	return informations
}

// Because the database account name is determined during provision step,
// it is possible the database account name has length 36 (the length of UUID).
// And in update step, a read location whose name is longer than 14 character
// may be provided, in which case will cause bug in updating.
// To avoid this case, we truncate the id of read locations to 50 characters.
// It is tested that truncating the read region id won't affect user's usage.
func generateIDForReadLocation(
	databaseAccountName string,
	location string,
) string {
	locationID := fmt.Sprintf("%s-%s", databaseAccountName, location)
	if len(locationID) > 50 {
		locationID = locationID[0:50]
	}
	return locationID
}

func isCreationSucceeded(
	desiredLocations []string,
	currentLocations []cosmosSDK.Location,
) bool {
	succeededLocations := make(map[string]bool)
	for i := range currentLocations {
		state := *(currentLocations[i].ProvisioningState)
		// If the status of any region is not succeeded, we haven't finished
		// the process and directly return false
		if strings.ToLower(state) != succeeded {
			return false
		}
		locationName := *(currentLocations[i].LocationName)
		locationName = strings.Replace(locationName, " ", "", -1)
		locationName = strings.ToLower(locationName)
		succeededLocations[locationName] = true
	}

	for _, location := range desiredLocations {
		if !succeededLocations[location] {
			return false
		}
	}
	return true
}
