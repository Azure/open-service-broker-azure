package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/async"
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/satori/uuid"
)

func (s *server) provision(w http.ResponseWriter, r *http.Request) {
	instanceID := mux.Vars(r)["instance_id"]

	logFields := log.Fields{
		"instanceID": instanceID,
	}

	log.WithFields(logFields).Debug("received provisioning request")

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
	if !ok {
		logFields["serviceID"] = serviceID
		log.WithFields(logFields).Debug(
			"bad provisioning request: invalid serviceID",
		)
		s.writeResponse(w, http.StatusBadRequest, generateInvalidServiceIDResponse())
		return
	}

	plan, ok := svc.GetPlan(planID)
	if !ok {
		logFields["serviceID"] = serviceID
		logFields["planID"] = planID
		log.WithFields(logFields).Debug(
			"bad provisioning request: invalid planID for service",
		)
		s.writeResponse(w, http.StatusBadRequest, generateInvalidPlanIDResponse())
		return
	}

	serviceManager := svc.GetServiceManager()

	// Unpack the parameter map...

	// Location...
	location := ""
	locIface, ok := provisioningRequest.Parameters["location"]
	if ok {
		location, ok = locIface.(string)
		if !ok {
			s.handlePossibleValidationError(
				service.NewValidationError(
					"location",
					fmt.Sprintf(`"%v" is not a string`, locIface),
				),
				w,
				logFields,
			)
			return
		}
	}
	location = s.getLocation(location)

	// Resource group...
	requestedResourceGroup := ""
	rgIface, ok := provisioningRequest.Parameters["resourceGroup"]
	if ok {
		requestedResourceGroup, ok = rgIface.(string)
		if !ok {
			s.handlePossibleValidationError(
				service.NewValidationError(
					"resourceGroup",
					fmt.Sprintf(`"%v" is not a string`, rgIface),
				),
				w,
				logFields,
			)
			return
		}
	}
	resourceGroup := s.getResourceGroup(requestedResourceGroup)

	// Tags...
	var tags map[string]string
	tagsIface, ok := provisioningRequest.Parameters["tags"]
	if ok {
		mapTagsIfaces, ok := tagsIface.(map[string]interface{})
		if !ok {
			s.handlePossibleValidationError(
				service.NewValidationError(
					"tags",
					fmt.Sprintf(`"%v" is not a map[string]string`, tagsIface),
				),
				w,
				logFields,
			)
			return
		}
		decoderConfig := &mapstructure.DecoderConfig{
			Result: &tags,
		}
		decoder, err := mapstructure.NewDecoder(decoderConfig)
		if err != nil {
			logFields["error"] = err
			log.WithFields(logFields).Error(
				"error building tag map decoder",
			)
			s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
			return
		}
		err = decoder.Decode(mapTagsIfaces)
		if err != nil {
			log.WithFields(logFields).Debug(
				"bad provisioning request: error decoding tags into map[string]string",
			)
			// This scenario is bad request because it means the tags weren't
			// a map[string]string, as we expected.
			s.writeResponse(w, http.StatusBadRequest, generateMalformedTagsResponse())
			return
		}
	}

	// Alias
	alias := ""
	aliasIface, ok := provisioningRequest.Parameters["alias"]
	if ok {
		alias, ok = aliasIface.(string)
		if !ok {
			s.handlePossibleValidationError(
				service.NewValidationError(
					"alias",
					fmt.Sprintf(`"%v" is not a string`, locIface),
				),
				w,
				logFields,
			)
			return
		}
	}

	// Parent alias
	parentAlias := ""
	parentAliasIface, ok := provisioningRequest.Parameters["parentAlias"]
	if ok {
		parentAlias, ok = parentAliasIface.(string)
		if !ok {
			s.handlePossibleValidationError(
				service.NewValidationError(
					"parentAlias",
					fmt.Sprintf(`"%v" is not a string`, parentAlias),
				),
				w,
				logFields,
			)
			return
		}
	}

	// Alias
	alias := ""
	aliasIface, ok := provisioningRequest.Parameters["alias"]
	if ok {
		alias, ok = aliasIface.(string)
		if !ok {
			s.handlePossibleValidationError(
				service.NewValidationError(
					"alias",
					fmt.Sprintf(`"%v" is not a string`, locIface),
				),
				w,
				logFields,
			)
			return
		}
	}

	// Parent alias
	parentAlias := ""
	parentAliasIface, ok := provisioningRequest.Parameters["parentAlias"]
	if ok {
		parentAlias, ok = parentAliasIface.(string)
		if !ok {
			s.handlePossibleValidationError(
				service.NewValidationError(
					"parentAlias",
					fmt.Sprintf(`"%v" is not a string`, parentAlias),
				),
				w,
				logFields,
			)
			return
		}
	}

	// Now service-specific parameters...
	provisioningParameters := serviceManager.GetEmptyProvisioningParameters()
	decoderConfig := &mapstructure.DecoderConfig{
		TagName: "json",
		Result:  provisioningParameters,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"error building parameter map decoder",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}
	err = decoder.Decode(provisioningRequest.Parameters)
	if err != nil {
		log.WithFields(logFields).Debug(
			"bad provisioning request: error decoding parameter map into " +
				"service-specific parameters",
		)
		// krancour: Choosing to interpret this scenario as a bad request since the
		// probable cause would be disagreement between provided and expected types
		s.writeResponse(w, http.StatusBadRequest, generateInvalidRequestResponse())
		return
	}

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
			instance.Location == location &&
			// If resourceGroup wasn't specified, we know one would be generated, so
			// we're going to not take the equality of the requested resourceGroup
			// and the existing resourceGroup into account if the requested
			// resourceGroup is the empty string...
			(requestedResourceGroup == "" ||
				instance.ResourceGroup == resourceGroup) &&
			reflect.DeepEqual(instance.Tags, tags) &&
			reflect.DeepEqual(
				instance.ProvisioningParameters,
				provisioningParameters,
			) {
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

	// Start by validating the location
	err = s.validateLocation(svc, location)
	if err != nil {
		s.handlePossibleValidationError(err, w, logFields)
		return
	}

	// Validate alias (only applies if this service type has children)
	err = s.validateAlias(svc, alias)
	if err != nil {
		s.handlePossibleValidationError(err, w, logFields)
		return
	}

	// Validate parent alias (only applies if this service type has a parent)
	err = s.validateParentAlias(svc, parentAlias)
	if err != nil {
		s.handlePossibleValidationError(err, w, logFields)
		return
	}

	// Then validate service-specific provisioning parameters
	err = serviceManager.ValidateProvisioningParameters(provisioningParameters)
	if err != nil {
		s.handlePossibleValidationError(err, w, logFields)
		return
	}

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
		Location:               location,
		ResourceGroup:          resourceGroup,
		ParentAlias:            parentAlias,
		Tags:                   tags,
		Details:                serviceManager.GetEmptyInstanceDetails(),
		Created:                time.Now(),
	}

	waitForParent, err := s.isParentProvisioning(instance)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"provisioning error: error related to parent instance",
		)
		s.writeResponse(w, http.StatusBadRequest, generateParentInvalidResponse())
		return
	}

	if err = s.store.WriteInstance(instance); err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"provisioning error: error persisting new instance",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}

	var task async.Task
	if waitForParent {
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

func (s *server) validateLocation(svc service.Service, location string) error {
	// Validate location only if this is a "root" service type (i.e. has no
	// parent)
	if svc.GetParentServiceID() == "" {
		if (location == "" && s.defaultAzureLocation == "") ||
			(location != "" && !azure.IsValidLocation(location)) {
			return service.NewValidationError(
				"location",
				fmt.Sprintf(`invalid location: "%s"`, location),
			)
		}
	}
	return nil
}

func (s *server) validateAlias(svc service.Service, alias string) error {
	if svc.GetChildServiceID() != "" && alias == "" {
		return service.NewValidationError(
			"alias",
			fmt.Sprintf(`invalid alias: "%s"`, alias),
		)
	}
	return nil
}

func (s *server) validateParentAlias(
	svc service.Service,
	parentAlias string,
) error {
	if svc.GetParentServiceID() != "" && parentAlias == "" {
		return service.NewValidationError(
			"parentAlias",
			fmt.Sprintf(`invalid parentAlias: "%s"`, parentAlias),
		)
	}
	return nil
}

func (s *server) handlePossibleValidationError(
	err error,
	w http.ResponseWriter,
	logFields log.Fields,
) {
	validationErr, ok := err.(*service.ValidationError)
	if ok {
		logFields["field"] = validationErr.Field
		logFields["issue"] = validationErr.Issue
		log.WithFields(logFields).Debug(
			"bad provisioning request: validation error",
		)
		response := generateValidationFailedResponse(validationErr)
		s.writeResponse(w, http.StatusBadRequest, response)
		return
	}
	s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
}

func (s *server) getLocation(location string) string {
	if location != "" {
		return location
	}
	return s.defaultAzureLocation
}

func (s *server) getResourceGroup(resourceGroup string) string {
	if resourceGroup != "" {
		return resourceGroup
	}
	if s.defaultAzureResourceGroup != "" {
		return s.defaultAzureResourceGroup
	}
	return uuid.NewV4().String()
}
