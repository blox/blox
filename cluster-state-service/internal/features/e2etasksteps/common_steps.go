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

package e2etasksteps

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/blox/blox/cluster-state-service/internal/features/wrappers"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/models"
	. "github.com/gucumber/gucumber"
	"github.com/pkg/errors"
)

type Exception struct {
	exceptionType string
	exceptionMsg  string
}

var (
	// Lists to memorize results required for the subsequent steps
	EcsTaskList   = []ecs.Task{}
	cssTaskList   = []models.Task{}
	exceptionList = []Exception{}

	taskDefnARN         = ""
	roleName            = "E2ETestRole"
	instanceProfileName = "E2ETestInstance"
	policyARN           = "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"
)

const (
	errorCode            = 1
	numInstancesLaunched = 1
)

func init() {

	cssWrapper := wrappers.NewCSSWrapper()
	ecsWrapper := wrappers.NewECSWrapper()
	ec2Wrapper := wrappers.NewEC2Wrapper()
	iamWrapper := wrappers.NewIAMWrapper()

	// TODO: Change these os.Exit calls to T.Errorf. Currently unable to do so because T is not initialized until the first test.
	// (https://github.com/gucumber/gucumber/issues/28)
	BeforeAll(func() {
		clusterName := wrappers.GetClusterName()

		err := ecsWrapper.CreateCluster(&clusterName)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(errorCode)
		}

		err = terminateAllContainerInstances(ec2Wrapper, ecsWrapper, clusterName)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(errorCode)
		}

		taskDefnARN, err = ecsWrapper.RegisterSleep360TaskDefinition()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(errorCode)
		}

		err = createInstanceProfile(iamWrapper)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(errorCode)
		}

		amiID, err := wrappers.GetLatestECSOptimizedAMIID()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(errorCode)
		}

		keyPair := wrappers.GetKeyPairName()
		numInstances := int64(numInstancesLaunched)
		_, err = ec2Wrapper.RunInstance(&numInstances, &clusterName, &instanceProfileName, &keyPair, &amiID)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(errorCode)
		}

		err = wrappers.ValidateECSContainerInstance(ecsWrapper, clusterName)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(errorCode)
		}
	})

	// TODO: Change these os.Exit calls to T.Errorf. Currently unable to do so because T is not initialized until the first test.
	// (https://github.com/gucumber/gucumber/issues/28)
	Before("@task|@stream-instances", func() {
		clusterName := wrappers.GetClusterName()
		err := stopAllTasks(ecsWrapper, clusterName)
		if err != nil {
			if T != nil {
				T.Errorf(err.Error())
				return
			}
			fmt.Println(err.Error())
			os.Exit(errorCode)
		}
	})

	AfterAll(func() {
		clusterName := wrappers.GetClusterName()

		err := stopAllTasks(ecsWrapper, clusterName)
		if err != nil {
			T.Errorf(err.Error())
			return
		}

		err = ecsWrapper.DeregisterTaskDefinition(taskDefnARN)
		if err != nil {
			T.Errorf(err.Error())
			return
		}

		err = terminateAllContainerInstances(ec2Wrapper, ecsWrapper, clusterName)
		if err != nil {
			T.Errorf(err.Error())
			return
		}

		err = ecsWrapper.DeleteCluster(&clusterName)
		if err != nil {
			T.Errorf(err.Error())
			return
		}

		err = deleteInstanceProfile(iamWrapper)
		if err != nil {
			T.Errorf(err.Error())
			return
		}
	})

	Given(`^I start (\d+) task(?:|s) in the ECS cluster$`, func(numTasks int) {
		startNTasks(numTasks, "someone", ecsWrapper)
	})

	And(`^I stop the (\d+) task(?:|s) in the ECS cluster$`, func(numTasks int) {
		clusterName := wrappers.GetClusterName()
		if len(EcsTaskList) != numTasks {
			T.Errorf("Error memorizing tasks started using ECS client. ")
			return
		}
		for _, t := range EcsTaskList {
			err := ecsWrapper.StopTask(clusterName, *t.TaskArn)
			if err != nil {
				T.Errorf(err.Error())
				return
			}
		}
	})

	When(`^I get task with the cluster name and task ARN$`, func() {
		cssTaskList = nil

		clusterName := wrappers.GetClusterName()

		time.Sleep(15 * time.Second)
		if len(EcsTaskList) != 1 {
			T.Errorf("Error memorizing task started using ECS client. ")
			return
		}
		taskARN := *EcsTaskList[0].TaskArn
		cssTask, err := cssWrapper.GetTask(clusterName, taskARN)
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		cssTaskList = append(cssTaskList, *cssTask)
	})

	Then(`^I get a (.+?) task exception$`, func(exceptionType string) {
		if len(exceptionList) != 1 {
			T.Errorf("Error memorizing exception. ")
			return
		}
		if exceptionType != exceptionList[0].exceptionType {
			T.Errorf("Expected exception '%s' but got '%s'. ", exceptionType, exceptionList[0].exceptionType)
		}
	})

	And(`^the task exception message contains "(.+?)"$`, func(exceptionMsg string) {
		if len(exceptionList) != 1 {
			T.Errorf("Error memorizing exception. ")
			return
		}
		if !strings.Contains(exceptionList[0].exceptionMsg, exceptionMsg) {
			T.Errorf("Expected exception message returned '%s' to contain '%s'. ", exceptionList[0].exceptionMsg, exceptionMsg)
		}
	})
}

func stopAllTasks(ecsWrapper wrappers.ECSWrapper, clusterName string) error {
	taskARNList, err := ecsWrapper.ListTasks(clusterName)
	if err != nil {
		return err
	}
	for _, t := range taskARNList {
		err = ecsWrapper.StopTask(clusterName, *t)
		if err != nil {
			return err
		}
	}
	return nil
}

func startNTasks(numTasks int, startedBy string, ecsWrapper wrappers.ECSWrapper) {
	EcsTaskList = nil
	cssTaskList = nil

	clusterName := wrappers.GetClusterName()

	for i := 0; i < numTasks; i++ {
		ecsTask, err := ecsWrapper.StartTask(clusterName, taskDefinitionSleep300, startedBy)
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		EcsTaskList = append(EcsTaskList, ecsTask)
	}
}

func createInstanceProfile(iamWrapper wrappers.IAMWrapper) error {
	assumeRolePolicy := `{
		"Version": "2012-10-17",
		"Statement": [
		{
		"Effect": "Allow",
		"Principal": {
			"Service": "ec2.amazonaws.com"
		},
		"Action": "sts:AssumeRole"
		}
		]
	}`

	err := iamWrapper.GetRole(&roleName)
	if err != nil {
		if awsErr, ok := errors.Cause(err).(awserr.Error); ok && awsErr.Code() == "NoSuchEntity" {
			err = iamWrapper.CreateRole(&roleName, &assumeRolePolicy)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	err = iamWrapper.GetInstanceProfile(&instanceProfileName)
	if err != nil {
		if awsErr, ok := errors.Cause(err).(awserr.Error); ok && awsErr.Code() == "NoSuchEntity" {
			err = iamWrapper.CreateInstanceProfile(&instanceProfileName)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	err = iamWrapper.AttachRolePolicy(&policyARN, &roleName)
	if err != nil {
		return err
	}

	err = iamWrapper.AddRoleToInstanceProfile(&roleName, &instanceProfileName)
	if err != nil {
		if awsErr, ok := errors.Cause(err).(awserr.Error); ok {
			if awsErr.Code() != "EntityAlreadyExists" && awsErr.Code() != "LimitExceeded" {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func deleteInstanceProfile(iamWrapper wrappers.IAMWrapper) error {
	err := iamWrapper.DetachRolePolicy(&policyARN, &roleName)
	if err != nil {
		return err
	}

	err = iamWrapper.RemoveRoleFromInstanceProfile(&roleName, &instanceProfileName)
	if err != nil {
		return err
	}

	err = iamWrapper.DeleteRole(&roleName)
	if err != nil {
		return err
	}

	err = iamWrapper.DeleteInstanceProfile(&instanceProfileName)
	if err != nil {
		return err
	}

	return nil
}

func terminateAllContainerInstances(ec2Wrapper wrappers.EC2Wrapper, ecsWrapper wrappers.ECSWrapper, clusterName string) error {
	instanceARNs, err := ecsWrapper.ListContainerInstances(clusterName)
	if err != nil {
		return errors.Wrapf(err, "Failed to list container instances from cluster '%v'.", clusterName)
	}

	if len(instanceARNs) == 0 {
		return nil
	}

	err = ecsWrapper.DeregisterContainerInstances(&clusterName, instanceARNs)
	if err != nil {
		return errors.Wrapf(err, "Failed to deregister container instances '%v'.", instanceARNs)
	}

	ec2InstanceIDs := make([]*string, 0, len(instanceARNs))
	for _, v := range instanceARNs {
		containerInstance, err := ecsWrapper.DescribeContainerInstance(clusterName, *v)
		if err != nil {
			return errors.Wrapf(err, "Failed to describe container instance '%v'.", v)
		}
		ec2InstanceIDs = append(ec2InstanceIDs, containerInstance.Ec2InstanceId)
	}

	err = ec2Wrapper.TerminateInstances(ec2InstanceIDs)
	if err != nil {
		return errors.Wrapf(err, "Failed to terminate container instances '%v'.", ec2InstanceIDs)
	}

	return nil
}
