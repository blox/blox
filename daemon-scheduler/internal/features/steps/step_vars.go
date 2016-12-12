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

import "github.com/blox/blox/daemon-scheduler/generated/v1/models"

var (
	taskDefinition  string
	cluster         string
	clusterARN      string
	environment     string
	asg             string
	deploymentID    string
	deploymentToken string
	environmentList []*models.Environment
	err             error
)

var deploymentIDs = make(map[string]*models.Deployment)
