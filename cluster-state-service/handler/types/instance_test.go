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

package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersionValidInstance(t *testing.T) {
	ins := ContainerInstance{}

	version := int64(1)
	instance := ContainerInstance{
		Detail: &InstanceDetail{
			Version: &version,
		},
	}
	instanceJSON := marshalInstance(t, instance)

	extractedVersion, err := ins.GetVersion(instanceJSON)
	assert.Nil(t, err, "Unexpected error getting instance version")
	assert.Equal(t, version, extractedVersion, "Invalid version extracted from instance JSON")
}

func TestGetVersionEmptyInstanceDetail(t *testing.T) {
	ins := ContainerInstance{}

	instance := ContainerInstance{}
	instanceJSON := marshalInstance(t, instance)

	_, err := ins.GetVersion(instanceJSON)
	assert.NotNil(t, err, "Expected an error getting instance version for an instance with no detail")
}

func TestGetVersionEmptyInstanceVersion(t *testing.T) {
	ins := ContainerInstance{}

	instance := ContainerInstance{
		Detail: &InstanceDetail{},
	}
	instanceJSON := marshalInstance(t, instance)

	_, err := ins.GetVersion(instanceJSON)
	assert.NotNil(t, err, "Expected an error getting instance version for an instance with no version")
}

func TestGerVersionInvalidInstance(t *testing.T) {
	ins := ContainerInstance{}
	_, err := ins.GetVersion("invalidInstanceJSON")
	assert.NotNil(t, err, "Expected an error getting instance version for an invalid instance")
}

func marshalInstance(t *testing.T, ins ContainerInstance) string {
	insJSON, err := json.Marshal(ins)
	assert.Nil(t, err, "Unexpected error marshaling instance")
	return string(insJSON)
}
