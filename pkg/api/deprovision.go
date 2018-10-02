package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	"github.com/deis/async"
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
			`bad deprovisioning request: query parameter has invalid value; only ` +
				`"true" is accepted`,
		)
		s.writeResponse(
			w,
			http.StatusUnprocessableEntity,
			generateAsyncRequiredResponse(),
		)
		return
	}

	instance, ok, err := s.store.GetInstance(instanceID)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-deprovisioning error: error retrieving instance by id",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}
	if !ok {
		log.WithFields(logFields).Debug(
			"no such instance remains to be deprovisioned",
		)
		// No instance was found-- per spec, we return a 410
		s.writeResponse(w, http.StatusGone, generateEmptyResponse())
		return
	}
	if instance.Details == nil {
		// If we get to here, we're dealing with an orphan -- the instance
		// detail is nil for some reason. We simply delete the record from
		// the store and return 200 OK.
		instanceFound, err := s.store.DeleteInstance(instanceID)
		if err != nil {
			logFields["error"] = err
			log.WithFields(logFields).Error(
				"pre-deprovisioning error: error deleting instance will nil detail",
			)
			s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
			return
		}
		if !instanceFound {
			// We should never land in here, since "instanceFound, err = false, nil"
			// only happens when the instanceID can't be found in the store, but in L63
			// we have done the validation and we can make sure the instanceID is found
			// in the store.
			s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
			log.Fatal(
				fmt.Errorf("pre-deprovision fatal error, store inconsistency %s", err),
			)
		}
		s.writeResponse(w, http.StatusOK, generateEmptyResponse())
		return
	}
	switch instance.Status {
	case service.InstanceStateDeprovisioning:
		log.WithFields(logFields).Debug(
			"deprovisioning is already in progress",
		)
		s.writeResponse(w, http.StatusAccepted, generateDeprovisionAcceptedResponse())
		return
	case service.InstanceStateProvisioned:
	case service.InstanceStateProvisioningFailed:
	case service.InstanceStateUpdatingFailed:
	default:
		// This is going to handle the case where we cannot deprovision because
		// the instance isn't in a terminal state-- i.e. it's still provisioning
		logFields["status"] = instance.Status
		log.WithFields(logFields).Debug(
			"cannot deprovision instance in its current state",
		)
		s.writeResponse(w, http.StatusConflict, generateEmptyResponse())
		return
	}

	// If we get to here, we're dealing with an instance that is fully provisioned
	// or has failed provisioning. We need to kick off asynchronous
	// deprovisioning.

	serviceManager := instance.Service.GetServiceManager()

	deprovisioner, err := serviceManager.GetDeprovisioner(instance.Plan)
	if err != nil {
		logFields["serviceID"] = instance.ServiceID
		logFields["planID"] = instance.PlanID
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-deprovisioning error: error retrieving deprovisioner for service " +
				"and plan",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
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
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}

	var task async.Task
	if childCount, err :=
		s.store.GetInstanceChildCountByAlias(instance.Alias); err != nil {
		logFields["step"] = firstStepName
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"deprovisioning error: error determining child count",
		)
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	} else if childCount > 0 {
		instance.Status = service.InstanceStateDeprovisioningDeferred
		logFields["provisionedChildren"] = childCount
		task = async.NewDelayedTask(
			"checkChildrenStatuses",
			map[string]string{
				"instanceID": instanceID,
			},
			time.Minute*1,
		)
		log.WithFields(logFields).Debug("children not deprovisioned, waiting")
	} else {
		instance.Status = service.InstanceStateDeprovisioning
		task = async.NewTask(
			"executeDeprovisioningStep",
			map[string]string{
				"stepName":   firstStepName,
				"instanceID": instanceID,
			},
		)
		log.WithFields(logFields).Debug(
			"no provisioned children, starting deprovision",
		)
	}

	if err = s.store.WriteInstance(instance); err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"deprovisioning error: error persisting updated instance",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}

	if err = s.asyncEngine.SubmitTask(task); err != nil {
		logFields["step"] = firstStepName
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"deprovisioning error: error submitting deprovisioning task",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}

	// If we get all the way to here, we've been successful!
	s.writeResponse(w, http.StatusAccepted, generateDeprovisionAcceptedResponse())

	log.WithFields(logFields).Debug("asynchronous deprovisioning initiated")
}
