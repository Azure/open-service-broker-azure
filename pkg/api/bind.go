package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/Azure/azure-service-broker/pkg/service"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func (s *server) bind(w http.ResponseWriter, r *http.Request) {
	instanceID := mux.Vars(r)["instance_id"]
	bindingID := mux.Vars(r)["binding_id"]

	logFields := log.Fields{
		"instanceID": instanceID,
		"bindingID":  bindingID,
	}

	log.WithFields(logFields).Debug("received binding request")

	instance, ok, err := s.store.GetInstance(instanceID)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-binding error: error retrieving instance by id",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	if !ok {
		log.WithFields(logFields).Debug(
			"bad binding request: the instance does not exist",
		)
		// The instance to bind to does not exist
		// krancour: Choosing to interpret this scenario as a bad request
		// TODO: Write a more detailed response
		s.writeResponse(w, http.StatusBadRequest, responseEmptyJSON)
		return
	}

	if instance.Status != service.InstanceStateProvisioned {
		log.WithFields(logFields).Debug(
			"bad binding request: the instance to bind to is not in a provisioned state",
		)
		// The instance to bind to does not exist
		// krancour: Choosing to interpret this scenario as unprocessable
		// TODO: Write a more detailed response
		s.writeResponse(w, http.StatusUnprocessableEntity, responseEmptyJSON)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-binding error: error reading request body",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	defer r.Body.Close() // nolint: errcheck

	bindingRequest := &service.BindingRequest{}
	err = service.GetBindingRequestFromJSON(bodyBytes, bindingRequest)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Debug(
			"bad binding request: error unmarshaling request body",
		)
		// This scenario is a bad request, as a valid request obviously must contain
		// valid, well-formed JSON
		s.writeResponse(w, http.StatusBadRequest, responseMalformedRequestBody)
		return
	}

	// Our broker doesn't actually require the serviceID and planID that, per
	// spec, are passed to us in the request body (since this broker is stateful,
	// we can get these details from the instance we already retrieved), BUT if
	// serviceID and planID were provided, they BETTER be the same as what's in
	// the instance-- or else we obviously have a conflict.
	if (bindingRequest.ServiceID != "" &&
		bindingRequest.ServiceID != instance.ServiceID) ||
		(bindingRequest.PlanID != "" &&
			bindingRequest.PlanID != instance.PlanID) {
		logFields["serviceID"] = instance.ServiceID
		logFields["requestServiceID"] = bindingRequest.ServiceID
		logFields["planID"] = instance.PlanID
		logFields["requestPlanID"] = bindingRequest.PlanID
		log.WithFields(logFields).Debug(
			"bad binding request: serviceID or planID does not match serviceID or " +
				"planID on the instance",
		)
		// TODO: Write a more detailed response
		s.writeResponse(w, http.StatusConflict, responseEmptyJSON)
		return
	}

	// At this point, there's absolute agreement on what service we're dealing
	// with. We can go ahead and find the module for this service.
	module, ok := s.modules[instance.ServiceID]
	if !ok {
		// If we don't find a module that handles the service, something is really
		// wrong. (It should exist, because an instance with this serviceID exists.)
		logFields["serviceID"] = instance.ServiceID
		log.WithFields(logFields).Error(
			"pre-binding error: no module found for service",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	// Now that we know what module we're dealing with, we can get an instance
	// of the module-specific type for bindingParameters and take a second
	// pass at parsing the request body
	bindingRequest.Parameters = module.GetEmptyBindingParameters()
	err = service.GetBindingRequestFromJSON(bodyBytes, bindingRequest)
	if err != nil {
		log.WithFields(logFields).Debug(
			"bad binding request: error unmarshaling request body",
		)
		// This scenario is a bad request, as a valid request obviously must contain
		// valid, well-formed JSON
		s.writeResponse(w, http.StatusBadRequest, responseMalformedRequestBody)
		return
	}
	if bindingRequest.Parameters == nil {
		bindingRequest.Parameters = module.GetEmptyBindingParameters()
	}

	binding, ok, err := s.store.GetBinding(bindingID)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-binding error: error retrieving binding by id",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	if ok {
		// We land in here if an existing binding was found-- the OSB spec
		// obligates us to compare this binding to the one that was requested and
		// respond with 200 if they're identical or 409 otherwise. It actually seems
		// best to compare instanceIDs to ensure there's no conflict and then
		// compare binding request parameters (not bindings) because binding objects
		// also contain binding context and other status information.
		if instanceID != binding.InstanceID {
			logFields["existingInstanceID"] = binding.InstanceID
			log.WithFields(logFields).Debug(
				"bad binding request: instanceID to bind to does not match " +
					"instanceID of existing binding",
			)
			// TODO: Write a more detailed response
			s.writeResponse(w, http.StatusConflict, responseEmptyJSON)
			return
		}
		previousBindingRequestParams := module.GetEmptyBindingParameters()
		if err = binding.GetBindingParameters(
			previousBindingRequestParams,
			s.codec,
		); err != nil {
			logFields["error"] = err
			log.WithFields(logFields).Error(
				"pre-binding error: error decoding persisted bindingParameters",
			)
			s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
			return
		}

		if reflect.DeepEqual(
			bindingRequest.Parameters,
			previousBindingRequestParams,
		) {
			// Per the spec, if bound, respond with a 200
			// Filling in a gap in the spec-- if the status is anything else, we'll
			// choose to respond with a 409
			switch binding.Status {
			case service.BindingStateBound:
				credentials := module.GetEmptyCredentials()
				if err = binding.GetCredentials(credentials, s.codec); err != nil {
					logFields["error"] = err
					log.WithFields(logFields).Error(
						"binding error: error decoding persisted credentials",
					)
					s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
					return
				}
				bindingResponse := &service.BindingResponse{
					Credentials: credentials,
				}
				var bindingResponseJSON []byte
				bindingResponseJSON, err = bindingResponse.ToJSON()
				if err != nil {
					logFields["error"] = err
					log.WithFields(logFields).Error(
						"binding error: error marshaling binding response",
					)
					s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
					return
				}
				// TODO: krancour: Is this a vulnerability? If I am interpreting the
				// spec correctly, this is the "right" thing to do, but it also means
				// any client can steal credentials just by emulating a binding requet
				// for an existing binding.
				s.writeResponse(w, http.StatusOK, bindingResponseJSON)
				return
			default:
				// TODO: Write a more detailed response
				s.writeResponse(w, http.StatusConflict, responseEmptyJSON)
				return
			}
		}

		// We land in here if an existing binding was found, but its atrributes
		// vary from what was requested. The spec requires us to respond with a
		// 409
		s.writeResponse(w, http.StatusConflict, responseEmptyJSON)
		return
	}

	// If we get to here, we need to provision a new binding.
	// Start by carrying out module-specific request validation
	err = module.ValidateBindingParameters(bindingRequest.Parameters)
	if err != nil {
		validationErr, ok := err.(*service.ValidationError)
		if ok {
			logFields["field"] = validationErr.Field
			logFields["issue"] = validationErr.Issue
			log.WithFields(logFields).Debug(
				"bad binding request: validation error",
			)
			// TODO: Send the correct response body-- this is a placeholder
			s.writeResponse(w, http.StatusBadRequest, responseEmptyJSON)
			return
		}
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	provisioningContext := module.GetEmptyProvisioningContext()
	err = instance.GetProvisioningContext(provisioningContext, s.codec)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"binding error: error decoding persisted provisioningContext",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	binding = &service.Binding{
		InstanceID: instanceID,
		BindingID:  bindingID,
	}

	// Starting here, if something goes wrong, we don't know what state module-
	// specific code has left us in, so we'll attempt to record the error in
	// the datastore.
	bindingContext, credentials, err := module.Bind(
		provisioningContext,
		bindingRequest.Parameters,
	)
	if err != nil {
		s.handleBindingError(
			binding,
			err,
			"error executing module-specific binding logic",
			w,
		)
		return
	}

	if err = binding.SetBindingContext(bindingContext, s.codec); err != nil {
		s.handleBindingError(
			binding,
			err,
			"error encoding bindingContext",
			w,
		)
		return
	}

	if err = binding.SetCredentials(credentials, s.codec); err != nil {
		s.handleBindingError(
			binding,
			err,
			"error encoding credentials",
			w,
		)
		return
	}

	binding.Status = service.BindingStateBound
	if err = s.store.WriteBinding(binding); err != nil {
		s.handleBindingError(
			binding,
			err,
			"error persisting binding",
			w,
		)
		return
	}

	// The binding is completed at this point. The only remaining errors that can
	// occur are errors in preparing or sending the response. Such errors do not
	// need to affect the binding's state.

	bindingResponse := &service.BindingResponse{
		Credentials: credentials,
	}
	bindingJSON, err := bindingResponse.ToJSON()
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"post-binding error: error marshaling bindingResponse",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	// If we get all the way to here, we've been successful!
	s.writeResponse(w, http.StatusCreated, bindingJSON)

	log.WithFields(logFields).Debug("binding complete")
}

// handleBindingError tries to handle the most serious binding errors. The
// binding status is updated and an attempt is made to persist the binding with
// updated status. If this fails, we have a very serious problem on our hands,
// so we log that failure and kill the process. Barring such a failure, a nicely
// formatted error message is logged.
func (s *server) handleBindingError(
	binding *service.Binding,
	e error,
	msg string,
	w http.ResponseWriter,
) {
	binding.Status = service.BindingStateBindingFailed
	if e == nil {
		binding.StatusReason = fmt.Sprintf(`binding error: %s`, msg)
	} else {
		binding.StatusReason = fmt.Sprintf(`binding error: %s: %s`, msg, e)
	}
	logFields := log.Fields{
		"bindingID":  binding.BindingID,
		"instanceID": binding.InstanceID,
		"status":     binding.Status,
	}
	if err := s.store.WriteBinding(binding); err != nil {
		logFields["originalError"] = binding.StatusReason
		logFields["persistenceError"] = err
		log.WithFields(logFields).Fatal(
			"binding error: error persisting binding with updated status",
		)
	}
	if e != nil {
		logFields["error"] = e
	}
	log.WithFields(logFields).Error(
		fmt.Sprintf(`binding error: %s`, msg),
	)
	s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
}
