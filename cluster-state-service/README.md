
# Cluster State Service

### Description

The Cluster State Service consumes the Amazon ECS event stream and builds the ECS cluster state locally. It also handles state reconciliation with ECS in order to handle scenarios where the events could be lost. The Cluster State Service provides simple REST API operations in order to operate on the local view of the cluster state.

### REST API

The Cluster State Service API operations:  
*	Lists and describes container instances and tasks
*	Filters container instances and tasks by status or cluster
*	Listens to streaming container instance and task state changes

### Building Cluster State Service

The Cluster State Service depends on golang and go-swagger. Install and configure [golang](https://golang.org/doc/). For more information about installing go-swagger, see the [go-swagger documentation](https://github.com/go-swagger/go-swagger).

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

We provide an AWS CloudFormation template to set up the necessary prerequisites for the Cluster State Service. After the prerequisites are ready, you can launch the Cluster State Service via the Docker compose file, if you prefer. For more information, see the Blox [Deployment Guide](../deploy).

To launch the Cluster State Service manually, use the following steps.

#### Prerequisites

In order to use the Cluster State Service, you need to set up an Amazon SQS queue, configure CloudWatch Events, and add the queue as a target for ECS events.

The Cluster State Service also depends on etcd to store the cluster state locally. To set up etcd manually, see the [etcd documentation](https://github.com/coreos/etcd).

#### Quick Start - Launching the Cluster State Service

The Cluster State Service is provided as a Docker image for your convenience. You can launch it with the following code. Use appropriate values for AWS_REGION, etcd IP, and port and queue names.

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

After you launch the Cluster State Service, you can interact with and use the REST API by using the endpoint at port 3000. Identify the Cluster State Service container IP address and connect to port 3000. For more information about the API definitions, see the [swagger specification](handler/api/v1/swagger/swagger.json).
