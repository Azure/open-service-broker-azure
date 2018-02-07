package api

import (
	"fmt"
	"net/http"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func (s *server) poll(
	w http.ResponseWriter,
	r *http.Request,
) {
	instanceID := mux.Vars(r)["instance_id"]

	logFields := log.Fields{
		"instanceID": instanceID,
	}

	log.WithFields(logFields).Debug("received polling request")

	operation := r.URL.Query().Get("operation")
	if operation == "" {
		logFields["parameter"] = "operation"
		log.WithFields(logFields).Debug(
			"bad polling request: request is missing required query parameter",
		)
		s.writeResponse(w, http.StatusBadRequest, generateOperationRequiredResponse())
		return
	}
	if operation != OperationProvisioning &&
		operation != OperationDeprovisioning &&
		operation != OperationUpdating {
		logFields["operation"] = operation
		log.WithFields(logFields).Debug(
			fmt.Sprintf(
				`bad polling request: query parameter has invalid value; only "%s",`+
					` %s, and "%s" are accepted`,
				OperationProvisioning,
				OperationDeprovisioning,
				OperationUpdating,
			),
		)
		s.writeResponse(w, http.StatusBadRequest, generateOperationInvalidResponse())
		return
	}

	logFields["operation"] = operation

	instance, ok, err := s.store.GetInstance(instanceID)
	if err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error(
			"polling error: error retrieving instance by id",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		return
	}
	if !ok {
		if operation == OperationDeprovisioning {
			s.writeResponse(w, http.StatusGone, generateEmptyResponse())
			return
		}
		s.writeResponse(w, http.StatusNotFound, generateEmptyResponse())
		return
	}

	logFields["status"] = instance.Status

	if operation == OperationProvisioning {
		switch instance.Status {
		case service.InstanceStateProvisioning:
			log.WithFields(logFields).Debug(
				"provisioning is in progress",
			)
			s.writeResponse(w, http.StatusOK, generateOperationInProgressResponse())
		case service.InstanceStateProvisioned:
			log.WithFields(logFields).Debug(
				"provisioning is complete",
			)
			s.writeResponse(w, http.StatusOK, generateOperationSucceededResponse())
		case service.InstanceStateProvisioningFailed:
			log.WithFields(logFields).Debug(
				"provisioning has failed",
			)
			s.writeResponse(w, http.StatusOK, generateOperationFailedResponse())
		default:
			log.WithFields(logFields).Error(
				"polling error: instance is in an unknown or invalid state",
			)
			s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		}
		return
	}

	if operation == OperationUpdating {
		switch instance.Status {
		case service.InstanceStateUpdating:
			log.WithFields(logFields).Debug(
				"updating is in progress",
			)
			s.writeResponse(w, http.StatusOK, generateOperationInProgressResponse())
		case service.InstanceStateUpdated:
			log.WithFields(logFields).Debug(
				"updating is complete",
			)
			s.writeResponse(w, http.StatusOK, generateOperationSucceededResponse())
		case service.InstanceStateUpdatingFailed:
			log.WithFields(logFields).Debug(
				"updating has failed",
			)
			s.writeResponse(w, http.StatusOK, generateOperationFailedResponse())
		default:
			log.WithFields(logFields).Error(
				"polling error: instance is in an unknown or invalid state",
			)
			s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
		}
		return
	}

	switch instance.Status {
	case service.InstanceStateDeprovisioning:
		log.WithFields(logFields).Debug(
			"deprovisioning is in progress",
		)
		s.writeResponse(w, http.StatusOK, generateOperationInProgressResponse())
	case service.InstanceStateDeprovisioningFailed:
		log.WithFields(logFields).Debug(
			"deprovisioning has failed",
		)
		s.writeResponse(w, http.StatusOK, generateOperationFailedResponse())
	default:
		log.WithFields(logFields).Error(
			"polling error: instance is in an unknown or invalid state",
		)
		s.writeResponse(w, http.StatusInternalServerError, generateEmptyResponse())
	}

}
