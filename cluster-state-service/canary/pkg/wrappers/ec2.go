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

package wrappers

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
)

// EC2Wrapper defines methods to access wrapper methods to call EC2 APIs
type EC2Wrapper interface {
	LaunchInstances(count *int64, clusterName *string) ([]*string, error)
	TerminateInstances(instanceIDs []*string) error
}

type ec2ClientWrapper struct {
	client ec2iface.EC2API
}

// NewEC2Wrapper returns a new EC2Wrapper for the canary
func NewEC2Wrapper(sess *session.Session) (EC2Wrapper, error) {
	if sess == nil {
		return nil, errors.New("AWS session has to be initialized to initialize the EC2 client. ")
	}
	ec2Client := ec2.New(sess)
	return ec2ClientWrapper{
		client: ec2Client,
	}, nil
}

func (wrapper ec2ClientWrapper) LaunchInstances(count *int64, clusterName *string) ([]*string, error) {
	// TODO: Self terminate ec2 instances launched using user data in case
	// test cleanup fails for any reason
	userData := `#!/bin/bash
echo ECS_CLUSTER=` + *clusterName + ` >> /etc/ecs/ecs.config`
	encodedUserData := base64.StdEncoding.EncodeToString([]byte(userData))

	ecsAMI := wrapper.getECSAMI()
	instanceType := "t2.micro"
	// TODO: Make the key-name configurable
	keyName := "blox-canary"
	// TODO: Make the instance-role configurable
	instanceRole := "blox-canary-ecs-role"
	in := ec2.RunInstancesInput{
		ImageId:      &ecsAMI,
		UserData:     &encodedUserData,
		InstanceType: &instanceType,
		MinCount:     count,
		MaxCount:     count,
		KeyName:      &keyName,
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: &instanceRole,
		},
	}

	resp, err := wrapper.client.RunInstances(&in)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to launch '%d' EC2 instances with "+
			"container instances registered to ECS cluster with name '%s'. ",
			*count, *clusterName)
	}

	instanceIDs := make([]*string, len(resp.Instances))
	for _, instance := range resp.Instances {
		instanceIDs = append(instanceIDs, instance.InstanceId)
	}
	return instanceIDs, nil
}

func (wrapper ec2ClientWrapper) TerminateInstances(instanceIDs []*string) error {
	in := ec2.TerminateInstancesInput{
		InstanceIds: instanceIDs,
	}

	_, err := wrapper.client.TerminateInstances(&in)
	if err != nil {
		return errors.Wrapf(err, "Failed to terminate EC2 instances with instance IDs '%v'. ", instanceIDs)
	}
	return nil
}

func (wrapper ec2ClientWrapper) getECSAMI() string {
	// TODO: Use ec2 DescribeImages API to get the latest ECS AMI.
	// The following AMI id corresponds to us-east-1
	return "ami-d69c74c0"
}
