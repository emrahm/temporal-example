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
// Schedule creates a schedule that runs a workflow recursively
func (w *Workflow) Schedule(ctx context.Context, config ScheduleConfig) error {
	// Validate workflow exists
	workflowFunc, ok := w.workflows[config.WorkflowName]
	if !ok {
		return fmt.Errorf("workflow %s not found", config.WorkflowName)
	}

	// Validate interval
	if config.Interval < time.Second {
		return fmt.Errorf("interval must be at least 1 second")
	}

	// Generate a schedule ID if not provided
	if config.ScheduleID == "" {
		id, err := gonanoid.New()
		if err != nil {
			return fmt.Errorf("failed to generate schedule ID: %w", err)
		}
		config.ScheduleID = id
	}

	// Create consistent naming based on your worker setup
	scheduleID := fmt.Sprintf("schedule-%s-%s", config.WorkflowName, config.ScheduleID)
	workflowID := fmt.Sprintf("sw-%s-%s", config.WorkflowName, uuid.New().String())
	taskQueue := fmt.Sprintf("schedule-%s", config.WorkflowName) // Match your worker naming

	// Create the schedule
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
			Workflow:  workflowFunc,
			TaskQueue: taskQueue, // Use the specific queue for this workflow
			Args:      []interface{}{config.Args},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}

	// Verify schedule creation (optional debugging)
	_, err = scheduleHandle.Describe(ctx)
	if err != nil {
		return fmt.Errorf("failed to verify schedule creation: %w", err)
	}

	return nil
}
