# Blox Installation

#### Introduction

Blox is an open source cluster manager and orchestration framework that enables developers to easily build, test, and run containerized applications locally and in production on Amazon ECS using a common toolset. This document describes how to install the Blox framework locally and on top of Amazon AWS.

#### Installation Instructions

- [Install Blox on Local Docker](#local-installation)
- [Install Blox on Amazon AWS](#amazon-aws-installation)

#### Framework Components

- Daemon Scheduler
- Cluster State Service
- Etcd

#### Required Amazon AWS Components

- Amazon CloudWatch
- Amazon ECS
- Amazon IAM
- Amazon SQS

#### Additional AWS Components (When Installing Blox on AWS)

- Amazon API Gateway
- Amazon Application Load Balancer
- Amazon Lambda
- Amazon VPC

## Prerequisites

#### Amazon AWS CLI

You will need to have the Amazon AWS CLI installed locally to create the required Amazon AWS components. Follow the instructions at [Installing the AWS Command Line Interface](http://docs.aws.amazon.com/cli/latest/userguide/installing.html) to install the CLI before proceeding. 

#### IAM Permissions

The AWS profile you use with the AWS CLI will need appropriate permissions to create the required Amazon AWS components. To make this easier, we have created IAM Policy Documents for both the Local and Amazon AWS installations. Follow the steps on the [Creating a New IAM Policy](http://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_create.html) guide to create this policy. In the `Create Policy` wizard, select `Create Your Own Policy` and paste the contents of the appropriate policy file in the `Policy Document` text area. Once you have created the IAM Policy, you will need to attach the policy to the AWS user that you will use to deploy Blox via the installation instructions below. If your AWS user already has full administrator permissions to your AWS account, you can ignore this step.

[Local Installation IAM Policy Document](deploy/docker/conf/cloudformation_policy.json)
[Amazon AWS Installation IAM Policy Document](deploy/aws/conf/cloudformation_policy.json)

**Warning**: Attaching the `Amazon AWS Installation IAM Policy Document` to a user grants the user IAM administrator privileges, which the CloudFormation template uses to create new IAM roles and policies required by the Blox framework. You should only attach this policy to users that you would trust with full administrator access to your AWS account.


## Local Installation

Our recommended way for getting started with Blox is to deploy the framework on your local Docker installation. We provide a pre-built AWS CloudFormation template that will deploy the required Amazon AWS components, and a Docker Compose file to launch the Blox framework in your local Docker environment. Please ensure that you have installed [Docker](https://docs.docker.com/engine/installation/) and [Docker Compose](https://docs.docker.com/compose/install/) locally before proceeding.

#### Create AWS Components

Run the following AWS CLI command to create the required Amazon AWS components needed to run the Blox framework locally.

```
$ cd <GitRepoBase>
$ aws --region <region> cloudformation create-stack --stack-name BloxLocal --template-body file://./deploy/docker/conf/cloudformation_template.json

$ Sample Output:
{
    "StackId": "arn:aws:cloudformation:us-east-1:123456789012:stack/BloxLocal/abcdefgh-abcd-abcd-abcd-abcdefghijkl"
}
```

#### Monitor Progress

The CloudFormation command above can take several minutes to create the required Amazon AWS components. You can monitor the progress with the following command. When the StackStatus shows `CREATE_COMPLETE`, proceed to the next step.

```
$ aws --region <region> cloudformation describe-stacks --stack-name BloxLocal

$ Sample Output:
{
  "Stacks": [
    {
      "StackId": "arn:aws:cloudformation:us-east-1:123456789012:stack/BloxLocal/abcdefgh-abcd-abcd-abcd-abcdefghijkl",
      "Description": "Template to deploy Blox framework locally",
      "StackName": "BloxLocal",
      "StackStatus": "CREATE_COMPLETE"
    }
  ]
}
```

#### Deploy Blox Locally to Docker

Before launching Blox, you will first need to update `<GitRepoBase>/deploy/docker/conf/docker-compose.yml` with the following changes:

- Update the `AWS_REGION` value with the region of your ECS and SQS resources.
- Update the `AWS_PROFILE` value with your profile name in ~/.aws/credentials. You can skip this step if you are using the `default` profile.

After you have updated `<GitRepoBase>/deploy/docker/conf/docker-compose.yml`, you can use the following commands to launch the Blox containers on your local Docker environment.

```
$ cd <GitRepoBase>/deploy/docker/conf/
$ docker-compose up -d
$ docker-compose ps

$ Sample Output:
Name             Command                          State   Ports
----------------------------------------------------------------------------------------
etcd_1        /usr/local/bin/etcd --data ...   Up      2379/tcp, 2380/tcp
scheduler_1   --bind 0.0.0.0:2000 --css- ...   Up      0.0.0.0:2000->2000/tcp
css_1         --bind 0.0.0.0:3000 --etcd ...   Up      3000/tcp
```

You have now completed the local installation of Blox. You can begin consuming the Scheduler API at http://localhost:2000/.


## Amazon AWS Installation

If you would prefer to run the Blox framework securely on Amazon AWS instead of locally, we provide a pre-built AWS CloudFormation template that will deploy Blox and the required Amazon AWS components with a single command. Deploying Blox on AWS adds TLS via Amazon API Gateway, and authentication through IAM security policies. This installation option is only recommended for advanced users who have already tested the Blox framework locally, and now wish to run it securely on AWS with a public HTTPS endpoint.

#### Create Custom Parameters File

Create a custom CloudFormation parameters file at `/tmp/blox_parameters.json` with the following content. You can remove any lines for ParameterKeys that you do not wish to override, and replace the ParameterValues for any settings that you do with to override. At the very least, you will need to override the `KeyName` parameter value with a valid EC2 Key Pair to enable SSH login to your EC2 instance. After the CloudFormation setup completes, you can follow the steps on the [Connecting to Your Linux Instance Using SSH](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AccessingInstancesLinux.html) to connect to the EC2 instance where Blox is installed. Make sure to perform the steps listed under `Enable inbound SSH traffic from your IP address to your instance`, as tcp/22 will be blocked by default into your VPC.

```
[
  {"ParameterKey":"EcsAmiId", "ParameterValue":""},
  {"ParameterKey":"InstanceType", "ParameterValue":"t2.micro"},
  {"ParameterKey":"KeyName", "ParameterValue":"<KeyPair>"},
  {"ParameterKey":"EcsClusterName", "ParameterValue":"Blox"},
  {"ParameterKey":"QueueName", "ParameterValue":"blox_queue"},
  {"ParameterKey":"ApiStageName", "ParameterValue":"blox"}
]
```

#### Deploy Blox via AWS CloudFormation

Run the following AWS CLI command to deploy Blox and the required Amazon AWS components.

```
$ cd <GitRepoBase>
$ aws --region <region> cloudformation create-stack --stack-name BloxAws --template-body file://./deploy/aws/conf/cloudformation_template.json --capabilities CAPABILITY_NAMED_IAM --parameters file:///tmp/blox_parameters.json

$ Sample Output:
{
    "StackId": "arn:aws:cloudformation:us-east-1:123456789012:stack/BloxAws/abcdefgh-abcd-abcd-abcd-abcdefghijkl"
}
```

You can monitor the progress via the same [Monitor Progress](#monitor-progress) steps above. Make sure to replace the --stack-name with 'BloxAws'. After the CloudFormation setup completes, you can retrieve the URL for your secure Daemon Scheduler REST API endpoint via the following command.

```
$ aws --region <region> cloudformation describe-stacks --stack-name BloxAws --query 'Stacks[0].Outputs[0].OutputValue' --output text

$ Sample Output:
https://<api-gateway-id>.execute-api.us-east-1.amazonaws.com/blox
```

#### Authentication

When deploying Blox on Amazon AWS, we use AWS IAM Authentication with [Signature Version 4](http://docs.aws.amazon.com/general/latest/gr/sigv4_signing.html) signing. The AWS user that you are using to authenticate with the Amazon API Gateway REST URL will need to have the following IAM Policy applied. Choose the appropriate Resource pattern based upon the permissions you want the user to have.
 
```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "execute-api:Invoke"
      ],
      "Resource": [
        "arn:aws:execute-api:<region>:<account-id>:<api-gateway-id>/*/*/*", <- Allow all API calls
        "arn:aws:execute-api:<region>:<account-id>:<api-gateway-id>/*/GET/*", <- Allow only GET calls
        "arn:aws:execute-api:<region>:<account-id>:<api-gateway-id>/*/POST/v1/environments/*/deployments" <- Allow only POST deployment calls
      ]
    }
  ]
}
```

Replace `region`, `account-id`, and `api-gateway-id` with the appropriate values. You can retrieve the `api-gateway-id` via the following command.

```
$ aws --region <region> cloudformation describe-stack-resource --stack-name BloxAws --logical-resource-id RestApi --query 'StackResourceDetail.PhysicalResourceId' --output text

$ Sample Output:
abcdef1234
```

## Upgrade Process

You can use the AWS CLI and CloudFormation to upgrade the required Amazon AWS components to the latest versions required by the Blox framework. Before starting the upgrade, you will need to export the CloudFormation parameters that are assigned to your existing CloudFormation Stack via the following commands. For all commands below, you should use the --stack-name of `BloxLocal` for local installations, and `BloxAws` for AWS installations.

```
$ aws --region <region> cloudformation describe-stacks --stack-name <BloxLocal|BloxAws> --query 'Stacks[0].Parameters' > /tmp/blox_parameters.json
$ cat /tmp/blox_parameters.json
```

If you installed Blox on AWS, and the `EcsAmiId` ParameterKey is empty in the output above, you may want to update this value to the AMI Id of your current EC2 instance. Failure to do this may result in your EC2 instance getting rebuilt on a different AMI, which would cause you to lose all Etcd state. You can use the following commands to retrieve the AMI Id from the AWS CLI. Afterwards, manually set the `EcsAmiId` ParameterValue in `/tmp/blox_parameters.json`. You can skip this step if you used the local installation.

```
$ AWS_INSTANCE=`aws --region <region> cloudformation describe-stack-resource --stack-name Blox --logical-resource-id Instance --query 'StackResourceDetail.PhysicalResourceId' --output text`
$ aws --region <region> ec2 describe-instances --instance-ids "$AWS_INSTANCE" --query 'Reservations[0].Instances[0].ImageId' --output text
```

Once you have created and updated `/tmp/blox_parameters.json` via the above commands, run the following AWS CLI commands to upgrade Blox to the latest version.

```
$ cd <GitRepoBase>

# Replace <docker|aws> with 'docker' for local installations, and 'aws' for AWS installations.
$ aws --region <region> cloudformation update-stack --stack-name <BloxLocal|BloxAws> --template-body file://./deploy/<docker|aws>/conf/cloudformation_template.json --capabilities CAPABILITY_NAMED_IAM --parameters file:///tmp/blox_parameters.json

$ Sample Output:
{
    "StackId": "arn:aws:cloudformation:us-east-1:123456789012:stack/BloxLocal/abcdefgh-abcd-abcd-abcd-abcdefghijkl"
}
```

You can monitor the upgrade progress via the same [Monitor Progress](#monitor-progress) steps above. When the StackStatus shows `UPDATE_COMPLETE`, you are running the latest versions of the required Amazon AWS components. If you are running the local installation, you will then need to launch the Docker Compose via the same [Deploy Blox Locally to Docker](#deploy-blox-locally-to-docker) steps above.

## Delete Process

If you installed Blox locally, you should stop the running Blox containers. You can skip this step if you installed Blox on AWS.

```
$ cd <GitRepoBase>/deploy/docker/conf/
$ docker-compose stop
```

You should now delete the attached Amazon AWS components via the following command. You should use the --stack-name of `BloxLocal` if you installed locally, and `BloxAws` if you installed on AWS.

```
$ aws --region <region> cloudformation delete-stack --stack-name <BloxLocal|BloxAws>
```

You can monitor the deletion progress via the same [Monitor Progress](#monitor-progress) steps above. When the AWS CLI response shows `Stack with id <BloxLocal|BloxAws> does not exist`, the deletion is complete.

**Note**: There is a known issue with the CloudFormation deletion failing on AWS installations if there is an active Lambda ENI. If the CloudFormation delete command fails, you can retry after a couple hours, or run the following AWS CLI commands to expedite the deletion.

```
$ aws --region <region> cloudformation describe-stack-resource --stack-name BloxAws --logical-resource-id Vpc --query 'StackResourceDetail.PhysicalResourceId' --output text
(Replace <VpcId> in the next command with the output returned) 

$ aws --region <region> ec2 describe-network-interfaces --filters "Name=vpc-id,Values=<VpcId>" | egrep "(Description|AttachmentId|NetworkInterfaceId)"
(Replace the <AttachmentId> and <NetworkInterfaceId> of the 'AWS Lambda VPC ENI' in the next two commands)

$ aws --region <region> ec2 detach-network-interface --attachment-id <AttachmentId>
$ aws --region <region> ec2 delete-network-interface --network-interface-id <NetworkInterfaceId>

$ aws --region <region> cloudformation delete-stack --stack-name BloxAws
```