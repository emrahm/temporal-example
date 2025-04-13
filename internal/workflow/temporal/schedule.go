// temporal/scheduler.go

package temporal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.temporal.io/sdk/client"
)

// ScheduleConfig holds configuration for scheduling workflows
type ScheduleConfig struct {
	WorkflowName string
	Args         []string
	Interval     time.Duration
	ScheduleID   string
}

// Schedule creates a schedule that runs a workflow recursively
func (w *Workflow) Schedule(ctx context.Context, config ScheduleConfig) error {
	if _, ok := w.workflows[config.WorkflowName]; !ok {
		return fmt.Errorf("workflow %s not found", config.WorkflowName)
	}

	if config.Interval < time.Second {
		return fmt.Errorf("interval must be at least 1 second")
	}

	// Generate a schedule ID if not provided
	if config.ScheduleID == "" {
		id, err := gonanoid.New()
		if err != nil {
			return err
		}

		config.ScheduleID = id
	}

	scheduleID := fmt.Sprintf("schedule-%s-%s", config.WorkflowName, config.ScheduleID)
	workflowID := fmt.Sprintf("sw-%s-%s", config.WorkflowName, uuid.New().String())
	// Create the schedule with the interval
	scheduleHandle, err := w.serviceClient.ScheduleClient().Create(ctx, client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			Intervals: []client.ScheduleIntervalSpec{
				{
					Every: config.Interval,
				},
			},
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  w.workflows[config.WorkflowName],
			TaskQueue: "schedule",
			Args:      []interface{}{config.Args},
		},
	})

	internalDescription, err := scheduleHandle.Describe(ctx)
	if err != nil {
		return err
	}
	_ = internalDescription

	return nil
}
