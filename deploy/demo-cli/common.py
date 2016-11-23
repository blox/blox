#!/usr/bin/env python
import argparse, json, subprocess, sys, urllib2

# Create blank object class.
class Object():
	pass

# Parse and return the CLI args.
def parse_cli_args(desc, args):
	obj = Object()

	# Print header.
	print "== Blox Demo CLI - %s ==\n" % desc

	# Create argument parser.
	parser = argparse.ArgumentParser(description=desc)
	parser.add_argument('--region', dest='region', default=None, help='AWS region')

	# Append each arg to the cli arguments.
	for arg in args:
		if 'type' in arg and arg['type'] == 'boolean':
			parser.add_argument(arg['arg'], dest=arg['dest'], action='store_true', help=arg['help'])
		else:
			parser.add_argument(arg['arg'], dest=arg['dest'], default=arg['default'], help=arg['help'])

	# Parse CLI arguments.
	arguments = parser.parse_args()
	obj.region = arguments.region

	# Loop through and process each arg.
	for arg in args:
		value = getattr(arguments, arg['dest'])

		# If value is None, prompt user for a value.
		if value == None:
			value = raw_input('- Enter %s: ' % arg['help'])
			value = value.strip()

		# Exit if argument value is not supplied.
		if value == None or (isinstance(value, basestring) and len(value) == 0):
			print 'Error: %s is required.' % arg['help']
			sys.exit(1)
		else:
			setattr(obj, arg['dest'], value)
	print "\n"

	# Return an object containing all argument values.
	return obj

# Run shell command and return object.
def run_shell_command(region, cmd):
	obj = Object()

	# Create the shell command to execute
	command = ["aws", "--output", "json"]
	if region != None:
		command.extend(["--region", region])
	command.extend(cmd)

	# Execute shell command.
	process = subprocess.Popen(command, shell=False, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
	result = process.communicate()

	if process.returncode != 0:
		print "Error: Command [%s] returned exit code %d." % (command, process.returncode)
		print result
		sys.exit(1)

	try:
		obj = json.loads(result[0])
	except Exception as e:
		print "Error: Could not parse response - %s" % e
		print result
		sys.exit(1)

	return obj

# Make HTTP call and return object.
def call_api(api):
	obj = Object()

	# URL to call.
	url = 'http://%s%s' % (api.host, api.uri)
	request = urllib2.Request(url=url)

	# Add any given headers.
	for header in api.headers:
		request.add_header(header, api.headers[header])

	# On POST requests, add a body.
	if api.method == 'POST':
		request.add_data( "" if api.data == None else json.dumps(api.data) )

	try:
		response = urllib2.urlopen(request, timeout=60)
		obj.status = response.getcode()
		obj.body = response.read()
	except urllib2.HTTPError as e:
		print "Error: Invalid response on '%s %s' - %s" % (api.method, url, str(e))
		print e.read()
		sys.exit(1)
	except Exception as e:
		print "Error: Could not call '%s %s' - %s" % (api.method, url, str(e))
		sys.exit(1)

	return obj
