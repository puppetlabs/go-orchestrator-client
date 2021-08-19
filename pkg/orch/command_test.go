package orch

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandTask(t *testing.T) {

	// Test success
	setupPostResponder(t, orchCommandTask, "command-task-request.json", "command-task-response.json")
	taskRequest := &TaskRequest{
		Description: "Sent from go-pe-client",
		Environment: "test-env-1",
		Task:        "package",
		Params: map[string]string{
			"action": "install",
			"name":   "httpd",
		},
		Scope: Scope{
			Application: "Wordpress_app[demo]",
			Nodes:       []string{"node1.example.com"},
			Query:       []interface{}{"from", "nodes", []interface{}{"~", "certname", ".*"}},
			NodeGroup:   "00000000-0000-4000-8000-000000000000",
		},
	}

	actual, err := orchClient.CommandTask(taskRequest)
	require.Nil(t, err)
	require.Equal(t, expectedCommandTaskResponse, actual)

	// Test Orchestrator error
	setupErrorResponder(t, orchCommandTask)
	actual, err = orchClient.CommandTask(taskRequest)
	require.Nil(t, actual)
	require.Equal(t, expectedError, err)

	//Test HTTP error
	setupResponderWithStatusCodeAndBody(t, orchCommandTask, http.StatusNotFound, []byte(`{"StatusCode": 400}`))
	actual, err = orchClient.CommandTask(taskRequest)
	testHTTPError(t, actual, err, http.StatusNotFound)

}

var expectedCommandTaskResponse = &JobID{Job: struct {
	ID   string "json:\"id\""
	Name string "json:\"name\""
}{ID: "https://orchestrator.example.com:8143/orchestrator/v1/jobs/1234", Name: "1234"}}

func TestCommandScheduleTask(t *testing.T) {

	// Test success
	setupPostResponder(t, orchCommandScheduleTask, "command-schedule_task-request.json", "command-schedule_task-response.json")
	scheduleTaskRequest := &ScheduleTaskRequest{
		Environment: "test-env-1",
		Task:        "package",
		Params: map[string]string{
			"action":  "install",
			"package": "httpd",
		},
		Scope: Scope{
			Nodes: []string{"node1.example.com"},
		},
		ScheduledTime: "2027-05-05T19:50:08Z",
	}

	actual, err := orchClient.CommandScheduleTask(scheduleTaskRequest)
	require.Nil(t, err)
	require.Equal(t, expectedCommandScheduleTaskResponse, actual)

	// Test error
	setupErrorResponder(t, orchCommandScheduleTask)
	actual, err = orchClient.CommandScheduleTask(scheduleTaskRequest)
	require.Nil(t, actual)
	require.Equal(t, expectedError, err)

	// Test HTTP error
	setupResponderWithStatusCodeAndBody(t, orchCommandScheduleTask, http.StatusBadRequest, []byte(`{"StatusCode": 400}`))
	actual, err = orchClient.CommandScheduleTask(scheduleTaskRequest)
	assert.Error(t, err)
	require.Nil(t, actual)
	testExpectError := getExpectedHTTPError(http.StatusBadRequest, "ignorefornow")
	httpErr, ok := err.(*HTTPError)
	if !ok {
		t.Error("Error returned is not of type HTTP error.")
	}
	require.Equal(t, httpErr.StatusCode, testExpectError.StatusCode)

	//Test Orchestrator error
	setupErrorResponder(t, orchCommandScheduleTask)
	actual, err = orchClient.CommandScheduleTask(scheduleTaskRequest)
	require.Nil(t, actual)
	assert.Error(t, err)
	require.Equal(t, expectedError, err)

	//Test HTTP error
	setupResponderWithStatusCodeAndBody(t, orchCommandScheduleTask, http.StatusNotFound, []byte(`{"StatusCode": 400}`))
	actual, err = orchClient.CommandScheduleTask(scheduleTaskRequest)
	testHTTPError(t, actual, err, http.StatusNotFound)
}

func TestCommandScheduleTaskWithScheduleOptions(t *testing.T) {

	var options = NewScheduleTaskOptions(time.Duration(24) * time.Hour)

	// Test success
	setupPostResponder(t, orchCommandScheduleTask, "command-schedule-interval_task-request.json", "command-schedule_task-response.json")
	scheduleTaskRequest := &ScheduleTaskRequest{
		Environment: "test-env-1",
		Task:        "package",
		Params: map[string]string{
			"action":  "install",
			"package": "httpd",
		},
		Scope: Scope{
			Nodes: []string{"node1.example.com"},
		},
		ScheduledTime:   "2027-05-05T19:50:08Z",
		ScheduleOptions: options,
	}

	actual, err := orchClient.CommandScheduleTask(scheduleTaskRequest)
	require.Nil(t, err)
	require.Equal(t, expectedCommandScheduleTaskResponse, actual)

	// Test error
	setupErrorResponder(t, orchCommandScheduleTask)
	actual, err = orchClient.CommandScheduleTask(scheduleTaskRequest)
	require.Nil(t, actual)
	require.Equal(t, expectedError, err)

	// Test HTTP error
	setupResponderWithStatusCodeAndBody(t, orchCommandScheduleTask, http.StatusBadRequest, []byte(`{"StatusCode": 400}`))
	actual, err = orchClient.CommandScheduleTask(scheduleTaskRequest)
	assert.Error(t, err)
	require.Nil(t, actual)
	testExpectError := getExpectedHTTPError(http.StatusBadRequest, "ignorefornow")
	httpErr, ok := err.(*HTTPError)
	if !ok {
		t.Error("Error returned is not of type HTTP error.")
	}
	require.Equal(t, httpErr.StatusCode, testExpectError.StatusCode)

	//Test Orchestrator error
	setupErrorResponder(t, orchCommandScheduleTask)
	actual, err = orchClient.CommandScheduleTask(scheduleTaskRequest)
	require.Nil(t, actual)
	assert.Error(t, err)
	require.Equal(t, expectedError, err)

	//Test HTTP error
	setupResponderWithStatusCodeAndBody(t, orchCommandScheduleTask, http.StatusNotFound, []byte(`{"StatusCode": 400}`))
	actual, err = orchClient.CommandScheduleTask(scheduleTaskRequest)
	testHTTPError(t, actual, err, http.StatusNotFound)
}

var expectedCommandScheduleTaskResponse = &ScheduledJobID{ScheduledJob: struct {
	ID   string "json:\"id\""
	Name string "json:\"name\""
}{ID: "https://orchestrator.example.com:8143/orchestrator/v1/scheduled_jobs/2", Name: "1234"}}

func TestNewScheduleTaskOptions(t *testing.T) {

	var optionsHour = NewScheduleTaskOptions(time.Duration(1) * time.Hour)
	var optionsMinutes = NewScheduleTaskOptions(time.Duration(60) * time.Minute)

	// Test Success from hour to seconds conversion
	require.Equal(t, expectedScheduleTaskOptions, optionsHour)

	// Test Success from minutes to seconds conversion
	require.Equal(t, expectedScheduleTaskOptions, optionsMinutes)
}

var expectedScheduleTaskOptions = &ScheduleOptions{
	Interval: Interval{
		Units: "seconds",
		Value: 3600,
	},
}

func TestCommandTaskTarget(t *testing.T) {

	// Test success
	setupPostResponder(t, orchCommandTaskTarget, "command-task_target-request.json", "command-task_target-response.json")
	taskTargetRequest := &TaskTargetRequest{
		DisplayName: "1",
		NodeGroups:  []string{"3c4df64f-7609-4d31-9c2d-acfa52ed66ec", "4932bfe7-69c4-412f-b15c-ac0a7c2883f1"},
		Nodes:       []string{"wss6c3w9wngpycg", "jjj2h5w8gpycgwn"},
		PQLQuery:    "nodes[certname] { catalog_environment = \"production\" }",
		Tasks:       []string{"package::install", "exec"},
	}

	actual, err := orchClient.CommandTaskTarget(taskTargetRequest)
	require.Nil(t, err)
	require.Equal(t, expectedCommandTaskTargetResponse, actual)

	// Test error
	setupErrorResponder(t, orchCommandTaskTarget)
	actual, err = orchClient.CommandTaskTarget(taskTargetRequest)
	require.Nil(t, actual)
	require.Equal(t, expectedError, err)

	//Test HTTP error
	setupResponderWithStatusCodeAndBody(t, orchCommandTaskTarget, http.StatusNotFound, []byte(`{"StatusCode": 400}`))
	actual, err = orchClient.CommandTaskTarget(taskTargetRequest)
	testHTTPError(t, actual, err, http.StatusNotFound)
}

var expectedCommandTaskTargetResponse = &TaskTargetJobID{TaskTargetJob: struct {
	ID   string "json:\"id\""
	Name string "json:\"name\""
}{ID: "https://orchestrator.example.com:8143/orchestrator/v1/scopes/task_targets/1", Name: "1"}}

func TestCommandPlanRun(t *testing.T) {

	// Test success
	setupPostResponder(t, orchCommandPlanRun, "command-plan_run-request.json", "command-plan_run-response.json")
	planRunRequest := &PlanRunRequest{
		Description: "Start the canary plan on node1 and node2",
		Params: map[string]interface{}{
			"canary":  1,
			"command": "whoami",
			"nodes":   []string{"node1.example.com", "node2.example.com"},
		},
		Name: "canary"}

	actual, err := orchClient.CommandPlanRun(planRunRequest)
	require.Nil(t, err)
	require.Equal(t, expectedCommandPlanRunResponse, actual)

	// Test error
	setupErrorResponder(t, orchCommandPlanRun)
	actual, err = orchClient.CommandPlanRun(planRunRequest)
	require.Nil(t, actual)
	require.Equal(t, expectedError, err)

	//Test HTTP error
	setupResponderWithStatusCodeAndBody(t, orchCommandPlanRun, http.StatusNotFound, []byte(`{"StatusCode": 400}`))
	actual, err = orchClient.CommandPlanRun(planRunRequest)
	testHTTPError(t, actual, err, http.StatusNotFound)
}

var expectedCommandPlanRunResponse = &PlanRunJobID{
	Name: "1234",
}

func TestCommandStop(t *testing.T) {

	// Test success
	setupPostResponder(t, orchCommandStop, "command-stop-request.json", "command-stop-response.json")
	stopRequest := &StopRequest{
		Job: "1234",
	}

	actual, err := orchClient.CommandStop(stopRequest)
	require.Nil(t, err)
	require.Equal(t, expectedCommandStopResponse, actual)

	// Test error
	setupErrorResponder(t, orchCommandStop)
	actual, err = orchClient.CommandStop(stopRequest)
	require.Nil(t, actual)
	require.Equal(t, expectedError, err)

	//Test HTTP error
	setupResponderWithStatusCodeAndBody(t, orchCommandStop, http.StatusNotFound, []byte(`{"StatusCode": 400}`))
	actual, err = orchClient.CommandStop(stopRequest)
	testHTTPError(t, actual, err, http.StatusNotFound)
}

var expectedCommandStopResponse = &StopJobID{Job: struct {
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	Nodes map[string]int `json:"nodes"`
}{ID: "https://orchestrator.example.com:8143/orchestrator/v1/jobs/1234", Name: "1234",
	Nodes: map[string]int{"new": 5, "errored": 1, "failed": 3, "finished": 5, "running": 8, "skipped": 2}}}

func TestCommandDeploy(t *testing.T) {

	// Test success
	setupPostResponder(t, orchCommandDeploy, "command-deploy-request.json", "command-deploy-response.json")
	deployRequest := &DeployRequest{
		Environment:        "production",
		EnforceEnvironment: true,
		Noop:               true,
		Scope: Scope{
			Nodes: []string{"node1.example.com"},
		},
	}

	actual, err := orchClient.CommandDeploy(deployRequest)
	require.Nil(t, err)
	require.Equal(t, expectedCommandDeployResponse, actual)

	// Test error
	setupErrorResponder(t, orchCommandDeploy)
	actual, err = orchClient.CommandDeploy(deployRequest)
	require.Nil(t, actual)
	require.Equal(t, expectedError, err)
}

var expectedCommandDeployResponse = &JobID{Job: struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}{ID: "https://orchestrator.example.com:8143/orchestrator/v1/jobs/1234", Name: "1234"}}
