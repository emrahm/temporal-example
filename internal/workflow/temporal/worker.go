package temporal

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func RunWorker(serviceClient client.Client) {

	workflow := NewWorkflow(serviceClient)
	workers := workflow.Register()

	for _, w := range workers {
		go func(w worker.Worker) {
			if err := w.Run(worker.InterruptCh()); err != nil {
				log.Fatalln("Unable to start worker:", err)
			}
		}(w)
	}
}
