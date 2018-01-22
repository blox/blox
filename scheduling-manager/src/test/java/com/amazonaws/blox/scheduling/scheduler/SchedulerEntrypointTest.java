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
package com.amazonaws.blox.scheduling.scheduler;

import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.mockito.Matchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.model.DeploymentConfiguration;
import com.amazonaws.blox.dataservicemodel.v1.model.Environment;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentHealth;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentRevision;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentStatus;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentType;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRevisionRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRevisionResponse;
import com.amazonaws.blox.scheduling.LambdaHandlerTestCase;
import com.amazonaws.blox.scheduling.scheduler.SchedulerEntrypointTest.TestConfig;
import com.amazonaws.blox.scheduling.scheduler.engine.SchedulerFactory;

import java.time.Instant;
import java.util.Collections;
import org.junit.Test;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Import;
import org.springframework.test.context.ContextConfiguration;
import software.amazon.awssdk.services.ecs.ECSAsyncClient;

@ContextConfiguration(classes = TestConfig.class)
public class SchedulerEntrypointTest extends LambdaHandlerTestCase {

  @Test
  public void convertsInputsAndOutputsFromJson() throws Exception {
    String result = callHandler(fixture("handlers/Scheduler.input.json"));
    assertThat(result, is(fixtureAsString("handlers/Scheduler.output.json")));
  }

  @Configuration
  @Import(SchedulerApplication.class)
  public static class TestConfig {
    private static final String ACCOUNT_ID = "123456789012";
    private static final String CLUSTER_NAME = "default";
    private static final String ENVIRONMENT_NAME = "SomeEnvironment";
    private static final String ACTIVE_ENVIRONMENT_REVISION_ID = "1";
    private static final String DEPLOYMENT_METHOD = "ReplaceAfterTerminate";
    private static final String taskDefinition = "taskDefinition:1";

    @Bean
    public DataService dataService() throws Exception {
      final EnvironmentId environmentId =
          EnvironmentId.builder()
              .accountId(ACCOUNT_ID)
              .cluster(CLUSTER_NAME)
              .environmentName(ENVIRONMENT_NAME)
              .build();
      final DataService dataService = mock(DataService.class);

      when(dataService.describeEnvironment(
              DescribeEnvironmentRequest.builder().environmentId(environmentId).build()))
          .thenReturn(
              DescribeEnvironmentResponse.builder()
                  .environment(
                      Environment.builder()
                          .environmentId(environmentId)
                          .role("")
                          .environmentType(EnvironmentType.SingleTask)
                          .createdTime(Instant.now())
                          .lastUpdatedTime(Instant.now())
                          .environmentHealth(EnvironmentHealth.HEALTHY)
                          .environmentStatus(EnvironmentStatus.ACTIVE)
                          .deploymentMethod(DEPLOYMENT_METHOD)
                          .deploymentConfiguration(DeploymentConfiguration.builder().build())
                          .activeEnvironmentRevisionId(ACTIVE_ENVIRONMENT_REVISION_ID)
                          .build())
                  .build());
      when(dataService.describeEnvironmentRevision(
              DescribeEnvironmentRevisionRequest.builder()
                  .environmentId(environmentId)
                  .environmentRevisionId(ACTIVE_ENVIRONMENT_REVISION_ID)
                  .build()))
          .thenReturn(
              DescribeEnvironmentRevisionResponse.builder()
                  .environmentRevision(
                      EnvironmentRevision.builder()
                          .environmentId(environmentId)
                          .environmentRevisionId(ACTIVE_ENVIRONMENT_REVISION_ID)
                          .taskDefinition(taskDefinition)
                          .createdTime(Instant.now())
                          .build())
                  .build());

      return dataService;
    }

    @Bean
    public ECSAsyncClient ecs() {
      return mock(ECSAsyncClient.class);
    }

    @Bean
    public SchedulerFactory schedulerFactory() throws Exception {
      return when(mock(SchedulerFactory.class).schedulerFor(any()))
          .thenReturn((snapshot, deploymentConfiguration) -> Collections.emptyList())
          .getMock();
    }
  }
}
