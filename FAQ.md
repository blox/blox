# Frequently Asked Questions

##### What is a Blox scheduler?
Blox schedulers are responsible for making the decision of “what” to run and “when” to run a container or set of containers based on a customers specification of tasks and services. The scheduling decision may be based on: task or service health, event response, schedule, or other custom business logic. The Blox scheduler is enabled to take advantage of the Amazon ECS task placement options, such as: availability spread, binpacking, attribute-based placement etc.  

##### What is a daemon scheduler?
A daemon scheduler runs one-task-per node on set of nodes. As new nodes join the cluster, the daemon scheduler identifies the new node and will automatically start a task on the node. The daemon scheduler also monitors the daemons running on every node and can restart the task as required. The most common use cases for this scheduler is to run logging or monitoring agents.    

##### How can I use Blox?
You can use Blox’s scheduling capabilities using Amazon ECS console or using the CLI or the SDK. The workflow will be very similar to the current workflow of using Amazon ECS.

##### What are the goals of Blox?
* Deliver a highly available, scalable, and performant set of schedulers for customers;
* Deliver schedulers as open source to enable additional extensibility and customization;
* Deliver schedulers as a managed service so customers do not need to manage or host critical infrastructure on their own;
* Enable customers to build and run the same schedulers that power Amazon ECS;
* Encourage transparent product feedback and community contributions.

##### Is there any cost to using Blox?
If using the managed Blox service, there is no cost as Blox is provided as part of Amazon ECS.  

##### What is the goal of Blox v1.0 compared to Blox v0.3?  
Customers wanted to see greater convergence across the open source Blox project and Amazon ECS schedulers, and asked for a managed or hosted option for Blox. Therefore, we are providing Blox as a managed option and making it available for use with ECS without having to deploy, run and manage it.

The code for v0.3 lives in the v0.3 branch. See [pull request](https://github.com/blox/blox/pull/214) for details.

##### How can I extend Blox for my specific scheduling needs?   
We encourage contributions to Blox both in the form of features requests (see [Issues](https://github.com/blox/blox/issues)) and pull requests. All projects under Blox are released under Apache 2.0 and contributions are accepted under individual Apache Contributor Agreements.

##### What is the status of the Blox 0.3.0 release?
Blox v0.3.0 is now deprecated. If you are actively using v0.3.0, please [email us](mailto:ecs-blox-team@amazon.com) to work out a migration plan. If you have a critical issue in v0.3.0, please open a [new issue](https://github.com/blox/blox/issues/new). We are currently working on and prioritizing the new v1.0 release that incorporates the feedback we received from customers using v0.3.0 pre-release. For more details on the Blox v1.0 release, see [here](https://github.com/blox/blox/blob/dev/docs/daemon_design.md)
