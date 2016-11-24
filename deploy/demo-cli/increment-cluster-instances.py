#!/usr/bin/env python
import base64, json, os, re, sys
import common

def main():
	args = [
		{'arg':'--cluster', 'dest':'cluster', 'default':None, 'help':'ECS cluster name'},
		{'arg':'--num', 'dest':'num', 'default':None, 'help':'number of instances to increment by'}
	]
	params = common.parse_cli_args('Increment Cluster Instances', args)
	params.num = int(params.num)

	command = ["ecs", "describe-clusters", "--clusters", params.cluster]
	clusterResult = common.run_shell_command(params.region, command)
	if len(clusterResult['clusters']) == 0:
		print "Error: Cluster '%s' does not exist." % params.cluster
		sys.exit(1)

	command = ["ecs", "list-container-instances", "--cluster", params.cluster]
	containersResult = common.run_shell_command(params.region, command)
	if len(containersResult['containerInstanceArns']) == 0:
		print "Error: Cluster '%s' does not contain any instances." % params.cluster
		sys.exit(1)

	command = ["ecs", "describe-container-instances", "--cluster", params.cluster, "--container-instances", containersResult['containerInstanceArns'][0]]
	containerResult = common.run_shell_command(params.region, command)
	if len(containerResult['containerInstances']) == 0:
		print "Error: Could not retrieve container instance '%s'." % containersResult['containerInstanceArns'][0]
		sys.exit(1)

	command = ["ec2", "describe-instances", "--instance-ids", containerResult['containerInstances'][0]['ec2InstanceId']]
	instanceResult = common.run_shell_command(params.region, command)
	if len(instanceResult['Reservations']) == 0:
		print "Error: Could not retrieve ec2 instance '%s'." % containerResult['containerInstances'][0]['ec2InstanceId']
		sys.exit(1)

	amiId = instanceResult['Reservations'][0]['Instances'][0]['ImageId']
	keyName = instanceResult['Reservations'][0]['Instances'][0]['KeyName']
	securityGroup = instanceResult['Reservations'][0]['Instances'][0]['SecurityGroups'][0]['GroupId']
	subnetId = instanceResult['Reservations'][0]['Instances'][0]['SubnetId']
	instanceType = instanceResult['Reservations'][0]['Instances'][0]['InstanceType']
	instanceProfileArn = instanceResult['Reservations'][0]['Instances'][0]['IamInstanceProfile']['Arn']
	instanceProfile = None

	instanceProfileRegex = re.match('(.+)\/(.+)', instanceProfileArn)
	if instanceProfileRegex != None:
		instanceProfile = instanceProfileRegex.group(2)

	userDataRaw = "#!/bin/bash\necho ECS_CLUSTER=%s >> /etc/ecs/ecs.config\n" % params.cluster
	userData = base64.b64encode(userDataRaw)

	for i in range(params.num):
		command = ["ec2", "run-instances", "--image-id", amiId, "--instance-type", instanceType, "--user-data", userData, "--iam-instance-profile", "Name=\"%s\"" % instanceProfile, "--security-group-ids", securityGroup, "--key-name", keyName]
		result = common.run_shell_command(params.region, command)
		print "Created instance '%s'." % result['Instances'][0]['InstanceId']

if __name__ == "__main__":
	main()
