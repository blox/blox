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

package e2etasksteps

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/models"
	"github.com/pkg/errors"
)

func ValidateTasksMatch(ecsTask ecs.Task, cssTask models.Task) error {
	if *ecsTask.TaskArn != *cssTask.Entity.TaskARN ||
		*ecsTask.ClusterArn != *cssTask.Entity.ClusterARN ||
		*ecsTask.ContainerInstanceArn != *cssTask.Entity.ContainerInstanceARN {
		return errors.New("Tasks don't match.")
	}
	return nil
}

func ValidateListContainsTask(ecsTask ecs.Task, cssTaskList []models.Task) error {
	taskARN := *ecsTask.TaskArn
	var cssTask models.Task
	for _, t := range cssTaskList {
		if *t.Entity.TaskARN == taskARN {
			cssTask = t
			break
		}
	}
	if cssTask.Entity.TaskARN == nil {
		return errors.Errorf("Task with ARN '%s' not found in response. ", taskARN)
	}
	return ValidateTasksMatch(ecsTask, cssTask)
}

func ValidateListContainsTaskWithDesiredStatus(ecsTask ecs.Task, cssTaskList []models.Task, desiredStatus string) error {
	taskARN := *ecsTask.TaskArn
	var cssTask models.Task
	for _, t := range cssTaskList {
		if *t.Entity.TaskARN == taskARN && *t.Entity.DesiredStatus == desiredStatus {
			cssTask = t
			break
		}
	}
	if cssTask.Entity == nil || cssTask.Entity.TaskARN == nil {
		return errors.Errorf("Task with ARN '%s' and desired status '%s' not found in response. ", taskARN, desiredStatus)
	}
	return ValidateTasksMatch(ecsTask, cssTask)
}
