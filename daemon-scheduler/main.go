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

package main

import (
	"fmt"

	"github.com/blox/blox/daemon-scheduler/logger"
	"github.com/blox/blox/daemon-scheduler/pkg/cmd"
	"github.com/blox/blox/daemon-scheduler/pkg/config"
	"github.com/blox/blox/daemon-scheduler/pkg/scheduler"
	log "github.com/cihub/seelog"

	"os"
)

func main() {
	defer log.Flush()
	err := logger.InitLogger()
	if err != nil {
		fmt.Printf("Could not initialize logger: %+v", err)
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		log.Criticalf("Error getting command line arguments: %v", err)
		os.Exit(1)
	}

	if err := scheduler.Run(config.SchedulerBindAddr, config.ClusterStateServiceEndpoint); err != nil {
		log.Criticalf("Error running scheduler: %v", err)
		os.Exit(1)
	}
}
