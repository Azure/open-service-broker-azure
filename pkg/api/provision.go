package api

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"github.com/Azure/azure-service-broker/pkg/async/model"
	"github.com/Azure/azure-service-broker/pkg/service"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func (s *server) provision(w http.ResponseWriter, r *http.Request) {
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

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	provisioningRequest := &service.ProvisioningRequest{}
	err = service.GetProvisioningRequestFromJSONString(
		string(bodyBytes),
		provisioningRequest,
	)
	if err != nil {
		log.Debug("error parsing request body")
		// krancour: Choosing to interpret this scenario as a bad request, as a
		// valid request, obviously contains valid, well-formed JSON
		w.WriteHeader(http.StatusBadRequest)
		// TODO: Write a more detailed response
		w.Write(responseEmptyJSON)
		return
	}

	instanceID := mux.Vars(r)["instance_id"]
	log.WithFields(log.Fields{
		"instanceID": instanceID,
		"serviceID":  provisioningRequest.ServiceID,
		"planID":     provisioningRequest.PlanID,
	}).Debug("received provisioning request")

	if provisioningRequest.ServiceID == "" {
		log.Debug("request body parameter service_id is a required field")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseServiceIDRequired)
		return
	}
	if provisioningRequest.PlanID == "" {
		log.Debug("request body parameter plan_id is a required field")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responsePlanIDRequired)
		return
	}

	svc, ok := s.catalog.GetService(provisioningRequest.ServiceID)
	if !ok {
		log.WithFields(log.Fields{
			"instanceID": instanceID,
			"serviceID":  provisioningRequest.ServiceID,
			"planID":     provisioningRequest.PlanID,
		}).Debug("invalid serviceID")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseInvalidServiceID)
		return
	}

	_, ok = svc.GetPlan(provisioningRequest.PlanID)
	if !ok {
		log.WithFields(log.Fields{
			"instanceID": instanceID,
			"serviceID":  provisioningRequest.ServiceID,
			"planID":     provisioningRequest.PlanID,
		}).Debug("invalid planID for service")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseInvalidPlanID)
		return
	}

	module, ok := s.modules[provisioningRequest.ServiceID]
	if !ok {
		// We already validated that the serviceID and planID are legitimate. If
		// we don't find a module that handles the service, something is really
		// wrong.
		log.WithField(
			"serviceID",
			provisioningRequest.ServiceID,
		).Error("no module found for service")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}

	// Now that we know what module we're dealing with, we can get an instance
	// of the module-specific type for provisioningParameters and take a second
	// pass at parsing the request body
	provisioningRequest.Parameters = module.GetEmptyProvisioningParameters()
	err = service.GetProvisioningRequestFromJSONString(
		string(bodyBytes),
		provisioningRequest,
	)
	if err != nil {
		log.Debug("error parsing request body")
		// krancour: Choosing to interpret this scenario as a bad request, as a
		// valid request, obviously contains valid, well-formed JSON
		w.WriteHeader(http.StatusBadRequest)
		// TODO: Write a more detailed response
		w.Write(responseEmptyJSON)
		return
	}
	if provisioningRequest.Parameters == nil {
		provisioningRequest.Parameters = module.GetEmptyProvisioningParameters()
	}

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
	if ok {
		// We land in here if an existing instance was found-- the OSB spec
		// obligates us to compare this instance to the one that was requested and
		// respond with 200 if they're identical or 409 otherwise. It actually seems
		// best to compare REQUESTS instead because instance objects also contain
		// provisioning context and other status information. So, let's reverse
		// engineer a request from the existing instance then compare it to the
		// current request.
		previousProvisioningRequestParams := module.GetEmptyProvisioningParameters()
		err = instance.GetProvisioningParameters(
			previousProvisioningRequestParams,
			s.codec,
		)
		if err != nil {
			log.WithFields(log.Fields{
				"instanceID": instanceID,
				"error":      err,
			}).Error("error decoding persisted provisioningParameters")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(responseEmptyJSON)
			return
		}
		previousProvisioningRequest := &service.ProvisioningRequest{
			ServiceID:  instance.ServiceID,
			PlanID:     instance.PlanID,
			Parameters: previousProvisioningRequestParams,
		}
		if reflect.DeepEqual(provisioningRequest, previousProvisioningRequest) {
			// Per the spec, if fully provisioned, respond with a 200, else a 202.
			// Filling in a gap in the spec-- if the status is anything else, we'll
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
			// TODO: Send the correct response body-- this is a placeholder
			w.Write(responseEmptyJSON)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}

	provisioner, err := module.GetProvisioner()
	if err != nil {
		log.WithFields(log.Fields{
			"serviceID": instance.ServiceID,
			"error":     err,
		}).Error("error retrieving provisioner for service")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}
	firstStepName, ok := provisioner.GetFirstStepName()
	if !ok {
		log.WithField(
			"serviceID",
			instance.ServiceID,
		).Error("no steps found for provisioning service")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}

	instance = &service.Instance{
		InstanceID: instanceID,
		ServiceID:  provisioningRequest.ServiceID,
		PlanID:     provisioningRequest.PlanID,
		Status:     service.InstanceStateProvisioning,
	}
	err = instance.SetProvisioningParameters(
		provisioningRequest.Parameters,
		s.codec,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"instanceID": instanceID,
			"error":      err,
		}).Error("error encoding provisioningParameters")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}
	err = instance.SetProvisioningContext(
		module.GetEmptyProvisioningContext(),
		s.codec,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"instanceID": instanceID,
			"error":      err,
		}).Error("error encoding empty provisioningContext")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}
	err = s.store.WriteInstance(instance)
	if err != nil {
		log.WithFields(log.Fields{
			"instanceID": instanceID,
			"error":      err,
		}).Error("error storing new instance")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
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
		log.WithFields(log.Fields{
			"step":       firstStepName,
			"instanceID": instanceID,
			"error":      err,
		}).Error("error submitting provisioning task")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseEmptyJSON)
		return
	}

	// If we get all the way to here, we've been successful!
	w.WriteHeader(http.StatusAccepted)
	w.Write(responseEmptyJSON)
}
