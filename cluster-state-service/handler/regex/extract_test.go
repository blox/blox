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

package regex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClusterNameFromARNEmptyARN(t *testing.T) {
	_, err := GetClusterNameFromARN("")
	assert.NotNil(t, err, "Expected an error when retrieving cluster name from empty ARN")
}

func TestGetClusterNameFromARNWithNoName(t *testing.T) {
	_, err := GetClusterNameFromARN(invalidClusterARNWithNoName)
	assert.NotNil(t, err, "Expected an error when retrieving cluster name from ARN with no name")
}

func TestGetClusterNameFromARNWithInvalidName(t *testing.T) {
	_, err := GetClusterNameFromARN(invalidClusterARNWithInvalidName)
	assert.NotNil(t, err, "Expected an error when retrieving cluster name from ARN with invalid name")
}

func TestGetClusterNameFromARNWithInvalidPrefix(t *testing.T) {
	_, err := GetClusterNameFromARN(invalidClusterARNWithInvalidPrefix)
	assert.NotNil(t, err, "Expected an error when retrieving cluster name from ARN with invalid prefix")
}

func TestGetClusterNameFromARN(t *testing.T) {
	c, err := GetClusterNameFromARN(validClusterARN)
	assert.Nil(t, err, "Unexpected error when retrieving cluster name from ARN")
	assert.NotNil(t, c, "Expected cluster name to be retrieved from ARN")
	assert.Equal(t, validClusterName, c, "Invalid cluster name retrieved from ARN")
}
