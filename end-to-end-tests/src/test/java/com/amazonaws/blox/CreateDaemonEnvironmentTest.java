/*
 * Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"). You may
 * not use this file except in compliance with the License. A copy of the
 * License is located at
 *
 *     http://aws.amazon.com/apache2.0/
 *
 * or in the "LICENSE" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */
package com.amazonaws.blox;

import static org.assertj.core.api.Assertions.assertThat;

import com.amazonaws.blox.model.CreateEnvironmentRequest;
import com.amazonaws.blox.model.DeploymentConfiguration;
import com.amazonaws.blox.model.DescribeEnvironmentRequest;
import com.amazonaws.blox.model.DescribeEnvironmentRevisionRequest;
import com.amazonaws.blox.model.Environment;
import com.amazonaws.blox.model.EnvironmentRevision;
import com.amazonaws.blox.model.StartDeploymentRequest;
import org.junit.Test;

public class CreateDaemonEnvironmentTest extends AbstractEndToEndTest {

  @Test
  public void creatingDaemonEnvironmentShouldLaunchTaskPerInstance() throws Exception {
    final String taskDefinition = stack.getTransientTaskDefinition();
    // Create environment
    final String revisionId =
        stack
            .getBlox()
            .createEnvironment(
                new CreateEnvironmentRequest()
                    .environmentName(environmentName)
                    .taskDefinition(taskDefinition)
                    .deploymentConfiguration(new DeploymentConfiguration())
                    .cluster(stack.getCluster())
                    .environmentType("Daemon")
                    .role("Test")
                    .deploymentMethod("ReplaceAfterTerminate"))
            .getEnvironmentRevisionId();

    // Then when I describe that environment ...
    Environment environment =
        stack
            .getBlox()
            .describeEnvironment(
                new DescribeEnvironmentRequest()
                    .cluster(stack.getCluster())
                    .environmentName(environmentName))
            .getEnvironment();

    assertThat(environment.getEnvironmentName()).as("environment name").isEqualTo(environmentName);
    assertThat(environment.getActiveEnvironmentRevisionId()).isNull();

    // And when I describe the revision
    EnvironmentRevision revision =
        stack
            .getBlox()
            .describeEnvironmentRevision(
                new DescribeEnvironmentRevisionRequest()
                    .cluster(stack.getCluster())
                    .environmentName(environmentName)
                    .environmentRevisionId(revisionId))
            .getEnvironmentRevision();

    // The task definition should match
    assertThat(revision.getTaskDefinition()).isEqualTo(stack.getTransientTaskDefinition());

    // Now start deployment
    final String deploymentId =
        stack
            .getBlox()
            .startDeployment(
                new StartDeploymentRequest()
                    .cluster(stack.getCluster())
                    .environmentName(environmentName)
                    .revisionId(revisionId))
            .getDeploymentId();
    assertThat(deploymentId).isNotEmpty();

    assertAllRunningTasksMatch(environmentName, taskDefinition);
  }
}
