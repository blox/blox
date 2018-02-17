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

import com.amazonaws.blox.model.CreateEnvironmentRequest;
import com.amazonaws.blox.model.DeploymentConfiguration;
import com.amazonaws.blox.model.StartDeploymentRequest;
import com.amazonaws.blox.model.UpdateEnvironmentRequest;
import org.junit.Test;

public class UpdateDaemonEnvironmentTest extends AbstractEndToEndTest {

  @Test
  public void updatingEnvironmentCreatesNewRevision() throws Exception {
    // Create environment
    final Blox blox = stack.getBlox();

    final String firstTaskDefinition = stack.getPersistentTaskDefinition();
    final String secondTaskDefinition = stack.getTransientTaskDefinition();

    final String firstRevisionId =
        blox.createEnvironment(
                new CreateEnvironmentRequest()
                    .environmentName(environmentName)
                    .taskDefinition(firstTaskDefinition)
                    .deploymentConfiguration(new DeploymentConfiguration())
                    .cluster(stack.getCluster())
                    .environmentType("Daemon")
                    .role("Test")
                    .deploymentMethod("ReplaceAfterTerminate"))
            .getEnvironmentRevisionId();

    // Now start deployment
    blox.startDeployment(
        new StartDeploymentRequest()
            .cluster(stack.getCluster())
            .environmentName(environmentName)
            .revisionId(firstRevisionId));

    assertAllRunningTasksMatch(environmentName, firstTaskDefinition);

    final String secondRevisionId =
        blox.updateEnvironment(
                new UpdateEnvironmentRequest()
                    .cluster(stack.getCluster())
                    .environmentName(environmentName)
                    .taskDefinition(secondTaskDefinition))
            .getEnvironmentRevisionId();

    blox.startDeployment(
        new StartDeploymentRequest()
            .cluster(stack.getCluster())
            .environmentName(environmentName)
            .revisionId(secondRevisionId));

    assertAllRunningTasksMatch(environmentName, secondTaskDefinition);
  }
}
