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

package e2etasksteps

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/blox/blox/cluster-state-service/internal/features/wrappers"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/models"
	. "github.com/gucumber/gucumber"
)

type Exception struct {
	exceptionType string
	exceptionMsg  string
}

var (
	// Lists to memorize results required for the subsequent steps
	ecsTaskList   = []ecs.Task{}
	cssTaskList   = []models.Task{}
	exceptionList = []Exception{}

	taskDefnARN = ""
)

func init() {

	ecsWrapper := wrappers.NewECSWrapper()

	BeforeAll(func() {
		var err error
		taskDefnARN, err = ecsWrapper.RegisterSleep360TaskDefinition()
		if err != nil {
			T.Errorf(err.Error())
		}
	})

	Before("@task|@stream-instances", func() {
		clusterName, err := wrappers.GetClusterName()
		if err != nil {
			T.Errorf(err.Error())
		}
		err = stopAllTasks(ecsWrapper, clusterName)
		if err != nil {
			T.Errorf("Failed to stop all ECS tasks. Error: %s", err)
		}
	})

	AfterAll(func() {
		clusterName, err := wrappers.GetClusterName()
		if err != nil {
			T.Errorf(err.Error())
		}
		err = stopAllTasks(ecsWrapper, clusterName)
		if err != nil {
			T.Errorf("Failed to stop all ECS tasks. Error: %s", err)
		}
		err = ecsWrapper.DeregisterTaskDefinition(taskDefnARN)
		if err != nil {
			T.Errorf("Failed to deregister task definition. Error: %s", err)
		}
	})

	Given(`^I start (\d+) task(?:|s) in the ECS cluster$`, func(numTasks int) {
		ecsTaskList = nil
		cssTaskList = nil

		clusterName, err := wrappers.GetClusterName()
		if err != nil {
			T.Errorf(err.Error())
		}

		for i := 0; i < numTasks; i++ {
			ecsTask, err := ecsWrapper.StartTask(clusterName, taskDefinitionSleep300)
			if err != nil {
				T.Errorf(err.Error())
			}
			ecsTaskList = append(ecsTaskList, ecsTask)
		}
	})

	Then(`^I get a (.+?) task exception$`, func(exceptionType string) {
		if len(exceptionList) != 1 {
			T.Errorf("Error memorizing exception. ")
		}
		if exceptionType != exceptionList[0].exceptionType {
			T.Errorf("Expected exception '%s' but got '%s'. ", exceptionType, exceptionList[0].exceptionType)
		}
	})

	And(`^the task exception message contains "(.+?)"$`, func(exceptionMsg string) {
		if len(exceptionList) != 1 {
			T.Errorf("Error memorizing exception. ")
		}
		if !strings.Contains(exceptionList[0].exceptionMsg, exceptionMsg) {
			T.Errorf("Expected exception message returned '%s' to contain '%s'. ", exceptionList[0].exceptionMsg, exceptionMsg)
		}
	})
}

func stopAllTasks(ecsWrapper wrappers.ECSWrapper, clusterName string) error {
	taskARNList, err := ecsWrapper.ListTasks(clusterName)
	if err != nil {
		return err
	}
	for _, t := range taskARNList {
		err = ecsWrapper.StopTask(clusterName, *t)
		if err != nil {
			return err
		}
	}
	return nil
}
