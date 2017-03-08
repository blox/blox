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
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

const (
	ec2InstanceLaunchRetrySleepSecond = 2 * time.Second
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

func (ec2Wrapper EC2Wrapper) RunInstance(count *int64, clusterName *string, instanceProfileName *string, keyName *string, amiID *string) ([]*string, error) {

	// The user data is run when the instance is launched. It registers the instance to the cluster
	// and makes the instance to self-terminate after one hour in case the test fails to clean up.
	currentDate := `$(date "+%s")`
	utmpStartDate := `$(last -t 19701111111111 | cut -d" " -f3-)`
	bootTime := `$(date -d "` + utmpStartDate + `" "+%s")`
	userData := `#!/bin/bash
echo ECS_CLUSTER=` + *clusterName + ` >> /etc/ecs/ecs.config
#cloud-boothook
#!/bin/sh -x
echo -e '#!/bin/bash' >> /etc/cron.hourly/shutdowninstance
echo -e '[ "$(( ` + currentDate + ` - ` + bootTime + ` ))" -gt 3600 ] && /sbin/shutdown -h now' >> /etc/cron.hourly/shutdowninstance
chmod +x /etc/cron.hourly/shutdowninstance`
	encodedUserData := base64.StdEncoding.EncodeToString([]byte(userData))

	instanceType := "t2.micro"
	shutdownBehavior := "terminate"

	in := ec2.RunInstancesInput{
		ImageId:      amiID,
		UserData:     &encodedUserData,
		InstanceType: &instanceType,
		MinCount:     count,
		MaxCount:     count,
		InstanceInitiatedShutdownBehavior: &shutdownBehavior,
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: instanceProfileName,
		},
	}

	if keyName != nil && *keyName != "" {
		in.SetKeyName(*keyName)
	}

	success := false
	var err error
	var resp *ec2.Reservation
	for i := 0; i < 5; i++ {
		resp, err = ec2Wrapper.client.RunInstances(&in)
		if err == nil {
			success = true
			break
		}
		time.Sleep(ec2InstanceLaunchRetrySleepSecond)
	}

	if !success {
		return nil, errors.Wrapf(err, "Failed to launch '%d' EC2 instances with "+
			"container instances registered to ECS cluster with name '%s'. ",
			*count, *clusterName)
	}

	instanceIDs := make([]*string, 0, len(resp.Instances))
	for _, instance := range resp.Instances {
		instanceIDs = append(instanceIDs, instance.InstanceId)
	}
	return instanceIDs, nil
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
