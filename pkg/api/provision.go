package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/Azure/azure-service-broker/pkg/async/model"
	"github.com/Azure/azure-service-broker/pkg/azure"
	"github.com/Azure/azure-service-broker/pkg/service"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/satori/uuid"
)

func (s *server) provision(w http.ResponseWriter, r *http.Request) {
	instanceID := mux.Vars(r)["instance_id"]

	logFields := log.Fields{
		"instanceID": instanceID,
	}

	log.WithFields(logFields).Debug("received provisioning request")

	// This broker provisions everything asynchronously. If a client doesn't
	// explicitly indicate that they will accept an incomplete result, the
	// spec says to respond with a 422
	acceptsIncompleteStr := r.URL.Query().Get("accepts_incomplete")
	if acceptsIncompleteStr == "" {
		logFields["parameter"] = "accepts_incomplete=true" // nolint: goconst
		log.WithFields(logFields).Debug(
			"bad provisioning request: request is missing required query parameter",
		)
		s.writeResponse(w, http.StatusUnprocessableEntity, responseAsyncRequired)
		return
	}
	acceptsIncomplete, err := strconv.ParseBool(acceptsIncompleteStr)
	if err != nil || !acceptsIncomplete {
		logFields["accepts_incomplete"] = acceptsIncompleteStr
		log.WithFields(logFields).Debug(
			`bad provisioning request: query parameter has invalid value; only ` +
				`"true" is accepted`,
		)
		s.writeResponse(w, http.StatusUnprocessableEntity, responseAsyncRequired)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-provisioning error: error reading request body",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	defer r.Body.Close() // nolint: errcheck

	provisioningRequest, err := NewProvisioningRequestFromJSON(bodyBytes)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Debug(
			"bad provisioning request: error unmarshaling request body",
		)
		// krancour: Choosing to interpret this scenario as a bad request, as a
		// valid request, obviously contains valid, well-formed JSON
		// TODO: Write a more detailed response
		s.writeResponse(w, http.StatusBadRequest, responseEmptyJSON)
		return
	}

	serviceID := provisioningRequest.ServiceID
	if serviceID == "" {
		logFields["field"] = "service_id"
		log.WithFields(logFields).Debug(
			"bad provisioning request: required request body field is missing",
		)
		s.writeResponse(w, http.StatusBadRequest, responseServiceIDRequired)
		return
	}

	planID := provisioningRequest.PlanID
	if planID == "" {
		logFields["field"] = "plan_id"
		log.WithFields(logFields).Debug(
			"bad provisioning request: required request body field is missing",
		)
		s.writeResponse(w, http.StatusBadRequest, responsePlanIDRequired)
		return
	}

	svc, ok := s.catalog.GetService(serviceID)
	if !ok {
		logFields["serviceID"] = serviceID
		log.WithFields(logFields).Debug(
			"bad provisioning request: invalid serviceID",
		)
		s.writeResponse(w, http.StatusBadRequest, responseInvalidServiceID)
		return
	}

	_, ok = svc.GetPlan(planID)
	if !ok {
		logFields["serviceID"] = serviceID
		logFields["planID"] = planID
		log.WithFields(logFields).Debug(
			"bad provisioning request: invalid planID for service",
		)
		s.writeResponse(w, http.StatusBadRequest, responseInvalidPlanID)
		return
	}

	module, ok := s.modules[serviceID]
	if !ok {
		// We already validated that the serviceID and planID are legitimate. If
		// we don't find a module that handles the service, something is really
		// wrong.
		logFields["serviceID"] = serviceID
		log.WithFields(logFields).Error(
			"pre-provisioning error: no module found for service",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	// Unpack the parameter map in the request to structs

	// Standard params (those common to all services) first
	standardProvisioningParameters := service.StandardProvisioningParameters{}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  &standardProvisioningParameters,
	})
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"error building parameter map decoder",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	err = decoder.Decode(provisioningRequest.Parameters)
	if err != nil {
		log.WithFields(logFields).Debug(
			"bad provisioning request: error decoding parameter map",
		)
		// krancour: Choosing to interpret this scenario as a bad request since the
		// probable cause would be disagreement between provided and expected types
		s.writeResponse(w, http.StatusBadRequest, responseEmptyJSON)
		return
	}

	// Then service-specific parameters
	provisioningParameters := module.GetEmptyProvisioningParameters()
	decoderConfig := &mapstructure.DecoderConfig{
		TagName: "json",
		Result:  provisioningParameters,
	}
	decoder, err = mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"error building parameter map decoder",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	err = decoder.Decode(provisioningRequest.Parameters)
	if err != nil {
		log.WithFields(logFields).Debug(
			"bad provisioning request: error decoding parameter map",
		)
		// krancour: Choosing to interpret this scenario as a bad request since the
		// probable cause would be disagreement between provided and expected types
		s.writeResponse(w, http.StatusBadRequest, responseEmptyJSON)
		return
	}

	instance, ok, err := s.store.GetInstance(instanceID)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-provisioning error: error retrieving instance by id",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	if ok {
		// We land in here if an existing instance was found-- the OSB spec
		// obligates us to compare this instance to the one that was requested and
		// respond with 200 if they're identical or 409 otherwise. It actually seems
		// best to compare REQUESTS instead because instance objects also contain
		// provisioning context and other status information. So, let's reverse
		// engineer a request from the existing instance then compare it to the
		// current request.
		//
		// Two requests are the same if they are for the same serviceID, the same,
		// planID, and both standard and module-specific provisioning parameters
		// are deeply equal.
		previousProvisioningParameters := module.GetEmptyProvisioningParameters()
		if err = instance.GetProvisioningParameters(
			previousProvisioningParameters,
			s.codec,
		); err != nil {
			logFields["error"] = err
			log.WithFields(logFields).Error(
				"pre-provisioning error: error decoding persisted " +
					"provisioningParameters",
			)
			s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
			return
		}
		if instance.ServiceID == serviceID &&
			instance.PlanID == planID &&
			reflect.DeepEqual(
				instance.StandardProvisioningParameters,
				standardProvisioningParameters,
			) &&
			reflect.DeepEqual(
				previousProvisioningParameters,
				provisioningParameters,
			) {
			// Per the spec, if fully provisioned, respond with a 200, else a 202.
			// Filling in a gap in the spec-- if the status is anything else, we'll
			// choose to respond with a 409
			switch instance.Status {
			case service.InstanceStateProvisioning:
				s.writeResponse(w, http.StatusAccepted, responseProvisioningAccepted)
				return
			case service.InstanceStateProvisioned:
				s.writeResponse(w, http.StatusOK, responseEmptyJSON)
				return
			default:
				// TODO: Write a more detailed response
				s.writeResponse(w, http.StatusConflict, responseEmptyJSON)
				return
			}
		}

		// We land in here if an existing instance was found, but its atrributes
		// vary from what was requested. The spec requires us to respond with a
		// 409
		s.writeResponse(w, http.StatusConflict, responseEmptyJSON)
		return
	}

	// If we get to here, we need to provision a new instance.

	// Start by validating all the standard provisioning parameters
	err = validateStandardProvisioningParameters(standardProvisioningParameters)
	if err != nil {
		s.handlePossibleValidationError(err, w, logFields)
		return
	}

	// Then validate module-specific provisioning parameters
	err = module.ValidateProvisioningParameters(provisioningParameters)
	if err != nil {
		s.handlePossibleValidationError(err, w, logFields)
		return
	}

	provisioner, err := module.GetProvisioner(
		serviceID,
		planID,
	)
	if err != nil {
		logFields["serviceID"] = serviceID
		logFields["planID"] = provisioningRequest.PlanID
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-provisioning error: error retrieving provisioner for service and " +
				"plan",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	firstStepName, ok := provisioner.GetFirstStepName()
	if !ok {
		logFields["serviceID"] = provisioningRequest.ServiceID
		logFields["planID"] = provisioningRequest.PlanID
		log.WithFields(logFields).Error(
			"pre-provisioning error: no steps found for provisioning service and " +
				"plan",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	standardProvisioningContext := service.StandardProvisioningContext{
		Location: standardProvisioningParameters.Location,
		Tags:     standardProvisioningParameters.Tags,
	}
	if standardProvisioningParameters.ResourceGroup == "" {
		standardProvisioningContext.ResourceGroup = uuid.NewV4().String()
	} else {
		standardProvisioningContext.ResourceGroup =
			standardProvisioningParameters.ResourceGroup
	}

	instance = &service.Instance{
		InstanceID: instanceID,
		ServiceID:  provisioningRequest.ServiceID,
		PlanID:     provisioningRequest.PlanID,
		StandardProvisioningParameters: standardProvisioningParameters,
		Status: service.InstanceStateProvisioning,
		StandardProvisioningContext: standardProvisioningContext,
		Created:                     time.Now(),
	}
	if err = instance.SetProvisioningParameters(
		provisioningRequest.Parameters,
		s.codec,
	); err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"provisioning error: error encoding provisioningParameters",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	if err = instance.SetProvisioningContext(
		module.GetEmptyProvisioningContext(),
		s.codec,
	); err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"provisioning error: error encoding provisioningContext",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	if err = s.store.WriteInstance(instance); err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"provisioning error: error persisting new instance",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	task := model.NewTask(
		"provisionStep",
		map[string]string{
			"stepName":   firstStepName,
			"instanceID": instanceID,
		},
	)
	if err = s.asyncEngine.SubmitTask(task); err != nil {
		logFields["step"] = firstStepName
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"provisioning error: error submitting provisioning task",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	// If we get all the way to here, we've been successful!
	s.writeResponse(w, http.StatusAccepted, responseProvisioningAccepted)

	log.WithFields(logFields).Debug("asynchronous provisioning initiated")
}

func validateStandardProvisioningParameters(
	spp service.StandardProvisioningParameters,
) error {
	if !azure.IsValidLocation(spp.Location) {
		return service.NewValidationError(
			"location",
			fmt.Sprintf(`invalid location: "%s"`, spp.Location),
		)
	}
	return nil
}

func (s *server) handlePossibleValidationError(
	err error,
	w http.ResponseWriter,
	logFields log.Fields,
) {
	validationErr, ok := err.(*service.ValidationError)
	if ok {
		logFields["field"] = validationErr.Field
		logFields["issue"] = validationErr.Issue
		log.WithFields(logFields).Debug(
			"bad provisioning request: validation error",
		)
		// TODO: Send the correct response body-- this is a placeholder
		s.writeResponse(w, http.StatusBadRequest, responseEmptyJSON)
		return
	}
	s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
}
