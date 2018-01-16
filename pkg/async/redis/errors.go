package redis

import "fmt"

type errCleaning struct {
	workerID string
	err      error
}

func (e *errCleaning) Error() string {
	if e.workerID == "" {
		return fmt.Sprintf(
			`error cleaning up after dead workers: %s`,
			e.err,
		)
	}
	return fmt.Sprintf(
		`error cleaning up after dead worker "%s": %s`,
		e.workerID,
		e.err,
	)
}

type errCleanerStopped struct {
	err error
}

func (e *errCleanerStopped) Error() string {
	baseMsg := "cleaner stopped"
	if e.err == nil {
		return baseMsg
	}
	return fmt.Sprintf("%s: %s", baseMsg, e.err)
}

type errWorkerStopped struct {
	workerID string
	err      error
}

func (e *errWorkerStopped) Error() string {
	baseMsg := fmt.Sprintf(`worker "%s" stopped`, e.workerID)
	if e.err == nil {
		return baseMsg
	}
	return fmt.Sprintf("%s: %s", baseMsg, e.err)
}

type errHeartStopped struct {
	workerID string
	err      error
}

func (e *errHeartStopped) Error() string {
	baseMsg := fmt.Sprintf(`worker "%s" heart stopped`, e.workerID)
	if e.err == nil {
		return baseMsg
	}
	return fmt.Sprintf("%s: %s", baseMsg, e.err)
}

type errReceiverStopped struct {
	workerID  string
	queueName string
	err       error
}

func (e *errReceiverStopped) Error() string {
	baseMsg := fmt.Sprintf(
		`worker "%s" receiver for queue "%s" stopped`,
		e.workerID,
		e.queueName,
	)
	if e.err == nil {
		return baseMsg
	}
	return fmt.Sprintf("%s: %s", baseMsg, e.err)
}

type errTaskExecutorStopped struct {
	workerID string
	err      error
}

func (e *errTaskExecutorStopped) Error() string {
	baseMsg := fmt.Sprintf(`worker "%s" task executor stopped`, e.workerID)
	if e.err == nil {
		return baseMsg
	}
	return fmt.Sprintf("%s: %s", baseMsg, e.err)
}

type errDeferredTaskWatcherStopped struct {
	workerID string
	err      error
}

func (e *errDeferredTaskWatcherStopped) Error() string {
	baseMsg := fmt.Sprintf(
		`worker "%s" deferred task watcher stopped`,
		e.workerID,
	)
	if e.err == nil {
		return baseMsg
	}
	return fmt.Sprintf("%s: %s", baseMsg, e.err)
}

type errDuplicateJob struct {
	name string
}

func (e *errDuplicateJob) Error() string {
	return fmt.Sprintf(`duplicate job name "%s"`, e.name)
}

type errHeartbeat struct {
	workerID string
	err      error
}

func (e *errHeartbeat) Error() string {
	return fmt.Sprintf(
		`error sending heartbeat for worker "%s": %s`,
		e.workerID,
		e.err,
	)
}
