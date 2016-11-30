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
	"time"

	"github.com/blox/blox/daemon-scheduler/pkg/deployment"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

const (
	inProgressMonitorTickerDuration = 10 * time.Second
)

type Monitor interface {
	InProgressMonitorLoop()
}

type monitor struct {
	ctx         context.Context
	environment deployment.Environment
	events      chan<- Event
}

func NewMonitor(
	ctx context.Context,
	environment deployment.Environment,
	events chan<- Event) Monitor {

	return monitor{
		ctx:         ctx,
		environment: environment,
		events:      events,
	}
}

func (m monitor) InProgressMonitorLoop() {
	ticker := time.NewTicker(inProgressMonitorTickerDuration)
	go func() {
		for {
			select {
			case <-ticker.C:
				m.runOnce()
			case <-m.ctx.Done():
				log.Infof("Shutting down in-progress monitor")
				ticker.Stop()
				return
			}
		}
	}()
}

func (m monitor) runOnce() error {
	environments, err := m.environment.ListEnvironments(m.ctx)
	if err != nil {
		return errors.New("Could not retrieve environments")
	}

	for _, environment := range environments {
		m.events <- UpdateInProgressDeploymentEvent{
			Environment: environment,
		}
	}

	return nil
}
