// Copyright 2016-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package main

import (
	"fmt"
	"os"

	"github.com/blox/blox/cluster-state-service/canary/pkg/cmd"
	"github.com/blox/blox/cluster-state-service/canary/pkg/logger"
	"github.com/blox/blox/cluster-state-service/canary/pkg/tests"
	"github.com/blox/blox/cluster-state-service/canary/pkg/wrappers"
	log "github.com/cihub/seelog"
)

const errorCode = 1

// TODO
// * Parallelize tests
// * Emit metric for time taken to run each test
func main() {
	defer log.Flush()
	err := logger.InitLogger()
	if err != nil {
		fmt.Printf("Could not initialize logger: %+v", err)
	}

	if err = cmd.RootCmd.Execute(); err != nil {
		log.Criticalf("Error getting command line arguments: %+v", err)
		os.Exit(errorCode)
	}

	awsSession, err := wrappers.NewAWSSession()
	if err != nil {
		log.Criticalf(err.Error())
		os.Exit(errorCode)
	}

	canary, err := tests.NewCanary(awsSession, cmd.ClusterStateServiceEndpoint)
	if err != nil {
		log.Criticalf(err.Error())
		os.Exit(errorCode)
	}

	canaryRunner, err := tests.NewCanaryRunner(awsSession)
	if err != nil {
		log.Criticalf(err.Error())
		os.Exit(errorCode)
	}

	canaryRunner.Run(canary.GetInstance, canary.GetInstanceMetric)
	canaryRunner.Run(canary.ListInstances, canary.ListInstancesMetric)
	canaryRunner.Run(canary.GetTask, canary.GetTaskMetric)
	canaryRunner.Run(canary.ListTasks, canary.ListTasksMetric)
}
