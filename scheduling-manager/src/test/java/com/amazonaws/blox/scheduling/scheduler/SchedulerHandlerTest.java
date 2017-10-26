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
import static org.mockito.Mockito.when;

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentType;
import com.amazonaws.blox.scheduling.FakeDataService;
import com.amazonaws.blox.scheduling.scheduler.engine.Scheduler;
import com.amazonaws.blox.scheduling.scheduler.engine.SchedulerFactory;
import com.amazonaws.blox.scheduling.scheduler.engine.SchedulingAction;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot;
import java.util.Arrays;
import java.util.Collections;
import java.util.concurrent.CompletableFuture;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.runners.MockitoJUnitRunner;
import software.amazon.awssdk.services.ecs.ECSAsyncClient;
import software.amazon.awssdk.services.ecs.model.StartTaskResponse;

@RunWith(MockitoJUnitRunner.class)
public class SchedulerHandlerTest {

  private static final String CLUSTER_ARN = "arn:::::cluster1";
  private static final String ENVIRONMENT_ID = "arn:::::environment1";
  private static final String TASK_DEFINITION = "arn:::::task:1";

  @Mock private SchedulerFactory schedulerFactory;

  private DataService dataService =
      FakeDataService.builder()
          .clusterArn(CLUSTER_ARN)
          .taskDefinition(TASK_DEFINITION)
          .environmentId(ENVIRONMENT_ID)
          .build();

  @Mock private ECSAsyncClient ecs;

  @Test
  public void invokesSchedulerCoreForDeploymentMethod() {
    ClusterSnapshot snapshot =
        new ClusterSnapshot(CLUSTER_ARN, Collections.emptyList(), Collections.emptyList());

    SchedulingAction succesfulAction = e -> CompletableFuture.completedFuture(true);
    SchedulingAction failedAction = e -> CompletableFuture.completedFuture(false);

    Scheduler fakeScheduler = (s, environment) -> Arrays.asList(succesfulAction, failedAction);

    when(schedulerFactory.schedulerFor(EnvironmentType.SingleTask)).thenReturn(fakeScheduler);

    StartTaskResponse successResponse = StartTaskResponse.builder().failures().build();

    SchedulerHandler handler = new SchedulerHandler(dataService, ecs, schedulerFactory);
    SchedulerOutput output =
        handler.handleRequest(new SchedulerInput(snapshot, ENVIRONMENT_ID), null);

    assertThat(output.getSuccessfulActions(), is(1L));
    assertThat(output.getFailedActions(), is(1L));
  }
}
