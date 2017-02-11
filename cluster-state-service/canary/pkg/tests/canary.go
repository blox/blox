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

package tests

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/blox/blox/cluster-state-service/canary/pkg/wrappers"
	"github.com/pkg/errors"
)

const (
	canaryClusterNamePrefix = "bloxCSSCanary"
)

// Canary defines methods for each test in the canary
type Canary interface {
	GetInstance() error
	ListInstances() error
	GetTask() error
	ListTasks() error

	GetInstanceMetric() string
	ListInstancesMetric() string
	GetTaskMetric() string
	ListTasksMetric() string
}

type canary struct {
	ecsWrapper wrappers.ECSWrapper
	ec2Wrapper wrappers.EC2Wrapper
	cssWrapper wrappers.CSSWrapper
}

// NewCanary generates a new canary
func NewCanary(sess *session.Session, clusterStateServiceEndpoint string) (Canary, error) {
	if sess == nil {
		return nil, errors.New("AWS session has to be initialized to initialize the canary. ")
	}
	if clusterStateServiceEndpoint == "" {
		return nil, errors.New("The address of the cluster-state-service endpoint had to be set to initialize the canary. ")
	}
	ecsWrapper, err := wrappers.NewECSWrapper(sess)
	if err != nil {
		return nil, err
	}
	ec2Wrapper, err := wrappers.NewEC2Wrapper(sess)
	if err != nil {
		return nil, err
	}
	cssWrapper, err := wrappers.NewCSSWrapper(clusterStateServiceEndpoint)
	if err != nil {
		return nil, err
	}
	return canary{
		ecsWrapper: ecsWrapper,
		ec2Wrapper: ec2Wrapper,
		cssWrapper: cssWrapper,
	}, nil
}
