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
import lombok.extern.log4j.Log4j2;
import software.amazon.awssdk.services.ecs.ECSAsyncClient;
import software.amazon.awssdk.services.ecs.model.StartTaskRequest;
import software.amazon.awssdk.services.ecs.model.StartTaskResponse;

@Value
@Builder
@Log4j2
public class StartTask implements SchedulingAction {
  private final String clusterArn;
  private final String containerInstanceArn;
  private final String taskDefinitionArn;
  private final String group;

  @Override
  public CompletableFuture<Boolean> execute(ECSAsyncClient ecs) {
    CompletableFuture<StartTaskResponse> pendingRequest =
        ecs.startTask(
            StartTaskRequest.builder()
                .cluster(clusterArn)
                .containerInstances(containerInstanceArn)
                .taskDefinition(taskDefinitionArn)
                .group(group)
                .startedBy("blox")
                .build());

    pendingRequest.thenAccept(r -> log.debug("ECS response: {}", r));

    // TODO: We probably need richer error reporting than a Boolean
    return pendingRequest.thenApply(startTaskResponse -> startTaskResponse.failures().size() > 0);
  }
}
