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

	// Start by carrying out plan-specific updating request parameters validation
	// This is the first of two validation stages. This one is schema-driven.
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

	serviceManager := svc.GetServiceManager()

	updatingParameters, secureUpdatingParameters, err :=
		serviceManager.SplitProvisioningParameters(updatingRequest.Parameters)
	if err != nil {
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}

	// This is the second of two validation stages. This weighs updating
	// parameters against current instance state. Service/plan-specific logic can
	// determine if a requested update isn't valid for some reason. Examples of
	// invalid updates might include, for instance, REDUCING the storage allocated
	// to a database.
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

	if instance.ServiceID == updatingRequest.ServiceID &&
		instance.PlanID == updatingRequest.PlanID &&
		reflect.DeepEqual(
			instance.UpdatingParameters,
			updatingParameters,
		) &&
		reflect.DeepEqual(
			instance.SecureUpdatingParameters,
			secureUpdatingParameters,
		) {
		// Per the spec, if fully provisioned, respond with a 200, else a 202.
		// Filling in a gap in the spec-- if the status is anything else, we'll
		// choose to respond with a 409
		switch instance.Status {
		case service.InstanceStateUpdating:
			s.writeResponse(w, http.StatusAccepted, generateUpdateAcceptedResponse())
			return
		case service.InstanceStateUpdated:
			s.writeResponse(w, http.StatusOK, generateEmptyResponse())
			return
		default:
			// TODO: Write a more detailed response
			s.writeResponse(w, http.StatusConflict, generateEmptyResponse())
			return
		}
	} else {
		switch instance.Status {
		case service.InstanceStateProvisioned:
		case service.InstanceStateUpdatingFailed:
		default:
			log.WithFields(logFields).Debug(
				"bad updating request: the instance to update to is not in a " +
					"provisioned state",
			)
			// The instance to update is not in a provisioned state
			// This could be from a previously failed update, which means
			// we will never allow a subsequent update to go through.
			// krancour: Choosing to interpret this scenario as unprocessable
			// TODO: Write a more detailed response
			s.writeResponse(w, http.StatusUnprocessableEntity, generateEmptyResponse())
			return
		}
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
