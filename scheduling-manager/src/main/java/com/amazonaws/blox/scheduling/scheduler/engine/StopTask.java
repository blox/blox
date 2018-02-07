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
package com.amazonaws.blox.scheduling.scheduler.engine;

import java.util.concurrent.CompletableFuture;
import lombok.Builder;
import lombok.Value;
import lombok.extern.slf4j.Slf4j;
import software.amazon.awssdk.services.ecs.ECSAsyncClient;
import software.amazon.awssdk.services.ecs.model.StopTaskRequest;
import software.amazon.awssdk.services.ecs.model.StopTaskResponse;

@Value
@Builder
@Slf4j
public class StopTask implements SchedulingAction {
  private final String clusterName;
  private final String task;
  private final String reason;

  @Override
  public CompletableFuture<Boolean> execute(ECSAsyncClient ecs) {
    CompletableFuture<StopTaskResponse> pendingRequest =
        ecs.stopTask(
            StopTaskRequest.builder().cluster(clusterName).task(task).reason(reason).build());

    pendingRequest.thenAccept(r -> log.debug("ECS response: {}", r));

    return pendingRequest.thenApply(stopTaskResponse -> stopTaskResponse.task() != null);
  }
}
