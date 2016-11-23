#!/usr/bin/env python
import json, os, sys
import common

def main():
	params = common.parse_cli_args('List Task Definitions', [])

	command = ["ecs", "list-task-definitions"]
	result = common.run_shell_command(params.region, command)

	print json.dumps(result, indent=2)

if __name__ == "__main__":
	main()
