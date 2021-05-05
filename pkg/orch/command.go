package orch

import (
	"fmt"
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
	if err != nil {
		return nil, FormatError(r, err.Error())
	}
	if r.IsError() {
		if r.Error() != nil {
			return nil, FormatError(r)
		}
		return nil, FormatError(r)
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
type TaskRequest struct {
	Environment string            `json:"environment,omitempty"`
	Task        string            `json:"task"`
	Params      map[string]string `json:"params"`
	Scope       Scope             `json:"scope"`
}

// CommandScheduleTask schedules a task to run at a future date and time (POST /command/schedule_task)
func (c *Client) CommandScheduleTask(scheduleTaskRequest *ScheduleTaskRequest) (*ScheduledJobID, error) {
	payload := ScheduledJobID{}
	r, err := c.resty.R().
		SetResult(&payload).
		SetBody(scheduleTaskRequest).
		Post(orchCommandScheduleTask)
	if err != nil {
		return nil, FormatError(r, err.Error())
	}
	if r.IsError() {
		if r.Error() != nil {
			return nil, FormatError(r)
		}
		return nil, FormatError(r)
	}
	return &payload, nil
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
	Environment   string            `json:"environment,omitempty"`
	Task          string            `json:"task"`
	Params        map[string]string `json:"params"`
	Scope         Scope             `json:"scope"`
	ScheduledTime string            `json:"scheduled_time"`
}

// CommandTaskTarget creates a new task-target (POST /command/task_target)
func (c *Client) CommandTaskTarget(taskTargetRequest *TaskTargetRequest) (*TaskTargetJobID, error) {
	payload := TaskTargetJobID{}
	r, err := c.resty.R().
		SetResult(&payload).
		SetBody(taskTargetRequest).
		Post(orchCommandTaskTarget)
	if err != nil {
		return nil, err
	}
	if r.IsError() {
		if r.Error() != nil {
			return nil, r.Error().(error)
		}
		return nil, fmt.Errorf("%s error: %s", orchCommandTaskTarget, r.Status())
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
	if err != nil {
		return nil, err
	}
	if r.IsError() {
		if r.Error() != nil {
			return nil, r.Error().(error)
		}
		return nil, fmt.Errorf("%s error: %s", orchCommandPlanRun, r.Status())
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
	if err != nil {
		return nil, err
	}
	if r.IsError() {
		if r.Error() != nil {
			return nil, r.Error().(error)
		}
		return nil, fmt.Errorf("%s error: %s", orchCommandStop, r.Status())
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
	if err != nil {
		return nil, err
	}
	if r.IsError() {
		if r.Error() != nil {
			return nil, r.Error().(error)
		}
		return nil, fmt.Errorf("%s error: %s", orchCommandDeploy, r.Status())
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
	EnforceEnvironment bool   `json:"enforce_environment,omitempty"`
	Debug              bool   `json:"debug,omitempty"`
	Trace              bool   `json:"trace,omitempty"`
	Evaltrace          bool   `json:"evaltrace,omitempty"`
	Target             string `json:"target,omitempty"`
}
