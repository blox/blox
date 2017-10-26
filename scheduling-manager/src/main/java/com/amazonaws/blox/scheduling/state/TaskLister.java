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

import com.amazonaws.blox.scheduling.state.ClusterSnapshot.Task;
import com.spotify.futures.CompletableFutures;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Collectors;
import java.util.stream.Stream;
import lombok.RequiredArgsConstructor;
import software.amazon.awssdk.services.ecs.ECSAsyncClient;
import software.amazon.awssdk.services.ecs.model.DescribeTasksRequest;
import software.amazon.awssdk.services.ecs.model.DescribeTasksResponse;
import software.amazon.awssdk.services.ecs.model.ListTasksRequest;
import software.amazon.awssdk.services.ecs.model.ListTasksResponse;

/** Wrapper around paginated ECS list/describe Task APIs that efficiently describes all tasks. */
@RequiredArgsConstructor
class TaskLister {
  private final ECSAsyncClient ecs;
  private final ListTasksRequest.Builder listRequest;
  private final DescribeTasksRequest.Builder describeRequest;

  public TaskLister(ECSAsyncClient ecs, String clusterArn) {
    this(
        ecs,
        ListTasksRequest.builder().cluster(clusterArn),
        DescribeTasksRequest.builder().cluster(clusterArn));
  }

  protected ListTasksResponse list(String nextToken) {
    return ecs.listTasks(listRequest.copy().nextToken(nextToken).build()).join();
  }

  public CompletableFuture<List<Task>> describe() {
    return new PaginatedResponseIterator<>(ListTasksResponse::nextToken, this::list)
        .stream()
        .map(ListTasksResponse::taskArns)
        .filter(arns -> !arns.isEmpty())
        .map(this::describeTasks)
        .collect(CompletableFutures.joinList())
        .thenApply(this::extractTasksFromResponses);
  }

  private CompletableFuture<DescribeTasksResponse> describeTasks(List<String> arns) {
    return ecs.describeTasks(describeRequest.copy().tasks(arns).build());
  }

  private List<Task> extractTasksFromResponses(List<DescribeTasksResponse> r) {
    return r.stream().flatMap(this::extractTasksFromResponse).collect(Collectors.toList());
  }

  private Stream<Task> extractTasksFromResponse(DescribeTasksResponse r) {
    return r.tasks().stream().map(Task::from);
  }
}
