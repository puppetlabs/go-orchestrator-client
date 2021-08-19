package orch

import (
	"fmt"
	"time"
)

const (
	orchCommandTask         = "/orchestrator/v1/command/task"
	orchCommandScheduleTask = "/orchestrator/v1/command/schedule_task"
	orchCommandTaskTarget   = "/orchestrator/v1/command/task_target"
	orchCommandPlanRun      = "/orchestrator/v1/command/plan_run"
	orchCommandStop         = "/orchestrator/v1/command/stop"
	orchCommandDeploy       = "/orchestrator/v1/command/deploy"
)

// CommandTask runs a permitted task job across a set of nodes (POST /command/task)
func (c *Client) CommandTask(taskRequest *TaskRequest) (*JobID, error) {
	payload := JobID{}
	r, err := c.resty.R().
		SetResult(&payload).
		SetBody(taskRequest).
		Post(orchCommandTask)

	if err = ProcessError(r, err, ""); err != nil {
		return nil, err
	}

	return &payload, nil
}

// JobID identifies a single Orchestrator job
type JobID struct {
	Job struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"job"`
}

// TaskRequest describes a task to be run
// Environment(string):	The environment to load the task from. The default is production.
// Scope	The PuppetDB query, list of nodes, or a node group ID. Application scopes are not allowed for task jobs. This key is required.
// Description	A description of the job.
// Noop	Whether to run the job in no-op mode. The default is false.
// Task	The task to run on the targets. This key is required.
// Params	The parameters to pass to the task. This key is required, but can be an empty object.
// Targets	A collection of target objects used for running the task on nodes through SSH or WinRM via Bolt server.
// Userdata	An object of arbitrary key/value data supplied to the job.
// REF: https://puppet.com/docs/pe/2019.8/orchestrator_api_commands_endpoint.html#orchestrator_api_post_command_task
type TaskRequest struct {
	Environment string            `json:"environment,omitempty"`
	Task        string            `json:"task"`
	Params      map[string]string `json:"params"`
	Scope       Scope             `json:"scope"`
	Description string            `json:"description"`
	Noop        bool              `json:"noop"`
}

// CommandScheduleTask schedules a task to run at a future date and time (POST /command/schedule_task)
func (c *Client) CommandScheduleTask(scheduleTaskRequest *ScheduleTaskRequest) (*ScheduledJobID, error) {
	payload := ScheduledJobID{}
	r, err := c.resty.R().
		SetResult(&payload).
		SetBody(scheduleTaskRequest).
		Post(orchCommandScheduleTask)

	if err = ProcessError(r, err, ""); err != nil {
		return nil, err
	}

	return &payload, nil
}

// NewScheduleTaskOptions will create task options based on time
func NewScheduleTaskOptions(interval time.Duration) *ScheduleOptions {
	return &ScheduleOptions{
		Interval: Interval{Units: "seconds", Value: int(interval.Seconds())},
	}
}

// ScheduledJobID identifies a single scheduled job
type ScheduledJobID struct {
	ScheduledJob struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"scheduled_job"`
}

// ScheduleTaskRequest describes a scheduled task
type ScheduleTaskRequest struct {
	Environment     string            `json:"environment,omitempty"`
	Task            string            `json:"task"`
	Params          map[string]string `json:"params"`
	Scope           Scope             `json:"scope"`
	ScheduledTime   string            `json:"scheduled_time"`
	ScheduleOptions *ScheduleOptions  `json:"schedule_options,omitempty"`
}

// CommandTaskTarget creates a new task-target (POST /command/task_target)
func (c *Client) CommandTaskTarget(taskTargetRequest *TaskTargetRequest) (*TaskTargetJobID, error) {
	payload := TaskTargetJobID{}
	r, err := c.resty.R().
		SetResult(&payload).
		SetBody(taskTargetRequest).
		Post(orchCommandTaskTarget)

	if err = ProcessError(r, err, fmt.Sprintf("%s error: %s", orchCommandTaskTarget, r.Status())); err != nil {
		return nil, err
	}

	return &payload, nil
}

// TaskTargetJobID identifies a task_target job
type TaskTargetJobID struct {
	TaskTargetJob struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"task_target"`
}

// TaskTargetRequest is a collection of tasks, nodes and node groups that define a permission group
type TaskTargetRequest struct {
	DisplayName string   `json:"display_name"`
	Tasks       []string `json:"tasks,omitempty"`
	AllTasks    bool     `json:"all_tasks,omitempty"`
	Nodes       []string `json:"nodes"`
	NodeGroups  []string `json:"node_groups"`
	PQLQuery    string   `json:"pql_query,omitempty"`
}

// CommandPlanRun runs a plan via the plan executor (POST /command/plan_run)
func (c *Client) CommandPlanRun(planRunRequest *PlanRunRequest) (*PlanRunJobID, error) {
	payload := PlanRunJobID{}
	r, err := c.resty.R().
		SetResult(&payload).
		SetBody(planRunRequest).
		Post(orchCommandPlanRun)

	if err = ProcessError(r, err, fmt.Sprintf("%s error: %s", orchCommandPlanRun, r.Status())); err != nil {
		return nil, err
	}

	return &payload, nil
}

// PlanRunJobID identifies a plan_run job
type PlanRunJobID struct {
	Name string `json:"name"`
}

// PlanRunRequest describes a plan_run request
type PlanRunRequest struct {
	Name        string                 `json:"plan_name"`
	Params      map[string]interface{} `json:"params"`
	Environment string                 `json:"environment,omitempty"`
	Description string                 `json:"description,omitempty"`
}

// CommandStop stops a orchestrator job that is currently in progress (POST /command/stop)
func (c *Client) CommandStop(stopRequest *StopRequest) (*StopJobID, error) {
	payload := StopJobID{}
	r, err := c.resty.R().
		SetResult(&payload).
		SetBody(stopRequest).
		Post(orchCommandStop)

	if err = ProcessError(r, err, fmt.Sprintf("%s error: %s", orchCommandStop, r.Status())); err != nil {
		return nil, err
	}

	return &payload, nil
}

// StopJobID describes jobs that were stopped successfully
type StopJobID struct {
	Job struct {
		ID    string         `json:"id"`
		Name  string         `json:"name"`
		Nodes map[string]int `json:"nodes"`
	} `json:"job"`
}

// StopRequest describes a stop request
type StopRequest struct {
	Job string `json:"job"`
}

// CommandDeploy runs the orchestrator across all nodes in an environment (POST /command/deploy)
func (c *Client) CommandDeploy(deployRequest *DeployRequest) (*JobID, error) {
	payload := JobID{}
	r, err := c.resty.R().
		SetResult(&payload).
		SetBody(deployRequest).
		Post(orchCommandDeploy)

	if err = ProcessError(r, err, fmt.Sprintf("%s error: %s", orchCommandDeploy, r.Status())); err != nil {
		return nil, err
	}

	return &payload, nil
}

// DeployRequest describes a deploy request
type DeployRequest struct {
	Environment        string `json:"environment"`
	Scope              Scope  `json:"scope,omitempty"`
	Description        string `json:"description,omitempty"`
	Noop               bool   `json:"noop,omitempty"`
	NoNoop             bool   `json:"no_noop,omitempty"`
	Concurrency        int    `json:"concurrency,omitempty"`
	EnforceEnvironment bool   `json:"enforce_environment"`
	Debug              bool   `json:"debug,omitempty"`
	Trace              bool   `json:"trace,omitempty"`
	Evaltrace          bool   `json:"evaltrace,omitempty"`
	Target             string `json:"target,omitempty"`
}
