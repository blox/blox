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

package cmd

import (
	"strings"
	"testing"

	"github.com/blox/blox/cluster-state-service/config"
	"github.com/stretchr/testify/assert"
)

func TestRootCommandFailsWithUnknownFlag(t *testing.T) {
	rootCmd := createRootCommand()
	rootCmd.SetArgs(strings.Split("--unknown", " "))
	assert.Error(t, rootCmd.Execute(), "Expected error processing an unknown flag")
}

func TestRootCommandWithSQSName(t *testing.T) {
	rootCmd := createRootCommand()
	rootCmd.SetArgs(strings.Split("--queue q", " "))
	assert.NoError(t, rootCmd.Execute(), "Error processing the --queue flag")
	assert.Equal(t, config.QueueName, "q", "Unexpected queue name set")
}

func TestRootCommandWithOneEtcdEndpoint(t *testing.T) {
	rootCmd := createRootCommand()
	rootCmd.SetArgs(strings.Split("--etcd-endpoint e1", " "))
	assert.NoError(t, rootCmd.Execute(), "Error processing the --etcd-endpoint flag")
	assert.Equal(t, config.EtcdEndpoints, []string{"e1"}, "Unexpected etcd endpoint set")
}

func TestRootCommandWithOneMultipleEndpoints(t *testing.T) {
	rootCmd := createRootCommand()
	rootCmd.SetArgs(strings.Split("--etcd-endpoint e1 --etcd-endpoint e2 --etcd-endpoint e3", " "))
	assert.NoError(t, rootCmd.Execute(), "Error processing the --etcd-endpoint flag")
	assert.Equal(t, config.EtcdEndpoints, []string{"e1", "e2", "e3"}, "Unexpected etcd endpoint set")
}
