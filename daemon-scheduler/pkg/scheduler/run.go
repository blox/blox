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

package scheduler

import (
	"context"

	"github.com/blox/blox/daemon-scheduler/pkg/api/v1"
	"github.com/blox/blox/daemon-scheduler/pkg/clients"
	"github.com/blox/blox/daemon-scheduler/pkg/config"
	"github.com/blox/blox/daemon-scheduler/pkg/deployment"
	"github.com/blox/blox/daemon-scheduler/pkg/engine"
	"github.com/blox/blox/daemon-scheduler/pkg/facade"
	"github.com/blox/blox/daemon-scheduler/pkg/store"
	log "github.com/cihub/seelog"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/pkg/errors"
	"github.com/urfave/negroni"

	"net/http"
	"time"
)

const (
	serverReadTimeout  = 10 * time.Second
	serverWriteTimeout = 10 * time.Second
)

// Run kickstarts the daemon scheduler service.
func Run(schedulerBindAddr string, clusterStateServiceEndpoint string) error {
	if schedulerBindAddr == "" {
		return errors.Errorf("The address for scheduler endpoint is not set")
	}
	if clusterStateServiceEndpoint == "" {
		return errors.Errorf("The address for cluster state service endpoint is not set")
	}

	etcdClient, err := clients.NewEtcdClient(config.EtcdEndpoints)
	if err != nil {
		log.Criticalf("Could not start etcd: %+v", err)
		return err
	}
	defer etcdClient.Close()

	// initialize the datastore
	datastore, err := store.NewDataStore(etcdClient)
	if err != nil {
		log.Criticalf("Could not initialize the datastore: %+v", err)
		return err
	}

	environmentStore, err := store.NewEnvironmentStore(datastore)
	if err != nil {
		log.Criticalf("Could not initialize the environment store: %+v", err)
		return err
	}

	ecsClient, err := clients.NewECSClient()
	if err != nil {
		log.Criticalf("Could not initialize ecs client: %+v", err)
		return err
	}

	cssClient := clients.NewCSSClient()
	cssTransport := httptransport.New(clusterStateServiceEndpoint, "/v1", []string{"http"})
	cssClient.SetTransport(cssTransport)

	ecs := facade.NewECS(ecsClient)
	css, err := facade.NewClusterState(cssClient)
	if err != nil {
		log.Criticalf("Could not initialize cluster state: %+v", err)
		return err
	}

	environment, err := deployment.NewEnvironment(environmentStore)
	if err != nil {
		log.Criticalf("Could not initialize environment: %+v", err)
		return err
	}

	deploymentWorker := deployment.NewDeploymentWorker(environment, ecs, css)
	deployment := deployment.NewDeployment(environment, css, ecs)

	ctx := context.Background()
	input := make(chan engine.Event)
	output := make(chan engine.Event)
	dispatcher := engine.NewDispatcher(ctx, environment, deployment, ecs, css, deploymentWorker, input, output)
	dispatcher.Start()
	scheduler := engine.NewScheduler(ctx, input, environment, deployment, css, ecs)
	scheduler.Start()

	monitor := engine.NewMonitor(ctx, environment, input)
	monitor.InProgressMonitorLoop()

	api := v1.NewAPI(environment, deployment, ecs)

	// start server
	router := v1.NewRouter(api)

	n := negroni.Classic()
	n.UseHandler(router)

	s := &http.Server{
		Addr:         schedulerBindAddr,
		Handler:      n,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
	}

	err = s.ListenAndServe()
	if err != nil {
		log.Criticalf("Could not start the server: %+v", err)
	}

	return err
}
