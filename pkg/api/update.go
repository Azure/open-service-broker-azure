package api

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"github.com/Azure/azure-service-broker/pkg/async/model"
	"github.com/Azure/azure-service-broker/pkg/service"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func (s *server) update(w http.ResponseWriter, r *http.Request) {
	instanceID := mux.Vars(r)["instance_id"]

	logFields := log.Fields{
		"instanceID": instanceID,
	}

	log.WithFields(logFields).Debug("received updating request")

	// This broker updates everything asynchronously. If a client doesn't
	// explicitly indicate that they will accept an incomplete result, the
	// spec says to respond with a 422
	acceptsIncompleteStr := r.URL.Query().Get("accepts_incomplete")
	if acceptsIncompleteStr == "" {
		logFields["parameter"] = "accepts_incomplete=true" // nolint: goconst
		log.WithFields(logFields).Debug(
			"bad updating request: request is missing required query parameter",
		)
		s.writeResponse(w, http.StatusUnprocessableEntity, responseAsyncRequired)
		return
	}
	acceptsIncomplete, err := strconv.ParseBool(acceptsIncompleteStr)
	if err != nil || !acceptsIncomplete {
		logFields["accepts_incomplete"] = acceptsIncompleteStr
		log.WithFields(logFields).Debug(
			`bad updating request: query paramater has invalid value; only ` +
				`"true" is accepted`,
		)
		s.writeResponse(w, http.StatusUnprocessableEntity, responseAsyncRequired)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-updating error: error reading request body",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	defer r.Body.Close() // nolint: errcheck

	updatingRequest := &UpdatingRequest{}
	err = GetUpdatingRequestFromJSON(bodyBytes, updatingRequest)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Debug(
			"bad updating request: error unmarshaling request body",
		)
		// krancour: Choosing to interpret this scenario as a bad request, as a
		// valid request, obviously contains valid, well-formed JSON
		// TODO: Write a more detailed response
		s.writeResponse(w, http.StatusBadRequest, responseEmptyJSON)
		return
	}

	if updatingRequest.ServiceID == "" {
		logFields["field"] = "service_id"
		log.WithFields(logFields).Debug(
			"bad updating request: required request body field is missing",
		)
		s.writeResponse(w, http.StatusBadRequest, responseServiceIDRequired)
		return
	}

	svc, ok := s.catalog.GetService(updatingRequest.ServiceID)
	if !ok {
		logFields["serviceID"] = updatingRequest.ServiceID
		log.WithFields(logFields).Debug(
			"bad updating request: invalid serviceID",
		)
		s.writeResponse(w, http.StatusBadRequest, responseInvalidServiceID)
		return
	}

	if updatingRequest.PlanID != "" {
		_, ok = svc.GetPlan(updatingRequest.PlanID)
		if !ok {
			logFields["serviceID"] = updatingRequest.ServiceID
			logFields["planID"] = updatingRequest.PlanID
			log.WithFields(logFields).Debug(
				"bad updating request: invalid planID for service",
			)
			s.writeResponse(w, http.StatusBadRequest, responseInvalidPlanID)
			return
		}
	}

	module, ok := s.modules[updatingRequest.ServiceID]
	if !ok {
		// We already validated that the serviceID and planID are legitimate. If
		// we don't find a module that handles the service, something is really
		// wrong.
		logFields["serviceID"] = updatingRequest.ServiceID
		log.WithFields(logFields).Error(
			"pre-provisioning error: no module found for service",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	// Now that we know what module we're dealing with, we can get an instance
	// of the module-specific type for updatingParameters and take a second
	// pass at parsing the request body
	updatingRequest.Parameters = module.GetEmptyUpdatingParameters()
	err = GetUpdatingRequestFromJSON(bodyBytes, updatingRequest)
	if err != nil {
		log.WithFields(logFields).Debug(
			"bad updating request: error unmarshaling request body",
		)
		// krancour: Choosing to interpret this scenario as a bad request, as a
		// valid request, obviously contains valid, well-formed JSON
		// TODO: Write a more detailed response
		s.writeResponse(w, http.StatusBadRequest, responseEmptyJSON)
		return
	}
	if updatingRequest.Parameters == nil {
		updatingRequest.Parameters = module.GetEmptyUpdatingParameters()
	}

	instance, ok, err := s.store.GetInstance(instanceID)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-updating error: error retrieving instance by id",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	if !ok {
		log.WithFields(logFields).Debug(
			"bad updating request: the instance does not exist",
		)
		// The instance to update does not exist
		// krancour: Choosing to interpret this scenario as a bad request
		// TODO: Write a more detailed response
		s.writeResponse(w, http.StatusBadRequest, responseEmptyJSON)
		return
	}

	// Our broker doesn't actually require the serviceID and previousValues that,
	// per spec, are passed to us in the request body (since this broker is
	// stateful, we can get these details from the instance we already
	// retrieved), BUT if serviceID and previousValues were provided, they BETTER
	// be the same as what's in the instance-- or else we obviously have a
	// conflict.
	if (updatingRequest.ServiceID != instance.ServiceID) ||
		(updatingRequest.PreviousValues.PlanID != "" &&
			updatingRequest.PreviousValues.PlanID != instance.PlanID) {
		logFields["serviceID"] = instance.ServiceID
		logFields["requestServiceID"] = updatingRequest.ServiceID
		logFields["previousPlanID"] = instance.PlanID
		logFields["requestPreviousPlanID"] = updatingRequest.PreviousValues.PlanID
		log.WithFields(logFields).Debug(
			"bad updating request: serviceID or previousPlanID does not match " +
				"serviceID or previousPlanID on the instance",
		)
		// TODO: Write a more detailed response
		s.writeResponse(w, http.StatusConflict, responseEmptyJSON)
		return
	}

	previousUpdatingRequestParams := module.GetEmptyUpdatingParameters()
	if err = instance.GetUpdatingParameters(
		previousUpdatingRequestParams,
		s.codec,
	); err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-updating error: error decoding persisted " +
				"updatingParameters",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	previousUpdatingRequest := &UpdatingRequest{
		ServiceID:  instance.ServiceID,
		PlanID:     instance.PlanID,
		Parameters: previousUpdatingRequestParams,
	}
	if reflect.DeepEqual(updatingRequest, previousUpdatingRequest) {
		// Per the spec, if fully provisioned, respond with a 200, else a 202.
		// Filling in a gap in the spec-- if the status is anything else, we'll
		// choose to respond with a 409
		switch instance.Status {
		case service.InstanceStateUpdating:
			s.writeResponse(w, http.StatusAccepted, responseUpdatingAccepted)
			return
		case service.InstanceStateUpdated:
			s.writeResponse(w, http.StatusOK, responseEmptyJSON)
			return
		default:
			// TODO: Write a more detailed response
			s.writeResponse(w, http.StatusConflict, responseEmptyJSON)
			return
		}
	} else if instance.Status != service.InstanceStateProvisioned {
		log.WithFields(logFields).Debug(
			"bad updating request: the instance to update to is not in a " +
				"provisioned state",
		)
		// The instance to update is not in a provisioned state
		// krancour: Choosing to interpret this scenario as unprocessable
		// TODO: Write a more detailed response
		s.writeResponse(w, http.StatusUnprocessableEntity, responseEmptyJSON)
		return
	}

	// If we get to here, we need to update the instance.
	// Start by carrying out module-specific request validation
	err = module.ValidateUpdatingParameters(updatingRequest.Parameters)
	if err != nil {
		validationErr, ok := err.(*service.ValidationError)
		if ok {
			logFields["field"] = validationErr.Field
			logFields["issue"] = validationErr.Issue
			log.WithFields(logFields).Debug(
				"bad updating request: validation error",
			)
			// TODO: Send the correct response body-- this is a placeholder
			s.writeResponse(w, http.StatusBadRequest, responseEmptyJSON)
			return
		}
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	updater, err := module.GetUpdater(
		updatingRequest.ServiceID,
		updatingRequest.PlanID,
	)
	if err != nil {
		logFields["serviceID"] = updatingRequest.ServiceID
		logFields["planID"] = updatingRequest.PlanID
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-updating error: error retrieving updater for service and plan",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	firstStepName, ok := updater.GetFirstStepName()
	if !ok {
		logFields["serviceID"] = updatingRequest.ServiceID
		logFields["planID"] = updatingRequest.PlanID
		log.WithFields(logFields).Error(
			"pre-updating error: no steps found for updating service and plan",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	if err := instance.SetUpdatingParameters(
		updatingRequest.Parameters,
		s.codec,
	); err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"updating error: error encoding updatingParameters",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	instance.Status = service.InstanceStateUpdating
	instance.PlanID = updatingRequest.PlanID
	if err := s.store.WriteInstance(instance); err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"updating error: error persisting updated instance",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	task := model.NewTask(
		"updateStep",
		map[string]string{
			"stepName":   firstStepName,
			"instanceID": instanceID,
		},
	)
	if err := s.asyncEngine.SubmitTask(task); err != nil {
		logFields["step"] = firstStepName
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"updating error: error submitting updating task",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	// If we get all the way to here, we've been successful!
	s.writeResponse(w, http.StatusAccepted, responseUpdatingAccepted)

	log.WithFields(logFields).Debug("asynchronous updating initiated")
}
