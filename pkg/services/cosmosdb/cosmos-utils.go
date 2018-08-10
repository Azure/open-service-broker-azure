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

// The deployment will return success once the write region is created,
// ignoring the status of read regions , so we must implement detection logic
// by ourselves.
func (c *cosmosAccountManager) waitForReadLocationsReady(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*cosmosdbInstanceDetails)
	resourceGroupName := instance.ProvisioningParameters.GetString("resourceGroup")
	accountName := dt.DatabaseAccountName
	databaseAccountClient := c.databaseAccountsClient

	err := pollingUntilReadLocationsReady(
		ctx,
		resourceGroupName,
		accountName,
		databaseAccountClient,
	)
	if err != nil {
		return nil, err
	}
	return dt, nil
}

// For sqlAllInOneManager, the real type of `instance.Details` is
// `*sqlAllInOneInstanceDetails`, so type assertion must be changed.
// Expect type assertion, this function is totally the same as previous one.
func (s *sqlAllInOneManager) waitForReadLocationsReady(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*sqlAllInOneInstanceDetails)
	resourceGroupName := instance.ProvisioningParameters.GetString("resourceGroup")
	accountName := dt.DatabaseAccountName
	databaseAccountClient := s.databaseAccountsClient

	err := pollingUntilReadLocationsReady(
		ctx,
		resourceGroupName,
		accountName,
		databaseAccountClient,
	)
	if err != nil {
		return nil, err
	}
	return dt, nil
}

const succeeded = "succeeded"

// For now, this method will return on either context is cancelled or
// every region's state is "succeeded" in seven consecutive check.
// The reason why we need seven consecutive check is that
// the read region is created one by one, there is a small gap between
// the finishment of previous creation and the start of the next creation.
// By this check, we can detect gaps shorter than 1 mintue,
// and report success within 70 seconds after completion.
func pollingUntilReadLocationsReady(
	ctx context.Context,
	resourceGroupName string,
	accountName string,
	databaseAccountClient cosmosSDK.DatabaseAccountsClient,
) error {
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Seven comsecutive check is needed to report success.
	const confirmNumberOfTimes = 7

	ticker := time.NewTicker(time.Second * 10)
	previousSucceededTimes := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			result, err := databaseAccountClient.Get(
				childCtx,
				resourceGroupName,
				accountName,
			)
			if err != nil {
				return err
			}

			//Check whether every read location's state is "Succeeded"
			allSucceed := true
			readLocations := *(result.DatabaseAccountProperties.ReadLocations)
			for i := range readLocations {
				state := *(readLocations[i].ProvisioningState)
				if strings.ToLower(state) != succeeded {
					allSucceed = false
					previousSucceededTimes = 0
					break
				}
			}
			if allSucceed && previousSucceededTimes >= confirmNumberOfTimes {
				return nil
			} else if allSucceed {
				previousSucceededTimes++
			}
		}
	}
}

func validateReadLocations(
	context string,
	regions []string,
) error {
	for i := range regions {
		region := regions[i]
		if !allowedReadLocations[region] {
			return service.NewValidationError(
				fmt.Sprintf("%s.allowedReadLocations", context),
				fmt.Sprintf("given region %s is not allowed", region),
			)
		}
	}
	return nil
}

// Allowed CosmosDB read locations, it is different from Azure regions.
// We use a map here to record all allowed regions.
var allowedReadLocations = map[string]bool{
	"westus2":            true,
	"westus":             true,
	"southcentralus":     true,
	"centralus":          true,
	"northcentralus":     true,
	"canadacentral":      true,
	"eastus":             true,
	"eastus2":            true,
	"canadaeast":         true,
	"brazilsouth":        true,
	"northeurope":        true,
	"ukwest":             true,
	"uksouth":            true,
	"francecentral":      true,
	"westeurope":         true,
	"westindia":          true,
	"centralindia":       true,
	"southindia":         true,
	"southeastasia":      true,
	"eastasia":           true,
	"koreacentral":       true,
	"koreasouth":         true,
	"japaneast":          true,
	"japanwest":          true,
	"australiasoutheast": true,
	"australiaeast":      true,
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
	p["readLocations"] = readLocations
	if pp.GetString("autoFailoverEnabled") == enabled {
		p["enableAutomaticFailover"] = true
	} else {
		p["enableAutomaticFailover"] = false
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
