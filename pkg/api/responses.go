package api

import "fmt"

var responseAsyncRequired = []byte(
	`{ "error": "AsyncRequired", "description": "This service plan requires ` +
		`client support for asynchronous service operations." }`,
)

var responseServiceIDRequired = []byte(
	`{ "error": "ServiceIdRequired", "description": "service_id is a required ` +
		`field." }`,
)

var responsePlanIDRequired = []byte(
	`{ "error": "PlanIdRequired", "description": "plan_id is a required ` +
		`field." }`)

var responseInvalidServiceID = []byte(
	`{ "error": "InvalidServiceId", "description": "The provided service_id is ` +
		`invalid." }`,
)

var responseInvalidPlanID = []byte(
	`{ "error": "InvalidPlanId", "description": "The provided plan_id is ` +
		`invalid." }`,
)

var responseProvisioningAccepted = []byte(
	fmt.Sprintf(`{ "operation": "%s" }`, OperationProvisioning),
)

var responseDeprovisioningAccepted = []byte(
	fmt.Sprintf(`{ "operation": "%s" }`, OperationDeprovisioning),
)

var responseInProgress = []byte(
	fmt.Sprintf(`{ "state": "%s" }`, OperationStateInProgress),
)

var responseSucceeded = []byte(
	fmt.Sprintf(`{ "state": "%s" }`, OperationStateSucceeded),
)

var responseFailed = []byte(
	fmt.Sprintf(`{ "state": "%s" }`, OperationStateGone),
)

var responseEmptyJSON = []byte("{}")

// The following are custom to this broker-- i.e. not explicitly declared by
// the OSB spec

var responseMalformedRequestBody = []byte(
	`{ "error": "MalformedRequestBody", "description": "The request body did ` +
		`not contain valid, well-formed JSON" }`,
)

var responseOperationRequired = []byte(
	`{ "error": "OperationRequired", "description": "The polling request did ` +
		`not include the required operation query parameter" }`,
)

var responseOperationInvalid = []byte(
	`{ "error": "OperationInvalid", "description": "The polling request ` +
		`included an invalid value for the required operation query parameter" }`,
)
