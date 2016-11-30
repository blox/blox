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

package wrappers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type AutoScalingWrapper struct {
	client *autoscaling.AutoScaling
}

func NewAutoScalingWrapper() AutoScalingWrapper {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	return AutoScalingWrapper{
		client: autoscaling.New(sess),
	}
}

func (wrapper AutoScalingWrapper) SetDesiredCapacity(asg string, count int64) error {
	in := autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(asg),
		DesiredCapacity:      aws.Int64(count),
	}
	_, err := wrapper.client.SetDesiredCapacity(&in)
	if err != nil {
		return err
	}
	return nil
}
