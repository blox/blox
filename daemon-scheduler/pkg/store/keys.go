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

package store

import (
	"github.com/pkg/errors"
)

const (
	environmentKeyPrefix = "ecs/environment/"
	taskKeyConnector     = "/task/"
)

func GenerateEnvironmentKey(envName string) (string, error) {
	if envName == "" {
		return "", errors.New("Environment name cannot be empty while generating environment key")
	}
	return environmentKeyPrefix + envName, nil
}

func GenerateTaskKey(envName string, taskARN string) (string, error) {
	if envName == "" {
		return "", errors.New("Environment name cannot be empty while generating task key")
	}
	if taskARN == "" {
		return "", errors.New("Task ARN cannot be empty while generating task key")
	}
	envKey, err := GenerateEnvironmentKey(envName)
	if err != nil {
		return "", errors.Wrap(err, "Error generating task key")
	}
	return envKey + taskKeyConnector + taskARN, nil
}
