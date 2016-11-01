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
	assert.Equal(t, clusterName, c, "Invalid cluster name retrieved from ARN")
}
