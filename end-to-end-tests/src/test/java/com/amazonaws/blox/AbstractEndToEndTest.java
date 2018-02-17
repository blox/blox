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

import com.amazonaws.blox.integ.BloxTestStack;
import com.amazonaws.blox.model.DeleteEnvironmentRequest;
import java.util.UUID;
import org.junit.After;
import org.junit.Before;

public abstract class AbstractEndToEndTest {

  private static final long RECONCILIATION_INTERVAL = 60_000;
  String environmentName;
  BloxTestStack stack;

  @Before
  public void setupBase() {
    final String bloxEndpoint = System.getProperty("blox.tests.apiUrl");
    stack = new BloxTestStack(bloxEndpoint);

    environmentName = "EndToEndTestEnvironment_" + UUID.randomUUID();
  }

  @After
  public void tearDownBase() {
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

  void assertAllRunningTasksMatch(final String environmentName, final String taskDefinition)
      throws InterruptedException {
    waitOrTimeout(
        RECONCILIATION_INTERVAL * 3 / 2,
        () ->
            assertThat(stack.describeTasks())
                .as("Tasks launched by blox")
                .allSatisfy(
                    t ->
                        assertThat(t)
                            .hasFieldOrPropertyWithValue("group", environmentName)
                            .hasFieldOrPropertyWithValue("taskDefinitionArn", taskDefinition)));
  }
}
