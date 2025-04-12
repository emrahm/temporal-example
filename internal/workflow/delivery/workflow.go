package delivery

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

type Workflow struct {
}

func NewWorkflow() Workflow {
	return Workflow{}
}

func (w Workflow) Workflow(ctx workflow.Context, args []string) (interface{}, error) {

	logger := workflow.GetLogger(ctx)
	logger.Info("Workflow started")

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)
	var result string

	name := args[0]
	logger.Info(fmt.Sprintf("Workflow name: %v", name))

	err := workflow.ExecuteActivity(ctx, w.ChargeStripe, name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity Failed", err)
	}
	return result, nil
}
