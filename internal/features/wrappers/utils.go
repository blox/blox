// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the License). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the license file accompanying this file. This file is distributed
// on an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package wrappers

import (
	"fmt"
	"os"
)

const (
	clusterNameEnvVarName = "ECS_CLUSTER"
	ecsEndpointEnvVarName = "ECS_ENDPOINT"
)

func GetClusterName() (string, error) {
	cluster := os.Getenv(clusterNameEnvVarName)
	if cluster == "" {
		return "", fmt.Errorf("Empty cluster name. Please specify the ECS cluster name using the '%s' environment variable", clusterNameEnvVarName)
	}

	return cluster, nil
}

func getECSEndpoint() (string, error) {
	endpoint := os.Getenv(ecsEndpointEnvVarName)
	if endpoint == "" {
		return "", fmt.Errorf("Empty endpoint. Please specify the ECS endpoint name using the '%s' environment variable", ecsEndpointEnvVarName)
	}

	return endpoint, nil
}
