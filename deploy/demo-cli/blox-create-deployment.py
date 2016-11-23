#!/usr/bin/env python
import json, os, sys
import common

def main(argv):
	# Command Line Arguments
	args = [{'arg':'--apigateway', 'dest':'apigateway', 'default':None, 'type':'boolean', 'help':'Call API Gateway endpoint'}]
	if '--apigateway' in argv:
		args.extend([{'arg':'--stack', 'dest':'stack', 'default':None, 'help':'CloudFormation stack name'}])
	else:
		args.extend([{'arg':'--host', 'dest':'host', 'default':'localhost:2000', 'help':'Blox Scheduler <Host>:<Port>'}])
	args.extend([{'arg':'--environment', 'dest':'environment', 'default':None, 'help':'Blox environment name'}])

	# Parse Command Line Arguments
	params = common.parse_cli_args('Create Blox Deployment', args)

	if params.apigateway:
		run_apigateway(params)
	else:
		run_local(params)

# Call Blox Scheduler API Gateway Endpoint
def run_apigateway(params):
	command = ["cloudformation", "describe-stack-resource", "--stack-name", params.stack, "--logical-resource-id", "RestApi"]
	restApi = common.run_shell_command(params.region, command)

	command = ["cloudformation", "describe-stack-resource", "--stack-name", params.stack, "--logical-resource-id", "ApiResource"]
	restResource = common.run_shell_command(params.region, command)

	command = ["apigateway", "test-invoke-method", "--rest-api-id", restApi['StackResourceDetail']['PhysicalResourceId'], "--resource-id", restResource['StackResourceDetail']['PhysicalResourceId'], "--http-method", "POST", "--headers", "{}", "--path-with-query-string", "/v1/environments/%s/deployments" % params.environment, "--body", ""]
	response = common.run_shell_command(params.region, command)

	print "HTTP Response Code: %d" % response['status']

	try:
		obj = json.loads(response['body'])
		print json.dumps(obj, indent=2)
	except Exception as e:
		print "Error: Could not parse response - %s" % e
		print json.dumps(response, indent=2)
		sys.exit(1)

# Call Blox Scheduler Local Endpoint
def run_local(params):
	api = common.Object()
	api.method = 'POST'
	api.headers = {}
	api.host = params.host
	api.uri = '/v1/environments/%s/deployments' % params.environment
	api.data = None

	response = common.call_api(api)

	print "HTTP Response Code: %d" % response.status

	try:
		obj = json.loads(response.body)
		print json.dumps(obj, indent=2)
	except Exception as e:
		print "Error: Could not parse response - %s" % e
		print response.body
		sys.exit(1)

if __name__ == "__main__":
	main(sys.argv[1:])
