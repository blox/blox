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

package reconcile

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ReconcilerTestSuite struct {
	suite.Suite
	taskLoader     *mocks.MockTaskLoader
	instanceLoader *mocks.MockContainerInstanceLoader
}

func (suite *ReconcilerTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	suite.taskLoader = mocks.NewMockTaskLoader(mockCtrl)
	suite.instanceLoader = mocks.NewMockContainerInstanceLoader(mockCtrl)
}

func TestReconcilerTestSuite(t *testing.T) {
	suite.Run(t, new(ReconcilerTestSuite))
}

func (suite *ReconcilerTestSuite) TestRunLoadTasksReturnsError() {
	reconciler := Reconciler{
		taskLoader:     suite.taskLoader,
		instanceLoader: suite.instanceLoader,
	}

	suite.taskLoader.EXPECT().LoadTasks().Return(errors.New("Error while loading tasks"))
	err := reconciler.RunOnce()
	assert.Error(suite.T(), err, "Expected an error when load tasks returns an error")
}

func (suite *ReconcilerTestSuite) TestRunLoadInstancesReturnsError() {
	reconciler := Reconciler{
		taskLoader:     suite.taskLoader,
		instanceLoader: suite.instanceLoader,
	}
	suite.taskLoader.EXPECT().LoadTasks().Return(nil)
	suite.instanceLoader.EXPECT().LoadContainerInstances().Return(errors.New("Error while loading instance"))

	err := reconciler.RunOnce()
	assert.Error(suite.T(), err, "Expected an error when load instances returns an error")
}

func (suite *ReconcilerTestSuite) TestRun() {
	reconciler := Reconciler{
		taskLoader:     suite.taskLoader,
		instanceLoader: suite.instanceLoader,
	}
	verifyInProgress := func() {
		assert.True(suite.T(), reconciler.isInProgress(), "Reconcile operation should be in progress")
	}
	suite.taskLoader.EXPECT().LoadTasks().Do(verifyInProgress).Return(nil)
	suite.instanceLoader.EXPECT().LoadContainerInstances().Do(verifyInProgress).Return(nil)

	err := reconciler.RunOnce()
	assert.Nil(suite.T(), err, "Unexpected error when performing bootstrapping")
	assert.False(suite.T(), reconciler.isInProgress(), "Reconcile operation should not be in progress")
}

func (suite *ReconcilerTestSuite) TestOverlappingRunInvocationsAreSkipped() {
	ctx, cancel := context.WithCancel(context.TODO())
	tickerDuration := 10 * time.Millisecond
	reconciler := Reconciler{
		taskLoader:     suite.taskLoader,
		instanceLoader: suite.instanceLoader,
		ctx:            ctx,
		tickerDuration: tickerDuration,
	}

	// verifyInProgress will be invoked by the LoadContainerInstances, in reconciler.Run()
	// This will cause reconciler.RunOnce() to be blocked because of the time.Sleep() call in
	// this method, which should result in reconciler.ticker's ticks being missed.
	// If there was a bug and the ticks were processed and resulted in reconciler.RunOnce() to
	// be invoked, the tests should fail as there are no matching EXPECT statements for
	// those calls.
	verifyInProgress := func() {
		assert.True(suite.T(), reconciler.isInProgress(), "Reconcile operation should be in progress")
		time.Sleep(3 * tickerDuration)
		cancel()
	}
	suite.taskLoader.EXPECT().LoadTasks().Return(nil)
	suite.instanceLoader.EXPECT().LoadContainerInstances().Do(verifyInProgress).Return(nil)
	reconciler.Run()
	select {
	case <-ctx.Done():
	}
}

func (suite *ReconcilerTestSuite) TestMultipleRunInvocations() {
	ctx, cancel := context.WithCancel(context.TODO())
	tickerDuration := 10 * time.Millisecond
	reconciler := Reconciler{
		taskLoader:     suite.taskLoader,
		instanceLoader: suite.instanceLoader,
		ctx:            ctx,
		tickerDuration: tickerDuration,
	}

	verifyInProgress := func() {
		assert.True(suite.T(), reconciler.isInProgress(), "Reconcile operation should be in progress")
		cancel()
	}
	gomock.InOrder(
		suite.taskLoader.EXPECT().LoadTasks().Return(nil),
		suite.instanceLoader.EXPECT().LoadContainerInstances().Return(nil),
		suite.taskLoader.EXPECT().LoadTasks().Return(nil),
		// Stop the Run() method by cancelling the context during its second invocation
		suite.instanceLoader.EXPECT().LoadContainerInstances().Do(verifyInProgress).Return(nil),
	)
	reconciler.Run()
	select {
	case <-ctx.Done():
	}
}
