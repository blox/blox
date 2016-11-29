# Daemon scheduler

### Description

The Daemon Scheduler ensures that one copy of a task is run on every specified instance in an Amazon ECS cluster. It also monitors the cluster state and launches the tasks on any new nodes joining the cluster. The Daemon Scheduler also restarts any of the stopped or failed tasks.

### Concepts

The Daemon Scheduler defines and depends on the following concepts:  

* An `Environment` represents the configuration for desired state of the tasks to be maintained. For the Daemon Scheduler, the environment indicates the task definition to launch in a specific cluster.
* `Deployment` is the operation that brings the environment into existence. Deployment indicates to the scheduler that the desired configuration state in `Environment` should be established in the cluster.

### REST API

The Daemon Scheduler API:  
* Creates and lists environments
* Creates and lists deployments

### Building the Daemon Scheduler

The Daemon Scheduler depends on golang and go-swagger. Install and configure [golang](https://golang.org/doc/). For more information about installing go-swagger, see the [go-swagger documentation](https://github.com/go-swagger/go-swagger).

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

The Daemon Scheduler depends on the Cluster State Service. We provide an AWS CloudFormation template to set up the necessary prerequisites for the Cluster state service. After the prerequisites are ready, you can launch the Daemon Scheduler via the Docker compose file. For more information, see the Blox [Deployment Guide](../deploy).

To launch the Daemon Scheduler manually, use the following steps.

#### Quick Start - Launching Daemon Scheduler

The Daemon Scheduler is provided as a Docker image for your convenience. You can launch it using the following command. Use the appropriate values for AWS_REGION, AWS_PROFILE, etcd IP address and port, and the Cluster State Service IP address and port.

```
docker run -e AWS_REGION=us-west-2 \
    AWS_PROFILE=default \
    -v ~/.aws:/.aws \
    -v /tmp/ds-logs:/var/output/logs \
    bloxoss/daemon-scheduler:0.1.0 \
    --etcd-endpoint $ETCD_IP:$ETCD_PORT \
    --css-endpoint $CSS_IP:$CS_PORT
```

#### API endpoint

After you launch the Daemon Scheduler, you can interact with and use the REST API by using the endpoint at port 2000. Identify the Daemon Scheduler container IP address and connect to port 2000. For more information about the API definitions, see the [swagger specification](generated/v1/swagger.json).
