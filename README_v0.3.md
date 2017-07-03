# ![Logo](blox-logo.png)


[![Build Status](https://travis-ci.org/blox/blox.svg?branch=master)](https://travis-ci.org/blox/blox)

### Description
Blox is a collection of open source projects for container management and orchestration. Blox gives you more control over how your containerized applications run on Amazon ECS. It enables you to build schedulers and integrate third-party schedulers on top of ECS, while leveraging Amazon ECS to fully manage and scale your clusters.

The *blox* project provides a scheduling framework to help you easily build custom tooling, such as schedulers, on top of Amazon ECS. The framework makes it easy to consume events from Amazon ECS, store the cluster state locally, and query the local data store though APIs. The *blox* project currently consists of two components:  

* *cluster-state-service*
* *daemon-scheduler*

The *cluster-state-service* consumes events from a stream of all changes to containers and instances across your Amazon ECS clusters, persists the events in a local data store, and provides APIs (e.g., search, filter, list, etc.) that enable you to query the state of your cluster so you can respond to changes in real-time. The *cluster-state-service* tracks your Amazon ECS cluster state locally, and manages any drift in state by periodically reconciling state with Amazon ECS.

The *daemon-scheduler* is a scheduler that allows you to run exactly one task per host across all nodes in a cluster. The scheduler monitors the cluster state and launches tasks as new nodes join the cluster, and it is ideal for running monitoring agents, log collectors, etc. The scheduler can be used as a reference for how to use the *cluster-state-service* to build custom scheduling logic, and we plan to add additional scheduling capabilities for different use cases.


### Interested in learning more?

If you are interested in learning more about the components, please read the [cluster-state-service](cluster-state-service) and [daemon-scheduler](daemon-scheduler) README files.

### Deploying Blox

We provide two methods for deploying *blox*:  
* Local deployment
* AWS deployment

#### Local Deployment

You can deploy locally and quickly try out Blox using our Docker Compose file. This allows you to get started with building custom schedulers using the cluster-state-service. The Docker Compose file launches the *blox* components, *cluster-state-service* and *daemon-scheduler*, along with a backing state store, etcd. Please see [Blox Deployment Guide](deploy) to launch *blox* using the Docker Compose file.

#### AWS Deployment

We also provide an AWS CloudFormation template to launch the *blox* stack easily on AWS. The AWS deployed *blox* stack makes use of AWS services designed to provide a secure public facing scheduler endpoint.

##### Creating a Blox stack on AWS

Deploying Blox using the AWS CloudFormation template in AWS sets up a stack consisting of the following components:
* An Amazon SQS queue is created for you, and Amazon CloudWatch is configured to deliver ECS events to your queue.
* *blox* components are set up as a service running on an Amazon ECS cluster. The *cluster-state-service*, *daemon-scheduler* , and etcd containers run as a single task on a container instance. The scheduler endpoint is made reachable, which allows you to securely interact with the endpoint.
* An Application Load Balancer (ALB) is created in front of your scheduler endpoint.
* An Amazon API Gateway endpoint is set up as the public facing frontend and provides an authentication mechanism for the *blox* stack. The API Gateway endpoint can be used to reach the scheduler and manage tasks on the ECS cluster.
* An AWS Lambda function acts as a simple proxy which enables the public facing API Gateway endpoint to forward requests onto the ALB listener in the VPC.

For more information about deployment instructions, see [Blox Deployment Guide](deploy).

### Building Blox

For more information about how to build these components, see [cluster-state-service](cluster-state-service) and [daemon-scheduler](daemon-scheduler).

### Contributing to Blox

All projects under Blox are released under Apache 2.0 and the usual Apache Contributor Agreements apply for individual contributors. All projects are maintained in public on GitHub, issues and pull requests use GitHub, and discussions use our [Gitter channel](https://gitter.im/blox). We look forward to collaborating with the community.
