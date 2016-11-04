// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the License). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the license file accompanying this file. This file is distributed
// on an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

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
