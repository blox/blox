# Blox

### Description
Blox is a collection of open source projects for container management and orchestration. Blox gives you more control over how your containerized applications run on Amazon ECS. It enables you to build custom schedulers and integrate third-party schedulers on top of ECS, all the while leveraging ECS to fully manage and scale your clusters

Blox currently consists of two components:  
* Cluster State Service
* Daemon Scheduler

The Cluster State Service provides a local materialized view of the ECS cluster state by consuming the ECS event stream. The ECS event stream provides the ability to listen to cluster state changes in near real-time and is delivered via CloudWatch events. Customers who want to build scheduling workflows often need to consume the events generated in the ECS cluster, persist this state locally, and operate on the local cluster state. Cluster State Service implements this functionality and provides APIs (e.g., search, filter, list, etc.) that enable you to query the state of your cluster so you can respond to changes in real-time. The Cluster State Service utilizes etcd as the data store to track your Amazon ECS cluster state locally, and it also manages any drift in state by periodically reconciling state with Amazon ECS.

The Daemon Scheduler allows you to run exactly one task per host across all nodes in a cluster. The Daemon Scheduler monitors the cluster state and launches tasks as new nodes join the cluster. It is ideal for running monitoring agents, log collectors, etc. The Daemon Scheduler can be used a reference for how to use the Cluster State Service to build custom scheduling logic, and we plan to add more scheduling logic for additional use cases.

### Interested in learning more?

If you are interested in learning more, see the [cluster-state-service](cluster-state-service) and [daemon-scheduler](daemon-scheduler).

### Deploying Blox

We provide two methods for deploying Blox:  
* Local deployment
* AWS deployment

#### Local Deployment

You can deploy locally using our Docker compose file in order to try out Blox quickly, or using the cluster-state-service during development of custom local schedulers. Local deployment launches the Blox components — the Cluster State Service and Daemon Scheduler, along with a backing state store, etcd. Please see [Blox Deployment Guide](deploy) to launch Blox using the docker compose file.

#### AWS Deployment

We also provide an AWS CloudFormation template to launch the Blox stack easily on AWS. The AWS-deployed Blox stack makes use of AWS services designed to provide a secure public facing Daemon Scheduler endpoint.

##### Creating a Blox stack in AWS

Deploying Blox using the AWS CloudFormation template in AWS sets up a stack consisting of the following components:
* An SQS queue is created and CloudWatch is configured to deliver ECS events to the queue.
* Blox components are set up as a service running on an Amazon ECS cluster. The Cluster State Service, Daemon Scheduler, and etcd containers making up a single task are run on a container instance. The Daemon Scheduler endpoint is then made reachable for customers to interact with securely.
* An Application Load Balancer (ALB) instance is created in front of the scheduler endpoint.
* An Amazon API Gateway endpoint is set up as the public facing front end to Blox and provides the authentication mechanism. This API Gateway can be used to reach the scheduler and manage tasks on the ECS cluster.
* An AWS Lambda function that acts as a simple proxy to enable the public-facing API gateway endpoint to forward requests onto the ELB listener in the VPC.

For more information about deployment instructions, see [Blox Deployment Guide](deploy).

### Building Blox

For more information about how to build these components, see [cluster-state-service](cluster-state-service) and [daemon-scheduler](daemon-scheduler).

### Contributing to Blox

All projects under Blox are released under Apache 2.0 and the usual Apache Contributor Agreements apply for individual contributors. All projects are maintained in public on GitHub, issues and pull requests use GitHub, and discussions use our Gitter channel. We look forward to collaborating with the community.
