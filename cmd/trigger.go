package cmd

import (
	"fmt"
	"log"

	"github.com/emrahm/temporal-example/internal/workflow/temporal"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/client"
)

func init() {
	rootCmd.AddCommand(triggerCmd)
}

var triggerCmd = &cobra.Command{
	Use:   "run [workflow_name] [args...]",
	Short: "Trigger a Temporal workflow",
	Args:  cobra.MinimumNArgs(1), // Require at least workflow_name
	Run: func(cmd *cobra.Command, args []string) {
		workflowName := args[0]
		var workflowArgs []string
		for _, arg := range args[1:] {
			workflowArgs = append(workflowArgs, arg)
		}

		serviceClient, err := client.NewLazyClient(client.Options{})
		if err != nil {
			log.Fatalln("Unable to create client", err)
		}
		defer serviceClient.Close()

		worker := temporal.NewWorkflow(serviceClient)

		fmt.Printf("Triggering workflow %s with args: %v\n", workflowName, workflowArgs)
		err = worker.Execute(workflowName, workflowArgs)
		if err != nil {
			log.Fatalln("Unable to schedule", err)
		}
		fmt.Printf("Workflow %s triggered successfully\n", workflowName)
	},
}
