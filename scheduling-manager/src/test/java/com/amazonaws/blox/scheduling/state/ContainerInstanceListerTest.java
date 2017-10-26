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
import software.amazon.awssdk.services.ecs.model.ContainerInstance;
import software.amazon.awssdk.services.ecs.model.DescribeContainerInstancesRequest;
import software.amazon.awssdk.services.ecs.model.DescribeContainerInstancesResponse;
import software.amazon.awssdk.services.ecs.model.ListContainerInstancesRequest;
import software.amazon.awssdk.services.ecs.model.ListContainerInstancesResponse;

public class ContainerInstanceListerTest {

  private static final String CLUSTER_1 = "cluster1";

  // TODO: Add tests for synchronous/asynchronous failures

  @Test
  public void returnsEmptyListWhenNoContainerInstances() {
    ECSAsyncClient ecs = mock(ECSAsyncClient.class);
    when(ecs.listContainerInstances(any()))
        .thenReturn(
            CompletableFuture.completedFuture(
                ListContainerInstancesResponse.builder().containerInstanceArns().build()));

    ContainerInstanceLister containerInstances = new ContainerInstanceLister(ecs, CLUSTER_1);

    List<ClusterSnapshot.ContainerInstance> describe = containerInstances.describe().join();

    assertThat(describe, empty());
  }

  @Test
  public void returnsContainerInstancesFromAllPages() {
    ECSAsyncClient ecs = mock(FakeECSAsyncClient.class, Mockito.CALLS_REAL_METHODS);

    ContainerInstanceLister containerInstances = new ContainerInstanceLister(ecs, CLUSTER_1);

    List<ClusterSnapshot.ContainerInstance> describe = containerInstances.describe().join();

    assertThat(
        describe
            .stream()
            .map(ClusterSnapshot.ContainerInstance::getArn)
            .collect(Collectors.toList()),
        containsInAnyOrder("1", "2", "3", "4", "5", "6"));
  }

  public abstract class FakeECSAsyncClient implements ECSAsyncClient {
    @Override
    public CompletableFuture<DescribeContainerInstancesResponse> describeContainerInstances(
        DescribeContainerInstancesRequest request) {
      return CompletableFuture.completedFuture(
          DescribeContainerInstancesResponse.builder()
              .containerInstances(
                  request
                      .containerInstances()
                      .stream()
                      .map(arn -> ContainerInstance.builder().containerInstanceArn(arn).build())
                      .collect(Collectors.toList()))
              .build());
    }

    @Override
    public CompletableFuture<ListContainerInstancesResponse> listContainerInstances(
        ListContainerInstancesRequest request) {
      if (request.nextToken() == null) {
        return CompletableFuture.completedFuture(
            ListContainerInstancesResponse.builder()
                .nextToken("1")
                .containerInstanceArns("1", "2")
                .build());
      }
      if (request.nextToken() == "1") {
        return CompletableFuture.completedFuture(
            ListContainerInstancesResponse.builder()
                .nextToken("2")
                .containerInstanceArns("3", "4")
                .build());
      }
      if (request.nextToken() == "2") {
        return CompletableFuture.completedFuture(
            ListContainerInstancesResponse.builder().containerInstanceArns("5", "6").build());
      }
      throw new UnsupportedOperationException("invalid nexttoken " + request.nextToken());
    }
  }
}
