package regex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
