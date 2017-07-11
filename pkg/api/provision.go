package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"reflect"

	"github.com/Azure/azure-service-broker/pkg/async/model"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/gorilla/mux"
)

func (s *server) provision(w http.ResponseWriter, r *http.Request) {
	// This broker provisions everything asynchronously. If a client doesn't
	// explicitly indicate that they will accept an incomplete result, the
	// spec says to respond with a 422
	acceptsIncompleteStr := r.URL.Query().Get("accepts_incomplete")
	if acceptsIncompleteStr == "" {
		log.Println("request is missing required query parameter accepts_incomplete=true")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(responseAsyncRequired)
		return
	}
	acceptsIncomplete, err := strconv.ParseBool(acceptsIncompleteStr)
	if err != nil || !acceptsIncomplete {
		log.Printf("query paramater accepts_incomplete has invalid value '%s'", acceptsIncompleteStr)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(responseAsyncRequired)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	provisioningRequest := &service.ProvisioningRequest{}
	err = service.GetProvisioningRequestFromJSONString(
		string(bodyBytes),
		provisioningRequest,
	)
	if err != nil {
		log.Println("error parsing request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if provisioningRequest.ServiceID == "" {
		log.Println("request body parameter service_id is a required field")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseServiceIDRequired)
		return
	}
	if provisioningRequest.PlanID == "" {
		log.Println("request body parameter plan_id is a required field")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responsePlanIDRequired)
		return
	}

	svc, ok := s.catalog.GetService(provisioningRequest.ServiceID)
	if !ok {
		log.Printf("invalid service %s", provisioningRequest.ServiceID)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseInvalidServiceID)
		return
	}

	_, ok = svc.GetPlan(provisioningRequest.PlanID)
	if !ok {
		log.Printf(
			"invalid plan %s for service %s",
			provisioningRequest.PlanID,
			provisioningRequest.ServiceID,
		)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseInvalidPlanID)
		return
	}

	module, ok := s.modules[provisioningRequest.ServiceID]
	if !ok {
		// We already validated that the serviceID and planID are legitimate. If
		// we don't find a module that handles the service, something is really
		// wrong.
		log.Printf(
			"error finding module for service %s",
			provisioningRequest.ServiceID,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Now that we know what module we're dealing with, we can get an instance
	// of the module-specific type for provisioning parameters and take a second
	// pass at parsing the request body
	provisioningRequest.Parameters = module.GetEmptyProvisioningParameters()
	err = service.GetProvisioningRequestFromJSONString(
		string(bodyBytes),
		provisioningRequest,
	)
	if err != nil {
		log.Println("error parsing request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	instanceID := mux.Vars(r)["instance_id"]
	instance, ok, err := s.store.GetInstance(instanceID)
	if err != nil {
		log.Printf("error retrieving instance with id %s", instanceID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if ok {
		// We land in here if an existing instance was found-- the OSB spec
		// obligates us to compare this instance to the one that was requested and
		// respond with 200 if they're identical or 409 otherwise. It actually seems
		// best to compare REQUESTS instead because instance objects also contain
		// provisioning results and other status information. So, let's reverse
		// engineer a request from the existing instance then compare it to the
		// current request.
		previousProvisioningRequestParams := module.GetEmptyProvisioningParameters()
		err = instance.GetProvisioningParameters(previousProvisioningRequestParams)
		if err != nil {
			log.Println("error decoding persisted provisioning parameters")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		previousProvisioningRequest := service.ProvisioningRequest{
			ServiceID:  instance.ServiceID,
			PlanID:     instance.PlanID,
			Parameters: previousProvisioningRequestParams,
		}

		if reflect.DeepEqual(provisioningRequest, previousProvisioningRequest) {
			// Per the spec, if fully provisioned, respond with a 200, else a 202.
			// Filling in a gap in teh spec-- if the status is anything else, we'll
			// choose to respond with a 409
			switch instance.Status {
			case service.InstanceStateProvisioning:
				w.WriteHeader(http.StatusAccepted)
			case service.InstanceStateProvisioned:
				w.WriteHeader(http.StatusOK)
			default:
				w.WriteHeader(http.StatusConflict)
			}
			w.Write(responseEmptyJSON)
			return
		}

		// We land in here if an existing instance was found, but its atrributes
		// vary from what was requested. The spec requires us to respond with a
		// 409
		w.WriteHeader(http.StatusConflict)
		w.Write(responseEmptyJSON)
		return
	}

	// If we get to here, we need to provision a new instance.
	// Start by carrying out module-specific request validation
	err = module.ValidateProvisioningParameters(provisioningRequest.Parameters)
	if err != nil {
		_, ok := err.(*service.ValidationError)
		if ok {
			w.WriteHeader(http.StatusBadRequest)
			// TODO: Send the correct response body
			// w.Write(responseValidationError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Start by persisting it
	instance = &service.Instance{
		InstanceID: instanceID,
		ServiceID:  provisioningRequest.ServiceID,
		PlanID:     provisioningRequest.PlanID,
		Status:     service.InstanceStateProvisioning,
	}
	err = instance.SetProvisioningResult(module.GetEmptyProvisioningResult())
	if err != nil {
		log.Println("error encoding empty provisioning result")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.store.WriteInstance(instance)
	if err != nil {
		log.Println("error storing new instance")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	provisioner, ok := s.provisioners[provisioningRequest.ServiceID]
	if !ok {
		log.Printf(
			`no provisioner found for service "%s"`,
			provisioningRequest.ServiceID,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	firstStepName, ok := provisioner.GetFirstStepName()
	if !ok {
		log.Printf(
			`no steps found for provisioning service "%s"`,
			provisioningRequest.ServiceID,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	task := model.NewTask(
		"provisionStep",
		map[string]string{
			"stepName":   firstStepName,
			"instanceID": instanceID,
		},
	)
	err = s.asyncEngine.SubmitTask(task)
	if err != nil {
		log.Println("error submitting provisioning task")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If we get all the way to here, we've been successful!
	w.WriteHeader(http.StatusAccepted)
}
