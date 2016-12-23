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

	mocks "github.com/blox/blox/daemon-scheduler/pkg/mocks"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MonitorTestSuite struct {
	suite.Suite
	ctx         context.Context
	environment *mocks.MockEnvironment
	monitor     Monitor
	env1        *types.Environment
	env2        *types.Environment
}

func (suite *MonitorTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	suite.ctx = context.Background()
	suite.environment = mocks.NewMockEnvironment(mockCtrl)

	var err error
	suite.env1, err = types.NewEnvironment(environmentName1, taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Cannot initialize MonitorTestSuite")

	suite.env2, err = types.NewEnvironment(environmentName2, taskDefinition, cluster2)
	assert.Nil(suite.T(), err, "Cannot initialize MonitorTestSuite")
}

func TestMonitorTestSuite(t *testing.T) {
	suite.Run(t, new(MonitorTestSuite))
}

func (suite *MonitorTestSuite) TestInProgressListEnvironmentsFails() {
	ctx, cancel := context.WithCancel(suite.ctx)
	defer cancel()
	events := make(chan Event)

	monitor := NewMonitor(ctx, suite.environment, events)
	suite.environment.EXPECT().ListEnvironments(ctx).Return(nil, errors.New("Could not retrieve environments"))

	monitor.InProgressMonitorLoop(1 * time.Millisecond)
	monitorErrorEvent := (<-events).(MonitorErrorEvent)
	assert.Error(suite.T(), errors.Cause(monitorErrorEvent.Error), "Expected a monitorErrorEvent")
}

func (suite *MonitorTestSuite) TestInProgressListMultipleEnvironments() {
	ctx, cancel := context.WithCancel(suite.ctx)
	defer cancel()
	events := make(chan Event)

	environments := []types.Environment{*suite.env1, *suite.env2}
	environmentsMap := map[string]types.Environment{
		suite.env1.Name: *suite.env1,
		suite.env2.Name: *suite.env2,
	}

	monitor := NewMonitor(ctx, suite.environment, events)
	suite.environment.EXPECT().ListEnvironments(ctx).Return(environments, nil)

	monitor.InProgressMonitorLoop(1 * time.Millisecond)

	for i := 0; i < len(environments); i++ {
		inProgressEvent, ok := (<-events).(UpdateInProgressDeploymentEvent)
		if !ok {
			assert.Fail(suite.T(), "Expected only in-progress deployment events")
		}

		_, ok = environmentsMap[inProgressEvent.Environment.Name]
		if !ok {
			assert.Fail(suite.T(), "Unexpected environment in in-progress deployment event")
		}
	}
}

func (suite *MonitorTestSuite) TestPendingListEnvironmentsFails() {
	ctx, cancel := context.WithCancel(suite.ctx)
	defer cancel()
	events := make(chan Event)

	monitor := NewMonitor(ctx, suite.environment, events)
	suite.environment.EXPECT().ListEnvironments(ctx).Return(nil, errors.New("Could not retrieve environments"))

	monitor.PendingMonitorLoop(1 * time.Millisecond)
	monitorErrorEvent := (<-events).(MonitorErrorEvent)
	assert.Error(suite.T(), errors.Cause(monitorErrorEvent.Error), "Expected a monitorErrorEvent")
}

func (suite *MonitorTestSuite) TestPendingListMultipleEnvironments() {
	ctx, cancel := context.WithCancel(suite.ctx)
	defer cancel()
	events := make(chan Event)

	environments := []types.Environment{*suite.env1, *suite.env2}
	environmentsMap := map[string]types.Environment{
		suite.env1.Name: *suite.env1,
		suite.env2.Name: *suite.env2,
	}

	monitor := NewMonitor(ctx, suite.environment, events)
	suite.environment.EXPECT().ListEnvironments(ctx).Return(environments, nil)

	monitor.PendingMonitorLoop(1 * time.Millisecond)

	for i := 0; i < len(environments); i++ {
		pendingEvent, ok := (<-events).(UpdatePendingDeploymentEvent)
		if !ok {
			assert.Fail(suite.T(), "Expected only pending deployment events")
		}

		_, ok = environmentsMap[pendingEvent.Environment.Name]
		if !ok {
			assert.Fail(suite.T(), "Unexpected environment in pending deployment event")
		}
	}
}
