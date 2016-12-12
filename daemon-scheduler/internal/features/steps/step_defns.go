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

package steps

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/blox/blox/daemon-scheduler/generated/v1/client/operations"
	"github.com/blox/blox/daemon-scheduler/generated/v1/models"
	"github.com/blox/blox/daemon-scheduler/internal/features/wrappers"
	. "github.com/gucumber/gucumber"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	deploymentCompleted           = "completed"
	taskRunning                   = "RUNNING"
	deploymentCompleteWaitSeconds = 50
)

func init() {
	asgWrapper := wrappers.NewAutoScalingWrapper()
	ecsWrapper := wrappers.NewECSWrapper()
	edsWrapper := wrappers.NewEDSWrapper()
	ctx := context.Background()

	css, err := wrappers.NewClusterState()
	if err != nil {
		T.Errorf("Error creating CSS client: %v", err)
		return
	}

	When(`^I make a Ping call$`, func() {
		err = edsWrapper.Ping()
	})

	Then(`^the Ping response indicates that the service is healthy$`, func() {
		if err != nil {
			T.Errorf(err.Error())
		}
	})

	Given(`^A cluster "env.(.+?)" and asg "env.(.+?)"$`, func(cEnv string, aEnv string) {
		c := os.Getenv(cEnv)
		if len(c) == 0 {
			T.Errorf("ECS_CLUSTER env-var is not defined")
		}
		a := os.Getenv(aEnv)
		if len(a) == 0 {
			T.Errorf("ECS_CLUSTER_ASG env-var is not defined")
		}
		cluster = c
		asg = a
	})

	Given(`^A cluster named "env.(.+?)"$`, func(cEnv string) {
		c := os.Getenv(cEnv)
		if len(c) == 0 {
			T.Errorf("ECS_CLUSTER env-var is not defined")
		}
		cluster = c
	})

	Given(`^(?:a|another) cluster "(.+?)"$`, func(c string) {
		cARN, err := ecsWrapper.CreateCluster(c)
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		cluster = c
		clusterARN = *cARN
	})

	When(`^I update the desired-capacity of cluster to (\d+) instances and wait for a max of (\d+) seconds$`, func(count int, seconds int) {
		err := asgWrapper.SetDesiredCapacity(asg, int64(count))
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		ok, err := doSomething(time.Duration(seconds)*time.Second, 1*time.Second, func() (bool, error) {
			instances, err := css.ListInstances(cluster)
			if err != nil {
				return false, errors.Wrapf(err, "Error calling ListInstances for cluster %s", cluster)
			}
			activeCount := 0
			for _, instance := range instances {
				if "ACTIVE" == aws.StringValue(instance.Status) {
					activeCount++
				}
			}
			return count == activeCount, nil
		})
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		if !ok {
			T.Errorf("Expected %d instances in cluster %s", count, cluster)
			return
		}

	})

	And(`^a registered "(.+?)" task-definition$`, func(td string) {
		resp, err := ecsWrapper.RegisterTaskDefinition(td)
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		taskDefinition = resp
	})

	And(`^I deregister task-definition$`, func() {
		err := ecsWrapper.DeregisterTaskDefinition(taskDefinition)
		if err != nil {
			T.Errorf(err.Error())
			return
		}
	})

	When(`^I create an environment with name "(.+?)" in the cluster using the task-definition$`,
		func(e string) {
			environment = e
			err := edsWrapper.DeleteEnvironment(&environment)
			if err != nil {
				T.Errorf("Was not able to delete environment %v: %v", environment, err.Error())
			}

			createEnvReq := &models.CreateEnvironmentRequest{
				InstanceGroup: &models.InstanceGroup{
					Cluster: cluster,
				},
				Name:           &environment,
				TaskDefinition: &taskDefinition,
			}
			env, err := edsWrapper.CreateEnvironment(createEnvReq)
			if err != nil {
				_, ok := err.(*operations.CreateEnvironmentBadRequest)
				if !ok {
					T.Errorf(err.Error())
					return
				}
			} else {
				deploymentToken = env.DeploymentToken
			}
		})

	Then(`^creating the same environment should fail with BadRequest$`,
		func() {
			createEnvReq := &models.CreateEnvironmentRequest{
				InstanceGroup: &models.InstanceGroup{
					Cluster: cluster,
				},
				Name:           &environment,
				TaskDefinition: &taskDefinition,
			}
			_, err := edsWrapper.CreateEnvironment(createEnvReq)
			if err != nil {
				_, ok := err.(*operations.CreateEnvironmentBadRequest)
				if !ok {
					T.Errorf(err.Error())
					return
				}
			} else {
				T.Errorf("Expecting CreateEnvironmentBadRequest error")
				return
			}
		})

	Then(`^I create an environment with name "(.+?)" it should fail with NotFound$`, func(e string) {
		err := edsWrapper.DeleteEnvironment(&e)
		if err != nil {
			T.Errorf("Was not able to delete environment %v: %v", e, err.Error())
		}

		createEnvReq := &models.CreateEnvironmentRequest{
			InstanceGroup: &models.InstanceGroup{
				Cluster: cluster,
			},
			Name:           &e,
			TaskDefinition: &taskDefinition,
		}
		_, err = edsWrapper.CreateEnvironment(createEnvReq)
		if err != nil {
			_, ok := err.(*operations.CreateEnvironmentNotFound)
			if !ok {
				T.Errorf(err.Error())
				return
			}
		} else {
			T.Errorf("Expecting CreateEnvironmentNotFound error")
			return
		}
	})

	Then(`^I create an environment with name "(.+?)" it should fail with BadRequest$`,
		func(e string) {
			err := edsWrapper.DeleteEnvironment(&e)
			if err != nil {
				T.Errorf("Was not able to delete environment %v: %v", e, err.Error())
			}

			createEnvReq := &models.CreateEnvironmentRequest{
				InstanceGroup: &models.InstanceGroup{
					Cluster: cluster,
				},
				Name:           &e,
				TaskDefinition: &taskDefinition,
			}
			_, err = edsWrapper.CreateEnvironment(createEnvReq)
			if err != nil {
				_, ok := err.(*operations.CreateEnvironmentBadRequest)
				if !ok {
					T.Errorf(err.Error())
					return
				}
			} else {
				T.Errorf("Expecting CreateEnvironmentBadRequest error")
				return
			}
		})

	And(`^I delete cluster$`,
		func() {
			_, err := ecsWrapper.DeleteCluster(cluster)
			if err != nil {
				T.Errorf(err.Error())
				return
			}
		})

	Then(`^GetEnvironment should succeed$`,
		func() {
			_, err := edsWrapper.GetEnvironment(&environment)
			if err != nil {
				T.Errorf(err.Error())
				return
			}
		})

	Then(`^GetEnvironment with name "(.+?)" should fail with NotFound$`,
		func(e string) {
			_, err := edsWrapper.GetEnvironment(&e)
			if err != nil {
				_, ok := err.(*operations.GetEnvironmentNotFound)
				if !ok {
					T.Errorf(err.Error())
					return
				}
			} else {
				T.Errorf("Expecting GetEnvironmentNotFound error")
				return
			}
		})

	Then(`^GetDeployment with created deployment should succeed$`,
		func() {
			deploymentGet, err := edsWrapper.GetDeployment(&environment, &deploymentID)
			if err != nil {
				T.Errorf(err.Error())
				return
			}
			assert.Equal(T, deploymentID, *deploymentGet.ID, "DeploymentID should match")
		})

	Then(`^the environment should be returned in ListEnvironments call$`,
		func() {
			environments, err := edsWrapper.ListEnvironments()
			if err != nil {
				T.Errorf(err.Error())
				return
			}
			found := false
			for _, env := range environments {
				if *env.Name == environment {
					found = true
					break
				}
			}
			assert.Equal(T, true, found, "Did not find environment with name "+environment)
		})

	Then(`^there should be at least (\d+) environment returned when I call ListEnvironments with cluster filter set to the second cluster$`, func(numEnvs int) {
		environments, err := edsWrapper.FilterEnvironments(clusterARN)
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		assert.True(T, len(environments) >= numEnvs,
			"Number of environments in the response should be at least "+string(numEnvs))
		environmentList = environments
	})

	And(`^all the environments in the response should correspond to the second cluster$`,
		func() {
			for _, env := range environmentList {
				if env.InstanceGroup.Cluster != clusterARN {
					T.Errorf("Environment in list environments response with cluster filter set to '" +
						clusterARN + "' belongs to cluster + '" + env.InstanceGroup.Cluster + "'")
				}
			}
		})

	And(`^second environment should be one of the environments in the response$`,
		func() {
			found := false
			for _, env := range environmentList {
				if *env.Name == environment {
					found = true
					break
				}
			}
			assert.True(T, found, "Did not find environment with name "+environment)
		})

	Then(`^I call CreateDeployment API$`, func() {
		createDeployment(deploymentToken, ctx, edsWrapper)
	})

	Then(`^creating another deployment with the same token should fail$`, func() {
		_, err := edsWrapper.CreateDeployment(context.TODO(), &environment, &deploymentToken)
		if err != nil {
			_, ok := err.(*operations.CreateDeploymentBadRequest)
			if !ok {
				T.Errorf(err.Error())
				return
			}
		} else {
			T.Errorf("Expecting CreateDeploymentBadRequest error")
			return
		}
	})

	Then(`^Deployment should be returned in ListDeployments call$`,
		func() {
			deployments, err := edsWrapper.ListDeployments(aws.String(environment))
			if err != nil {
				T.Errorf(err.Error())
				return
			}
			found := false
			for _, d := range deployments {
				if *d.ID == deploymentID {
					found = true
					break
				}
			}
			assert.Equal(T, true, found, fmt.Sprintf("Did not find deployment with id:%s under environment:%s", deploymentID, environment))
		})

	Then(`^the deployment should have (\d+) task(?:|s) running within (\d+) seconds$`, func(count int, seconds int) {
		ok, err := doSomething(time.Duration(seconds)*time.Second, 1*time.Second, func() (bool, error) {
			tasks, err := ecsWrapper.ListTasks(cluster, aws.String(deploymentID))
			if err != nil {
				return false, errors.Wrapf(err, "Error calling ListTasks for cluster %s and deployment %s", cluster, deploymentID)
			}

			runningTasks := filterTasksByStatusRunning(aws.String(cluster), tasks, ecsWrapper)
			return count == len(runningTasks), nil
		})

		if err != nil {
			T.Errorf(err.Error())
			return
		}

		if !ok {
			T.Errorf("Expecting at least %d tasks to be launched in the cluster %v", count, cluster)
			return
		}
	})

	Then(`^the deployment should complete in (\d+) seconds$`, func(seconds int) {
		ok, err := doSomething(time.Duration(seconds)*time.Second, 1*time.Second, func() (bool, error) {
			deployment, err := edsWrapper.GetDeployment(aws.String(environment), aws.String(deploymentID))
			if err != nil {
				return false, errors.Wrapf(err, "Error calling GetDeployment for environment %s and deployment %s", environment, deploymentID)
			}

			return aws.StringValue(deployment.Status) == models.DeploymentStatusCompleted, nil
		})

		if err != nil {
			T.Errorf(err.Error())
			return
		}

		if !ok {
			T.Errorf("Expecting the deployment status to be %v", taskRunning)
			return
		}
	})

	And(`^Deployment should be marked as completed$`, func() {
		deployment, err := edsWrapper.GetDeployment(aws.String(environment), aws.String(deploymentID))
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		if aws.StringValue(deployment.Status) != deploymentCompleted {
			T.Errorf("Expected deployment %s to be completed but was %s", deploymentID, *deployment.Status)
			return
		}
	})

	And(`^I stop the tasks running in cluster$`, func() {
		tasks, err := ecsWrapper.ListTasks(cluster, nil)
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		for _, task := range tasks {
			err := ecsWrapper.StopTask(cluster, *task)
			if err != nil {
				T.Errorf(err.Error())
				return
			}
		}
	})

	When(`^I call GetDeployment with environment "(.+?)", it should fail with NotFound$`,
		func(e string) {
			_, err := edsWrapper.GetDeployment(aws.String(e), aws.String(deploymentID))
			if err != nil {
				_, ok := err.(*operations.GetDeploymentNotFound)
				if !ok {
					T.Errorf(err.Error())
					return
				}
			} else {
				T.Errorf("Expecting GetDeploymentNotFound error")
				return
			}
		})

	When(`^I call GetDeployment with id "(.+?)", it should fail with NotFound$`,
		func(d string) {
			_, err := edsWrapper.GetDeployment(aws.String(environment), aws.String(d))
			if err != nil {
				_, ok := err.(*operations.GetDeploymentNotFound)
				if !ok {
					T.Errorf(err.Error())
					return
				}
			} else {
				T.Errorf("Expecting GetDeploymentNotFound error")
				return
			}
		})

	When(`^I call ListDeployments with environment "(.+?)", it should fail with NotFound$`,
		func(e string) {
			_, err := edsWrapper.ListDeployments(aws.String(e))
			if err != nil {
				_, ok := err.(*operations.ListDeploymentsNotFound)
				if !ok {
					T.Errorf(err.Error())
					return
				}
			} else {
				T.Errorf("Expecting ListDeploymentsNotFound error")
				return
			}
		})

	And(`^I call CreateDeployment API (.+?) times$`, func(count int) {
		for i := 0; i < count; i++ {
			deployment := createDeployment("", ctx, edsWrapper)
			waitForDeploymentToComplete(deploymentCompleteWaitSeconds, edsWrapper)
			deploymentIDs[*deployment.ID] = deployment
		}
	})

	And(`^ListDeployments should return (.+?) deployment$`, func(count int) {
		deployments, err := edsWrapper.ListDeployments(aws.String(environment))
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		assert.True(T, count <= len(deployments), "Wrong number of deployments returned")
		deploymentsFromResponse := make(map[string]bool)
		for _, d := range deployments {
			deploymentsFromResponse[*d.ID] = true
		}
		for key, _ := range deploymentIDs {
			if !deploymentsFromResponse[key] {
				T.Errorf("Did not find deployment with id:%s under environment:%s", key, environment)
			}
		}
	})

	When(`^I delete the environment$`, func() {
		deleteEnvironment(environment, edsWrapper)
	})

	And(`^deleting the environment again should succeed$`, func() {
		deleteEnvironment(environment, edsWrapper)
	})

	Then(`^get environment should return empty$`, func() {
		_, err := edsWrapper.GetEnvironment(&environment)
		if err != nil {
			_, ok := err.(*operations.GetEnvironmentNotFound)
			if !ok {
				T.Errorf(err.Error())
				return
			}
		} else {
			T.Errorf("Expecting GetEnvironmentNotFound error")
			return
		}
	})
}

func deleteEnvironment(environment string, edsWrapper wrappers.EDSWrapper) {
	err := edsWrapper.DeleteEnvironment(&environment)
	if err != nil {
		T.Errorf("Was not able to delete environment %v: %v", environment, err.Error())
		return
	}
}

func createDeployment(deploymentToken string, ctx context.Context, edsWrapper wrappers.EDSWrapper) *models.Deployment {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	deployment, err := edsWrapper.CreateDeployment(opCtx, &environment, &deploymentToken)
	if err != nil {
		T.Errorf(err.Error())
		return nil
	}
	deploymentID = *deployment.ID
	return deployment
}

func waitForDeploymentToComplete(seconds int, edsWrapper wrappers.EDSWrapper) {
	ok, err := doSomething(time.Duration(seconds)*time.Second, 1*time.Second, func() (bool, error) {
		deployment, err := edsWrapper.GetDeployment(aws.String(environment), aws.String(deploymentID))
		if err != nil {
			return false, errors.Wrapf(err, "Error calling GetDeployment for environment %s and deployment %s", environment, deploymentID)
		}

		return strings.ToLower(taskRunning) == strings.ToLower(aws.StringValue(deployment.Status)), nil
	})

	if err != nil {
		T.Errorf(err.Error())
		return
	}

	if !ok {
		T.Errorf("Expecting the deployment status to be %v", taskRunning)
		return
	}
}

func doSomething(ttl time.Duration, tickTime time.Duration, fn func() (bool, error)) (bool, error) {
	timeout := time.After(ttl)
	tick := time.Tick(tickTime)
	for {
		select {
		case <-timeout:
			return false, errors.New("timed out")
		case <-tick:
			ok, err := fn()
			if err != nil {
				return false, err
			} else if ok {
				return true, nil
			}
		}
	}
}

func filterTasksByStatusRunning(cluster *string, taskARNs []*string, ecsWrapper wrappers.ECSWrapper) []*string {
	runningTasks := make([]*string, len(taskARNs))
	if len(taskARNs) == 0 {
		return runningTasks
	}
	tasks, err := ecsWrapper.DescribeTasks(cluster, taskARNs)
	if err != nil {
		T.Errorf(err.Error())
	}
	for _, t := range tasks {
		if aws.StringValue(t.LastStatus) == taskRunning {
			runningTasks = append(runningTasks, t.TaskArn)
		}
	}
	return runningTasks
}
