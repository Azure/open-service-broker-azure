package model

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testTask     Task
	testTaskJSON string
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

	testTaskJSON = fmt.Sprintf(
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
	testTaskJSON = strings.Replace(testTaskJSON, " ", "", -1)
	testTaskJSON = strings.Replace(testTaskJSON, "\n", "", -1)
	testTaskJSON = strings.Replace(testTaskJSON, "\t", "", -1)
}

func TestNewTaskFromJSONString(t *testing.T) {
	task, err := NewTaskFromJSONString(testTaskJSON)
	assert.Nil(t, err)
	assert.Equal(t, testTask, task)
}

func TestTaskToJSON(t *testing.T) {
	jsonStr, err := testTask.ToJSONString()
	assert.Nil(t, err)
	assert.Equal(t, testTaskJSON, jsonStr)
}
