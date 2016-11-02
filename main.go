package main

import (
	"fmt"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/api/v1"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/clients"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/event"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/store"
	"github.com/aws/amazon-ecs-event-stream-handler/logger"

	log "github.com/cihub/seelog"
	"github.com/urfave/negroni"
	"golang.org/x/net/context"

	"net/http"
	"os"
	"time"
)

const (
	serverAddress     = "localhost:3000"
	serverReadTimeout = 10 * time.Second
	errorCode         = 1
)

func main() {
	os.Exit(_main())
}

func _main() int {
	defer log.Flush()
	err := logger.InitLogger()
	if err != nil {
		fmt.Printf("Could not initialize logger: %+v", err)
	}

	etcdClient, err := clients.NewEtcdClient()
	if err != nil {
		log.Criticalf("Could not start etcd: %+v", err)
		return errorCode
	}
	defer etcdClient.Close()

	// initialize the datastore
	datastore, err := store.NewDataStore(etcdClient)
	if err != nil {
		log.Criticalf("Could not initialize the datastore: %+v", err)
		return errorCode
	}

	// initialize services
	stores, err := store.NewStores(datastore)
	if err != nil {
		log.Criticalf("Could not initialize stores: %+v", err)
		return errorCode
	}

	// initialize apis
	apis := v1.NewAPIs(stores)

	// start event processor
	processor := event.NewProcessor(stores)

	awsSession, err := clients.NewAWSSession()
	if err != nil {
		log.Criticalf("Could not load aws session: %+v", err)
		panic(err)
	}

	sqsClient := clients.NewSQSClient(awsSession)

	// start event consumer
	consumer, err := event.NewConsumer(sqsClient, processor)
	if err != nil {
		log.Criticalf("Could not start the consumer: %+v", err)
		return errorCode
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

	err = s.ListenAndServe()
	if err != nil {
		log.Criticalf("Could not start the server: %+v", err)
		return errorCode
	}

	return 0
}
