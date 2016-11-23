#!/usr/bin/env python
import json, os, sys
import common

def main():
	args = [
		{'arg':'--file', 'dest':'file', 'default':None, 'help':'path to task definition file'}
	]
	params = common.parse_cli_args('Register Task Definition', args)

	if not os.path.isfile(params.file):
		print "Error: File path '%s' does not exist." % params.file
		sys.exit(1)

	command = ["ecs", "register-task-definition", "--cli-input-json", "file://%s" % params.file]
	result = common.run_shell_command(params.region, command)

	print json.dumps(result, indent=2)

if __name__ == "__main__":
	main()
