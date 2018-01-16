package async

import "fmt"

const (
	workerSetName         = "workers"
	aliveIndicator        = "alive"
	pendingTaskQueueName  = "pendingTasks"
	deferredTaskQueueName = "deferredTasks"
)

func getActiveTaskQueueName(workerID string) string {
	return fmt.Sprintf("active-tasks:%s", workerID)
}

func getWatchedTaskQueueName(workerID string) string {
	return fmt.Sprintf("watched-tasks:%s", workerID)
}
