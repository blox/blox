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

package main

import (
	"fmt"

	"github.com/aws/amazon-ecs-event-stream-handler/cmd"
	"github.com/aws/amazon-ecs-event-stream-handler/config"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/run"
	"github.com/aws/amazon-ecs-event-stream-handler/logger"

	log "github.com/cihub/seelog"

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
		log.Criticalf("Error executing: %v", err)
		os.Exit(errorCode)
	}
	if err := run.StartEventStreamHandler(config.SQSQueueName, config.EtcdEndpoints); err != nil {
		log.Criticalf("Error starting event stream handler: %v", err)
		os.Exit(errorCode)
	}
}
