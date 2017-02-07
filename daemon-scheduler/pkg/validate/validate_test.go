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

package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	validClusterName = "clust_er-1"
	validClusterARN  = "arn:aws:ecs:us-east-1:123456789123:cluster/" + validClusterName

	invalidClusterName                 = "cluster1/cluster1"
	invalidClusterARNWithNoName        = "arn:aws:ecs:us-east-1:123456789123:cluster/"
	invalidClusterARNWithInvalidName   = "arn:aws:ecs:us-east-1:123456789123:cluster/" + invalidClusterName
	invalidClusterARNWithInvalidPrefix = "arn/cluster"
)

func TestIsClusterNameEmptyName(t *testing.T) {
	isValid := IsClusterName("")
	assert.False(t, isValid, "Empty cluster name should not satisfy regex")
}

func TestIsClusterNameInvalidName(t *testing.T) {
	isValid := IsClusterName(invalidClusterName)
	assert.False(t, isValid, "Invalid cluster name should not satisfy regex")
}

func TestIsClusterName(t *testing.T) {
	isValid := IsClusterName(validClusterName)
	assert.True(t, isValid, "Valid cluster name should satisfy regex")
}

func TestIsClusterARNEmptyARN(t *testing.T) {
	isValid := IsClusterARN("")
	assert.False(t, isValid, "Empty cluster ARN should not satisfy regex")
}

func TestIsClusterARNNoNameInARN(t *testing.T) {
	isValid := IsClusterARN(invalidClusterARNWithNoName)
	assert.False(t, isValid, "Invalid cluster ARN with no name should not satisfy regex")
}

func TestIsClusterARNInvalidNameInARN(t *testing.T) {
	isValid := IsClusterARN(invalidClusterARNWithInvalidName)
	assert.False(t, isValid, "Invalid cluster ARN with invalid name should not satisfy regex")
}

func TestIsClusterARNInvalidPrefixInARN(t *testing.T) {
	isValid := IsClusterARN(invalidClusterARNWithInvalidPrefix)
	assert.False(t, isValid, "Invalid cluster ARN with invalid prefix should not satisfy regex")
}

func TestIsClusterARN(t *testing.T) {
	isValid := IsClusterARN(validClusterARN)
	assert.True(t, isValid, "Valid cluster ARN should satisfy regex")
}
