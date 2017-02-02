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

package loader

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/blox/blox/cluster-state-service/handler/mocks"
	storetypes "github.com/blox/blox/cluster-state-service/handler/store/types"
	"github.com/blox/blox/cluster-state-service/handler/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	instanceClusterARN1           = "arn:aws:ecs:us-east-1:123456789012:cluster/cluster1"
	instanceClusterARN2           = "arn:aws:ecs:us-east-1:123456789012:cluster/cluster2"
	instanceARN1                  = "arn:aws:ecs:us-east-1:123456789012:container-instance/b6b9eace-958e-4f2a-a09c-8cf43b76cf97"
	instanceARN2                  = "arn:aws:ecs:us-east-1:123456789012:container-instance/ab345dfe-6578-2eab-c671-72847ffe8122"
	redundantClusterARNOfInstance = "arn:aws:ecs:us-east-1:123456789012:cluster/red-un-da-nt"
	redundantInstanceARN          = "arn:aws:ecs:us-east-1:123456789012:container-instance/"
)

type InstanceLoaderTestSuite struct {
	suite.Suite
	instanceStore              *mocks.MockContainerInstanceStore
	ecsWrapper                 *mocks.MockECSWrapper
	instanceLoader             ContainerInstanceLoader
	clusterARNList             []*string
	instance                   types.ContainerInstance
	versionedInstance          storetypes.VersionedContainerInstance
	redundantInstance          types.ContainerInstance
	redundantVersionedInstance storetypes.VersionedContainerInstance
	instanceJSON               string
}

func (suite *InstanceLoaderTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())

	suite.instanceStore = mocks.NewMockContainerInstanceStore(mockCtrl)

	suite.ecsWrapper = mocks.NewMockECSWrapper(mockCtrl)

	suite.instanceLoader = instanceLoader{
		instanceStore: suite.instanceStore,
		ecsWrapper:    suite.ecsWrapper,
	}

	suite.clusterARNList = []*string{&instanceClusterARN1, &instanceClusterARN2}

	agentConnected := true
	containerStatus := "ACTIVE"
	instanceVersion := version
	suite.instance = types.ContainerInstance{
		Detail: &types.InstanceDetail{
			AgentConnected:       &agentConnected,
			Attributes:           []*types.Attribute{},
			ClusterARN:           &instanceClusterARN1,
			ContainerInstanceARN: &instanceARN1,
			RegisteredResources:  []*types.Resource{},
			RemainingResources:   []*types.Resource{},
			Status:               &containerStatus,
			Version:              &instanceVersion,
			VersionInfo:          &types.VersionInfo{},
		},
	}
	suite.versionedInstance = storetypes.VersionedContainerInstance{
		ContainerInstance: suite.instance,
		Version: "123",
	}

	ins, err := json.Marshal(suite.instance)
	assert.Nil(suite.T(), err, "Cannot setup testSuite: Unexpected error when marshaling instance")
	suite.instanceJSON = string(ins)

	suite.redundantInstance = types.ContainerInstance{
		Detail: &types.InstanceDetail{
			ClusterARN:           &redundantClusterARNOfInstance,
			ContainerInstanceARN: &redundantInstanceARN,
		},
	}
	suite.redundantVersionedInstance = storetypes.VersionedContainerInstance{
		ContainerInstance: suite.redundantInstance,
		Version: "123",
	}
}

func TestInstanceLoaderTestSuite(t *testing.T) {
	suite.Run(t, new(InstanceLoaderTestSuite))
}

func (suite *InstanceLoaderTestSuite) TestLoadContainerInstancesListAllClustersReturnsError() {
	gomock.InOrder(
		suite.instanceStore.EXPECT().ListContainerInstances().Return(make([]storetypes.VersionedContainerInstance, 0), nil),
		suite.ecsWrapper.EXPECT().ListAllClusters().Return(nil, errors.New("Error while listing all clusters")),
		suite.ecsWrapper.EXPECT().ListAllContainerInstances(gomock.Any()).Times(0),
		suite.ecsWrapper.EXPECT().DescribeContainerInstances(gomock.Any(), gomock.Any()).Times(0),
	)

	err := suite.instanceLoader.LoadContainerInstances()
	assert.Error(suite.T(), err, "Expected an error when ecs returns an error when listing all clusters")
}

func (suite *InstanceLoaderTestSuite) TestLoadContainerInstancesListAllContainerInstancesReturnsError() {
	gomock.InOrder(
		suite.instanceStore.EXPECT().ListContainerInstances().Return(make([]storetypes.VersionedContainerInstance, 0), nil),
		suite.ecsWrapper.EXPECT().ListAllClusters().Return(suite.clusterARNList, nil),
		suite.ecsWrapper.EXPECT().ListAllContainerInstances(suite.clusterARNList[0]).Return(nil, errors.New("Error while listing all container instances")),
		suite.ecsWrapper.EXPECT().ListAllContainerInstances(suite.clusterARNList[1]).Times(0),
		suite.ecsWrapper.EXPECT().DescribeContainerInstances(gomock.Any(), gomock.Any()).Times(0),
	)

	err := suite.instanceLoader.LoadContainerInstances()
	assert.Error(suite.T(), err, "Expected an error when ecs returns an error when listing all container instances in a cluster")
}

func (suite *InstanceLoaderTestSuite) TestLoadContainerInstancesDescribeContainerInstancesReturnsError() {
	instanceARNList := []*string{&instanceARN1}

	gomock.InOrder(
		suite.instanceStore.EXPECT().ListContainerInstances().Return(make([]storetypes.VersionedContainerInstance, 0), nil),
		suite.ecsWrapper.EXPECT().ListAllClusters().Return(suite.clusterARNList, nil),
		suite.ecsWrapper.EXPECT().ListAllContainerInstances(suite.clusterARNList[0]).Return(instanceARNList, nil),
		suite.ecsWrapper.EXPECT().ListAllContainerInstances(suite.clusterARNList[1]).Times(0),
		suite.ecsWrapper.EXPECT().DescribeContainerInstances(suite.clusterARNList[0], instanceARNList).Return(nil, nil, errors.New("Error while desribing container instance")),
	)
	err := suite.instanceLoader.LoadContainerInstances()
	assert.Error(suite.T(), err, "Expected an error when ecs returns an error when describing container instances")
}

func (suite *InstanceLoaderTestSuite) TestLoadContainerInstancesStoreReturnsError() {
	instanceARNList := []*string{&instanceARN1}
	instanceList := []types.ContainerInstance{suite.instance}
	emptyInstanceARNList := []*string{}

	gomock.InOrder(
		suite.instanceStore.EXPECT().ListContainerInstances().Return(make([]storetypes.VersionedContainerInstance, 0), nil),
		suite.ecsWrapper.EXPECT().ListAllClusters().Return(suite.clusterARNList, nil),
		suite.ecsWrapper.EXPECT().ListAllContainerInstances(suite.clusterARNList[0]).Return(instanceARNList, nil),
		suite.ecsWrapper.EXPECT().DescribeContainerInstances(suite.clusterARNList[0], instanceARNList).Return(instanceList, nil, nil),
		suite.ecsWrapper.EXPECT().ListAllContainerInstances(suite.clusterARNList[1]).Return(emptyInstanceARNList, nil),
		suite.ecsWrapper.EXPECT().DescribeContainerInstances(suite.clusterARNList[1], gomock.Any()).Times(0),
		suite.instanceStore.EXPECT().AddUnversionedContainerInstance(suite.instanceJSON).Return(errors.New("Error while adding container instance to store")),
	)
	err := suite.instanceLoader.LoadContainerInstances()
	assert.Error(suite.T(), err, "Expected an error when store returns an error when adding container instance")
}

func (suite *InstanceLoaderTestSuite) TestLoadContainerInstancesEmptyLocalStore() {
	instanceARNList := []*string{&instanceARN1}
	instanceList := []types.ContainerInstance{suite.instance}
	emptyInstanceARNList := []*string{}
	gomock.InOrder(
		suite.instanceStore.EXPECT().ListContainerInstances().Return(make([]storetypes.VersionedContainerInstance, 0), nil),
		suite.ecsWrapper.EXPECT().ListAllClusters().Return(suite.clusterARNList, nil),
		suite.ecsWrapper.EXPECT().ListAllContainerInstances(suite.clusterARNList[0]).Return(instanceARNList, nil),
		suite.ecsWrapper.EXPECT().DescribeContainerInstances(suite.clusterARNList[0], instanceARNList).Return(instanceList, nil, nil),
		suite.ecsWrapper.EXPECT().ListAllContainerInstances(suite.clusterARNList[1]).Return(emptyInstanceARNList, nil),
		suite.ecsWrapper.EXPECT().DescribeContainerInstances(suite.clusterARNList[1], gomock.Any()).Times(0),
		suite.instanceStore.EXPECT().AddUnversionedContainerInstance(suite.instanceJSON).Return(nil),
	)
	err := suite.instanceLoader.LoadContainerInstances()
	assert.Nil(suite.T(), err, "Unexpected error when loading container instances")
}

func (suite *InstanceLoaderTestSuite) TestLoadContainerInstancesLocalStoreSameAsECS() {
	instanceARNList := []*string{&instanceARN1}
	emptyInstanceARNList := []*string{}
	instanceListInStore := []storetypes.VersionedContainerInstance{suite.versionedInstance}
	instanceList := []types.ContainerInstance{suite.instance}
	// instanceListInStore == instanceList, which should mean that there shouldn't
	// be a call to DeleteContainerInstance()
	suite.instanceStore.EXPECT().DeleteContainerInstance(gomock.Any(), gomock.Any()).Return(nil).Times(0)
	gomock.InOrder(
		suite.instanceStore.EXPECT().ListContainerInstances().Return(instanceListInStore, nil),
		suite.ecsWrapper.EXPECT().ListAllClusters().Return(suite.clusterARNList, nil),
		suite.ecsWrapper.EXPECT().ListAllContainerInstances(suite.clusterARNList[0]).Return(instanceARNList, nil),
		suite.ecsWrapper.EXPECT().DescribeContainerInstances(suite.clusterARNList[0], instanceARNList).Return(instanceList, nil, nil),
		suite.ecsWrapper.EXPECT().ListAllContainerInstances(suite.clusterARNList[1]).Return(emptyInstanceARNList, nil),
		suite.ecsWrapper.EXPECT().DescribeContainerInstances(suite.clusterARNList[1], gomock.Any()).Times(0),
		suite.instanceStore.EXPECT().AddUnversionedContainerInstance(suite.instanceJSON).Return(nil),
	)
	err := suite.instanceLoader.LoadContainerInstances()
	assert.Nil(suite.T(), err, "Unexpected error when loading container instances")
}

func (suite *InstanceLoaderTestSuite) TestLoadContainerInstancesRedundantEntriesInLocalStore() {
	instanceARNList := []*string{&instanceARN1}
	emptyInstanceARNList := []*string{}
	instanceListInStore := []storetypes.VersionedContainerInstance{suite.versionedInstance, suite.redundantVersionedInstance}
	instanceList := []types.ContainerInstance{suite.instance}
	gomock.InOrder(
		suite.instanceStore.EXPECT().ListContainerInstances().Return(instanceListInStore, nil),
		suite.ecsWrapper.EXPECT().ListAllClusters().Return(suite.clusterARNList, nil),
		suite.ecsWrapper.EXPECT().ListAllContainerInstances(suite.clusterARNList[0]).Return(instanceARNList, nil),
		suite.ecsWrapper.EXPECT().DescribeContainerInstances(suite.clusterARNList[0], instanceARNList).Return(instanceList, nil, nil),
		suite.ecsWrapper.EXPECT().ListAllContainerInstances(suite.clusterARNList[1]).Return(emptyInstanceARNList, nil),
		suite.ecsWrapper.EXPECT().DescribeContainerInstances(suite.clusterARNList[1], gomock.Any()).Times(0),
		suite.instanceStore.EXPECT().AddUnversionedContainerInstance(suite.instanceJSON).Return(nil),
		// Expect delete container instance for the redundant instance
		suite.instanceStore.EXPECT().DeleteContainerInstance(redundantClusterARNOfInstance, redundantInstanceARN).Return(nil),
	)
	err := suite.instanceLoader.LoadContainerInstances()
	assert.Nil(suite.T(), err, "Unexpected error when loading container instances")
}
