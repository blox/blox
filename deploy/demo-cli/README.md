# Blox Demo CLI Tools

The following CLI tools are included in the `<GitRepoBase>/deploy/demo-cli/` directory for interacting with the Blox daemon-scheduler API.

## List Blox Environments

```
$ ./blox-list-environments.py --help
== Blox Demo CLI - List Blox Environments ==

usage: blox-list-environments.py [-h] [--region REGION] [--apigateway]
                                 [--host HOST]

List Blox Environments

optional arguments:
  -h, --help       show this help message and exit
  --region REGION  AWS region
  --apigateway     Call API Gateway endpoint
  --host HOST      Blox Scheduler <Host>:<Port>
```

## Create Blox Environment

```
$ ./blox-create-environment.py --help
== Blox Demo CLI - Create Blox Environment ==

usage: blox-create-environment.py [-h] [--region REGION] [--apigateway]
                                  [--host HOST] [--environment ENVIRONMENT]
                                  [--cluster CLUSTER]
                                  [--task-definition TASKDEF]

Create Blox Environment

optional arguments:
  -h, --help            show this help message and exit
  --region REGION       AWS region
  --apigateway          Call API Gateway endpoint
  --host HOST           Blox Scheduler <Host>:<Port>
  --environment ENVIRONMENT
                        Blox environment name
  --cluster CLUSTER     ECS cluster name
  --task-definition TASKDEF
                        ECS task definition arn
```

## List Blox Deployments

```
$ ./blox-list-deployments.py --help
== Blox Demo CLI - List Blox Deployments ==

usage: blox-list-deployments.py [-h] [--region REGION] [--apigateway]
                                [--host HOST] [--environment ENVIRONMENT]

List Blox Deployments

optional arguments:
  -h, --help            show this help message and exit
  --region REGION       AWS region
  --apigateway          Call API Gateway endpoint
  --host HOST           Blox Scheduler <Host>:<Port>
  --environment ENVIRONMENT
                        Blox environment name
```

## Create Blox Deployment

```
$ ./blox-create-deployment.py --help
== Blox Demo CLI - Create Blox Deployment ==

usage: blox-create-deployment.py [-h] [--region REGION] [--apigateway]
                                 [--host HOST] [--environment ENVIRONMENT]

Create Blox Deployment

optional arguments:
  -h, --help            show this help message and exit
  --region REGION       AWS region
  --apigateway          Call API Gateway endpoint
  --host HOST           Blox Scheduler <Host>:<Port>
  --environment ENVIRONMENT
                        Blox environment name
```

## List cluster-state-service Instances

```
$ ./css-list-instances.py --help
== Blox Demo CLI - List Blox Instances ==

usage: css-list-instances.py [-h] [--region REGION] [--host HOST]
                             [--cluster CLUSTER] [--status STATUS]
                             [--instance-arn INSTANCE]

List Blox Instances

optional arguments:
  -h, --help            show this help message and exit
  --region REGION       AWS region
  --host HOST           Blox CSS <Host>:<Port>
  --cluster CLUSTER     ECS cluster name
  --status STATUS       EC2 instance status
  --instance-arn INSTANCE
                        EC2 instance Arn
```

## List cluster-state-service Tasks

```
$ ./css-list-tasks.py --help
== Blox Demo CLI - List Blox Tasks ==

usage: css-list-tasks.py [-h] [--region REGION] [--host HOST]
                         [--cluster CLUSTER] [--status STATUS]
                         [--task-arn TASK]

List Blox Tasks

optional arguments:
  -h, --help         show this help message and exit
  --region REGION    AWS region
  --host HOST        Blox CSS <Host>:<Port>
  --cluster CLUSTER  ECS cluster name
  --status STATUS    ECS task status
  --task-arn TASK    ECS task Arn
```

## List ECS Task Definitions

```
$ ./list-task-definitions.py --help
== Blox Demo CLI - List Task Definitions ==

usage: list-task-definitions.py [-h] [--region REGION]

List Task Definitions

optional arguments:
  -h, --help       show this help message and exit
  --region REGION  AWS region
```

## Register ECS Task Definition

```
$ ./register-task-definition.py --help
== Blox Demo CLI - Register Task Definition ==

usage: register-task-definition.py [-h] [--region REGION] [--file FILE]

Register Task Definition

optional arguments:
  -h, --help       show this help message and exit
  --region REGION  AWS region
  --file FILE      path to task definition file
```

## Increment ECS Cluster Instances

```
$ ./increment-cluster-instances.py --help
== Blox Demo CLI - Increment Cluster Instances ==

usage: increment-cluster-instances.py [-h] [--region REGION]
                                      [--cluster CLUSTER] [--num NUM]

Increment Cluster Instances

optional arguments:
  -h, --help         show this help message and exit
  --region REGION    AWS region
  --cluster CLUSTER  ECS cluster name
  --num NUM          number of instances to increment by
```
