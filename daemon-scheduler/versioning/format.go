// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package versioning

import "fmt"

type Versioner interface {
	Version() (string, error)
}

// PrintVersions prints the version information on stdout as a multi-line
// string. The output will look similar to the following:
//
//    Blox Daemon Scheduler:
//        Version: 0.0.1
//        Commit: 55347bc
func PrintVersion(extra ...Versioner) {
	cleanliness := ""
	if GitDirty {
		cleanliness = "\tDirty: true\n"
	}

	fmt.Printf(`Blox Daemon Scheduler:
	Version: %v
	Commit: %v
%v`, Version, GitShortHash, cleanliness)

	for _, versioner := range extra {
		if str, err := versioner.Version(); err == nil {
			fmt.Printf("\t%v\n", str)
		}
	}

}

// String produces a human-readable string showing the agent version.
func String() string {
	ret := "Blox Daemon Scheduler - v" + Version + " ("
	if GitDirty {
		ret += "*"
	}
	return ret + GitShortHash + ")"
}

func GitHashString() string {
	if GitDirty {
		return "*" + GitShortHash
	}
	return GitShortHash
}
