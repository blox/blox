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

package engine

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/models"
	"github.com/blox/blox/daemon-scheduler/pkg/deployment/types"
	environmenttypes "github.com/blox/blox/daemon-scheduler/pkg/environment/types"
	"github.com/blox/blox/daemon-scheduler/pkg/facade"
	mocks "github.com/blox/blox/daemon-scheduler/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SchedulerTestSuite struct {
	suite.Suite
	environmentService *mocks.MockEnvironmentService
	deploymentService  *mocks.MockDeploymentService
	css                *facade.MockClusterState
	ecs                *mocks.MockECS
}

func (suite *SchedulerTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	suite.environmentService = mocks.NewMockEnvironmentService(mockCtrl)
	suite.deploymentService = mocks.NewMockDeploymentService(mockCtrl)
	suite.css = facade.NewMockClusterState(mockCtrl)
	suite.ecs = mocks.NewMockECS(mockCtrl)
}

func TestSchedulerTestSuite(t *testing.T) {
	suite.Run(t, new(SchedulerTestSuite))
}

func (suite *SchedulerTestSuite) TestRunInProgress() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*schedulerTickerDuration)
	defer cancel()
	events := make(chan Event)
	scheduler := NewScheduler(ctx, events, suite.environmentService, suite.deploymentService, suite.css, suite.ecs)
	scheduler.setInProgress(true)
	scheduler.Start()
	schedulerErrorEvent := (<-events).(SchedulerErrorEvent)
	assert.Error(suite.T(), schedulerErrorEvent.Error, "Expected error due to in-progress scheduler run")
}

func (suite *SchedulerTestSuite) TestRunListEnvironmentsReturnsError() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*schedulerTickerDuration)
	defer cancel()
	var err error
	err = errors.New("Error calling ListEnvironments")
	suite.environmentService.EXPECT().ListEnvironments(ctx).Return(nil, err)
	events := make(chan Event)
	scheduler := NewScheduler(ctx, events, suite.environmentService, suite.deploymentService, suite.css, suite.ecs)
	scheduler.Start()
	schedulerErrorEvent := (<-events).(SchedulerErrorEvent)
	assert.Equal(suite.T(), err, errors.Cause(schedulerErrorEvent.Error))

	//next run of scheduler should occur after ticker and do the same thing
	err = errors.New("Error calling ListEnvironments")
	suite.environmentService.EXPECT().ListEnvironments(ctx).Return(nil, err)
	schedulerErrorEvent = (<-events).(SchedulerErrorEvent)
	assert.Equal(suite.T(), err, errors.Cause(schedulerErrorEvent.Error))
}

func (suite *SchedulerTestSuite) TestRunGetCurrentDeploymentReturnsError() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*schedulerTickerDuration)
	defer cancel()
	err := errors.New("Error calling GetCurrentDeployment")
	environment := environmenttypes.Environment{
		Name: "TestRunGetCurrentDeploymentReturnsError",
	}
	environments := []environmenttypes.Environment{environment}
	suite.environmentService.EXPECT().ListEnvironments(ctx).Return(environments, nil)
	suite.deploymentService.EXPECT().GetCurrentDeployment(ctx, environment.Name).Return(nil, err)
	events := make(chan Event)
	scheduler := NewScheduler(ctx, events, suite.environmentService, suite.deploymentService, suite.css, suite.ecs)
	scheduler.Start()
	schedulerErrorEvent := (<-events).(SchedulerErrorEvent)
	assert.Equal(suite.T(), err, errors.Cause(schedulerErrorEvent.Error))
}

func (suite *SchedulerTestSuite) TestRunGetCurrentDeploymentReturnsNil() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*schedulerTickerDuration)
	defer cancel()
	environment := environmenttypes.Environment{
		Name: "TestRunGetCurrentDeploymentReturnsNil",
	}
	environments := []environmenttypes.Environment{environment}
	suite.environmentService.EXPECT().ListEnvironments(ctx).Return(environments, nil)
	suite.deploymentService.EXPECT().GetCurrentDeployment(ctx, environment.Name).Return(nil, nil)
	events := make(chan Event)
	scheduler := NewScheduler(ctx, events, suite.environmentService, suite.deploymentService, suite.css, suite.ecs)
	scheduler.Start()
	schedulerEnvironmentEvent := (<-events).(SchedulerEnvironmentEvent)
	assert.Equal(suite.T(), environment.Name, schedulerEnvironmentEvent.Environment.Name)
}

func (suite *SchedulerTestSuite) TestRunListInstancesReturnsError() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*schedulerTickerDuration)
	defer cancel()
	environment := environmenttypes.Environment{
		Name:    "TestRunListInstancesReturnsError",
		Cluster: "testCluster",
	}
	environments := []environmenttypes.Environment{environment}
	currentDeployment := types.Deployment{
		ID:     "dep-id",
		Status: types.DeploymentInProgress,
		Health: types.DeploymentHealthy,
	}
	suite.environmentService.EXPECT().ListEnvironments(ctx).Return(environments, nil)
	suite.deploymentService.EXPECT().GetCurrentDeployment(ctx, environment.Name).Return(&currentDeployment, nil)

	err := errors.New("Error getting instances from css")
	suite.css.EXPECT().ListInstances(environment.Cluster).Return(nil, err)

	events := make(chan Event)
	scheduler := NewScheduler(ctx, events, suite.environmentService, suite.deploymentService, suite.css, suite.ecs)
	scheduler.Start()
	schedulerErrorEvent := (<-events).(SchedulerErrorEvent)
	assert.Equal(suite.T(), err, errors.Cause(schedulerErrorEvent.Error))
}

func (suite *SchedulerTestSuite) TestRunListTasksReturnsError() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*schedulerTickerDuration)
	defer cancel()

	environment := environmenttypes.Environment{
		Name:    "TestRunListTasksReturnsError",
		Cluster: "testCluster",
	}
	environments := []environmenttypes.Environment{environment}
	suite.environmentService.EXPECT().ListEnvironments(ctx).Return(environments, nil)

	currentDeployment := types.Deployment{
		ID:     "dep-id",
		Status: types.DeploymentInProgress,
		Health: types.DeploymentHealthy,
	}
	suite.deploymentService.EXPECT().GetCurrentDeployment(ctx, environment.Name).Return(&currentDeployment, nil)

	instance := &models.ContainerInstance{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.ContainerInstanceDetail{
			ClusterARN:           aws.String(environment.Cluster),
			ContainerInstanceARN: aws.String("instance-arn"),
			Status:               aws.String("ACTIVE"),
		},
	}
	instances := []*models.ContainerInstance{instance}
	suite.css.EXPECT().ListInstances(environment.Cluster).Return(instances, nil)

	err := errors.New("Error getting tasks from css")
	suite.css.EXPECT().ListTasks(environment.Cluster).Return(nil, err)

	events := make(chan Event)
	scheduler := NewScheduler(ctx, events, suite.environmentService, suite.deploymentService, suite.css, suite.ecs)
	scheduler.Start()
	schedulerErrorEvent := (<-events).(SchedulerErrorEvent)
	assert.Equal(suite.T(), err, errors.Cause(schedulerErrorEvent.Error))
}

func (suite *SchedulerTestSuite) TestRunListDeploymentsReturnsError() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*schedulerTickerDuration)
	defer cancel()

	environment := environmenttypes.Environment{
		Name:    "TestRunListDeploymentsReturnsError",
		Cluster: "testCluster",
	}
	environments := []environmenttypes.Environment{environment}
	suite.environmentService.EXPECT().ListEnvironments(ctx).Return(environments, nil)

	currentDeployment := types.Deployment{
		ID:     "dep-id",
		Status: types.DeploymentInProgress,
		Health: types.DeploymentHealthy,
	}
	suite.deploymentService.EXPECT().GetCurrentDeployment(ctx, environment.Name).Return(&currentDeployment, nil)

	instance := &models.ContainerInstance{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.ContainerInstanceDetail{
			ClusterARN:           aws.String(environment.Cluster),
			ContainerInstanceARN: aws.String("instance-arn"),
			Status:               aws.String("ACTIVE"),
		},
	}
	instances := []*models.ContainerInstance{instance}
	suite.css.EXPECT().ListInstances(environment.Cluster).Return(instances, nil)

	task := &models.Task{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.TaskDetail{
			ClusterARN:           instance.Entity.ClusterARN,
			ContainerInstanceARN: instance.Entity.ContainerInstanceARN,
			TaskARN:              aws.String("task-arn"),
		},
	}
	tasks := []*models.Task{task}
	suite.css.EXPECT().ListTasks(environment.Cluster).Return(tasks, nil)

	err := errors.New("Error getting deployments for environment")
	suite.deploymentService.EXPECT().ListDeploymentsSortedReverseChronologically(ctx, environment.Name).Return(nil, err)

	events := make(chan Event)
	scheduler := NewScheduler(ctx, events, suite.environmentService, suite.deploymentService, suite.css, suite.ecs)
	scheduler.Start()
	schedulerErrorEvent := (<-events).(SchedulerErrorEvent)
	assert.Equal(suite.T(), err, errors.Cause(schedulerErrorEvent.Error))
}

func (suite *SchedulerTestSuite) TestRunAllInstancesDeployed() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*schedulerTickerDuration)
	defer cancel()

	environment := environmenttypes.Environment{
		Name:    "TestRunAllInstancesDeployed",
		Cluster: "testCluster",
	}
	environments := []environmenttypes.Environment{environment}
	suite.environmentService.EXPECT().ListEnvironments(ctx).Return(environments, nil)

	currentDeployment := types.Deployment{
		ID:     "dep-id",
		Status: types.DeploymentInProgress,
		Health: types.DeploymentHealthy,
	}
	suite.deploymentService.EXPECT().GetCurrentDeployment(ctx, environment.Name).Return(&currentDeployment, nil)

	instance1 := &models.ContainerInstance{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.ContainerInstanceDetail{
			ClusterARN:           aws.String(environment.Cluster),
			ContainerInstanceARN: aws.String("instance-arn-1"),
			Status:               aws.String("ACTIVE"),
		},
	}
	instance2 := &models.ContainerInstance{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.ContainerInstanceDetail{
			ClusterARN:           aws.String(environment.Cluster),
			ContainerInstanceARN: aws.String("instance-arn-2"),
			Status:               aws.String("ACTIVE"),
		},
	}
	instances := []*models.ContainerInstance{instance1, instance2}
	suite.css.EXPECT().ListInstances(environment.Cluster).Return(instances, nil)

	task1 := &models.Task{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.TaskDetail{
			ClusterARN:           instance1.Entity.ClusterARN,
			ContainerInstanceARN: instance1.Entity.ContainerInstanceARN,
			TaskARN:              aws.String("task-arn-1"),
			StartedBy:            currentDeployment.ID,
			DesiredStatus:        aws.String(runningTaskStatus),
		},
	}
	task2 := &models.Task{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.TaskDetail{
			ClusterARN:           instance2.Entity.ClusterARN,
			ContainerInstanceARN: instance2.Entity.ContainerInstanceARN,
			TaskARN:              aws.String("task-arn-2"),
			StartedBy:            currentDeployment.ID,
			DesiredStatus:        aws.String(runningTaskStatus),
		},
	}
	tasks := []*models.Task{task1, task2}
	suite.css.EXPECT().ListTasks(environment.Cluster).Return(tasks, nil)

	deployments := []types.Deployment{currentDeployment}
	suite.deploymentService.EXPECT().ListDeploymentsSortedReverseChronologically(ctx, environment.Name).Return(deployments, nil)

	events := make(chan Event)
	scheduler := NewScheduler(ctx, events, suite.environmentService, suite.deploymentService, suite.css, suite.ecs)
	scheduler.Start()

	schedulerEnvironmentEvent := (<-events).(SchedulerEnvironmentEvent)
	assert.Equal(suite.T(), environment.Name, schedulerEnvironmentEvent.Environment.Name)
}

func (suite *SchedulerTestSuite) TestRunEnvironmentStateInProgress() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*schedulerTickerDuration)
	defer cancel()

	environment := environmenttypes.Environment{
		Name:    "TestRunEnvironmentStateInProgress",
		Cluster: "testCluster",
	}
	environments := []environmenttypes.Environment{environment}
	suite.environmentService.EXPECT().ListEnvironments(ctx).Return(environments, nil)

	events := make(chan Event)
	scheduler := NewScheduler(ctx, events, suite.environmentService, suite.deploymentService, suite.css, suite.ecs)
	trackingInfo := make(map[string]time.Time)
	previousState := make(map[string]environmentExecutionState)
	previousState[environment.Name] = environmentExecutionState{
		environment:  environment,
		trackingInfo: trackingInfo,
		inProgress:   true,
	}
	scheduler.setExecutionState(previousState)

	scheduler.Start()

	schedulerEnvironmentEvent := (<-events).(SchedulerEnvironmentEvent)
	assert.Equal(suite.T(), environment.Name, schedulerEnvironmentEvent.Environment.Name)
}

func (suite *SchedulerTestSuite) TestRunNewInstance() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*schedulerTickerDuration)
	defer cancel()

	environment := environmenttypes.Environment{
		Name:    "TestRunNewInstance",
		Cluster: "testCluster",
	}
	environments := []environmenttypes.Environment{environment}
	suite.environmentService.EXPECT().ListEnvironments(ctx).Return(environments, nil)

	currentDeployment := types.Deployment{
		ID:     "dep-id",
		Status: types.DeploymentInProgress,
		Health: types.DeploymentHealthy,
	}
	suite.deploymentService.EXPECT().GetCurrentDeployment(ctx, environment.Name).Return(&currentDeployment, nil)

	instance1 := &models.ContainerInstance{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.ContainerInstanceDetail{
			ClusterARN:           aws.String(environment.Cluster),
			ContainerInstanceARN: aws.String("instance-arn-1"),
			Status:               aws.String("ACTIVE"),
		},
	}
	instance2 := &models.ContainerInstance{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.ContainerInstanceDetail{
			ClusterARN:           aws.String(environment.Cluster),
			ContainerInstanceARN: aws.String("instance-arn-2"),
			Status:               aws.String("INACTIVE"),
		},
	}
	newInstance := &models.ContainerInstance{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.ContainerInstanceDetail{
			ClusterARN:           aws.String(environment.Cluster),
			ContainerInstanceARN: aws.String("instance-arn-new"),
			Status:               aws.String("ACTIVE"),
		},
	}
	instances := []*models.ContainerInstance{instance1, instance2, newInstance}
	suite.css.EXPECT().ListInstances(environment.Cluster).Return(instances, nil)

	task1 := &models.Task{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.TaskDetail{
			ClusterARN:           instance1.Entity.ClusterARN,
			ContainerInstanceARN: instance1.Entity.ContainerInstanceARN,
			TaskARN:              aws.String("task-arn-1"),
			StartedBy:            currentDeployment.ID,
			DesiredStatus:        aws.String(runningTaskStatus),
		},
	}
	task2 := &models.Task{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.TaskDetail{
			ClusterARN:           instance2.Entity.ClusterARN,
			ContainerInstanceARN: instance2.Entity.ContainerInstanceARN,
			TaskARN:              aws.String("task-arn-2"),
			StartedBy:            currentDeployment.ID,
			DesiredStatus:        aws.String(runningTaskStatus),
		},
	}
	// task on newInstance which is not related to environment
	task3 := &models.Task{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.TaskDetail{
			ClusterARN:           newInstance.Entity.ClusterARN,
			ContainerInstanceARN: newInstance.Entity.ContainerInstanceARN,
			TaskARN:              aws.String("task-arn-3"),
			StartedBy:            "non-scheduler",
			DesiredStatus:        aws.String(runningTaskStatus),
		},
	}

	tasks := []*models.Task{task1, task2, task3}
	suite.css.EXPECT().ListTasks(environment.Cluster).Return(tasks, nil)

	deployments := []types.Deployment{currentDeployment}
	suite.deploymentService.EXPECT().ListDeploymentsSortedReverseChronologically(ctx, environment.Name).Return(deployments, nil)

	events := make(chan Event)
	scheduler := NewScheduler(ctx, events, suite.environmentService, suite.deploymentService, suite.css, suite.ecs)
	scheduler.Start()

	startDeploymentEvent := (<-events).(StartDeploymentEvent)
	assert.Equal(suite.T(), environment.Name, startDeploymentEvent.Environment.Name)
	assert.Equal(suite.T(), aws.StringValue(newInstance.Entity.ContainerInstanceARN),
		aws.StringValue(startDeploymentEvent.Instances[0]))

	schedulerEnvironmentEvent := (<-events).(SchedulerEnvironmentEvent)
	assert.Equal(suite.T(), environment.Name, schedulerEnvironmentEvent.Environment.Name)
}

func (suite *SchedulerTestSuite) TestRunInstancesWithOldDeployments() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*schedulerTickerDuration)
	defer cancel()

	environment := environmenttypes.Environment{
		Name:    "TestRunInstancesWithOldDeployments",
		Cluster: "testCluster",
	}
	environments := []environmenttypes.Environment{environment}
	suite.environmentService.EXPECT().ListEnvironments(ctx).Return(environments, nil)

	currentDeployment := types.Deployment{
		ID:     "dep-id",
		Status: types.DeploymentInProgress,
	}
	suite.deploymentService.EXPECT().GetCurrentDeployment(ctx, environment.Name).Return(&currentDeployment, nil)

	instance1 := &models.ContainerInstance{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.ContainerInstanceDetail{
			ClusterARN:           aws.String(environment.Cluster),
			ContainerInstanceARN: aws.String("instance-arn-1"),
			Status:               aws.String("ACTIVE"),
		},
	}
	instance2 := &models.ContainerInstance{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.ContainerInstanceDetail{
			ClusterARN:           aws.String(environment.Cluster),
			ContainerInstanceARN: aws.String("instance-arn-2"),
			Status:               aws.String("ACTIVE"),
		},
	}
	instances := []*models.ContainerInstance{instance1, instance2}
	suite.css.EXPECT().ListInstances(environment.Cluster).Return(instances, nil)

	oldDeployment := types.Deployment{
		ID:     "old-dep-id",
		Status: types.DeploymentCompleted,
	}
	deployments := []types.Deployment{currentDeployment, oldDeployment}

	task1 := &models.Task{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.TaskDetail{
			ClusterARN:           instance1.Entity.ClusterARN,
			ContainerInstanceARN: instance1.Entity.ContainerInstanceARN,
			TaskARN:              aws.String("task-arn-1"),
			StartedBy:            currentDeployment.ID,
			DesiredStatus:        aws.String(runningTaskStatus),
		},
	}
	task2 := &models.Task{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.TaskDetail{
			ClusterARN:           instance2.Entity.ClusterARN,
			ContainerInstanceARN: instance2.Entity.ContainerInstanceARN,
			TaskARN:              aws.String("task-arn-2"),
			StartedBy:            oldDeployment.ID,
			DesiredStatus:        aws.String(runningTaskStatus),
		},
	}
	tasks := []*models.Task{task1, task2}
	suite.css.EXPECT().ListTasks(environment.Cluster).Return(tasks, nil)

	suite.deploymentService.EXPECT().ListDeploymentsSortedReverseChronologically(ctx, environment.Name).Return(deployments, nil)

	events := make(chan Event)
	scheduler := NewScheduler(ctx, events, suite.environmentService, suite.deploymentService, suite.css, suite.ecs)
	scheduler.Start()

	stopTasksEvent := (<-events).(StopTasksEvent)
	assert.Equal(suite.T(), environment.Name, stopTasksEvent.Environment.Name)
	assert.Equal(suite.T(), aws.StringValue(task2.Entity.TaskARN), stopTasksEvent.Tasks[0])

	startDeploymentEvent := (<-events).(StartDeploymentEvent)
	assert.Equal(suite.T(), environment.Name, startDeploymentEvent.Environment.Name)
	assert.Equal(suite.T(), aws.StringValue(instance2.Entity.ContainerInstanceARN),
		aws.StringValue(startDeploymentEvent.Instances[0]))

	schedulerEnvironmentEvent := (<-events).(SchedulerEnvironmentEvent)
	assert.Equal(suite.T(), environment.Name, schedulerEnvironmentEvent.Environment.Name)
}

func (suite *SchedulerTestSuite) TestRunTrackedInstance() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*schedulerTickerDuration)
	defer cancel()

	environment := environmenttypes.Environment{
		Name:    "TestRunTrackedInstance",
		Cluster: "testCluster",
	}
	environments := []environmenttypes.Environment{environment}
	suite.environmentService.EXPECT().ListEnvironments(ctx).Return(environments, nil)

	currentDeployment := types.Deployment{
		ID:     "dep-id",
		Status: types.DeploymentInProgress,
		Health: types.DeploymentHealthy,
	}
	suite.deploymentService.EXPECT().GetCurrentDeployment(ctx, environment.Name).Return(&currentDeployment, nil)

	instance := &models.ContainerInstance{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.ContainerInstanceDetail{
			ClusterARN:           aws.String(environment.Cluster),
			ContainerInstanceARN: aws.String("instance-arn"),
			Status:               aws.String("ACTIVE"),
		},
	}
	instances := []*models.ContainerInstance{instance}
	suite.css.EXPECT().ListInstances(environment.Cluster).Return(instances, nil)

	// ListTasks from CSS returns empty due to lag
	tasks := []*models.Task{}
	suite.css.EXPECT().ListTasks(environment.Cluster).Return(tasks, nil)

	deployments := []types.Deployment{currentDeployment}
	suite.deploymentService.EXPECT().ListDeploymentsSortedReverseChronologically(ctx, environment.Name).Return(deployments, nil)

	suite.ecs.EXPECT().ListTasksByInstance(environment.Cluster, aws.StringValue(instance.Entity.ContainerInstanceARN)).Times(0)

	events := make(chan Event)
	scheduler := NewScheduler(ctx, events, suite.environmentService, suite.deploymentService, suite.css, suite.ecs)
	trackingInfo := make(map[string]time.Time)
	trackingInfo[aws.StringValue(instance.Entity.ContainerInstanceARN)] = time.Now().UTC()
	previousState := make(map[string]environmentExecutionState)
	previousState[environment.Name] = environmentExecutionState{
		environment:  environment,
		trackingInfo: trackingInfo,
		inProgress:   false,
	}
	scheduler.setExecutionState(previousState)
	scheduler.Start()

	schedulerEnvironmentEvent := (<-events).(SchedulerEnvironmentEvent)
	assert.Equal(suite.T(), environment.Name, schedulerEnvironmentEvent.Environment.Name)
}

func (suite *SchedulerTestSuite) TestRunTrackedInstanceTTLExpired() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*schedulerTickerDuration)
	defer cancel()

	environment := environmenttypes.Environment{
		Name:    "TestRunTrackedInstanceTTLExpired",
		Cluster: "testCluster",
	}
	environments := []environmenttypes.Environment{environment}
	suite.environmentService.EXPECT().ListEnvironments(ctx).Return(environments, nil)

	currentDeployment := types.Deployment{
		ID:     "dep-id",
		Status: types.DeploymentInProgress,
		Health: types.DeploymentHealthy,
	}
	suite.deploymentService.EXPECT().GetCurrentDeployment(ctx, environment.Name).Return(&currentDeployment, nil)

	instance := &models.ContainerInstance{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.ContainerInstanceDetail{
			ClusterARN:           aws.String(environment.Cluster),
			ContainerInstanceARN: aws.String("instance-arn"),
			Status:               aws.String("ACTIVE"),
		},
	}
	instances := []*models.ContainerInstance{instance}
	suite.css.EXPECT().ListInstances(environment.Cluster).Return(instances, nil)

	// ListTasks from CSS returns empty due to lag
	tasks := []*models.Task{}
	suite.css.EXPECT().ListTasks(environment.Cluster).Return(tasks, nil)

	deployments := []types.Deployment{currentDeployment}
	suite.deploymentService.EXPECT().ListDeploymentsSortedReverseChronologically(ctx, environment.Name).Return(deployments, nil)

	taskARNFromECS := []*string{aws.String("task-arn")}
	suite.ecs.EXPECT().ListTasksByInstance(environment.Cluster, aws.StringValue(instance.Entity.ContainerInstanceARN)).Return(taskARNFromECS, nil)
	tasksFromECS := &ecs.DescribeTasksOutput{
		Tasks: []*ecs.Task{
			&ecs.Task{
				TaskArn:              aws.String("task-arn"),
				ContainerInstanceArn: instance.Entity.ContainerInstanceARN,
				ClusterArn:           instance.Entity.ClusterARN,
				DesiredStatus:        aws.String("RUNNING"),
				StartedBy:            aws.String(currentDeployment.ID),
			},
		},
	}
	suite.ecs.EXPECT().DescribeTasks(environment.Cluster, taskARNFromECS).Return(tasksFromECS, nil)

	events := make(chan Event)
	scheduler := NewScheduler(ctx, events, suite.environmentService, suite.deploymentService, suite.css, suite.ecs)
	trackingInfo := make(map[string]time.Time)
	trackingInfo[aws.StringValue(instance.Entity.ContainerInstanceARN)] = time.Now().UTC().Add(-2 * trackingInfoTTL)
	previousState := make(map[string]environmentExecutionState)
	previousState[environment.Name] = environmentExecutionState{
		environment:  environment,
		trackingInfo: trackingInfo,
		inProgress:   false,
	}
	scheduler.setExecutionState(previousState)
	scheduler.Start()

	schedulerEnvironmentEvent := (<-events).(SchedulerEnvironmentEvent)
	assert.Equal(suite.T(), environment.Name, schedulerEnvironmentEvent.Environment.Name)
}

func (suite *SchedulerTestSuite) TestRunTrackedInstanceDescribeTasksReturnsError() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*schedulerTickerDuration)
	defer cancel()

	environment := environmenttypes.Environment{
		Name:    "TestRunTrackedInstanceDescribeTasksReturnsError",
		Cluster: "testCluster",
	}
	environments := []environmenttypes.Environment{environment}
	suite.environmentService.EXPECT().ListEnvironments(ctx).Return(environments, nil)

	currentDeployment := types.Deployment{
		ID:     "dep-id",
		Status: types.DeploymentInProgress,
		Health: types.DeploymentHealthy,
	}
	suite.deploymentService.EXPECT().GetCurrentDeployment(ctx, environment.Name).Return(&currentDeployment, nil)

	instance := &models.ContainerInstance{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.ContainerInstanceDetail{
			ClusterARN:           aws.String(environment.Cluster),
			ContainerInstanceARN: aws.String("instance-arn"),
			Status:               aws.String("ACTIVE"),
		},
	}
	instances := []*models.ContainerInstance{instance}
	suite.css.EXPECT().ListInstances(environment.Cluster).Return(instances, nil)

	// ListTasks from CSS returns empty due to lag
	tasks := []*models.Task{}
	suite.css.EXPECT().ListTasks(environment.Cluster).Return(tasks, nil)

	deployments := []types.Deployment{currentDeployment}
	suite.deploymentService.EXPECT().ListDeploymentsSortedReverseChronologically(ctx, environment.Name).Return(deployments, nil)

	taskARNFromECS := []*string{aws.String("task-arn")}
	suite.ecs.EXPECT().ListTasksByInstance(environment.Cluster, aws.StringValue(instance.Entity.ContainerInstanceARN)).Return(taskARNFromECS, nil)

	err := errors.Errorf("Error from ecs.DescribeTasks")
	suite.ecs.EXPECT().DescribeTasks(environment.Cluster, taskARNFromECS).Return(nil, err)

	events := make(chan Event)
	scheduler := NewScheduler(ctx, events, suite.environmentService, suite.deploymentService, suite.css, suite.ecs)
	trackingInfo := make(map[string]time.Time)
	trackingInfo[aws.StringValue(instance.Entity.ContainerInstanceARN)] = time.Now().UTC().Add(-2 * trackingInfoTTL)
	previousState := make(map[string]environmentExecutionState)
	previousState[environment.Name] = environmentExecutionState{
		environment:  environment,
		trackingInfo: trackingInfo,
		inProgress:   false,
	}
	scheduler.setExecutionState(previousState)
	scheduler.Start()
	schedulerErrorEvent := (<-events).(SchedulerErrorEvent)
	assert.Equal(suite.T(), err, errors.Cause(schedulerErrorEvent.Error))
}

func (suite *SchedulerTestSuite) TestRunTrackedInstanceListTasksReturnsError() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*schedulerTickerDuration)
	defer cancel()

	environment := environmenttypes.Environment{
		Name:    "TestRunTrackedInstanceListTasksReturnsError",
		Cluster: "testCluster",
	}
	environments := []environmenttypes.Environment{environment}
	suite.environmentService.EXPECT().ListEnvironments(ctx).Return(environments, nil)

	currentDeployment := types.Deployment{
		ID:     "dep-id",
		Status: types.DeploymentInProgress,
		Health: types.DeploymentHealthy,
	}
	suite.deploymentService.EXPECT().GetCurrentDeployment(ctx, environment.Name).Return(&currentDeployment, nil)

	instance := &models.ContainerInstance{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.ContainerInstanceDetail{
			ClusterARN:           aws.String(environment.Cluster),
			ContainerInstanceARN: aws.String("instance-arn"),
			Status:               aws.String("ACTIVE"),
		},
	}
	instances := []*models.ContainerInstance{instance}
	suite.css.EXPECT().ListInstances(environment.Cluster).Return(instances, nil)

	// ListTasks from CSS returns empty due to lag
	tasks := []*models.Task{}
	suite.css.EXPECT().ListTasks(environment.Cluster).Return(tasks, nil)

	deployments := []types.Deployment{currentDeployment}
	suite.deploymentService.EXPECT().ListDeploymentsSortedReverseChronologically(ctx, environment.Name).Return(deployments, nil)

	err := errors.Errorf("Error from ecs.ListTasks")
	suite.ecs.EXPECT().ListTasksByInstance(environment.Cluster, aws.StringValue(instance.Entity.ContainerInstanceARN)).Return(nil, err)

	events := make(chan Event)
	scheduler := NewScheduler(ctx, events, suite.environmentService, suite.deploymentService, suite.css, suite.ecs)
	trackingInfo := make(map[string]time.Time)
	trackingInfo[aws.StringValue(instance.Entity.ContainerInstanceARN)] = time.Now().UTC().Add(-2 * trackingInfoTTL)
	previousState := make(map[string]environmentExecutionState)
	previousState[environment.Name] = environmentExecutionState{
		environment:  environment,
		trackingInfo: trackingInfo,
		inProgress:   false,
	}
	scheduler.setExecutionState(previousState)
	scheduler.Start()
	schedulerErrorEvent := (<-events).(SchedulerErrorEvent)
	assert.Equal(suite.T(), err, errors.Cause(schedulerErrorEvent.Error))
}

func (suite *SchedulerTestSuite) TestRunTrackedInstanceListTasksReturnsEmpty() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*schedulerTickerDuration)
	defer cancel()

	environment := environmenttypes.Environment{
		Name:    "TestRunTrackedInstanceListTasksReturnsEmpty",
		Cluster: "testCluster",
	}
	environments := []environmenttypes.Environment{environment}
	suite.environmentService.EXPECT().ListEnvironments(ctx).Return(environments, nil)

	currentDeployment := types.Deployment{
		ID:     "dep-id",
		Status: types.DeploymentInProgress,
		Health: types.DeploymentHealthy,
	}
	suite.deploymentService.EXPECT().GetCurrentDeployment(ctx, environment.Name).Return(&currentDeployment, nil)

	instance := &models.ContainerInstance{
		Metadata: &models.Metadata{EntityVersion: aws.String("123")},
		Entity: &models.ContainerInstanceDetail{
			ClusterARN:           aws.String(environment.Cluster),
			ContainerInstanceARN: aws.String("instance-arn"),
			Status:               aws.String("ACTIVE"),
		},
	}
	instances := []*models.ContainerInstance{instance}
	suite.css.EXPECT().ListInstances(environment.Cluster).Return(instances, nil)

	// ListTasks from CSS returns empty due to lag
	tasks := []*models.Task{}
	suite.css.EXPECT().ListTasks(environment.Cluster).Return(tasks, nil)

	deployments := []types.Deployment{currentDeployment}
	suite.deploymentService.EXPECT().ListDeploymentsSortedReverseChronologically(ctx, environment.Name).Return(deployments, nil)

	taskARNFromECS := []*string{}
	suite.ecs.EXPECT().ListTasksByInstance(environment.Cluster, aws.StringValue(instance.Entity.ContainerInstanceARN)).Return(taskARNFromECS, nil)

	events := make(chan Event)
	scheduler := NewScheduler(ctx, events, suite.environmentService, suite.deploymentService, suite.css, suite.ecs)
	trackingInfo := make(map[string]time.Time)
	trackingInfo[aws.StringValue(instance.Entity.ContainerInstanceARN)] = time.Now().UTC().Add(-2 * trackingInfoTTL)
	previousState := make(map[string]environmentExecutionState)
	previousState[environment.Name] = environmentExecutionState{
		environment:  environment,
		trackingInfo: trackingInfo,
		inProgress:   false,
	}
	scheduler.setExecutionState(previousState)
	scheduler.Start()

	startDeploymentEvent := (<-events).(StartDeploymentEvent)
	assert.Equal(suite.T(), environment.Name, startDeploymentEvent.Environment.Name)
	assert.Equal(suite.T(), aws.StringValue(instance.Entity.ContainerInstanceARN),
		aws.StringValue(startDeploymentEvent.Instances[0]))

	schedulerEnvironmentEvent := (<-events).(SchedulerEnvironmentEvent)
	assert.Equal(suite.T(), environment.Name, schedulerEnvironmentEvent.Environment.Name)
}
