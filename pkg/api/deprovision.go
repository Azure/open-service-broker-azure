package api

import (
	"net/http"
	"strconv"

	"github.com/Azure/azure-service-broker/pkg/async/model"
	"github.com/Azure/azure-service-broker/pkg/service"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func (s *server) deprovision(w http.ResponseWriter, r *http.Request) {
	instanceID := mux.Vars(r)["instance_id"]

	logFields := log.Fields{
		"instanceID": instanceID,
	}

	log.WithFields(logFields).Debug("received deprovisioning request")

	// This broker provisions everything asynchronously. If a client doesn't
	// explicitly indicate that they will accept an incomplete result, the
	// spec says to respond with a 422
	acceptsIncompleteStr := r.URL.Query().Get("accepts_incomplete")
	if acceptsIncompleteStr == "" {
		logFields["parameter"] = "accepts_incomplete=true" // nolint: goconst
		log.WithFields(logFields).Debug(
			"bad deprovisioning request: request is missing required query parameter",
		)
		s.writeResponse(w, http.StatusUnprocessableEntity, responseAsyncRequired)
		return
	}
	acceptsIncomplete, err := strconv.ParseBool(acceptsIncompleteStr)
	if err != nil || !acceptsIncomplete {
		logFields["accepts_incomplete"] = acceptsIncompleteStr
		log.WithFields(logFields).Debug(
			`bad deprovisioning request: query parameter has invalid value; only ` +
				`"true" is accepted`,
		)
		s.writeResponse(w, http.StatusUnprocessableEntity, responseAsyncRequired)
		return
	}

	instance, ok, err := s.store.GetInstance(instanceID)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-deprovisioning error: error retrieving instance by id",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	if !ok {
		log.WithFields(logFields).Debug(
			"no such instance remains to be deprovisioned",
		)
		// No instance was found-- per spec, we return a 410
		s.writeResponse(w, http.StatusGone, responseEmptyJSON)
		return
	}
	switch instance.Status {
	case service.InstanceStateDeprovisioning:
		log.WithFields(logFields).Debug(
			"deprovisioning is already in progress",
		)
		s.writeResponse(w, http.StatusAccepted, responseDeprovisioningAccepted)
		return
	case service.InstanceStateProvisioned:
	case service.InstanceStateProvisioningFailed:
	default:
		// This is going to handle the case where we cannot deprovision because
		// the instance isn't in a terminal state-- i.e. it's still provisioning
		logFields["status"] = instance.Status
		log.WithFields(logFields).Debug(
			"cannot deprovision instance in its current state",
		)
		s.writeResponse(w, http.StatusConflict, responseEmptyJSON)
		return
	}

	// If we get to here, we're dealing with an instance that is fully provisioned
	// or has failed provisioning. We need to kick off asynchronous
	// deprovisioning.

	module, ok := s.modules[instance.ServiceID]
	if !ok {
		logFields["serviceID"] = instance.ServiceID
		log.WithFields(logFields).Error(
			"pre-deprovisioning error: no module found for service",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	deprovisioner, err := module.GetDeprovisioner(
		instance.ServiceID,
		instance.PlanID,
	)
	if err != nil {
		logFields["serviceID"] = instance.ServiceID
		logFields["planID"] = instance.PlanID
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-deprovisioning error: error retrieving deprovisioner for service " +
				"and plan",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	firstStepName, ok := deprovisioner.GetFirstStepName()
	if !ok {
		logFields["serviceID"] = instance.ServiceID
		logFields["planID"] = instance.PlanID
		log.WithFields(logFields).Error(
			"pre-deprovisioning error: no steps found for deprovisioning service " +
				"and plan",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	instance.Status = service.InstanceStateDeprovisioning
	if err = s.store.WriteInstance(instance); err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"deprovisioning error: error persisting updated instance",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	task := model.NewTask(
		"deprovisionStep",
		map[string]string{
			"stepName":   firstStepName,
			"instanceID": instanceID,
		},
	)
	if err = s.asyncEngine.SubmitTask(task); err != nil {
		logFields["step"] = firstStepName
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"deprovisioning error: error submitting deprovisioning task",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}

	// If we get all the way to here, we've been successful!
	s.writeResponse(w, http.StatusAccepted, responseDeprovisioningAccepted)

	log.WithFields(logFields).Debug("asynchronous deprovisioning initiated")
}
