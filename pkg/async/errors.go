package async

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
	if e.err == nil {
		return "cleaner stopped"
	}
	return fmt.Sprintf("cleaner stopped: %s", e.err)
}

type errWorkerStopped struct {
	workerID string
	err      error
}

func (e *errWorkerStopped) Error() string {
	if e.err == nil {
		return fmt.Sprintf(`worker "%s" stopped`, e.workerID)
	}
	return fmt.Sprintf(`worker "%s" stopped: %s`, e.workerID, e.err)
}

type errHeartStopped struct {
	workerID string
	err      error
}

func (e *errHeartStopped) Error() string {
	if e.err == nil {
		return fmt.Sprintf(`worker "%s" heart stopped`, e.workerID)
	}
	return fmt.Sprintf(`worker "%s" heart stopped: %s`, e.workerID, e.err)
}

type errReceiveAndWorkStopped struct {
	workerID string
	err      error
}

type errWatchDelayedTasksStopped struct {
	workerID string
	err      error
}

func (e *errReceiveAndWorkStopped) Error() string {
	if e.err == nil {
		return fmt.Sprintf(`worker "%s" errReceiveAndWork stopped`, e.workerID)
	}
	return fmt.Sprintf(
		`worker "%s" errReceiveAndWork stopped: %s`,
		e.workerID,
		e.err,
	)
}

func (e *errWatchDelayedTasksStopped) Error() string {
	if e.err == nil {
		return fmt.Sprintf(`worker "%s" watchDelayedTasks stopped`, e.workerID)
	}
	return fmt.Sprintf(
		`worker "%s" watchDelayedTasks stopped: %s`,
		e.workerID,
		e.err,
	)
}

type errDuplicateJob struct {
	name string
}

func (e *errDuplicateJob) Error() string {
	return fmt.Sprintf(`duplicate job name "%s"`, e.name)
}

type errJobNotFound struct {
	name string
}

func (e *errJobNotFound) Error() string {
	return fmt.Sprintf(`no job named "%s" is registered with the worker`, e.name)
}
