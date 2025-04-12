package temporal

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

func (w *Workflow) Execute(workflowName string, args []string) error {

	if workflow, ok := w.workflows[workflowName]; ok {
		workflowOptions := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("%s-%s", workflowName, uuid.New().String()),
			TaskQueue: fmt.Sprintf("%s-%s", workflowName, "queue"),
		}

		workflowExecution, err := w.serviceClient.ExecuteWorkflow(context.Background(),
			workflowOptions,
			workflow,
			args)
		if err != nil {
			return fmt.Errorf("failed to execute workflow: %w", err)
		}
		var result string
		err = workflowExecution.Get(context.Background(), &result)
		if err != nil {
			return fmt.Errorf("failed to get workflow result: %w", err)
		}
		log.Println("Result: ", result)

	} else {
		return errors.New("workflow not found")
	}

	return nil

}
