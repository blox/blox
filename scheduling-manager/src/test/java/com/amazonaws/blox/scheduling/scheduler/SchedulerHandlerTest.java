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
import static org.junit.Assert.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verify;
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
import com.amazonaws.blox.scheduling.scheduler.engine.Scheduler;
import com.amazonaws.blox.scheduling.scheduler.engine.SchedulerFactory;
import com.amazonaws.blox.scheduling.scheduler.engine.SchedulingAction;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot;

import java.time.Instant;
import java.util.Arrays;
import java.util.Collections;
import java.util.concurrent.CompletableFuture;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.runners.MockitoJUnitRunner;
import software.amazon.awssdk.services.ecs.ECSAsyncClient;

@RunWith(MockitoJUnitRunner.class)
public class SchedulerHandlerTest {

  private static final String ACCOUNT_ID = "123456789012";
  private static final String CLUSTER_NAME = "cluster1";
  private static final String ENVIRONMENT_NAME = "environment1";
  private static final String ACTIVE_ENVIRONMENT_REVISION_ID = "1";
  private static final String DEPLOYMENT_METHOD = "ReplaceAfterTerminate";
  private static final String TASK_DEFINITION = "arn:::::task:1";

  @Mock private SchedulerFactory schedulerFactory;
  @Mock private DataService dataService;
  @Mock private ECSAsyncClient ecs;

  @Test
  public void invokesSchedulerCoreForDeploymentMethod() throws Exception {
    EnvironmentId environmentId =
        EnvironmentId.builder()
            .accountId(ACCOUNT_ID)
            .cluster(CLUSTER_NAME)
            .environmentName(ENVIRONMENT_NAME)
            .build();

    DescribeEnvironmentRequest describeEnvironmentRequest =
        DescribeEnvironmentRequest.builder().environmentId(environmentId).build();
    DescribeEnvironmentRevisionRequest describeEnvironmentRevisionRequest =
        DescribeEnvironmentRevisionRequest.builder()
            .environmentId(environmentId)
            .environmentRevisionId(ACTIVE_ENVIRONMENT_REVISION_ID)
            .build();

    when(dataService.describeEnvironment(describeEnvironmentRequest))
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
    when(dataService.describeEnvironmentRevision(describeEnvironmentRevisionRequest))
        .thenReturn(
            DescribeEnvironmentRevisionResponse.builder()
                .environmentRevision(
                    EnvironmentRevision.builder()
                        .environmentId(environmentId)
                        .environmentRevisionId(ACTIVE_ENVIRONMENT_REVISION_ID)
                        .taskDefinition(TASK_DEFINITION)
                        .createdTime(Instant.now())
                        .build())
                .build());

    ClusterSnapshot snapshot =
        new ClusterSnapshot(CLUSTER_NAME, Collections.emptyList(), Collections.emptyList());

    SchedulingAction successfulAction = e -> CompletableFuture.completedFuture(true);
    SchedulingAction failedAction = e -> CompletableFuture.completedFuture(false);

    Scheduler fakeScheduler = (s, environment) -> Arrays.asList(successfulAction, failedAction);

    when(schedulerFactory.schedulerFor(any())).thenReturn(fakeScheduler);

    SchedulerHandler handler = new SchedulerHandler(dataService, ecs, schedulerFactory);

    SchedulerOutput output =
        handler.handleRequest(new SchedulerInput(snapshot, environmentId), null);

    verify(dataService).describeEnvironment(describeEnvironmentRequest);
    verify(dataService).describeEnvironmentRevision(describeEnvironmentRevisionRequest);

    assertThat(output.getClusterName(), is(CLUSTER_NAME));
    assertThat(output.getEnvironmentId(), is(environmentId));
    assertThat(output.getSuccessfulActions(), is(1L));
    assertThat(output.getFailedActions(), is(1L));
  }
}
