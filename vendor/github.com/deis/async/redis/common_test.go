package redis

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

var errSome = errors.New("an error")

func getDisposableQueueName() string {
	return uuid.NewV4().String()
}

func getDisposableWorkerID() string {
	return uuid.NewV4().String()
}

func getDisposableWorkerSetName() string {
	return uuid.NewV4().String()
}
