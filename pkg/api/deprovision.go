package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Azure/azure-service-broker/pkg/async/model"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/gorilla/mux"
)

func (s *server) deprovision(w http.ResponseWriter, r *http.Request) {
	// This broker provisions everything asynchronously. If a client doesn't
	// explicitly indicate that they will accept an incomplete result, the
	// spec says to respond with a 422
	acceptsIncompleteStr := r.URL.Query().Get("accepts_incomplete")
	if acceptsIncompleteStr == "" {
		log.Println(
			"request is missing required query parameter accepts_incomplete=true",
		)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(responseAsyncRequired)
		return
	}
	acceptsIncomplete, err := strconv.ParseBool(acceptsIncompleteStr)
	if err != nil || !acceptsIncomplete {
		log.Printf(
			"query paramater accepts_incomplete has invalid value '%s'",
			acceptsIncompleteStr,
		)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(responseAsyncRequired)
		return
	}

	instanceID := mux.Vars(r)["instance_id"]
	instance, ok, err := s.store.GetInstance(instanceID)
	if err != nil {
		log.Printf("error retrieving instance with id %s", instanceID)
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
		log.Printf("error finding module for service %s", instance.ServiceID)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}

	deprovisioner, err := module.GetDeprovisioner()
	if err != nil {
		log.Printf(
			`error retrieving deprovisioner for service "%s"`,
			instance.ServiceID,
		)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}
	firstStepName, ok := deprovisioner.GetFirstStepName()
	if !ok {
		log.Printf(
			`no steps found for deprovisioning service "%s"`,
			instance.ServiceID,
		)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}

	instance.Status = service.InstanceStateDeprovisioning
	err = s.store.WriteInstance(instance)
	if err != nil {
		log.Println("error updating instance")
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
		log.Println("error submitting deprovisioning task")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}

	// If we get all the way to here, we've been successful!
	w.WriteHeader(http.StatusAccepted)
	w.Write(responseEmptyJSON)
}
