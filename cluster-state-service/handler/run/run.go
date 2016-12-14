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

package run

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/cihub/seelog"
	"github.com/pkg/errors"

	"github.com/blox/blox/cluster-state-service/handler/api/v1"
	"github.com/blox/blox/cluster-state-service/handler/clients"
	"github.com/blox/blox/cluster-state-service/handler/event"
	"github.com/blox/blox/cluster-state-service/handler/reconcile"
	"github.com/blox/blox/cluster-state-service/handler/store"
	"github.com/urfave/negroni"
	"strings"
)

const (
	serverReadTimeout = 10 * time.Second
	kinesisPrefix     = "kinesis://"
	sqsPrefix         = "sqs://"
)

// StartClusterStateService starts the Cluster State Service. It creates an ETCD
// client, a data store using this client and an event processor to process
// events from the provided queue. It also starts the RESTful server and blocks on
// the listen method of the same to listen to requests that query for task and
// instance state from the store.
func StartClusterStateService(queueNameURI string, bindAddr string, etcdEndpoints []string) error {
	if bindAddr == "" {
		return fmt.Errorf("The cluster state service listen address is not set")
	}

	etcdClient, err := clients.NewEtcdClient(etcdEndpoints)
	if err != nil {
		return errors.Wrapf(err, "Could not start etcd")
	}
	defer etcdClient.Close()

	// initialize the datastore
	datastore, err := store.NewDataStore(etcdClient)
	if err != nil {
		return errors.Wrapf(err, "Could not initialize the datastore")
	}

	etcdTXStore, err := store.NewEtcdTXStore(etcdClient)
	if err != nil {
		return errors.Wrapf(err, "Could not initialize the etcd transactional store")
	}

	// initialize services
	stores, err := store.NewStores(datastore, etcdTXStore)
	if err != nil {
		return errors.Wrapf(err, "Could not initialize stores")
	}

	awsSession, err := clients.NewAWSSession()
	if err != nil {
		return errors.Wrapf(err, "Could not load aws session")
	}

	ecsClient := clients.NewECSClient(awsSession)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	recon, err := reconcile.NewReconciler(ctx, stores, ecsClient, reconcile.ReconcileDuration)
	if err != nil {
		return errors.Wrapf(err, "Could not start reconciler")
	}
	err = recon.RunOnce()
	if err != nil {
		return errors.Wrapf(err, "Error bootstrapping")
	}
	log.Infof("Bootstrapping completed")
	go recon.Run()

	// initialize apis
	apis := v1.NewAPIs(stores)

	// start event processor
	processor := event.NewProcessor(stores)

	if strings.HasPrefix(queueNameURI, kinesisPrefix) {
		kinesisClient := clients.NewKinesisClient(awsSession)

		// start event consumer
		consumer, err := event.NewKinesisConsumer(kinesisClient, processor, strings.TrimPrefix(queueNameURI, kinesisPrefix))
		if err != nil {
			return errors.Wrapf(err, "Could not start the consumer")
		}

		go consumer.PollForEvents(ctx)
	} else {
		sqsClient := clients.NewSQSClient(awsSession)

		// start event consumer
		consumer, err := event.NewSQSConsumer(sqsClient, processor, strings.TrimPrefix(queueNameURI, sqsPrefix))
		if err != nil {
			return errors.Wrapf(err, "Could not start the consumer")
		}

		go consumer.PollForEvents(ctx)
	}

	// start server
	router := v1.NewRouter(apis)

	n := negroni.Classic()
	n.UseHandler(router)

	s := &http.Server{
		Addr:        bindAddr,
		Handler:     n,
		ReadTimeout: serverReadTimeout,
	}

	return s.ListenAndServe()
}
