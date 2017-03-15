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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"time"
	"fmt"
)

const (
	ASGConfigLaunchRetrySleepSecond = 2 * time.Second
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

func (wrapper AutoScalingWrapper) CreateAutoScalingGroup(asg string, configName string, availabilityZones []*string) error {
	minSize := int64(0)
	// MaxSize is set to 10 because it's the maximum number of instances launched by the e2e test
	maxSize := int64(10)
	in := autoscaling.CreateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(asg),
		MinSize: aws.Int64(minSize),
		MaxSize: aws.Int64(maxSize),
		LaunchConfigurationName: aws.String(configName),
		AvailabilityZones: availabilityZones,
	}

	_, err := wrapper.client.CreateAutoScalingGroup(&in)
	if err != nil {
		return err
	}
	return nil
}

func (wrapper AutoScalingWrapper) GetAutoScalingGroupStatus(asg string) (string, error) {
	in := autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{aws.String(asg)},
	}

	resp, err := wrapper.client.DescribeAutoScalingGroups(&in)
	if err != nil {
		return "", err
	}

	if len(resp.AutoScalingGroups) == 0 {
		return "", fmt.Errorf("Cannot get status of Autoscaling Group: Autoscaling Group not found.")
	}

	var status string
	if resp.AutoScalingGroups[0].Status != nil {
		status = *resp.AutoScalingGroups[0].Status
	} else {
		status = ""
	}

	return status, nil
}

func (wrapper AutoScalingWrapper) DeleteAutoScalingGroup(asg string, forceDelete bool) error {
	in := autoscaling.DeleteAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(asg),
		ForceDelete: aws.Bool(forceDelete),
	}

	_, err := wrapper.client.DeleteAutoScalingGroup(&in)
	if err != nil {
		return err
	}
	return nil
}

func (wrapper AutoScalingWrapper) CreateLaunchConfiguration(configName string, clusterName string, instanceProfileName string, keyName string, amiID string)  error {
	// The user data is run when the instance is launched. It registers the instance to the cluster
	// and makes the instance to self-stop after one hour in case the test fails to clean up.
	// TODO: Change the shutdown behavior from stop to terminate when shutdown behavior is supported by Autoscaling
	currentDate := `$(date "+%s")`
	utmpStartDate := `$(last -t 19701111111111 | cut -d" " -f3-)`
	bootTime := `$(date -d "` + utmpStartDate + `" "+%s")`
	userData := `#!/bin/bash
echo ECS_CLUSTER=` + clusterName + ` >> /etc/ecs/ecs.config
#cloud-boothook
#!/bin/sh -x
echo -e '#!/bin/bash' >> /etc/cron.hourly/shutdowninstance
echo -e '[ "$(( ` + currentDate + ` - ` + bootTime + ` ))" -gt 3600 ] && /sbin/shutdown -h now' >> /etc/cron.hourly/shutdowninstance
chmod +x /etc/cron.hourly/shutdowninstance`
	encodedUserData := base64.StdEncoding.EncodeToString([]byte(userData))

	instanceType := "t2.micro"

	in := autoscaling.CreateLaunchConfigurationInput{
		LaunchConfigurationName: aws.String(configName),
		IamInstanceProfile: aws.String(instanceProfileName),
		ImageId: aws.String(amiID),
		UserData: &encodedUserData,
		InstanceType: &instanceType,
	}

	if keyName != "" {
		in.SetKeyName(keyName)
	}

	success := false
	var err error
	for i := 0; i < 5; i++ {
		_, err := wrapper.client.CreateLaunchConfiguration(&in)
		if err == nil {
			success = true
			break
		}
		time.Sleep(ASGConfigLaunchRetrySleepSecond)
	}

	if !success {
		return err
	}

	return nil
}

func (wrapper AutoScalingWrapper) DeleteLaunchConfiguration(configName string) error {
	in := autoscaling.DeleteLaunchConfigurationInput{
		LaunchConfigurationName: aws.String(configName),
	}

	_, err := wrapper.client.DeleteLaunchConfiguration(&in)
	if err != nil {
		return err
	}
	return nil
}
