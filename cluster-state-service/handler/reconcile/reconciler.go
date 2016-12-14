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

package reconcile

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/blox/blox/cluster-state-service/handler/reconcile/loader"
	"github.com/blox/blox/cluster-state-service/handler/store"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

// ReconcileDuration specifies the interval between each reconcile loop
const ReconcileDuration = 20 * time.Minute

type Reconciler struct {
	taskLoader     loader.TaskLoader
	instanceLoader loader.ContainerInstanceLoader
	ticker         *time.Ticker
	tickerDuration time.Duration
	ctx            context.Context
	inProgress     bool
	inProgressLock sync.RWMutex
}

func NewReconciler(ctx context.Context, stores store.Stores, ecsClient *ecs.ECS, tickerDuration time.Duration) (*Reconciler, error) {
	var reconciler *Reconciler
	if ecsClient == nil {
		return reconciler, errors.New("Failed to initialize Reconciler. ECS client is not initialized.")
	}
	if tickerDuration <= 0 {
		return reconciler, fmt.Errorf("Invalid duration specified for running the reconciler: %s", tickerDuration.String())
	}
	return &Reconciler{
		taskLoader:     loader.NewTaskLoader(stores.TaskStore, ecsClient),
		instanceLoader: loader.NewContainerInstanceLoader(stores.ContainerInstanceStore, ecsClient),
		tickerDuration: tickerDuration,
		ctx:            ctx,
		inProgress:     false,
	}, nil
}

func (reconciler *Reconciler) Run() {
	reconciler.initTicker()
	for {
		select {
		case <-reconciler.ticker.C:
			if reconciler.isInProgress() {
				log.Info("Reconcile loop in progress, skipping")
				continue
			}
			go func() {
				err := reconciler.RunOnce()
				if err != nil {
					log.Warnf("Error reconciling: %v", err)
				}
			}()
		case <-reconciler.ctx.Done():
			reconciler.ticker.Stop()
			return
		}
	}
}

// RunOnce loads all existing ECS tasks and instances into the datastore
func (reconciler *Reconciler) RunOnce() error {
	reconciler.setInProgress(true)
	defer reconciler.setInProgress(false)

	log.Infof("Reconciler loading tasks and instances")
	// TODO: Pass in context everywhere so that cancelling the context cancels any outstanding
	// requests as well
	err := reconciler.taskLoader.LoadTasks()
	if err != nil {
		return errors.Wrapf(err, "Failed to reconcile. Could not load tasks.")
	}

	err = reconciler.instanceLoader.LoadContainerInstances()
	if err != nil {
		return errors.Wrapf(err, "Failed to reconcile. Could not load container instances.")
	}
	return nil
}

func (reconciler *Reconciler) setInProgress(val bool) {
	reconciler.inProgressLock.Lock()
	defer reconciler.inProgressLock.Unlock()

	reconciler.inProgress = val
}

func (reconciler *Reconciler) isInProgress() bool {
	reconciler.inProgressLock.RLock()
	defer reconciler.inProgressLock.RUnlock()

	return reconciler.inProgress
}

func (reconciler *Reconciler) initTicker() {
	if reconciler.ticker == nil {
		reconciler.ticker = time.NewTicker(reconciler.tickerDuration)
	}
}
