package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/Azure/azure-service-broker/pkg/async/model"
	"github.com/Azure/azure-service-broker/pkg/service"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
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
			`bad provisioning request: query paramater has invalid value; only ` +
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

	rawProvisioningRequest := map[string]interface{}{}
	err = json.Unmarshal(bodyBytes, &rawProvisioningRequest)
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

	serviceIDIface, ok := rawProvisioningRequest["service_id"]
	serviceID := fmt.Sprintf("%v", serviceIDIface)
	if !ok || serviceID == "" {
		logFields["field"] = "service_id"
		log.WithFields(logFields).Debug(
			"bad provisioning request: required request body field is missing",
		)
		s.writeResponse(w, http.StatusBadRequest, responseServiceIDRequired)
		return
	}

	planIDIface, ok := rawProvisioningRequest["plan_id"]
	planID := fmt.Sprintf("%v", planIDIface)
	if !ok || planID == "" {
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

	// Now that we know what module we're dealing with, we can get an instance
	// of the module-specific type for provisioningParameters and take a second
	// pass at parsing the request body
	provisioningRequest := &ProvisioningRequest{
		Parameters: module.GetEmptyProvisioningParameters(),
	}
	err = GetProvisioningRequestFromJSON(bodyBytes, provisioningRequest)
	if err != nil {
		log.WithFields(logFields).Debug(
			"bad provisioning request: error unmarshaling request body",
		)
		// krancour: Choosing to interpret this scenario as a bad request, as a
		// valid request, obviously contains valid, well-formed JSON
		// TODO: Write a more detailed response
		s.writeResponse(w, http.StatusBadRequest, responseEmptyJSON)
		return
	}
	if provisioningRequest.Parameters == nil {
		provisioningRequest.Parameters = module.GetEmptyProvisioningParameters()
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
		previousProvisioningRequestParams := module.GetEmptyProvisioningParameters()
		if err = instance.GetProvisioningParameters(
			previousProvisioningRequestParams,
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
		previousProvisioningRequest := &ProvisioningRequest{
			ServiceID:  instance.ServiceID,
			PlanID:     instance.PlanID,
			Parameters: previousProvisioningRequestParams,
		}
		if reflect.DeepEqual(provisioningRequest, previousProvisioningRequest) {
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
	// Start by carrying out module-specific request validation
	err = module.ValidateProvisioningParameters(provisioningRequest.Parameters)
	if err != nil {
		var validationErr *service.ValidationError
		validationErr, ok = err.(*service.ValidationError)
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

	instance = &service.Instance{
		InstanceID: instanceID,
		ServiceID:  provisioningRequest.ServiceID,
		PlanID:     provisioningRequest.PlanID,
		Status:     service.InstanceStateProvisioning,
		Created:    time.Now(),
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
