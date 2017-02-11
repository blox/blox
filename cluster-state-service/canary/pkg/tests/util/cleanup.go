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

package util

import "github.com/blox/blox/cluster-state-service/canary/pkg/wrappers"

// DeleteCluster deletes a cluster 'clusterName' after cleaning up all the resources within
func DeleteCluster(ecsWrapper wrappers.ECSWrapper, clusterName string) error {
	instanceARNs, err := ecsWrapper.ListContainerInstances(&clusterName)
	if err != nil {
		return NewCleanUpError(err)
	}
	err = ecsWrapper.DeregisterContainerInstances(&clusterName, instanceARNs)
	if err != nil {
		return NewCleanUpError(err)
	}
	return ecsWrapper.DeleteCluster(&clusterName)
}

// TerminateInstances terminates all the EC2 instances corresponding to 'instanceIDs'
func TerminateInstances(ec2Wrapper wrappers.EC2Wrapper, instanceIDs []*string) error {
	err := ec2Wrapper.TerminateInstances(instanceIDs)
	if err != nil {
		return NewCleanUpError(err)
	}
	return nil
}
