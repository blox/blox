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
		return "", fmt.Errorf("Empty endpoit. Please specify the ECS endpoint name using the '%s' environment variable", ecsEndpointEnvVarName)
	}

	return endpoint, nil
}
