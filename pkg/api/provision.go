package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/krancour/async"
)

func (s *server) provision(w http.ResponseWriter, r *http.Request) {
	instanceID := mux.Vars(r)["instance_id"]

	logFields := log.Fields{
		"instanceID": instanceID,
	}

	log.WithFields(logFields).Debug("received provisioning request")

	// TODO: krancour: Move all the accepts incomplete stuff into a filter.
	// This broker provisions everything asynchronously. If a client doesn't
	// explicitly indicate that they will accept an incomplete result, the
	// spec says to respond with a 422
	acceptsIncompleteStr := r.URL.Query().Get("accepts_incomplete")
	if acceptsIncompleteStr == "" {
		logFields["parameter"] = "accepts_incomplete=true" // nolint: goconst
		log.WithFields(logFields).Debug(
			"bad provisioning request: request is missing required query parameter",
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
			`bad provisioning request: query parameter has invalid value; only ` +
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
			"pre-provisioning error: error reading request body",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}
	defer r.Body.Close() // nolint: errcheck

	provisioningRequest, err := NewProvisioningRequestFromJSON(bodyBytes)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Debug(
			"bad provisioning request: error unmarshaling request body",
		)
		// krancour: Choosing to interpret this scenario as a bad request, as a
		// valid request, obviously contains valid, well-formed JSON
		// TODO: Write a more detailed response
		s.writeResponse(w, http.StatusBadRequest, generateMalformedRequestResponse())
		return
	}

	serviceID := provisioningRequest.ServiceID
	if serviceID == "" {
		logFields["field"] = "service_id"
		log.WithFields(logFields).Debug(
			"bad provisioning request: required request body field is missing",
		)
		s.writeResponse(w, http.StatusBadRequest, generateServiceIDRequiredResponse())
		return
	}

	planID := provisioningRequest.PlanID
	if planID == "" {
		logFields["field"] = "plan_id"
		log.WithFields(logFields).Debug(
			"bad provisioning request: required request body field is missing",
		)
		s.writeResponse(w, http.StatusBadRequest, generatePlanIDRequiredResponse())
		return
	}

	svc, ok := s.catalog.GetService(serviceID)
	if !ok || svc.IsEndOfLife() {
		logFields["serviceID"] = serviceID
		log.WithFields(logFields).Debug(
			"bad provisioning request: invalid serviceID",
		)
		s.writeResponse(w, http.StatusBadRequest, generateInvalidServiceIDResponse())
		return
	}

	plan, ok := svc.GetPlan(planID)
	if !ok || plan.IsEndOfLife() {
		logFields["serviceID"] = serviceID
		logFields["planID"] = planID
		log.WithFields(logFields).Debug(
			"bad provisioning request: invalid planID for service",
		)
		s.writeResponse(w, http.StatusBadRequest, generateInvalidPlanIDResponse())
		return
	}

	// Validate the provisioning parameters
	if err =
		plan.GetSchemas().ServiceInstances.ProvisioningParametersSchema.Validate(
			provisioningRequest.Parameters,
		); err != nil {
		var validationErr *service.ValidationError
		validationErr, ok = err.(*service.ValidationError)
		if ok {
			logFields["field"] = validationErr.Field
			logFields["issue"] = validationErr.Issue
			log.WithFields(logFields).Debug(
				"bad provisioning request: validation error",
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

	// Wrap the provisioning parameters with a "params" object that guides access
	// to the parameters using schema
	pps := plan.GetSchemas().ServiceInstances.ProvisioningParametersSchema
	provisioningParameters := &service.ProvisioningParameters{
		Parameters: service.Parameters{
			Schema: &pps,
			Data:   provisioningRequest.Parameters,
		},
	}

	// Unpack the generic bits of the parameter map...
	// Alias
	alias := provisioningParameters.GetString("alias")
	// Parent alias
	parentAlias := provisioningParameters.GetString("parentAlias")

	instance, ok, err := s.store.GetInstance(instanceID)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-provisioning error: error retrieving instance by id",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}
	if ok {
		// We land in here if an existing instance was found-- the OSB spec
		// obligates us to compare this instance to the one that was requested and
		// respond with 200 if they're identical or 409 otherwise. It actually seems
		// best to compare REQUESTS instead because instance objects also contain
		// instance details and other status information. So, let's reverse
		// engineer a request from the existing instance then compare it to the
		// current request.
		//
		// Two requests are the same if they are for the same serviceID, the same,
		// planID, and all other relevant fields are equal.
		if instance.ServiceID == serviceID &&
			instance.PlanID == planID &&
			((instance.ProvisioningParameters == nil && len(provisioningRequest.Parameters) == 0) || // nolint: lll
				(instance.ProvisioningParameters != nil && reflect.DeepEqual(instance.ProvisioningParameters.Data, provisioningRequest.Parameters))) { // nolint: lll
			// Per the spec, if fully provisioned, respond with a 200, else a 202.
			// Filling in a gap in the spec-- if the status is anything else, we'll
			// choose to respond with a 409
			switch instance.Status {
			case service.InstanceStateProvisioning:
				s.writeResponse(w, http.StatusAccepted, generateProvisionAcceptedResponse())
				return
			case service.InstanceStateProvisioned:
				s.writeResponse(w, http.StatusOK, generateEmptyResponse())
				return
			default:
				// TODO: Write a more detailed response
				s.writeResponse(w, http.StatusConflict, generateConflictResponse())
				return
			}
		}

		// We land in here if an existing instance was found, but its atrributes
		// vary from what was requested. The spec requires us to respond with a
		// 409
		s.writeResponse(w, http.StatusConflict, generateConflictResponse())
		return
	}

	// If we get to here, we need to provision a new instance.

	serviceManager := svc.GetServiceManager()

	provisioner, err := serviceManager.GetProvisioner(plan)
	if err != nil {
		logFields["serviceID"] = serviceID
		logFields["planID"] = provisioningRequest.PlanID
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"pre-provisioning error: error retrieving provisioner for service and " +
				"plan",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}

	firstStepName, ok := provisioner.GetFirstStepName()
	if !ok {
		logFields["serviceID"] = provisioningRequest.ServiceID
		logFields["planID"] = provisioningRequest.PlanID
		log.WithFields(logFields).Error(
			"pre-provisioning error: no steps found for provisioning service and " +
				"plan",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}

	instance = service.Instance{
		InstanceID:             instanceID,
		Alias:                  alias,
		ServiceID:              provisioningRequest.ServiceID,
		PlanID:                 provisioningRequest.PlanID,
		ProvisioningParameters: provisioningParameters,
		Status:                 service.InstanceStateProvisioning,
		ParentAlias:            parentAlias,
		Created:                time.Now(),
	}

	var task async.Task
	var waitForParent bool
	if waitForParent, err = s.isParentProvisioning(instance); err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"provisioning error: error related to parent instance",
		)
		s.writeResponse(w, http.StatusBadRequest, generateParentInvalidResponse())
		return
	} else if waitForParent {
		instance.Status = service.InstanceStateProvisioningDeferred
		task = async.NewDelayedTask(
			"checkParentStatus",
			map[string]string{
				"instanceID": instanceID,
			},
			time.Minute*1,
		)
		log.WithFields(logFields).Debug("parent not provisioned, waiting")
	} else {
		task = async.NewTask(
			"executeProvisioningStep",
			map[string]string{
				"stepName":   firstStepName,
				"instanceID": instanceID,
			},
		)
		log.WithFields(logFields).Debug(
			"no need to wait for parent, starting provision",
		)
	}

	if err = s.store.WriteInstance(instance); err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"provisioning error: error persisting new instance",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}

	if err = s.asyncEngine.SubmitTask(task); err != nil {
		logFields["step"] = firstStepName
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"provisioning error: error submitting provisioning task",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}

	// If we get all the way to here, we've been successful!
	s.writeResponse(w, http.StatusAccepted, generateProvisionAcceptedResponse())

	log.WithFields(logFields).Debug("asynchronous provisioning initiated")
}

func (s *server) isParentProvisioning(instance service.Instance) (bool, error) {
	//No parent, so no need to wait
	if instance.ParentAlias == "" {
		return false, nil
	}

	parent, parentFound, err := s.store.GetInstanceByAlias(instance.ParentAlias)

	if err != nil {
		log.WithFields(log.Fields{
			"error":       "waitforParent",
			"instanceID":  instance.InstanceID,
			"parentAlias": instance.ParentAlias,
		}).Error(
			"bad provision request: unable to retrieve parent",
		)
		return false, err
	}

	//Parent has was not found, so wait for that that to occur
	if !parentFound {
		return true, nil
	}

	//If parent failed, we should not even attempt to provision this
	if parent.Status == service.InstanceStateProvisioningFailed {
		log.WithFields(log.Fields{
			"error":      "waitforParent",
			"instanceID": instance.InstanceID,
			"parentID":   instance.Parent.InstanceID,
		}).Info(
			"bad provision request: parent failed provisioning",
		)
		return false, fmt.Errorf("error provisioning: parent provision failed")
	}

	//If parent is deprovisioning, we should not even attempt to provision this
	if parent.Status == service.InstanceStateDeprovisioning {
		log.WithFields(log.Fields{
			"error":      "waitforParent",
			"instanceID": instance.InstanceID,
			"parentID":   instance.Parent.InstanceID,
		}).Info(
			"bad provision request: parent is deprovisioning",
		)
		return false, fmt.Errorf("error provisioning: parent is deprovisioning")
	}

	//If parent is provisioned, then no need to wait.
	if parent.Status == service.InstanceStateProvisioned {
		return false, nil
	}

	return true, nil
}
