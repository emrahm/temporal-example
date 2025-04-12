package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/emrahm/temporal-example/internal/workflow/temporal"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/client"
)

func init() {
	rootCmd.AddCommand(workerCmd)
}

var workerCmd = &cobra.Command{
	Use:   "run-worker",
	Short: "Run Temporal workers",
	Run: func(cmd *cobra.Command, args []string) {
		serviceClient, err := client.NewLazyClient(client.Options{})
		if err != nil {
			log.Fatalln("Unable to create Temporal client:", err)
		}
		defer serviceClient.Close()

		temporal.RunWorker(serviceClient)

		// Keep process alive to listen for system signals
		log.Println("Worker(s) running. Press Ctrl+C to exit...")
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down worker...")
	},
}
