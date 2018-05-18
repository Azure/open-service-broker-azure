package api

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
)

type errorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}

var responseAsyncRequired = []byte(
	`{ "error": "AsyncRequired", "description": "This service plan requires ` +
		`client support for asynchronous service operations." }`,
)

func generateAsyncRequiredResponse() []byte {
	return responseAsyncRequired
}

var responseServiceIDRequired = []byte(
	`{ "error": "ServiceIdRequired", "description": "service_id is a required ` +
		`field." }`,
)

func generateServiceIDRequiredResponse() []byte {
	return responseServiceIDRequired
}

var responsePlanIDRequired = []byte(
	`{ "error": "PlanIdRequired", "description": "plan_id is a required ` +
		`field." }`)

func generatePlanIDRequiredResponse() []byte {
	return responsePlanIDRequired
}

var responseInvalidServiceID = []byte(
	`{ "error": "InvalidServiceId", "description": "The provided service_id is ` +
		`invalid." }`,
)

func generateInvalidServiceIDResponse() []byte {
	return responseInvalidServiceID
}

var responseInvalidPlanID = []byte(
	`{ "error": "InvalidPlanId", "description": "The provided plan_id is ` +
		`invalid." }`,
)

func generateInvalidPlanIDResponse() []byte {
	return responseInvalidPlanID
}

var responseProvisioningAccepted = []byte(
	fmt.Sprintf(`{ "operation": "%s" }`, OperationProvisioning),
)

func generateProvisionAcceptedResponse() []byte {
	return responseProvisioningAccepted
}

var responseUpdatingAccepted = []byte(
	fmt.Sprintf(`{ "operation": "%s" }`, OperationUpdating),
)

func generateUpdateAcceptedResponse() []byte {
	return responseUpdatingAccepted
}

var responseDeprovisioningAccepted = []byte(
	fmt.Sprintf(`{ "operation": "%s" }`, OperationDeprovisioning),
)

func generateDeprovisionAcceptedResponse() []byte {
	return responseDeprovisioningAccepted
}

var responseInProgress = []byte(
	fmt.Sprintf(`{ "state": "%s" }`, OperationStateInProgress),
)

func generateOperationInProgressResponse() []byte {
	return responseInProgress
}

var responseSucceeded = []byte(
	fmt.Sprintf(`{ "state": "%s" }`, OperationStateSucceeded),
)

func generateOperationSucceededResponse() []byte {
	return responseSucceeded
}

var responseFailed = []byte(
	fmt.Sprintf(`{ "state": "%s" }`, OperationStateFailed),
)

func generateOperationFailedResponse() []byte {
	return responseFailed
}

var responseEmptyJSON = []byte("{}")

func generateEmptyResponse() []byte {
	return responseEmptyJSON
}

var responseConflict = []byte(`{ "description": "A service instance exists ` +
	`with the specified service id" }`)

func generateConflictResponse() []byte {
	return responseConflict
}

// The following are custom to this broker-- i.e. not explicitly declared by
// the OSB spec

var responseMalformedRequestBody = []byte(
	`{ "error": "MalformedRequestBody", "description": "The request body did ` +
		`not contain valid, well-formed JSON" }`,
)

func generateMalformedRequestResponse() []byte {
	return responseMalformedRequestBody
}

var responseOperationRequired = []byte(
	`{ "error": "OperationRequired", "description": "The polling request did ` +
		`not include the required operation query parameter" }`,
)

func generateOperationRequiredResponse() []byte {
	return responseOperationRequired
}

var responseOperationInvalid = []byte(
	`{ "error": "OperationInvalid", "description": "The polling request ` +
		`included an invalid value for the required operation query parameter" }`,
)

func generateOperationInvalidResponse() []byte {
	return responseOperationInvalid
}

var responseTagsMalformed = []byte(`{ "error": "MalformedRequestTags", ` +
	`"description": The provided tags were not well-formed JSON" }`,
)

func generateMalformedTagsResponse() []byte {
	return responseTagsMalformed
}

// var responseIncorrectRequestBody = []byte(`{ "error": "MalformedRequest", ` +
// 	`"description": The provided request did not match what was expected ` +
// 	`for the service" }`,
// )

// func generateInvalidRequestResponse() []byte {
// 	return responseIncorrectRequestBody
// }

var validationFailedGenericResponse = []byte(
	`{ "error" : "ValidationFailure", ` +
		`"description" : "Failed to validate request" }`,
)
var responseValidationFailedTemplate = `The value provided for %s is ` +
	`invalid: %s`

func generateValidationFailedResponse(
	validationErr *service.ValidationError,
) []byte {

	response := errorResponse{
		Error: "ValidationError",
		Description: fmt.Sprintf(
			responseValidationFailedTemplate,
			validationErr.Field,
			validationErr.Issue,
		),
	}
	responseBody, err := json.Marshal(response)
	if err != nil {
		log.WithFields(
			log.Fields{
				"field": validationErr.Field,
				"issue": validationErr.Issue,
			},
		).Error("Error generating validation error response")
		// There was a failure marshalling the body, so return
		// a generic validation failed message in it's place
		return validationFailedGenericResponse
	}
	return responseBody
}

var responseParentInvalid = []byte(
	`{ "error": "InvalidParent", "description": "The parentAlias provided ` +
		`refers to a service instance that failed to provision or is currently ` +
		`deprovisioning"  }`,
)

func generateParentInvalidResponse() []byte {
	return responseParentInvalid
}
