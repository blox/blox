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

import static com.amazonaws.blox.integ.AsynchronousTestSupport.waitOrTimeout;
import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.groups.Tuple.tuple;

import com.amazonaws.blox.integ.BloxTestStack;
import com.amazonaws.blox.model.CreateEnvironmentRequest;
import com.amazonaws.blox.model.DeleteEnvironmentRequest;
import com.amazonaws.blox.model.DeploymentConfiguration;
import com.amazonaws.blox.model.DescribeEnvironmentRequest;
import com.amazonaws.blox.model.DescribeEnvironmentRevisionRequest;
import com.amazonaws.blox.model.Environment;
import com.amazonaws.blox.model.EnvironmentRevision;
import com.amazonaws.blox.model.StartDeploymentRequest;
import java.util.UUID;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

public class CreateDaemonEnvironmentTest {

  private String environmentName;
  private static final long RECONCILIATION_INTERVAL = 60_000;

  private BloxTestStack stack;

  @Before
  public void setUp() {
    final String bloxEndpoint = System.getProperty("blox.tests.apiUrl");
    stack = new BloxTestStack(bloxEndpoint);

    environmentName = "EndToEndTestEnvironment_" + UUID.randomUUID();
  }

  @After
  public void tearDown() {
    // Delete environment
    stack
        .getBlox()
        .deleteEnvironment(
            new DeleteEnvironmentRequest()
                .cluster(stack.getCluster())
                .environmentName(environmentName));
    // Cleanup ECS tasks
    stack.reset();
  }

  @Test
  public void creatingDaemonEnvironmentShouldLaunchTaskPerInstance() throws Exception {
    // Create environment
    final String revisionId =
        stack
            .getBlox()
            .createEnvironment(
                new CreateEnvironmentRequest()
                    .environmentName(environmentName)
                    .taskDefinition(stack.getTaskDefinition())
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

    // The names should match
    assertThat(environment.getEnvironmentName()).as("environment name").isEqualTo(environmentName);
    // and the active environment revision should not set
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
    assertThat(revision.getTaskDefinition()).isEqualTo(stack.getTaskDefinition());

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

    waitOrTimeout(
        RECONCILIATION_INTERVAL * 3 / 2,
        () -> {
          assertThat(stack.describeTasks())
              .as("Tasks launched by blox")
              .extracting("group", "taskDefinitionArn")
              .containsExactly(tuple(environment.getEnvironmentName(), stack.getTaskDefinition()));
        });
  }
}
