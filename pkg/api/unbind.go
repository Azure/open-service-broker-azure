package api

import (
	"fmt"
	"net/http"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func (s *server) unbind(w http.ResponseWriter, r *http.Request) {
	instanceID := mux.Vars(r)["instance_id"]
	bindingID := mux.Vars(r)["binding_id"]

	logFields := log.Fields{
		"instanceID": instanceID,
		"bindingID":  bindingID,
	}

	log.WithFields(logFields).Debug("received unbinding request")

	binding, ok, err := s.store.GetBinding(bindingID)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-unbinding error: error retrieving binding by id",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	if !ok {
		log.WithFields(logFields).Debug(
			"no such binding remains to be unbound",
		)
		// No binding was found-- per spec, we return a 410
		s.writeResponse(w, http.StatusGone, responseEmptyJSON)
		return
	}

	if binding.InstanceID != instanceID {
		logFields["instanceID"] = binding.InstanceID
		logFields["requestInstanceID"] = instanceID
		log.WithFields(logFields).Debug(
			"bad unbinding request: instanceID does not match instanceID on the " +
				"binding",
		)
		// TODO: Write a more detailed response
		s.writeResponse(w, http.StatusConflict, responseEmptyJSON)
		return
	}

	instance, ok, err := s.store.GetInstance(instanceID)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-unbinding error: error retrieving instance by id",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	if !ok {
		// The instance to unbind from does not exist!
		// krancour: Not totally sure what to do here. It seems within the realm
		// of possibility that an instance could be deprovisioned without all
		// bindings to that instance first having been unbound-- at least there
		// isn't any logic in this broker to prevent that and I don't believe the
		// spec is clear on whether that's permissible or not. So for now, we must
		// accept the possibility that orphaned bindings may exist. I'm choosing
		// to deal with this by skipping straight to deleting the binding from the
		// datastore without invoking any service-specific unbinding logic. (We
		// cannot, because with the instance no longer existing, we cannot identify
		// the service and plan of the instance, and therefore do not know which
		// serviceManager can successfully effect binding).
		// TODO: Re-evaluate this decision later.
		log.WithFields(logFields).Debug(
			"unbinding an orphaned binding",
		)
	} else {
		serviceManager := instance.Service.GetServiceManager()

		// Starting here, if something goes wrong, we don't know what state service-
		// specific code has left us in, so we'll attempt to record the error in
		// the datastore.
		err = serviceManager.Unbind(instance, binding.Details)
		if err != nil {
			s.handleUnbindingError(
				binding,
				err,
				"error executing service-specific unbinding logic",
				w,
			)
			return
		}

	}

	if _, err = s.store.DeleteBinding(bindingID); err != nil {
		s.handleUnbindingError(
			binding,
			err,
			"error deleting binding",
			w,
		)
		return
	}

	s.writeResponse(w, http.StatusOK, responseEmptyJSON)
}

// handleUnbindingError tries to handle the most serious unbinding errors. The
// binding status is updated and an attempt is made to persist the binding with
// updated status. If this fails, we have a very serious problem on our hands,
// so we log that failure and kill the process. Barring such a failure, a nicely
// formatted error message is logged.
func (s *server) handleUnbindingError(
	binding service.Binding,
	e error,
	msg string,
	w http.ResponseWriter,
) {
	binding.Status = service.BindingStateUnbindingFailed
	if e == nil {
		binding.StatusReason = fmt.Sprintf(`unbinding error: %s`, msg)
	} else {
		binding.StatusReason = fmt.Sprintf(`unbinding error: %s: %s`, msg, e)
	}
	logFields := log.Fields{
		"bindingID":  binding.BindingID,
		"instanceID": binding.InstanceID,
		"status":     binding.Status,
	}
	err := s.store.WriteBinding(binding)
	if err != nil {
		logFields["originalError"] = binding.StatusReason
		logFields["persistenceError"] = err
		log.WithFields(logFields).Fatal(
			"unbinding error: error persisting binding with updated status",
		)
	}
	if e != nil {
		logFields["error"] = e
	}
	log.WithFields(logFields).Error(
		fmt.Sprintf(`unbinding error: %s`, msg),
	)
	s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
}
