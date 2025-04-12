package delivery

import (
	"context"

	"go.temporal.io/sdk/activity"
)

func (w Workflow) ChargeStripe(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity ChargeStripe started")
	return "input: " + name, nil
}
