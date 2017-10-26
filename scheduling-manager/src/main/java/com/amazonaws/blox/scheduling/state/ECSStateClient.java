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
package com.amazonaws.blox.scheduling.state;

import java.util.List;
import java.util.concurrent.CompletableFuture;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Profile;
import org.springframework.stereotype.Component;
import software.amazon.awssdk.services.ecs.ECSAsyncClient;

@Component
@Profile("!test")
@RequiredArgsConstructor
public class ECSStateClient implements ECSState {
  private final ECSAsyncClient ecs;

  @Override
  public ClusterSnapshot snapshotState(String clusterArn) {
    CompletableFuture<List<ClusterSnapshot.Task>> tasks =
        new TaskLister(ecs, clusterArn).describe();
    CompletableFuture<List<ClusterSnapshot.ContainerInstance>> instances =
        new ContainerInstanceLister(ecs, clusterArn).describe();

    return new ClusterSnapshot(clusterArn, tasks.join(), instances.join());
  }
}
