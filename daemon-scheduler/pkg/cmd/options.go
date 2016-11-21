// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package cmd

import (
	"github.com/blox/blox/daemon-scheduler/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var etcdAddrList string

// RootCmd represents the base command when called without any subcommands
var RootCmd = createRootCommand()

// Init the CLI.
func init() {
	cobra.OnInitialize(initConfig)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}

func createRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ecs-daemon-scheduler",
		Short: "ecs-daemon-scheduler ",
		Long: `ecs-daemon-scheduler supports launching daemon like tasks that need to
be launched only once in every node of an ECS cluster. As new nodes join the cluster,
it launches the task in those nodes as well.`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	rootCmd.PersistentFlags().StringArrayVar(&config.EtcdEndpoints, "etcd-endpoint", make([]string, 0), "Etcd node addresses")
	rootCmd.PersistentFlags().StringVar(&config.SchedulerBindAddr, "bind", "", "Scheduler bind address")
	rootCmd.PersistentFlags().StringVar(&config.ClusterStateServiceEndpoint, "css-endpoint", "", "Cluster state service address")
	return rootCmd
}
