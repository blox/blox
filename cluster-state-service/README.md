
# cluster-state-service

### Description

The cluster-state-service consumes events from a stream of all changes to containers and instances across your Amazon ECS clusters, persists the events in a local data store, and provides APIs (e.g., search, filter, list, etc.) that enable you to query the state of your cluster so you can respond to changes in real-time. The cluster-state-service utilizes etcd as the data store to track your Amazon ECS cluster state locally, and it also manages any drift in state by periodically reconciling state with Amazon ECS.  

### REST API

The cluster-state-service API operations:  
*	Lists and describes container instances and tasks
*	Filters container instances and tasks by status or cluster
*	Listens to streaming container instance and task state changes

### Building cluster-state-service

The cluster-state-service depends on golang and go-swagger. Install and configure [golang](https://golang.org/doc/). For more information about installing go-swagger, see the [go-swagger documentation](https://github.com/go-swagger/go-swagger).

```
$ git clone https://github.com/blox/blox.git blox/blox
$ cd blox/blox/cluster-state-service
$ make get-deps
$ make

# Find the cluster-state-service binary in 'out' folder
$ ls out/
LICENSE                 cluster-state-service

```

### Usage

We provide an AWS CloudFormation template to set up the necessary prerequisites for the cluster-state-service. After the prerequisites are ready, you can launch the cluster-state-service via the Docker compose file, if you prefer. For more information, see the Blox [Deployment Guide](../deploy).

To launch the cluster-state-service manually, use the following steps.

#### Prerequisites

In order to use the cluster-state-service, you need to set up an Amazon SQS queue, configure CloudWatch Events, and add the queue as a target for ECS events.

The cluster-state-service also depends on etcd to store the cluster state locally. To set up etcd manually, see the [etcd documentation](https://github.com/coreos/etcd).

#### Quick Start - Launching the cluster-state-service

The cluster-state-service is provided as a Docker image for your convenience. You can launch it with the following code. Use appropriate values for AWS_REGION, etcd IP, and port and queue names.

```
docker run -e AWS_REGION=us-west-2 \
    AWS_PROFILE=default \
    -v ~/.aws:/.aws \
    -v /tmp/css-logs:/var/output/logs \
    bloxoss/cluster-state-service:0.1.0 \
    --etcd-endpoint $ETCD_IP:$ETCD_PORT \
    --queue_name $SQS_QUEUE_NAME
```

You can also override the logger configuration like the log file and log level.

```
docker run -e AWS_REGION=us-west-2 \
AWS_PROFILE=default \
    CSS_LOG_FILE=/var/output/logs/css.log \
    CSS_LOG_LEVEL=info \
    -v ~/.aws:/.aws \
    -v /tmp/css-logs:/var/output/logs \
    bloxoss/cluster-state-service:0.1.0 \
    --etcd-endpoint $ETCD_IP:$ETCD_PORT \
    --queue event_stream
```

#### API endpoint

After you launch the cluster-state-service, you can interact with and use the REST API by using the endpoint at port 3000. Identify the cluster-state-service container IP address and connect to port 3000. For more information about the API definitions, see the [swagger specification](swagger/v1/swagger.json).
