package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/emrahm/temporal-example/internal/workflow/temporal"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/client"
)

var (
	intervalStr string
	scheduleID  string
)

func init() {
	// Register the command
	rootCmd.AddCommand(scheduleCmd)

	// Register flags
	scheduleCmd.Flags().StringVarP(&intervalStr, "interval", "i", "10s", "Workflow interval (e.g., 10s, 1m, 1h, 1d)")
}

var scheduleCmd = &cobra.Command{
	Use:   "schedule [workflow_name] [schedule_id] [args...] --interval 1s",
	Short: "Schedule a Temporal workflow",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		workflowName := args[0]
		scheduleID := args[1]
		workflowArgs := args[1:]

		// Parse interval string to time.Duration
		interval, err := time.ParseDuration(intervalStr)
		if err != nil {
			log.Fatalf("Invalid interval: %v", err)
		}

		// Initialize Temporal client
		serviceClient, err := client.NewLazyClient(client.Options{})
		if err != nil {
			log.Fatalln("Unable to create Temporal client:", err)
		}
		defer serviceClient.Close()

		// Create schedule config from flags
		scheduleConfig := temporal.ScheduleConfig{
			WorkflowName: workflowName,
			Args:         workflowArgs,
			Interval:     interval,
			ScheduleID:   scheduleID,
		}

		worker := temporal.NewWorkflow(serviceClient)

		fmt.Printf("Triggering workflow %s with args: %v and interval: %v\n", workflowName, workflowArgs, interval)
		err = worker.Schedule(scheduleConfig)
		if err != nil {
			log.Fatalln("Unable to schedule workflow:", err)
		}
	},
}
