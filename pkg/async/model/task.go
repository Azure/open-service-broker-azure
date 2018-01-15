package model

import (
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Task is an interface to be implemented by types that represent a single
// asynchronous task
type Task interface {
	GetID() string
	GetJobName() string
	GetArgs() map[string]string
	GetWorkerRejectionCount() int
	IncrementWorkerRejectionCount() int
	ToJSON() ([]byte, error)
	GetExecuteTime() *time.Time
}

type task struct {
	ID                   string            `json:"id"`
	JobName              string            `json:"jobName"`
	Args                 map[string]string `json:"args"`
	WorkerRejectionCount int               `json:"workerRejectionCount"`
	ExecuteTime          *time.Time        `json:"executeTime"`
}

// NewTask returns a new task
func NewTask(jobName string, args map[string]string) Task {
	t := &task{
		JobName: jobName,
		Args:    args,
	}
	t.ID = uuid.NewV4().String()
	return t
}

// NewDelayedTask returns a new task that will fire after a specified duration
func NewDelayedTask(
	jobName string,
	args map[string]string,
	delay time.Duration,
) Task {
	tIface := NewTask(jobName, args)
	t := tIface.(*task)
	executeTime := time.Now().Add(delay)
	t.ExecuteTime = &executeTime
	return t
}

// NewScheduledTask returns a new task that will fire after a specified time
func NewScheduledTask(
	jobName string,
	args map[string]string,
	time time.Time,
) Task {
	tIface := NewTask(jobName, args)
	t := tIface.(*task)
	t.ExecuteTime = &time
	return t
}

// NewTaskFromJSON returns a new Task unmarshalled from the provided []byte
func NewTaskFromJSON(jsonBytes []byte) (Task, error) {
	t := &task{}
	if err := json.Unmarshal(jsonBytes, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *task) GetID() string {
	return t.ID
}

func (t *task) GetJobName() string {
	return t.JobName
}

func (t *task) GetArgs() map[string]string {
	return t.Args
}

func (t *task) GetWorkerRejectionCount() int {
	return t.WorkerRejectionCount
}

func (t *task) IncrementWorkerRejectionCount() int {
	t.WorkerRejectionCount++
	return t.WorkerRejectionCount
}

// ToJSON returns a []byte containing a JSON representation of the task
func (t *task) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

func (t *task) GetExecuteTime() *time.Time {
	return t.ExecuteTime
}
