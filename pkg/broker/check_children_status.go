package broker

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/async"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
)

func (b *broker) doCheckChildrenStatuses(
	_ context.Context,
	task async.Task,
) ([]async.Task, error) {
	instanceID, ok := task.GetArgs()["instanceID"]
	if !ok {
		return nil, errors.New(`missing required argument "instanceID"`)
	}
	instance, ok, err := b.store.GetInstance(instanceID)
	if !ok {
		return nil, b.handleDeprovisioningError(
			instanceID,
			"checkChildrenStatus",
			nil,
			"error loading persisted instance",
		)
	}
	if err != nil {
		return nil, b.handleDeprovisioningError(
			instanceID,
			"checkChildrenStatus",
			err,
			"error loading persisted instance",
		)
	}
	childCount, err := b.store.GetInstanceChildCountByAlias(instance.Alias)
	if err != nil {
		log.WithFields(log.Fields{
			"step":       "checkChildrenStatus",
			"instanceID": instanceID,
			"error":      err,
		}).Error(
			"deprovisioning error: error determining child count",
		)
		return nil, b.handleDeprovisioningError(
			instance,
			"checkChildrenStatus",
			err,
			"error determining child count",
		)
	}
	if childCount > 0 {
		//Put this task back into the queue
		log.WithFields(log.Fields{
			"instanceID":          instanceID,
			"provisionedChildren": childCount,
		}).Debug("children not deprovisioned, will wait again")
		return []async.Task{
			async.NewDelayedTask(
				"checkChildrenStatus",
				map[string]string{
					"instanceID": instanceID,
				},
				time.Minute*1,
			),
		}, nil
	}
	serviceManager := instance.Service.GetServiceManager()
	var deprovisioner service.Deprovisioner
	deprovisioner, err = serviceManager.GetDeprovisioner(instance.Plan)
	if err != nil {
		log.WithFields(log.Fields{
			"instanceID": instanceID,
			"serviceID":  instance.ServiceID,
			"planID":     instance.PlanID,
			"error":      err,
		}).Error(
			"deprovisioning error: error retrieving deprovisioner for " +
				"service and plan",
		)
		return nil, b.handleDeprovisioningError(
			instance,
			"checkChildrenStatus",
			err,
			"error retrieving deprovisioner for service and service",
		)
	}
	deprovisionFirstStep, ok := deprovisioner.GetFirstStepName()
	if !ok {
		log.WithFields(log.Fields{
			"instanceID": instanceID,
			"serviceID":  instance.ServiceID,
			"planID":     instance.PlanID,
		}).Error(
			"pre-deprovisioning error: no steps found for deprovisioning " +
				"service and plan",
		)
		return nil, b.handleDeprovisioningError(
			instance,
			"checkChildrenStatus",
			nil,
			"error: no steps found for deprovisioning service ance plan",
		)
	}

	//Put the real deprovision task into the queue
	log.WithFields(log.Fields{
		"step":       "checkChildrenStatus",
		"instanceID": instanceID,
	}).Debug("children deprovisioned,  sending start deprovision task")
	return []async.Task{
		async.NewTask(
			"executeDeprovisioningStep",
			map[string]string{
				"stepName":   deprovisionFirstStep,
				"instanceID": instanceID,
			},
		),
	}, nil
}
