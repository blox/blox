#!/usr/bin/env python
import json, os, sys
import common

def main():
	# Command Line Arguments
	args = [
		{'arg':'--host', 'dest':'host', 'default':'localhost:3000', 'help':'Blox CSS <Host>:<Port>'},
		{'arg':'--cluster', 'dest':'cluster', 'default':None, 'required':False, 'help':'ECS cluster name'},
		{'arg':'--status', 'dest':'status', 'default':None, 'required':False, 'help':'ECS task status'},
		{'arg':'--task-arn', 'dest':'task', 'default':None, 'required':False, 'help':'ECS task Arn'}
	]
	
	# Parse Command Line Arguments
	params = common.parse_cli_args('List Blox Tasks', args)

	run_local(params)

# Call Blox CSS Local Endpoint
def run_local(params):
	api = common.Object()
	api.method = 'GET'
	api.headers = {}
	api.host = params.host
	api.uri = '/v1/tasks'
	api.data = None

	if params.cluster != None and params.task != None:
		api.uri = '/v1/tasks/%s/%s' % (params.cluster, params.task)
	elif params.cluster != None:
		api.uri = '/v1/tasks/filter?cluster=%s' % params.cluster
	elif params.status != None:
		api.uri = '/v1/tasks/filter?status=%s' % params.status

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
	main()
