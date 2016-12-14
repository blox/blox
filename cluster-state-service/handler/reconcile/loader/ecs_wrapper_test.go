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

package loader

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/blox/blox/cluster-state-service/handler/mocks"
	"github.com/blox/blox/cluster-state-service/handler/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	ecsClusterARN1  = "arn:aws:ecs:us-east-1:123456789012:cluster/cluster1"
	ecsClusterARN2  = "arn:aws:ecs:us-east-1:123456789012:cluster/cluster2"
	ecsTaskARN1     = "arn:aws:ecs:us-east-1:123456789012:container-task/b6b9eace-958e-4f2a-a09c-8cf43b76cf97"
	ecsTaskARN2     = "arn:aws:ecs:us-east-1:123456789012:container-task/ab345dfe-6578-2eab-c671-72847ffe8122"
	ecsInstanceARN1 = "arn:aws:ecs:us-east-1:123456789012:container-instance/b6b9eace-958e-4f2a-a09c-8cf43b76cf97"
	ecsInstanceARN2 = "arn:aws:ecs:us-east-1:123456789012:container-instance/ab345dfe-6578-2eab-c671-72847ffe8122"
	ecsNextToken    = "eyJuZXh0VG9rZW4iOiBudWxsLCAiYm90b190cnVuY2F0ZV9hbW91bnQiOiAxfQ=="
)

type ECSWrapperTestSuite struct {
	suite.Suite
	mockECSClient *mocks.MockECSAPI
	ecsWrapper    ECSWrapper
	task          types.Task
	ecsTask       ecs.Task
	instance      types.ContainerInstance
	ecsInstance   ecs.ContainerInstance
}

func (suite *ECSWrapperTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())

	suite.mockECSClient = mocks.NewMockECSAPI(mockCtrl)

	suite.ecsWrapper = clientWrapper{
		client: suite.mockECSClient,
	}

	createdAt := "2016-11-07T15:30:00Z"
	startedAt := "2016-11-07T15:45:00Z"
	desiredStatus := "RUNNING"
	lastStatus := "PENDING"
	taskVersion := version
	ecsTaskInstanceARN := "arn:aws:ecs:us-east-1:123456789012:container-task/b6b9eace-958e-4f2a-a09c-8cf43b76cf97"
	ecsTaskDefinitionARN := "arn:aws:ecs:us-east-1:123456789012:task-definition/testTask:1"
	suite.task = types.Task{
		Detail: &types.TaskDetail{
			ClusterARN:           &ecsClusterARN1,
			ContainerInstanceARN: &ecsTaskInstanceARN,
			Containers:           []*types.Container{},
			CreatedAt:            &createdAt,
			DesiredStatus:        &desiredStatus,
			LastStatus:           &lastStatus,
			Overrides:            &types.Overrides{},
			StartedAt:            startedAt,
			TaskARN:              &ecsTaskARN1,
			TaskDefinitionARN:    &ecsTaskDefinitionARN,
			Version:              &taskVersion,
		},
	}

	ecsCreatedAt, err := time.Parse(timeLayout, createdAt)
	assert.Nil(suite.T(), err, "Cannot setup testSuite: Unexpected error when parsing time")
	ecsStartedAt, err := time.Parse(timeLayout, startedAt)
	assert.Nil(suite.T(), err, "Cannot setup testSuite: Unexpected error when parsing time")
	suite.ecsTask = ecs.Task{
		ClusterArn:           &ecsClusterARN1,
		ContainerInstanceArn: &ecsTaskInstanceARN,
		Containers:           []*ecs.Container{},
		CreatedAt:            &ecsCreatedAt,
		DesiredStatus:        &desiredStatus,
		LastStatus:           &lastStatus,
		Overrides:            &ecs.TaskOverride{},
		StartedAt:            &ecsStartedAt,
		TaskArn:              &ecsTaskARN1,
		TaskDefinitionArn:    &ecsTaskDefinitionARN,
	}

	agentConnected := true
	containerStatus := "ACTIVE"
	instanceVersion := version
	suite.instance = types.ContainerInstance{
		Detail: &types.InstanceDetail{
			AgentConnected:       &agentConnected,
			Attributes:           []*types.Attribute{},
			ClusterARN:           &ecsClusterARN1,
			ContainerInstanceARN: &ecsInstanceARN1,
			RegisteredResources:  []*types.Resource{},
			RemainingResources:   []*types.Resource{},
			Status:               &containerStatus,
			Version:              &instanceVersion,
			VersionInfo:          &types.VersionInfo{},
		},
	}

	suite.ecsInstance = ecs.ContainerInstance{
		AgentConnected:       &agentConnected,
		Attributes:           []*ecs.Attribute{},
		ContainerInstanceArn: &ecsInstanceARN1,
		RegisteredResources:  []*ecs.Resource{},
		RemainingResources:   []*ecs.Resource{},
		Status:               &containerStatus,
		VersionInfo:          &ecs.VersionInfo{},
	}
}

func TestECSWrapperTestSuite(t *testing.T) {
	suite.Run(t, new(ECSWrapperTestSuite))
}

func (suite *ECSWrapperTestSuite) TestListAllClustersECSListClustersWithoutTokenReturnsError() {
	in := ecs.ListClustersInput{}
	suite.mockECSClient.EXPECT().ListClusters(&in).Return(nil, errors.New("Error while listing clusters without next token"))

	_, err := suite.ecsWrapper.ListAllClusters()

	assert.Error(suite.T(), err, "Expected an error when ecs client returns an error when listing clusters without next token")
}

func (suite *ECSWrapperTestSuite) TestListAllClustersECSListClustersWithTokenReturnsError() {
	in1 := ecs.ListClustersInput{}
	resp := ecs.ListClustersOutput{
		ClusterArns: []*string{&ecsClusterARN1},
		NextToken:   &ecsNextToken,
	}
	listClustersWithoutTokenCall := suite.mockECSClient.EXPECT().ListClusters(&in1).Return(&resp, nil)

	in2 := ecs.ListClustersInput{
		NextToken: &ecsNextToken,
	}
	listClustersWithTokenCall := suite.mockECSClient.EXPECT().ListClusters(&in2).Return(nil, errors.New("Error while listing clusters with next token"))

	gomock.InOrder(listClustersWithoutTokenCall, listClustersWithTokenCall)

	_, err := suite.ecsWrapper.ListAllClusters()

	assert.Error(suite.T(), err, "Expected an error when ecs client returns an error when listing clusters with next token")
}

func (suite *ECSWrapperTestSuite) TestListAllClusters() {
	in1 := ecs.ListClustersInput{}
	resp1 := ecs.ListClustersOutput{
		ClusterArns: []*string{&ecsClusterARN1},
		NextToken:   &ecsNextToken,
	}
	listClustersWithoutTokenCall := suite.mockECSClient.EXPECT().ListClusters(&in1).Return(&resp1, nil)

	in2 := ecs.ListClustersInput{
		NextToken: &ecsNextToken,
	}
	resp2 := ecs.ListClustersOutput{
		ClusterArns: []*string{&ecsClusterARN2},
	}
	listClustersWithTokenCall := suite.mockECSClient.EXPECT().ListClusters(&in2).Return(&resp2, nil)

	gomock.InOrder(listClustersWithoutTokenCall, listClustersWithTokenCall)

	clusterARNs, err := suite.ecsWrapper.ListAllClusters()
	assert.Nil(suite.T(), err, "Unexpected error when listing clusters")
	expectedClusterARNs := []*string{&ecsClusterARN1, &ecsClusterARN2}
	assert.Equal(suite.T(), expectedClusterARNs, clusterARNs, "Cluster ARNs received using list clusters is not equal to the expected list")
}

func (suite *ECSWrapperTestSuite) TestListAllTasksECSListTasksWithoutTokenReturnsError() {
	in := ecs.ListTasksInput{
		Cluster: &ecsClusterARN1,
	}
	suite.mockECSClient.EXPECT().ListTasks(&in).Return(nil, errors.New("Error while listing tasks without next token"))

	_, err := suite.ecsWrapper.ListAllTasks(&ecsClusterARN1)

	assert.Error(suite.T(), err, "Expected an error when ecs client returns an error when listing tasks without next token")
}

func (suite *ECSWrapperTestSuite) TestListAllTasksECSListTasksWithTokenReturnsError() {
	in1 := ecs.ListTasksInput{
		Cluster: &ecsClusterARN1,
	}
	resp := ecs.ListTasksOutput{
		TaskArns:  []*string{&ecsTaskARN1},
		NextToken: &ecsNextToken,
	}
	listTasksWithoutTokenCall := suite.mockECSClient.EXPECT().ListTasks(&in1).Return(&resp, nil)

	in2 := ecs.ListTasksInput{
		Cluster:   &ecsClusterARN1,
		NextToken: &ecsNextToken,
	}
	listTasksWithTokenCall := suite.mockECSClient.EXPECT().ListTasks(&in2).Return(nil, errors.New("Error while listing tasks with next token"))

	gomock.InOrder(listTasksWithoutTokenCall, listTasksWithTokenCall)

	_, err := suite.ecsWrapper.ListAllTasks(&ecsClusterARN1)

	assert.Error(suite.T(), err, "Expected an error when ecs client returns an error when listing tasks with next token")
}

func (suite *ECSWrapperTestSuite) TestListAllTasks() {
	in1 := ecs.ListTasksInput{
		Cluster: &ecsClusterARN1,
	}
	resp1 := ecs.ListTasksOutput{
		TaskArns:  []*string{&ecsTaskARN1},
		NextToken: &ecsNextToken,
	}
	listTasksWithoutTokenCall := suite.mockECSClient.EXPECT().ListTasks(&in1).Return(&resp1, nil)

	in2 := ecs.ListTasksInput{
		Cluster:   &ecsClusterARN1,
		NextToken: &ecsNextToken,
	}
	resp2 := ecs.ListTasksOutput{
		TaskArns: []*string{&ecsTaskARN2},
	}
	listTasksWithTokenCall := suite.mockECSClient.EXPECT().ListTasks(&in2).Return(&resp2, nil)

	gomock.InOrder(listTasksWithoutTokenCall, listTasksWithTokenCall)

	taskARNs, err := suite.ecsWrapper.ListAllTasks(&ecsClusterARN1)
	assert.Nil(suite.T(), err, "Unexpected error when listing tasks")
	expectedTaskARNs := []*string{&ecsTaskARN1, &ecsTaskARN2}
	assert.Equal(suite.T(), expectedTaskARNs, taskARNs, "Task ARNs received using list tasks is not equal to the expected list")
}

func (suite *ECSWrapperTestSuite) DescribeTasksECSDescribeTasksReturnsError() {
	taskList := []*string{&ecsTaskARN1}
	in := ecs.DescribeTasksInput{
		Cluster: &ecsClusterARN1,
		Tasks:   taskList,
	}
	suite.mockECSClient.EXPECT().DescribeTasks(&in).Return(nil, errors.New("Error while describing tasks"))

	_, _, err := suite.ecsWrapper.DescribeTasks(&ecsClusterARN1, taskList)

	assert.Error(suite.T(), err, "Expected an error when ecs client returns an error when describing tasks")
}

func (suite *ECSWrapperTestSuite) DescribeTasks() {
	taskList := []*string{&ecsTaskARN1, &ecsTaskARN2}
	in := ecs.DescribeTasksInput{
		Cluster: &ecsClusterARN1,
		Tasks:   taskList,
	}
	resp := ecs.DescribeTasksOutput{
		Tasks: []*ecs.Task{
			&suite.ecsTask,
		},
		Failures: []*ecs.Failure{
			&ecs.Failure{
				Arn: &ecsTaskARN2,
			},
		},
	}
	suite.mockECSClient.EXPECT().DescribeTasks(&in).Return(resp, nil)

	tasks, failures, err := suite.ecsWrapper.DescribeTasks(&ecsClusterARN1, taskList)

	assert.Nil(suite.T(), err, "Unexpected error when describing tasks")
	expectedTasks := []types.Task{suite.task}
	assert.Equal(suite.T(), expectedTasks, tasks, "Tasks received on describing tasks does not match expected tasks")
	expectedFailures := []string{ecsTaskARN2}
	assert.Equal(suite.T(), expectedFailures, failures, "Failures received on describing tasks does not match expected failures")
}

func (suite *ECSWrapperTestSuite) TestListAllContainerInstancesECSListContainerInstancesWithoutTokenReturnsError() {
	in := ecs.ListContainerInstancesInput{
		Cluster: &ecsClusterARN1,
	}
	suite.mockECSClient.EXPECT().ListContainerInstances(&in).Return(nil, errors.New("Error while listing container instances without next token"))

	_, err := suite.ecsWrapper.ListAllContainerInstances(&ecsClusterARN1)

	assert.Error(suite.T(), err, "Expected an error when ecs client returns an error when listing container instances without next token")
}

func (suite *ECSWrapperTestSuite) TestListAllContainerInstancesECSListContainerInstancesWithTokenReturnsError() {
	in1 := ecs.ListContainerInstancesInput{
		Cluster: &ecsClusterARN1,
	}
	resp := ecs.ListContainerInstancesOutput{
		ContainerInstanceArns: []*string{&ecsInstanceARN1},
		NextToken:             &ecsNextToken,
	}
	listInstancesWithoutTokenCall := suite.mockECSClient.EXPECT().ListContainerInstances(&in1).Return(&resp, nil)

	in2 := ecs.ListContainerInstancesInput{
		Cluster:   &ecsClusterARN1,
		NextToken: &ecsNextToken,
	}
	listInstancesWithTokenCall := suite.mockECSClient.EXPECT().ListContainerInstances(&in2).Return(nil, errors.New("Error while listing container instances with next token"))

	gomock.InOrder(listInstancesWithoutTokenCall, listInstancesWithTokenCall)

	_, err := suite.ecsWrapper.ListAllContainerInstances(&ecsClusterARN1)

	assert.Error(suite.T(), err, "Expected an error when ecs client returns an error when listing container instances with next token")
}

func (suite *ECSWrapperTestSuite) TestListAllContainerInstances() {
	in1 := ecs.ListContainerInstancesInput{
		Cluster: &ecsClusterARN1,
	}
	resp1 := ecs.ListContainerInstancesOutput{
		ContainerInstanceArns: []*string{&ecsInstanceARN1},
		NextToken:             &ecsNextToken,
	}
	listInstancesWithoutTokenCall := suite.mockECSClient.EXPECT().ListContainerInstances(&in1).Return(&resp1, nil)

	in2 := ecs.ListContainerInstancesInput{
		Cluster:   &ecsClusterARN1,
		NextToken: &ecsNextToken,
	}
	resp2 := ecs.ListContainerInstancesOutput{
		ContainerInstanceArns: []*string{&ecsInstanceARN2},
	}
	listInstancesWithTokenCall := suite.mockECSClient.EXPECT().ListContainerInstances(&in2).Return(&resp2, nil)

	gomock.InOrder(listInstancesWithoutTokenCall, listInstancesWithTokenCall)

	instanceARNs, err := suite.ecsWrapper.ListAllContainerInstances(&ecsClusterARN1)
	assert.Nil(suite.T(), err, "Unexpected error when listing container instances")
	expectedContainerInstanceARNs := []*string{&ecsInstanceARN1, &ecsInstanceARN2}
	assert.Equal(suite.T(), expectedContainerInstanceARNs, instanceARNs, "ContainerInstance ARNs received using list instances is not equal to the expected list")
}

func (suite *ECSWrapperTestSuite) DescribeContainerInstancesECSDescribeContainerInstancesReturnsError() {
	instanceList := []*string{&ecsInstanceARN1}
	in := ecs.DescribeContainerInstancesInput{
		Cluster:            &ecsClusterARN1,
		ContainerInstances: instanceList,
	}
	suite.mockECSClient.EXPECT().DescribeContainerInstances(&in).Return(nil, errors.New("Error while describing container instances"))

	_, _, err := suite.ecsWrapper.DescribeContainerInstances(&ecsClusterARN1, instanceList)

	assert.Error(suite.T(), err, "Expected an error when ecs client returns an error when describing container instances")
}

func (suite *ECSWrapperTestSuite) DescribeContainerInstances() {
	instanceList := []*string{&ecsInstanceARN1, &ecsInstanceARN2}
	in := ecs.DescribeContainerInstancesInput{
		Cluster:            &ecsClusterARN1,
		ContainerInstances: instanceList,
	}
	resp := ecs.DescribeContainerInstancesOutput{
		ContainerInstances: []*ecs.ContainerInstance{
			&suite.ecsInstance,
		},
		Failures: []*ecs.Failure{
			&ecs.Failure{
				Arn: &ecsInstanceARN2,
			},
		},
	}
	suite.mockECSClient.EXPECT().DescribeContainerInstances(&in).Return(resp, nil)

	instances, failures, err := suite.ecsWrapper.DescribeContainerInstances(&ecsClusterARN1, instanceList)

	assert.Nil(suite.T(), err, "Unexpected error when describing container instances")
	expectedInstances := []types.ContainerInstance{suite.instance}
	assert.Equal(suite.T(), expectedInstances, instances, "Container instances received on describing container instances does not match expected container instances")
	expectedFailures := []string{ecsInstanceARN2}
	assert.Equal(suite.T(), expectedFailures, failures, "Failures received on describing container instances does not match expected failures")
}
