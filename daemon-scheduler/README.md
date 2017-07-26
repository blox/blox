# daemon-scheduler

### Description

The daemon-scheduler allows you to run exactly one task per host across all nodes in a cluster. It monitors the cluster state and launches tasks as new nodes join the cluster, and it is ideal for running monitoring agents, log collectors, etc. The daemon-scheduler can be used a reference for how to use the cluster-state-service to build custom scheduling logic.

### Concepts

The daemon-scheduler defines and depends on the following concepts:  

* An `Environment` represents the configuration for desired state of the tasks to be maintained. For the daemon-scheduler, the environment indicates the task definition to launch in a specific cluster.
* `Deployment` is the operation that brings the environment into existence. Deployment indicates to the scheduler that the desired configuration state in `Environment` should be established in the cluster.

### REST API

The daemon-scheduler API:  
* Creates and lists environments
* Creates and lists deployments

### Building the daemon-scheduler

The daemon-scheduler depends on golang and go-swagger. Install and configure [golang](https://golang.org/doc/). For more information about installing go-swagger, see the [go-swagger documentation](https://github.com/go-swagger/go-swagger). Also, make sure to clone this repo into your appropriate [$GOPATH](https://golang.org/doc/code.html#GOPATH) directory.

```
$ git clone https://github.com/blox/blox.git blox/blox
$ cd blox/blox/daemon-scheduler
$ make get-deps
$ make

# Find the daemon-scheduler binary in 'out' folder
$ ls out/
LICENSE                 daemon-scheduler

```

### Usage

The daemon-scheduler depends on the cluster-state-service. We provide an AWS CloudFormation template to set up the necessary prerequisites for the cluster-state-service. After the prerequisites are ready, you can launch the daemon-scheduler via the Docker compose file. For more information, see the Blox [Deployment Guide](../deploy).

To launch the daemon-scheduler manually, use the following steps.

#### Quick Start - Launching daemon-scheduler

The daemon-scheduler is provided as a Docker image for your convenience. You can launch it using the following command. Use the appropriate values for AWS_REGION, AWS_PROFILE, etcd IP address and port, and the cluster-state-service IP address and port.

```
docker run -e AWS_REGION=us-west-2 \
    AWS_PROFILE=default \
    -v ~/.aws:/.aws \
    -v /tmp/ds-logs:/var/output/logs \
    bloxoss/daemon-scheduler:0.3.0 \
    --etcd-endpoint $ETCD_IP:$ETCD_PORT \
    --css-endpoint $CSS_IP:$CS_PORT
```

#### API endpoint

After you launch the daemon-scheduler, you can interact with and use the REST API by using the endpoint at port 2000. Identify the daemon-scheduler container IP address and connect to port 2000. For more information about the API definitions, see the [swagger specification](swagger/v1/swagger.json).
