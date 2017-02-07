// Copyright 2016-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package regex

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// GetClusterNameFromARN extracts the cluster name from a cluster ARN
func GetClusterNameFromARN(clusterARN string) (string, error) {
	if len(clusterARN) == 0 {
		return "", errors.New("Cluster ARN cannot be empty")
	}

	if !IsClusterARN(clusterARN) {
		return "", fmt.Errorf("Invalid cluster ARN: %s", clusterARN)
	}

	re := regexp.MustCompile(ClusterNameAsARNSuffixRegex)
	matchedStrs := re.FindStringSubmatch(clusterARN)
	if len(matchedStrs) != 1 {
		return "", errors.New("Unable to extract cluster name from cluster ARN")
	}

	// matchedStrs[0]=/clusterName. Strip off "/" in the beginning.
	clusterName := matchedStrs[0][1:]
	return clusterName, nil
}

// GetEntityVersion extracts the entity version as an int.
func GetEntityVersion(entityVersion string) (int64, error) {
	if !IsEntityVersion(entityVersion) {
		return 0, fmt.Errorf("Invalid entity version: %s", entityVersion)
	}

	value, err := strconv.ParseInt(entityVersion, 10, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}