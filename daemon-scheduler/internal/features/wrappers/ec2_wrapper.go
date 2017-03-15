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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

type EC2Wrapper struct {
	client *ec2.EC2
}

func NewEC2Wrapper() EC2Wrapper {
	awsSession, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	return EC2Wrapper{
		client: ec2.New(awsSession),
	}
}

func (ec2Wrapper EC2Wrapper) TerminateInstances(instanceIDs []*string) error {
	in := ec2.TerminateInstancesInput{
		InstanceIds: instanceIDs,
	}

	_, err := ec2Wrapper.client.TerminateInstances(&in)
	if err != nil {
		return errors.Wrapf(err, "Failed to terminate EC2 instances with instance IDs '%v'. ", instanceIDs)
	}

	return nil
}

func (EC2Wrapper EC2Wrapper) DescribeAvailabilityZones() ([]*string, error) {
	// No need to pass any arguments since EC2 can obtain the region name itself
	in := ec2.DescribeAvailabilityZonesInput{}

	resp, err := EC2Wrapper.client.DescribeAvailabilityZones(&in)
	if err != nil {
		return nil, errors.Wrapf(err,"Fail to describe availability zones")
	}

	availabilityZones := make([]*string, 0, len(resp.AvailabilityZones))
	for _, v := range resp.AvailabilityZones {
		availabilityZones = append(availabilityZones, v.ZoneName)
	}

	return availabilityZones, nil
}
