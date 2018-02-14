package account

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"net/http"
)

// StorageAccountsClient is the creates an Azure Data Lake Analytics account management client.
type StorageAccountsClient struct {
	ManagementClient
}

// NewStorageAccountsClient creates an instance of the StorageAccountsClient client.
func NewStorageAccountsClient(subscriptionID string) StorageAccountsClient {
	return NewStorageAccountsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewStorageAccountsClientWithBaseURI creates an instance of the StorageAccountsClient client.
func NewStorageAccountsClientWithBaseURI(baseURI string, subscriptionID string) StorageAccountsClient {
	return StorageAccountsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Add updates the specified Data Lake Analytics account to add an Azure Storage account.
//
// resourceGroupName is the name of the Azure resource group that contains the Data Lake Analytics account. accountName
// is the name of the Data Lake Analytics account to which to add the Azure Storage account. storageAccountName is the
// name of the Azure Storage account to add parameters is the parameters containing the access key and optional suffix
// for the Azure Storage Account.
func (client StorageAccountsClient) Add(resourceGroupName string, accountName string, storageAccountName string, parameters AddStorageAccountParameters) (result autorest.Response, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.StorageAccountProperties", Name: validation.Null, Rule: true,
				Chain: []validation.Constraint{{Target: "parameters.StorageAccountProperties.AccessKey", Name: validation.Null, Rule: true, Chain: nil}}}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "account.StorageAccountsClient", "Add")
	}

	req, err := client.AddPreparer(resourceGroupName, accountName, storageAccountName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "Add", nil, "Failure preparing request")
		return
	}

	resp, err := client.AddSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "Add", resp, "Failure sending request")
		return
	}

	result, err = client.AddResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "Add", resp, "Failure responding to request")
	}

	return
}

// AddPreparer prepares the Add request.
func (client StorageAccountsClient) AddPreparer(resourceGroupName string, accountName string, storageAccountName string, parameters AddStorageAccountParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":        autorest.Encode("path", accountName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"storageAccountName": autorest.Encode("path", storageAccountName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/StorageAccounts/{storageAccountName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// AddSender sends the Add request. The method will close the
// http.Response Body if it receives an error.
func (client StorageAccountsClient) AddSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoRetryWithRegistration(client.Client))
}

// AddResponder handles the response to the Add request. The method always
// closes the http.Response Body.
func (client StorageAccountsClient) AddResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Delete updates the specified Data Lake Analytics account to remove an Azure Storage account.
//
// resourceGroupName is the name of the Azure resource group that contains the Data Lake Analytics account. accountName
// is the name of the Data Lake Analytics account from which to remove the Azure Storage account. storageAccountName is
// the name of the Azure Storage account to remove
func (client StorageAccountsClient) Delete(resourceGroupName string, accountName string, storageAccountName string) (result autorest.Response, err error) {
	req, err := client.DeletePreparer(resourceGroupName, accountName, storageAccountName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client StorageAccountsClient) DeletePreparer(resourceGroupName string, accountName string, storageAccountName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":        autorest.Encode("path", accountName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"storageAccountName": autorest.Encode("path", storageAccountName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/StorageAccounts/{storageAccountName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client StorageAccountsClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoRetryWithRegistration(client.Client))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client StorageAccountsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets the specified Azure Storage account linked to the given Data Lake Analytics account.
//
// resourceGroupName is the name of the Azure resource group that contains the Data Lake Analytics account. accountName
// is the name of the Data Lake Analytics account from which to retrieve Azure storage account details.
// storageAccountName is the name of the Azure Storage account for which to retrieve the details.
func (client StorageAccountsClient) Get(resourceGroupName string, accountName string, storageAccountName string) (result StorageAccountInfo, err error) {
	req, err := client.GetPreparer(resourceGroupName, accountName, storageAccountName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client StorageAccountsClient) GetPreparer(resourceGroupName string, accountName string, storageAccountName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":        autorest.Encode("path", accountName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"storageAccountName": autorest.Encode("path", storageAccountName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/StorageAccounts/{storageAccountName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client StorageAccountsClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client StorageAccountsClient) GetResponder(resp *http.Response) (result StorageAccountInfo, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// GetStorageContainer gets the specified Azure Storage container associated with the given Data Lake Analytics and
// Azure Storage accounts.
//
// resourceGroupName is the name of the Azure resource group that contains the Data Lake Analytics account. accountName
// is the name of the Data Lake Analytics account for which to retrieve blob container. storageAccountName is the name
// of the Azure storage account from which to retrieve the blob container. containerName is the name of the Azure
// storage container to retrieve
func (client StorageAccountsClient) GetStorageContainer(resourceGroupName string, accountName string, storageAccountName string, containerName string) (result StorageContainer, err error) {
	req, err := client.GetStorageContainerPreparer(resourceGroupName, accountName, storageAccountName, containerName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "GetStorageContainer", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetStorageContainerSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "GetStorageContainer", resp, "Failure sending request")
		return
	}

	result, err = client.GetStorageContainerResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "GetStorageContainer", resp, "Failure responding to request")
	}

	return
}

// GetStorageContainerPreparer prepares the GetStorageContainer request.
func (client StorageAccountsClient) GetStorageContainerPreparer(resourceGroupName string, accountName string, storageAccountName string, containerName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":        autorest.Encode("path", accountName),
		"containerName":      autorest.Encode("path", containerName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"storageAccountName": autorest.Encode("path", storageAccountName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/StorageAccounts/{storageAccountName}/Containers/{containerName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetStorageContainerSender sends the GetStorageContainer request. The method will close the
// http.Response Body if it receives an error.
func (client StorageAccountsClient) GetStorageContainerSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetStorageContainerResponder handles the response to the GetStorageContainer request. The method always
// closes the http.Response Body.
func (client StorageAccountsClient) GetStorageContainerResponder(resp *http.Response) (result StorageContainer, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByAccount gets the first page of Azure Storage accounts, if any, linked to the specified Data Lake Analytics
// account. The response includes a link to the next page, if any.
//
// resourceGroupName is the name of the Azure resource group that contains the Data Lake Analytics account. accountName
// is the name of the Data Lake Analytics account for which to list Azure Storage accounts. filter is the OData filter.
// Optional. top is the number of items to return. Optional. skip is the number of items to skip over before returning
// elements. Optional. selectParameter is oData Select statement. Limits the properties on each entry to just those
// requested, e.g. Categories?$select=CategoryName,Description. Optional. orderby is orderBy clause. One or more
// comma-separated expressions with an optional "asc" (the default) or "desc" depending on the order you'd like the
// values sorted, e.g. Categories?$orderby=CategoryName desc. Optional. count is the Boolean value of true or false to
// request a count of the matching resources included with the resources in the response, e.g. Categories?$count=true.
// Optional.
func (client StorageAccountsClient) ListByAccount(resourceGroupName string, accountName string, filter string, top *int32, skip *int32, selectParameter string, orderby string, count *bool) (result DataLakeAnalyticsAccountListStorageAccountsResult, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: top,
			Constraints: []validation.Constraint{{Target: "top", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "top", Name: validation.InclusiveMinimum, Rule: 1, Chain: nil}}}}},
		{TargetValue: skip,
			Constraints: []validation.Constraint{{Target: "skip", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "skip", Name: validation.InclusiveMinimum, Rule: 1, Chain: nil}}}}}}); err != nil {
		return result, validation.NewErrorWithValidationError(err, "account.StorageAccountsClient", "ListByAccount")
	}

	req, err := client.ListByAccountPreparer(resourceGroupName, accountName, filter, top, skip, selectParameter, orderby, count)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListByAccount", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByAccountSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListByAccount", resp, "Failure sending request")
		return
	}

	result, err = client.ListByAccountResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListByAccount", resp, "Failure responding to request")
	}

	return
}

// ListByAccountPreparer prepares the ListByAccount request.
func (client StorageAccountsClient) ListByAccountPreparer(resourceGroupName string, accountName string, filter string, top *int32, skip *int32, selectParameter string, orderby string, count *bool) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	}
	if skip != nil {
		queryParameters["$skip"] = autorest.Encode("query", *skip)
	}
	if len(selectParameter) > 0 {
		queryParameters["$select"] = autorest.Encode("query", selectParameter)
	}
	if len(orderby) > 0 {
		queryParameters["$orderby"] = autorest.Encode("query", orderby)
	}
	if count != nil {
		queryParameters["$count"] = autorest.Encode("query", *count)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/StorageAccounts/", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListByAccountSender sends the ListByAccount request. The method will close the
// http.Response Body if it receives an error.
func (client StorageAccountsClient) ListByAccountSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListByAccountResponder handles the response to the ListByAccount request. The method always
// closes the http.Response Body.
func (client StorageAccountsClient) ListByAccountResponder(resp *http.Response) (result DataLakeAnalyticsAccountListStorageAccountsResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByAccountNextResults retrieves the next set of results, if any.
func (client StorageAccountsClient) ListByAccountNextResults(lastResults DataLakeAnalyticsAccountListStorageAccountsResult) (result DataLakeAnalyticsAccountListStorageAccountsResult, err error) {
	req, err := lastResults.DataLakeAnalyticsAccountListStorageAccountsResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListByAccount", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListByAccountSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListByAccount", resp, "Failure sending next results request")
	}

	result, err = client.ListByAccountResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListByAccount", resp, "Failure responding to next results request")
	}

	return
}

// ListByAccountComplete gets all elements from the list without paging.
func (client StorageAccountsClient) ListByAccountComplete(resourceGroupName string, accountName string, filter string, top *int32, skip *int32, selectParameter string, orderby string, count *bool, cancel <-chan struct{}) (<-chan StorageAccountInfo, <-chan error) {
	resultChan := make(chan StorageAccountInfo)
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			close(resultChan)
			close(errChan)
		}()
		list, err := client.ListByAccount(resourceGroupName, accountName, filter, top, skip, selectParameter, orderby, count)
		if err != nil {
			errChan <- err
			return
		}
		if list.Value != nil {
			for _, item := range *list.Value {
				select {
				case <-cancel:
					return
				case resultChan <- item:
					// Intentionally left blank
				}
			}
		}
		for list.NextLink != nil {
			list, err = client.ListByAccountNextResults(list)
			if err != nil {
				errChan <- err
				return
			}
			if list.Value != nil {
				for _, item := range *list.Value {
					select {
					case <-cancel:
						return
					case resultChan <- item:
						// Intentionally left blank
					}
				}
			}
		}
	}()
	return resultChan, errChan
}

// ListSasTokens gets the SAS token associated with the specified Data Lake Analytics and Azure Storage account and
// container combination.
//
// resourceGroupName is the name of the Azure resource group that contains the Data Lake Analytics account. accountName
// is the name of the Data Lake Analytics account from which an Azure Storage account's SAS token is being requested.
// storageAccountName is the name of the Azure storage account for which the SAS token is being requested.
// containerName is the name of the Azure storage container for which the SAS token is being requested.
func (client StorageAccountsClient) ListSasTokens(resourceGroupName string, accountName string, storageAccountName string, containerName string) (result ListSasTokensResult, err error) {
	req, err := client.ListSasTokensPreparer(resourceGroupName, accountName, storageAccountName, containerName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListSasTokens", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSasTokensSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListSasTokens", resp, "Failure sending request")
		return
	}

	result, err = client.ListSasTokensResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListSasTokens", resp, "Failure responding to request")
	}

	return
}

// ListSasTokensPreparer prepares the ListSasTokens request.
func (client StorageAccountsClient) ListSasTokensPreparer(resourceGroupName string, accountName string, storageAccountName string, containerName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":        autorest.Encode("path", accountName),
		"containerName":      autorest.Encode("path", containerName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"storageAccountName": autorest.Encode("path", storageAccountName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/StorageAccounts/{storageAccountName}/Containers/{containerName}/listSasTokens", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListSasTokensSender sends the ListSasTokens request. The method will close the
// http.Response Body if it receives an error.
func (client StorageAccountsClient) ListSasTokensSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListSasTokensResponder handles the response to the ListSasTokens request. The method always
// closes the http.Response Body.
func (client StorageAccountsClient) ListSasTokensResponder(resp *http.Response) (result ListSasTokensResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListSasTokensNextResults retrieves the next set of results, if any.
func (client StorageAccountsClient) ListSasTokensNextResults(lastResults ListSasTokensResult) (result ListSasTokensResult, err error) {
	req, err := lastResults.ListSasTokensResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListSasTokens", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListSasTokensSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListSasTokens", resp, "Failure sending next results request")
	}

	result, err = client.ListSasTokensResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListSasTokens", resp, "Failure responding to next results request")
	}

	return
}

// ListSasTokensComplete gets all elements from the list without paging.
func (client StorageAccountsClient) ListSasTokensComplete(resourceGroupName string, accountName string, storageAccountName string, containerName string, cancel <-chan struct{}) (<-chan SasTokenInfo, <-chan error) {
	resultChan := make(chan SasTokenInfo)
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			close(resultChan)
			close(errChan)
		}()
		list, err := client.ListSasTokens(resourceGroupName, accountName, storageAccountName, containerName)
		if err != nil {
			errChan <- err
			return
		}
		if list.Value != nil {
			for _, item := range *list.Value {
				select {
				case <-cancel:
					return
				case resultChan <- item:
					// Intentionally left blank
				}
			}
		}
		for list.NextLink != nil {
			list, err = client.ListSasTokensNextResults(list)
			if err != nil {
				errChan <- err
				return
			}
			if list.Value != nil {
				for _, item := range *list.Value {
					select {
					case <-cancel:
						return
					case resultChan <- item:
						// Intentionally left blank
					}
				}
			}
		}
	}()
	return resultChan, errChan
}

// ListStorageContainers lists the Azure Storage containers, if any, associated with the specified Data Lake Analytics
// and Azure Storage account combination. The response includes a link to the next page of results, if any.
//
// resourceGroupName is the name of the Azure resource group that contains the Data Lake Analytics account. accountName
// is the name of the Data Lake Analytics account for which to list Azure Storage blob containers. storageAccountName
// is the name of the Azure storage account from which to list blob containers.
func (client StorageAccountsClient) ListStorageContainers(resourceGroupName string, accountName string, storageAccountName string) (result ListStorageContainersResult, err error) {
	req, err := client.ListStorageContainersPreparer(resourceGroupName, accountName, storageAccountName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListStorageContainers", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListStorageContainersSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListStorageContainers", resp, "Failure sending request")
		return
	}

	result, err = client.ListStorageContainersResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListStorageContainers", resp, "Failure responding to request")
	}

	return
}

// ListStorageContainersPreparer prepares the ListStorageContainers request.
func (client StorageAccountsClient) ListStorageContainersPreparer(resourceGroupName string, accountName string, storageAccountName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":        autorest.Encode("path", accountName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"storageAccountName": autorest.Encode("path", storageAccountName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/StorageAccounts/{storageAccountName}/Containers", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// ListStorageContainersSender sends the ListStorageContainers request. The method will close the
// http.Response Body if it receives an error.
func (client StorageAccountsClient) ListStorageContainersSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListStorageContainersResponder handles the response to the ListStorageContainers request. The method always
// closes the http.Response Body.
func (client StorageAccountsClient) ListStorageContainersResponder(resp *http.Response) (result ListStorageContainersResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListStorageContainersNextResults retrieves the next set of results, if any.
func (client StorageAccountsClient) ListStorageContainersNextResults(lastResults ListStorageContainersResult) (result ListStorageContainersResult, err error) {
	req, err := lastResults.ListStorageContainersResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListStorageContainers", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}

	resp, err := client.ListStorageContainersSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListStorageContainers", resp, "Failure sending next results request")
	}

	result, err = client.ListStorageContainersResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "ListStorageContainers", resp, "Failure responding to next results request")
	}

	return
}

// ListStorageContainersComplete gets all elements from the list without paging.
func (client StorageAccountsClient) ListStorageContainersComplete(resourceGroupName string, accountName string, storageAccountName string, cancel <-chan struct{}) (<-chan StorageContainer, <-chan error) {
	resultChan := make(chan StorageContainer)
	errChan := make(chan error, 1)
	go func() {
		defer func() {
			close(resultChan)
			close(errChan)
		}()
		list, err := client.ListStorageContainers(resourceGroupName, accountName, storageAccountName)
		if err != nil {
			errChan <- err
			return
		}
		if list.Value != nil {
			for _, item := range *list.Value {
				select {
				case <-cancel:
					return
				case resultChan <- item:
					// Intentionally left blank
				}
			}
		}
		for list.NextLink != nil {
			list, err = client.ListStorageContainersNextResults(list)
			if err != nil {
				errChan <- err
				return
			}
			if list.Value != nil {
				for _, item := range *list.Value {
					select {
					case <-cancel:
						return
					case resultChan <- item:
						// Intentionally left blank
					}
				}
			}
		}
	}()
	return resultChan, errChan
}

// Update updates the Data Lake Analytics account to replace Azure Storage blob account details, such as the access key
// and/or suffix.
//
// resourceGroupName is the name of the Azure resource group that contains the Data Lake Analytics account. accountName
// is the name of the Data Lake Analytics account to modify storage accounts in storageAccountName is the Azure Storage
// account to modify parameters is the parameters containing the access key and suffix to update the storage account
// with, if any. Passing nothing results in no change.
func (client StorageAccountsClient) Update(resourceGroupName string, accountName string, storageAccountName string, parameters *UpdateStorageAccountParameters) (result autorest.Response, err error) {
	req, err := client.UpdatePreparer(resourceGroupName, accountName, storageAccountName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "Update", nil, "Failure preparing request")
		return
	}

	resp, err := client.UpdateSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "Update", resp, "Failure sending request")
		return
	}

	result, err = client.UpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.StorageAccountsClient", "Update", resp, "Failure responding to request")
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client StorageAccountsClient) UpdatePreparer(resourceGroupName string, accountName string, storageAccountName string, parameters *UpdateStorageAccountParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":        autorest.Encode("path", accountName),
		"resourceGroupName":  autorest.Encode("path", resourceGroupName),
		"storageAccountName": autorest.Encode("path", storageAccountName),
		"subscriptionId":     autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeAnalytics/accounts/{accountName}/StorageAccounts/{storageAccountName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	if parameters != nil {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithJSON(parameters))
	}
	return preparer.Prepare(&http.Request{})
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client StorageAccountsClient) UpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client,
		req,
		azure.DoRetryWithRegistration(client.Client))
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client StorageAccountsClient) UpdateResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}
