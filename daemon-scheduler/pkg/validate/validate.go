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

package validate

import (
	"regexp"

	"github.com/blox/blox/daemon-scheduler/pkg/regex"
)

// IsClusterName validates a cluster name against the cluster name regex
func IsClusterName(clusterName string) bool {
	validClusterName := regexp.MustCompile(regex.ClusterNameRegex)
	if validClusterName.MatchString(clusterName) {
		return true
	}
	return false
}

// IsClusterARN validates a cluster ARN against the cluster ARN regex
func IsClusterARN(clusterARN string) bool {
	validClusterARN := regexp.MustCompile(regex.ClusterARNRegex)
	if validClusterARN.MatchString(clusterARN) {
		return true
	}
	return false
}
