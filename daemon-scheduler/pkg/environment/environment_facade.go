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

package environment

import (
	"github.com/blox/blox/daemon-scheduler/pkg/environment/types"
	"github.com/blox/blox/daemon-scheduler/pkg/facade"
	"github.com/pkg/errors"
)

type EnvironmentFacade interface {
	InstanceARNs(environment *types.Environment) ([]*string, error)
}

type environmentFacade struct {
	css facade.ClusterState
}

func NewEnvironmentFacade(css facade.ClusterState) (EnvironmentFacade, error) {
	if css == nil {
		return nil, errors.New("Cluster state service facade should not be nil")
	}

	return environmentFacade{
		css: css,
	}, nil
}

func (f environmentFacade) InstanceARNs(environment *types.Environment) ([]*string, error) {
	if environment.Cluster == "" {
		return nil, errors.New("Environment cluster name is required")
	}

	instances, err := f.css.ListInstances(environment.Cluster)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting instances for cluster %s in environment %s", environment.Cluster, environment.Name)
	}

	instanceARNs := make([]*string, 0, len(instances))
	for _, instance := range instances {
		instanceARNs = append(instanceARNs, instance.Entity.ContainerInstanceARN)
	}

	return instanceARNs, nil
}