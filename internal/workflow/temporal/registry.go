package temporal

import (
	"fmt"

	"github.com/emrahm/temporal-example/internal/workflow/delivery"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type ScheduleWorkflow func(ctx workflow.Context, args []string) (interface{}, error)

type Workflow struct {
	serviceClient client.Client
	workflows     map[string]ScheduleWorkflow
	activities    map[string]interface{}
	workers       map[string]worker.Options
}

func NewWorkflow(serviceClient client.Client) *Workflow {
	dw := delivery.NewWorkflow()

	var workflows = map[string]ScheduleWorkflow{
		"delivery": dw.Workflow,
	}

	var activities = map[string]interface{}{
		"delivery": dw.ChargeStripe,
	}

	var workers = map[string]worker.Options{
		"delivery": {},
	}

	return &Workflow{
		serviceClient: serviceClient,
		workflows:     workflows,
		activities:    activities,
		workers:       workers,
	}
}

func (w *Workflow) Register() map[string]worker.Worker {
	queueWorkers := make(map[string]worker.Worker, 0)
	for name, workflow := range w.workflows {
		if _, ok := w.workers[name]; ok {
			queueName := fmt.Sprintf("queue-%s", name)
			scheduleName := fmt.Sprintf("schedule-%s", name)
			mainWorker := worker.New(w.serviceClient, queueName, w.workers[name])
			scheduleWorker := worker.New(w.serviceClient, scheduleName, w.workers[name])
			mainWorker.RegisterWorkflow(workflow)
			scheduleWorker.RegisterWorkflow(workflow)
			queueWorkers[queueName] = mainWorker
			queueWorkers[scheduleName] = scheduleWorker
		} else {
			panic("worker not found")
		}
	}

	for name, activity := range w.activities {
		if _, ok := queueWorkers[name]; ok {
			queueName := fmt.Sprintf("queue-%s", name)
			scheduleName := fmt.Sprintf("schedule-%s", name)
			queueWorkers[queueName].RegisterActivity(activity)
			queueWorkers[scheduleName].RegisterActivity(activity)
		} else {
			panic("worker not found")
		}
	}
	return queueWorkers
}
