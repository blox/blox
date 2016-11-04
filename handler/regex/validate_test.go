// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the License). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the license file accompanying this file. This file is distributed
// on an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package regex

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestIsTaskARNEmptyARN(t *testing.T) {
	isValid := IsTaskARN("")
	assert.False(t, isValid, "Empty task ARN should not satisfy regex")
}

func TestIsTaskARNNoIDInARN(t *testing.T) {
	isValid := IsTaskARN(invalidTaskARNWithNoID)
	assert.False(t, isValid, "Invalid task ARN with no ID should not satisfy regex")
}

func TestIsTaskARNInvalidIDInARN(t *testing.T) {
	isValid := IsTaskARN(invalidTaskARNWithInvalidID)
	assert.False(t, isValid, "Invalid task ARN with invalid ID should not satisfy regex")
}

func TestIsTaskARNInvalidPrefixInARN(t *testing.T) {
	isValid := IsTaskARN(invalidTaskARNWithInvalidPrefix)
	assert.False(t, isValid, "Invalid task ARN with invalid prefix should not satisfy regex")
}

func TestIsTaskARN(t *testing.T) {
	isValid := IsTaskARN(validTaskARN)
	assert.True(t, isValid, "Valid task ARN should satisfy regex")
}

func TestIsInstanceARNEmptyARN(t *testing.T) {
	isValid := IsInstanceARN("")
	assert.False(t, isValid, "Empty instance ARN should not satisfy regex")
}

func TestIsInstanceARNNoIDInARN(t *testing.T) {
	isValid := IsInstanceARN(invalidInstanceARNWithNoID)
	assert.False(t, isValid, "Invalid instance ARN with no ID should not satisfy regex")
}

func TestIsInstanceARNInvalidIDInARN(t *testing.T) {
	isValid := IsInstanceARN(invalidInstanceARNWithInvalidID)
	assert.False(t, isValid, "Invalid instance ARN with invalid ID should not satisfy regex")
}

func TestIsInstanceARNInvalidPrefixInARN(t *testing.T) {
	isValid := IsInstanceARN(invalidInstanceARNWithInvalidPrefix)
	assert.False(t, isValid, "Invalid instance ARN with invalid prefix should not satisfy regex")
}

func TestIsInstanceARN(t *testing.T) {
	isValid := IsInstanceARN(validInstanceARN)
	assert.True(t, isValid, "Valid instance ARN should satisfy regex")
}
