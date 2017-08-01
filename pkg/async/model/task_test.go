package model

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testTask     Task
	testTaskJSON []byte
)

func init() {
	jobName := "test-job"
	argName := "foo"
	argValue := "FOO"

	testTask = NewTask(
		jobName,
		map[string]string{
			argName: argValue,
		},
	)

	testTaskJSONStr := fmt.Sprintf(
		`{
			"id":"%s",
			"jobName":"%s",
			"args":{"%s":"%s"},
			"workerRejectionCount": %d
		}`,
		testTask.GetID(),
		jobName,
		argName,
		argValue,
		0,
	)
	testTaskJSONStr = strings.Replace(testTaskJSONStr, " ", "", -1)
	testTaskJSONStr = strings.Replace(testTaskJSONStr, "\n", "", -1)
	testTaskJSONStr = strings.Replace(testTaskJSONStr, "\t", "", -1)
	testTaskJSON = []byte(testTaskJSONStr)
}

func TestNewTaskFromJSON(t *testing.T) {
	task, err := NewTaskFromJSON(testTaskJSON)
	assert.Nil(t, err)
	assert.Equal(t, testTask, task)
}

func TestTaskToJSON(t *testing.T) {
	json, err := testTask.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, testTaskJSON, json)
}
