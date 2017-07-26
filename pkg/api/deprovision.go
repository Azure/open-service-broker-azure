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
	// This broker provisions everything asynchronously. If a client doesn't
	// explicitly indicate that they will accept an incomplete result, the
	// spec says to respond with a 422
	acceptsIncompleteStr := r.URL.Query().Get("accepts_incomplete")
	if acceptsIncompleteStr == "" {
		log.WithField(
			"parameter",
			"accepts_incomplete=true",
		).Debug("request is missing required query parameter")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(responseAsyncRequired)
		return
	}
	acceptsIncomplete, err := strconv.ParseBool(acceptsIncompleteStr)
	if err != nil || !acceptsIncomplete {
		log.WithField(
			"accepts_incomplete",
			acceptsIncompleteStr,
		).Debug(`query paramater has invalid value; only "true" is accepted`)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(responseAsyncRequired)
		return
	}

	instanceID := mux.Vars(r)["instance_id"]
	instance, ok, err := s.store.GetInstance(instanceID)
	if err != nil {
		log.WithFields(log.Fields{
			"instanceID": instanceID,
			"error":      err,
		}).Error("error retrieving instance by id")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}
	if !ok {
		// No instance was found-- per spec, we return a 410
		w.WriteHeader(http.StatusGone)
		w.Write(responseEmptyJSON)
		return
	}
	switch instance.Status {
	case service.InstanceStateDeprovisioning:
		w.WriteHeader(http.StatusAccepted)
		w.Write(responseEmptyJSON)
		return
	case service.InstanceStateProvisioned:
	case service.InstanceStateProvisioningFailed:
	default:
		w.WriteHeader(http.StatusConflict)
		w.Write(responseEmptyJSON)
		return
	}

	// If we get to here, we're dealing with an instance that is fully provisioned
	// or has failed provisioning. We need to kick off asynchronous deprovisioning

	module, ok := s.modules[instance.ServiceID]
	if !ok {
		log.WithField(
			"serviceID",
			instance.ServiceID,
		).Error("could not find module for service")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}

	deprovisioner, err := module.GetDeprovisioner()
	if err != nil {
		log.WithFields(log.Fields{
			"serviceID": instance.ServiceID,
			"error":     err,
		}).Error("error retrieving deprovisioner for service")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}
	firstStepName, ok := deprovisioner.GetFirstStepName()
	if !ok {
		log.WithField(
			"serviceID",
			instance.ServiceID,
		).Error("no steps found for deprovisioning service")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}

	instance.Status = service.InstanceStateDeprovisioning
	err = s.store.WriteInstance(instance)
	if err != nil {
		log.WithFields(log.Fields{
			"instanceID": instanceID,
			"error":      err,
		}).Error("error updating instance")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}

	task := model.NewTask(
		"deprovisionStep",
		map[string]string{
			"stepName":   firstStepName,
			"instanceID": instanceID,
		},
	)
	err = s.asyncEngine.SubmitTask(task)
	if err != nil {
		log.WithFields(log.Fields{
			"step":       firstStepName,
			"instanceID": instanceID,
			"error":      err,
		}).Error("error submitting deprovisioning task")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}

	// If we get all the way to here, we've been successful!
	w.WriteHeader(http.StatusAccepted)
	w.Write(responseEmptyJSON)
}
