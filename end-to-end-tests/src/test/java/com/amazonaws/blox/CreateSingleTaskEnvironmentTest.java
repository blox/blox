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

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.model.DeploymentConfiguration;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentType;
import com.amazonaws.blox.dataservicemodel.v1.model.InstanceGroup;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateTargetEnvironmentRevisionRequest;
import com.amazonaws.blox.integ.BloxTestStack;
import java.util.UUID;
import lombok.extern.log4j.Log4j;
import org.junit.After;
import org.junit.Before;
import org.junit.BeforeClass;
import org.junit.Test;

@Log4j
public class CreateSingleTaskEnvironmentTest {
  private static final String ENVIRONMENT_NAME = "EndToEndTestEnvironment_" + UUID.randomUUID();
  private static final String ACCOUNT_ID = "012345789";
  private static final long RECONCILIATION_INTERVAL = 60_000;

  private static BloxTestStack stack;

  private DataService backend = stack.createDataService();

  @BeforeClass
  public static void initializeStack() {
    stack = new BloxTestStack();
  }

  @Before
  @After
  public void resetStack() {
    stack.reset();
  }

  @Test
  public void creatingSingleTaskEnvironmentShouldLaunchSingleTask() throws Exception {
    // TODO: Call the frontend APIs instead of DataService once it's wired up:
    CreateEnvironmentResponse environment =
        backend.createEnvironment(
            CreateEnvironmentRequest.builder()
                .environmentName(ENVIRONMENT_NAME)
                .environmentType(EnvironmentType.SingleTask)
                .taskDefinition(stack.getTaskDefinition())
                .instanceGroup(InstanceGroup.builder().cluster(stack.getCluster()).build())
                .deploymentConfiguration(DeploymentConfiguration.builder().build())
                .role("test-role")
                .accountId(ACCOUNT_ID)
                .build());

    backend.createTargetEnvironmentRevision(
        CreateTargetEnvironmentRevisionRequest.builder()
            .environmentId(environment.getEnvironmentId())
            .environmentVersion(environment.getEnvironmentVersion())
            .build());

    waitOrTimeout(
        RECONCILIATION_INTERVAL * 3 / 2,
        () -> {
          assertThat(stack.describeTasks())
              .as("Tasks launched by blox")
              .extracting("group", "taskDefinitionArn")
              .containsExactly(
                  tuple(environment.getEnvironmentId(), environment.getTaskDefinition()));
        });
  }
}
