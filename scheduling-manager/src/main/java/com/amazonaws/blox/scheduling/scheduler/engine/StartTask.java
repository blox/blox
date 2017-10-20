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
import lombok.Value;
import software.amazon.awssdk.services.ecs.ECSAsyncClient;
import software.amazon.awssdk.services.ecs.model.StartTaskRequest;
import software.amazon.awssdk.services.ecs.model.StartTaskResponse;

@Value
public class StartTask implements SchedulingAction {
  private final String clusterArn;
  private final String containerInstanceArn;
  private final String taskDefinitionArn;

  @Override
  public CompletableFuture<Boolean> execute(ECSAsyncClient ecs) {
    // TODO: This will probably require setting the task group to the environment ID too
    CompletableFuture<StartTaskResponse> pendingRequest =
        ecs.startTask(
            StartTaskRequest.builder()
                .cluster(clusterArn)
                .containerInstances(containerInstanceArn)
                .taskDefinition(taskDefinitionArn)
                .build());

    // TODO: We probably need richer error reporting than a Boolean
    return pendingRequest.thenApply(startTaskResponse -> startTaskResponse.failures().size() == 0);
  }
}
