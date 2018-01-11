package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/async/model"
	"github.com/Azure/open-service-broker-azure/pkg/service"
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

	svc, ok := s.catalog.GetService(instance.ServiceID)
	if !ok {
		// If we don't find the Service in the catalog, something is really wrong.
		// (It should exist, because an instance with this serviceID exists.)
		logFields["serviceID"] = instance.ServiceID
		log.WithFields(logFields).Error(
			"pre-deprovisioning error: no Service found for serviceID",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	plan, ok := svc.GetPlan(instance.PlanID)
	if !ok {
		// If we don't find the Service in the catalog, something is really wrong.
		// (It should exist, because an instance with this serviceID exists.)
		logFields["serviceID"] = instance.ServiceID
		logFields["planID"] = instance.PlanID
		log.WithFields(logFields).Error(
			"pre-deprovisioning error: no Plan found for planID in Service",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	serviceManager := svc.GetServiceManager()

	deprovisioner, err := serviceManager.GetDeprovisioner(plan)
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

	childCount, err := s.store.GetInstanceChildCountByAlias(instance.Alias)
	if err != nil {
		logFields["step"] = firstStepName
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"deprovisioning error: error determining child count",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
	}

	var task model.Task
	if childCount > 0 {
		task = model.NewDelayedTask(
			"waitForChildrenStep",
			map[string]string{
				"deprovisionFirstStep": firstStepName,
				"instanceID":           instanceID,
			},
			time.Minute*5,
		)
		log.WithFields(logFields).Debug("children not deprovisioned, waiting")
	} else {
		task = model.NewTask(
			"deprovisionStep",
			map[string]string{
				"stepName":   firstStepName,
				"instanceID": instanceID,
			},
		)
		log.WithFields(logFields).Debug(
			"no provisioned children, starting deprovision",
		)
	}
	if err = s.asyncEngine.SubmitTask(task); err != nil {
		logFields["step"] = firstStepName
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"deprovisioning error: error determining child count",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
	}

	if childCount > 0 {
		task := model.NewTask(
			"waitForChildrenStep",
			map[string]string{
				"deprovisionFirstStep": firstStepName,
				"instanceID":           instanceID,
			},
		)
		log.WithFields(logFields).Debug("children not deprovisioned, waiting")
		if err = s.asyncEngine.SubmitDelayedTask(instance.Alias, task); err != nil {
			logFields["step"] = "waitForChildrenStep"
			logFields["error"] = err
			log.WithFields(logFields).Error(
				"provisioning error: error submitting delayed provisioning task",
			)
			s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
			return
		}
	} else {
		task := model.NewTask(
			"deprovisionStep",
			map[string]string{
				"stepName":   firstStepName,
				"instanceID": instanceID,
			},
		)
		log.WithFields(logFields).Debug(
			"no provisioned children, starting deprovision",
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
	}
	// If we get all the way to here, we've been successful!
	s.writeResponse(w, http.StatusAccepted, responseDeprovisioningAccepted)

	log.WithFields(logFields).Debug("asynchronous deprovisioning initiated")
}
