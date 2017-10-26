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

import static org.hamcrest.Matchers.containsInAnyOrder;
import static org.hamcrest.Matchers.empty;
import static org.junit.Assert.assertThat;
import static org.mockito.Matchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Collectors;
import org.junit.Test;
import org.mockito.Mockito;
import software.amazon.awssdk.services.ecs.ECSAsyncClient;
import software.amazon.awssdk.services.ecs.model.DescribeTasksRequest;
import software.amazon.awssdk.services.ecs.model.DescribeTasksResponse;
import software.amazon.awssdk.services.ecs.model.ListTasksRequest;
import software.amazon.awssdk.services.ecs.model.ListTasksResponse;
import software.amazon.awssdk.services.ecs.model.Task;

public class TaskListerTest {

  private static final String CLUSTER_ARN = "cluster1";

  // TODO: Add tests for synchronous/asynchronous failures

  @Test
  public void returnsEmptyListWhenNoTasks() {
    ECSAsyncClient ecs = mock(ECSAsyncClient.class);
    when(ecs.listTasks(any()))
        .thenReturn(
            CompletableFuture.completedFuture(ListTasksResponse.builder().taskArns().build()));

    TaskLister tasks = new TaskLister(ecs, CLUSTER_ARN);

    List<ClusterSnapshot.Task> describe = tasks.describe().join();

    assertThat(describe, empty());
  }

  @Test
  public void returnsTasksFromAllPages() {
    ECSAsyncClient ecs = mock(FakeECSAsyncClient.class, Mockito.CALLS_REAL_METHODS);

    TaskLister tasks = new TaskLister(ecs, CLUSTER_ARN);

    List<ClusterSnapshot.Task> describe = tasks.describe().join();

    assertThat(
        describe.stream().map(ClusterSnapshot.Task::getArn).collect(Collectors.toList()),
        containsInAnyOrder("1", "2", "3", "4", "5", "6"));
  }

  public abstract class FakeECSAsyncClient implements ECSAsyncClient {
    @Override
    public CompletableFuture<DescribeTasksResponse> describeTasks(DescribeTasksRequest request) {
      return CompletableFuture.completedFuture(
          DescribeTasksResponse.builder()
              .tasks(
                  request
                      .tasks()
                      .stream()
                      .map(arn -> Task.builder().taskArn(arn).build())
                      .collect(Collectors.toList()))
              .build());
    }

    @Override
    public CompletableFuture<ListTasksResponse> listTasks(ListTasksRequest request) {
      if (request.nextToken() == null) {
        return CompletableFuture.completedFuture(
            ListTasksResponse.builder().nextToken("1").taskArns("1", "2").build());
      }
      if (request.nextToken() == "1") {
        return CompletableFuture.completedFuture(
            ListTasksResponse.builder().nextToken("2").taskArns("3", "4").build());
      }
      if (request.nextToken() == "2") {
        return CompletableFuture.completedFuture(
            ListTasksResponse.builder().taskArns("5", "6").build());
      }
      throw new UnsupportedOperationException("invalid nexttoken " + request.nextToken());
    }
  }
}
