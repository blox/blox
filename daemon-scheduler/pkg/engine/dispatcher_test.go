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

package engine

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/blox/blox/daemon-scheduler/pkg/clients/css/models"
	mocks "github.com/blox/blox/daemon-scheduler/pkg/mocks"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	clusterArn      = "cluster-arn"
	environmentName = "env-name"
)

type DispatcherTestSuite struct {
	suite.Suite
	environmentSvc   *mocks.MockEnvironment
	deploymentSvc    *mocks.MockDeployment
	css              *mocks.MockClusterState
	ecs              *mocks.MockECS
	deploymentWorker *mocks.MockDeploymentWorker
}

func (suite *DispatcherTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	suite.environmentSvc = mocks.NewMockEnvironment(mockCtrl)
	suite.deploymentSvc = mocks.NewMockDeployment(mockCtrl)
	suite.css = mocks.NewMockClusterState(mockCtrl)
	suite.ecs = mocks.NewMockECS(mockCtrl)
	suite.deploymentWorker = mocks.NewMockDeploymentWorker(mockCtrl)
}

func TestDispatcherTestSuite(t *testing.T) {
	suite.Run(t, new(DispatcherTestSuite))
}

func (suite *DispatcherTestSuite) TestUnknownEvent() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	input := make(chan Event)
	output := make(chan Event)
	dispatcher := NewDispatcher(ctx,
		suite.environmentSvc,
		suite.deploymentSvc,
		suite.ecs, suite.css,
		suite.deploymentWorker,
		input, output,
	)

	dispatcher.Start()
	input <- ErrorEvent{
		Error: errors.New("Unexpected Error"),
	}

	ticker := time.NewTicker(1 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			return
		case <-output:
			assert.Fail(suite.T(), "Received unexpected event from dispatcher")
			return
		}
	}
}

func (suite *DispatcherTestSuite) TestUpdateInProgressDeploymentEventReturnsError() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	input := make(chan Event)
	output := make(chan Event)
	dispatcher := NewDispatcher(ctx,
		suite.environmentSvc,
		suite.deploymentSvc,
		suite.ecs, suite.css,
		suite.deploymentWorker,
		input, output,
	)

	environment := types.Environment{
		Name:    environmentName,
		Cluster: clusterArn,
	}

	event := UpdateInProgressDeploymentEvent{
		Environment: environment,
	}

	err := errors.New("Error calling UpdateInProgressDeployment")
	suite.deploymentWorker.EXPECT().
		UpdateInProgressDeployment(ctx, event.Environment.Name).
		Return(nil, err).
		Times(1)

	dispatcher.Start()
	input <- event

	observedErr := errors.Cause((<-output).(ErrorEvent).Error)
	assert.Equal(suite.T(), err, observedErr)
}

func (suite *DispatcherTestSuite) TestUpdateInProgressDeploymentEvent() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	input := make(chan Event)
	output := make(chan Event)
	dispatcher := NewDispatcher(ctx,
		suite.environmentSvc,
		suite.deploymentSvc,
		suite.ecs, suite.css,
		suite.deploymentWorker,
		input, output,
	)

	environment := types.Environment{
		Name:    environmentName,
		Cluster: clusterArn,
	}

	event := UpdateInProgressDeploymentEvent{
		Environment: environment,
	}
	suite.deploymentWorker.EXPECT().
		UpdateInProgressDeployment(ctx, event.Environment.Name).
		Return(nil, nil).
		Times(1)

	dispatcher.Start()
	input <- event
}

func (suite *DispatcherTestSuite) TestStartDeploymentEventReturnsError() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	input := make(chan Event)
	output := make(chan Event)
	dispatcher := NewDispatcher(ctx,
		suite.environmentSvc,
		suite.deploymentSvc,
		suite.ecs, suite.css,
		suite.deploymentWorker,
		input, output,
	)

	environment := types.Environment{
		Name:    environmentName,
		Cluster: clusterArn,
	}
	instances := []*string{
		aws.String("instance-arn-1"),
		aws.String("instance-arn-2"),
	}
	event := StartDeploymentEvent{
		Environment: environment,
		Instances:   instances,
	}

	err := errors.New("Error creating sub-deployment")
	suite.deploymentSvc.EXPECT().
		CreateSubDeployment(ctx, event.Environment.Name, event.Instances).
		Return(nil, err)

	dispatcher.Start()
	input <- event

	observedErr := errors.Cause((<-output).(ErrorEvent).Error)
	assert.Equal(suite.T(), err, observedErr)
}

func (suite *DispatcherTestSuite) TestStartDeploymentEvent() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	input := make(chan Event)
	output := make(chan Event)
	dispatcher := NewDispatcher(ctx,
		suite.environmentSvc,
		suite.deploymentSvc,
		suite.ecs, suite.css,
		suite.deploymentWorker,
		input, output,
	)

	environment := types.Environment{
		Name:    environmentName,
		Cluster: clusterArn,
	}
	instances := []*string{
		aws.String("instance-arn-1"),
		aws.String("instance-arn-2"),
	}
	event := StartDeploymentEvent{
		Environment: environment,
		Instances:   instances,
	}

	deployment := types.Deployment{
		ID: uuid.NewRandom().String(),
	}
	suite.deploymentSvc.EXPECT().
		CreateSubDeployment(ctx, event.Environment.Name, event.Instances).
		Return(&deployment, nil).
		Times(1)

	dispatcher.Start()
	input <- event

	deploymentResult := (<-output).(StartDeploymentResult).Deployment
	assert.Equal(suite.T(), deployment.ID, deploymentResult.ID)
}

func (suite *DispatcherTestSuite) TestStopTasksEventListTasksReturnsError() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	input := make(chan Event)
	output := make(chan Event)
	dispatcher := NewDispatcher(ctx,
		suite.environmentSvc,
		suite.deploymentSvc,
		suite.ecs, suite.css,
		suite.deploymentWorker,
		input, output,
	)

	tasksToStop := []string{
		"task-arn-1",
		"task-arn-2",
		"unknown-task-arn-1",
	}
	event := StopTasksEvent{
		Cluster: "cluster-arn",
		Tasks:   tasksToStop,
	}

	err := errors.New("Error from css.ListTasks")
	suite.css.EXPECT().ListTasks(event.Cluster).Return(nil, err).Times(1)
	suite.ecs.EXPECT().StopTask(gomock.Any(), gomock.Any()).Times(0)

	dispatcher.Start()
	input <- event

	observedErr := errors.Cause((<-output).(ErrorEvent).Error)
	assert.Equal(suite.T(), err, observedErr)
}

func (suite *DispatcherTestSuite) TestStopTasksEventECSStopTaskReturnsError() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	input := make(chan Event)
	output := make(chan Event)
	dispatcher := NewDispatcher(ctx,
		suite.environmentSvc,
		suite.deploymentSvc,
		suite.ecs, suite.css,
		suite.deploymentWorker,
		input, output,
	)

	tasksToStop := []string{
		"task-arn-1",
	}
	event := StopTasksEvent{
		Cluster: "cluster-arn",
		Tasks:   tasksToStop,
	}

	tasksFromECS := []*models.Task{
		&models.Task{
			TaskARN:       aws.String("task-arn-1"),
			ClusterARN:    aws.String(event.Cluster),
			DesiredStatus: aws.String("RUNNING"),
		},
	}
	suite.css.EXPECT().ListTasks(event.Cluster).Return(tasksFromECS, nil).Times(1)

	err := errors.New("Error stopping task")
	suite.ecs.EXPECT().StopTask(event.Cluster, "task-arn-1").Return(err).Times(1)

	dispatcher.Start()
	input <- event

	stoppedTasks := (<-output).(StopTasksResult).StoppedTasks
	assert.Equal(suite.T(), []string{}, stoppedTasks)
}

func (suite *DispatcherTestSuite) TestStopTasksEvent() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	input := make(chan Event)
	output := make(chan Event)
	dispatcher := NewDispatcher(ctx,
		suite.environmentSvc,
		suite.deploymentSvc,
		suite.ecs, suite.css,
		suite.deploymentWorker,
		input, output,
	)

	tasksToStop := []string{
		"task-arn-1",
		"task-arn-2",
		"unknown-task-arn-1",
	}
	event := StopTasksEvent{
		Cluster: "cluster-arn",
		Tasks:   tasksToStop,
	}

	tasksFromECS := []*models.Task{
		&models.Task{
			TaskARN:       aws.String("task-arn-1"),
			ClusterARN:    aws.String(event.Cluster),
			DesiredStatus: aws.String("RUNNING"),
		},
		&models.Task{
			TaskARN:       aws.String("task-arn-2"),
			ClusterARN:    aws.String(event.Cluster),
			DesiredStatus: aws.String("STOPPED"),
		},
		&models.Task{
			TaskARN:       aws.String("task-arn-3"),
			ClusterARN:    aws.String(event.Cluster),
			DesiredStatus: aws.String("RUNNING"),
		},
	}
	suite.css.EXPECT().ListTasks(event.Cluster).Return(tasksFromECS, nil).Times(1)
	suite.ecs.EXPECT().StopTask(event.Cluster, "task-arn-1").Return(nil).Times(1)
	suite.ecs.EXPECT().StopTask(event.Cluster, "task-arn-2").Times(0)
	suite.ecs.EXPECT().StopTask(event.Cluster, "unknown-task-arn-1").Times(0)
	suite.ecs.EXPECT().StopTask(event.Cluster, "task-arn-3").Times(0)

	dispatcher.Start()
	input <- event

	stoppedTasks := (<-output).(StopTasksResult).StoppedTasks
	assert.Equal(suite.T(), []string{"task-arn-1", "task-arn-2"}, stoppedTasks)
}
