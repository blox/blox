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

package run

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/api/v1"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/clients"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/event"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/store"
	"github.com/urfave/negroni"
)

const (
	serverAddress     = "localhost:3000"
	serverReadTimeout = 10 * time.Second
)

// StartEventStreamHandler starts the event stream handler. It creates an ETCD
// client, a data store using this client and an event processor to process
// events from the SQS queue. It also starts the RESTful server and blocks on
// the listen method of the same to listen to requests that query for task and
// instance state from the store.
func StartEventStreamHandler(queueName string) error {
	etcdClient, err := clients.NewEtcdClient()
	if err != nil {
		return fmt.Errorf("Could not start etcd: %+v", err)
	}
	defer etcdClient.Close()

	// initialize the datastore
	datastore, err := store.NewDataStore(etcdClient)
	if err != nil {
		return fmt.Errorf("Could not initialize the datastore: %+v", err)
	}

	// initialize services
	stores, err := store.NewStores(datastore)
	if err != nil {
		return fmt.Errorf("Could not initialize stores: %+v", err)
	}

	// initialize apis
	apis := v1.NewAPIs(stores)

	// start event processor
	processor := event.NewProcessor(stores)

	awsSession, err := clients.NewAWSSession()
	if err != nil {
		return fmt.Errorf("Could not load aws session: %+v", err)
	}

	sqsClient := clients.NewSQSClient(awsSession)

	// start event consumer
	consumer, err := event.NewConsumer(sqsClient, processor, queueName)
	if err != nil {
		return fmt.Errorf("Could not start the consumer: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go consumer.PollForEvents(ctx)
	defer cancel()

	// start server
	router := v1.NewRouter(apis)

	n := negroni.Classic()
	n.UseHandler(router)

	s := &http.Server{
		Addr:        serverAddress,
		Handler:     n,
		ReadTimeout: serverReadTimeout,
	}

	return s.ListenAndServe()
}
