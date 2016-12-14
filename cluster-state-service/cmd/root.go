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
	"github.com/blox/blox/cluster-state-service/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	queueNameURIFlag = "queue"
	cssBindFlag      = "bind"
	etcdEndpointFlag = "etcd-endpoint"
	versionFlag      = "version"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd *cobra.Command

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd = createRootCommand()
}

func createRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		// TODO: Fix these messages
		Use:   "cluster-state-service",
		Short: "cluster-state-service consumes events from Amazon ECS and provides a local view of the cluster state",
		Long: `cluster-state-service processes EC2 Container Service events and  creates 
a localized data store, which provides you a near-real-time view of your cluster state.`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	// TODO: Fix the description
	rootCmd.PersistentFlags().StringVar(&config.QueueNameURI, queueNameURIFlag, "", "Queue name should be of the form sqs://name or kinesis://name")
	rootCmd.PersistentFlags().StringVar(&config.CSSBindAddr, cssBindFlag, "", "Cluster State Service listen address")
	rootCmd.PersistentFlags().StringArrayVar(&config.EtcdEndpoints, etcdEndpointFlag, make([]string, 0), "Etcd node addresses")
	rootCmd.PersistentFlags().BoolVar(&config.PrintVersion, versionFlag, false, "Print version and exit")
	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
