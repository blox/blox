package cmd

import (
	"fmt"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/run"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const sqsQueueNameFlag = "queue"

var sqsQueueName string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	// TODO: Fix these messages
	Use:   "amazon-ecs-event-stream-handler",
	Short: "amazon-ecs-event-stream-handler handles amazon ecs event stream",
	Long: `amazon-ecs-event-stream-handler handles amazon ecs event stream. It
processes EC2 Container Service events and creates a localized data store, which
provides you a near-real-time view of your cluster state.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if sqsQueueName == "" {
			return fmt.Errorf("SQS queue name cannot be empty")
		}

		return run.StartEventStreamHandler(sqsQueueName)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	// TODO: Fix the description
	RootCmd.PersistentFlags().StringVar(&sqsQueueName, sqsQueueNameFlag, "", "SQS queue name")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
