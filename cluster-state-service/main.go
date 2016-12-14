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

	"github.com/blox/blox/cluster-state-service/logger"
	log "github.com/cihub/seelog"

	"github.com/blox/blox/cluster-state-service/cmd"
	"github.com/blox/blox/cluster-state-service/config"
	"github.com/blox/blox/cluster-state-service/handler/run"
	"github.com/blox/blox/cluster-state-service/versioning"
	"os"
)

const errorCode = 1

func main() {
	defer log.Flush()
	err := logger.InitLogger()
	if err != nil {
		fmt.Printf("Could not initialize logger: %+v", err)
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Criticalf("Error executing: %+v", err)
		os.Exit(errorCode)
	}
	if config.PrintVersion {
		versioning.PrintVersion()
		os.Exit(0)
	}
	if err := run.StartClusterStateService(config.QueueNameURI, config.CSSBindAddr, config.EtcdEndpoints); err != nil {
		log.Criticalf("Error starting event stream handler: %+v", err)
		os.Exit(errorCode)
	}
}
