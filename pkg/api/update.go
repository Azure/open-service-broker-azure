package api

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"github.com/Azure/open-service-broker-azure/pkg/async"
	"github.com/Azure/open-service-broker-azure/pkg/service"
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
		s.writeResponse(
			w,
			http.StatusUnprocessableEntity,
			generateAsyncRequiredResponse(),
		)
		return
	}
	acceptsIncomplete, err := strconv.ParseBool(acceptsIncompleteStr)
	if err != nil || !acceptsIncomplete {
		logFields["accepts_incomplete"] = acceptsIncompleteStr
		log.WithFields(logFields).Debug(
			`bad updating request: query parameter has invalid value; only ` +
				`"true" is accepted`,
		)
		s.writeResponse(
			w,
			http.StatusUnprocessableEntity,
			generateAsyncRequiredResponse(),
		)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-updating error: error reading request body",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}
	defer r.Body.Close() // nolint: errcheck

	updatingRequest, err := NewUpdatingRequestFromJSON(bodyBytes)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Debug(
			"bad updating request: error unmarshaling request body",
		)
		// krancour: Choosing to interpret this scenario as a bad request, as a
		// valid request, obviously contains valid, well-formed JSON
		// TODO: Write a more detailed response
		s.writeResponse(w, http.StatusBadRequest, generateEmptyResponse())
		return
	}

	if updatingRequest.ServiceID == "" {
		logFields["field"] = "service_id"
		log.WithFields(logFields).Debug(
			"bad updating request: required request body field is missing",
		)
		s.writeResponse(w, http.StatusBadRequest, generateServiceIDRequiredResponse())
		return
	}

	svc, ok := s.catalog.GetService(updatingRequest.ServiceID)
	if !ok {
		logFields["serviceID"] = updatingRequest.ServiceID
		log.WithFields(logFields).Debug(
			"bad updating request: invalid serviceID",
		)
		s.writeResponse(w, http.StatusBadRequest, generateInvalidServiceIDResponse())
		return
	}

	var plan service.Plan
	if updatingRequest.PlanID != "" {
		plan, ok = svc.GetPlan(updatingRequest.PlanID)
		if !ok {
			logFields["serviceID"] = updatingRequest.ServiceID
			logFields["planID"] = updatingRequest.PlanID
			log.WithFields(logFields).Debug(
				"bad updating request: invalid planID for service",
			)
			s.writeResponse(w, http.StatusBadRequest, generateInvalidPlanIDResponse())
			return
		}
	}

	instance, ok, err := s.store.GetInstance(instanceID)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-updating error: error retrieving instance by id",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}
	if !ok {
		log.WithFields(logFields).Debug(
			"bad updating request: the instance does not exist",
		)
		// The instance to update does not exist
		// krancour: Choosing to interpret this scenario as a bad request
		// TODO: Write a more detailed response
		s.writeResponse(w, http.StatusBadRequest, generateEmptyResponse())
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
		s.writeResponse(w, http.StatusConflict, generateEmptyResponse())
		return
	}

	serviceManager := svc.GetServiceManager()

	updatingParameters, secureUpdatingParameters, err :=
		serviceManager.SplitProvisioningParameters(updatingRequest.Parameters)
	if err != nil {
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}

	// This determines whether the parameters of the update request are already
	// reflected in the provisioning parameters of a fully provisioned (or fully
	// updated) instance OR the parameters of the update request are already
	// reflected in the existing updating paramters of an in-progress update.
	var existingParams map[string]interface{}
	var existingSecureParams map[string]interface{}
	switch instance.Status {
	case service.InstanceStateProvisioned:
		existingParams = instance.ProvisioningParameters
		existingSecureParams = instance.SecureProvisioningParameters
	case service.InstanceStateUpdating:
		existingParams = instance.UpdatingParameters
		existingSecureParams = instance.SecureUpdatingParameters
	default:
		// If instance isn't fully provisioned (or updated) and there isn't an
		// update in-progress, we cannot handle this request. It's a conflict.
		s.writeResponse(w, http.StatusConflict, generateEmptyResponse())
		return
	}
	var needsUpdate bool
	for reqParamKey, reqParamVal := range updatingParameters {
		existingVal, ok := existingParams[reqParamKey]
		if !ok || !reflect.DeepEqual(reqParamVal, existingVal) {
			needsUpdate = true
			break
		}
	}
	for reqParamKey, reqParamVal := range secureUpdatingParameters {
		existingVal, ok := existingSecureParams[reqParamKey]
		if !ok || !reflect.DeepEqual(reqParamVal, existingVal) {
			needsUpdate = true
			break
		}
	}
	if needsUpdate {
		if instance.Status == service.InstanceStateUpdating {
			// We cannot handle two updates at once. This is a conflict.
			s.writeResponse(w, http.StatusConflict, generateEmptyResponse())
		}
	} else {
		if instance.Status == service.InstanceStateProvisioned {
			// In this case, the requested update is already completed
			s.writeResponse(w, http.StatusOK, generateEmptyResponse())
		}
		// In this case, the requested update is already in-progress
		s.writeResponse(w, http.StatusAccepted, generateUpdateAcceptedResponse())
	}

	// Only one scenario gets us to this point-- the instance is fully provisioned
	// (or fully updated) and the parameters of the update request indicate the
	// need for a new update.

	// Carry out schema-driven update request parameters validation.
	if instance.Plan.GetSchemas().ServiceInstances.UpdatingParametersSchema != nil { // nolint: lll
		if err :=
			instance.Plan.GetSchemas().ServiceInstances.UpdatingParametersSchema.Validate( // nolint: lll
				updatingRequest.Parameters,
			); err != nil {
			var validationErr *service.ValidationError
			validationErr, ok = err.(*service.ValidationError)
			if ok {
				logFields["field"] = validationErr.Field
				logFields["issue"] = validationErr.Issue
				log.WithFields(logFields).Debug(
					"bad updating request: validation error",
				)
				s.writeResponse(
					w,
					http.StatusBadRequest,
					generateValidationFailedResponse(validationErr),
				)
				return
			}
			s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
			return
		}
	}

	// This uses module-specific logic to weigh update parameters against current
	// instance state to detect any invalid state changes. An example of this
	// might be reducing the amound of storage allocated to a database.
	instance.UpdatingParameters = updatingParameters
	instance.SecureUpdatingParameters = secureUpdatingParameters
	if err := serviceManager.ValidateUpdatingParameters(instance); err != nil {
		var validationErr *service.ValidationError
		validationErr, ok = err.(*service.ValidationError)
		if ok {
			logFields["field"] = validationErr.Field
			logFields["issue"] = validationErr.Issue
			log.WithFields(logFields).Debug(
				"bad updating request: validation error",
			)
			s.writeResponse(
				w,
				http.StatusBadRequest,
				generateValidationFailedResponse(validationErr),
			)
			return
		}
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}

	// If we get to here, we need to update the instance.

	if plan == nil {
		plan, ok = svc.GetPlan(instance.PlanID)
		if !ok {
			logFields["serviceID"] = updatingRequest.ServiceID
			logFields["planID"] = instance.PlanID
			log.WithFields(logFields).Error(
				"pre-updating error: no Plan found for planID in Service",
			)
			s.writeResponse(
				w,
				http.StatusInternalServerError,
				generateInvalidPlanIDResponse(),
			)
			return
		}
	}
	updater, err := serviceManager.GetUpdater(plan)
	if err != nil {
		logFields["serviceID"] = updatingRequest.ServiceID
		logFields["planID"] = updatingRequest.PlanID
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-updating error: error retrieving updater for service and plan",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}
	firstStepName, ok := updater.GetFirstStepName()
	if !ok {
		logFields["serviceID"] = updatingRequest.ServiceID
		logFields["planID"] = updatingRequest.PlanID
		log.WithFields(logFields).Error(
			"pre-updating error: no steps found for updating service and plan",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}

	instance.Status = service.InstanceStateUpdating
	instance.PlanID = updatingRequest.PlanID
	if err := s.store.WriteInstance(instance); err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"updating error: error persisting updated instance",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}

	task := async.NewTask(
		"executeUpdatingStep",
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
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}

	// If we get all the way to here, we've been successful!
	s.writeResponse(w, http.StatusAccepted, generateUpdateAcceptedResponse())

	log.WithFields(logFields).Debug("asynchronous updating initiated")
}
