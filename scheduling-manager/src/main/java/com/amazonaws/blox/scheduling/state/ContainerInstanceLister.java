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

import com.amazonaws.blox.scheduling.state.ClusterSnapshot.ContainerInstance;
import com.spotify.futures.CompletableFutures;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Collectors;
import java.util.stream.Stream;
import lombok.RequiredArgsConstructor;
import software.amazon.awssdk.services.ecs.ECSAsyncClient;
import software.amazon.awssdk.services.ecs.model.DescribeContainerInstancesRequest;
import software.amazon.awssdk.services.ecs.model.DescribeContainerInstancesResponse;
import software.amazon.awssdk.services.ecs.model.ListContainerInstancesRequest;
import software.amazon.awssdk.services.ecs.model.ListContainerInstancesResponse;

/**
 * Wrapper around paginated ECS list/describe ContainerInstance APIs that efficiently describes all
 * tasks.
 */
@RequiredArgsConstructor
class ContainerInstanceLister {
  private final ECSAsyncClient ecs;
  private final ListContainerInstancesRequest.Builder listRequest;
  private final DescribeContainerInstancesRequest.Builder describeRequest;

  public ContainerInstanceLister(ECSAsyncClient ecs, String clusterArn) {
    this(
        ecs,
        ListContainerInstancesRequest.builder().cluster(clusterArn),
        DescribeContainerInstancesRequest.builder().cluster(clusterArn));
  }

  protected ListContainerInstancesResponse list(String nextToken) {
    CompletableFuture<ListContainerInstancesResponse> r =
        ecs.listContainerInstances(listRequest.copy().nextToken(nextToken).build());
    return r.join();
  }

  public CompletableFuture<List<ContainerInstance>> describe() {
    return new PaginatedResponseIterator<>(ListContainerInstancesResponse::nextToken, this::list)
        .stream()
        .map(ListContainerInstancesResponse::containerInstanceArns)
        .filter(arns -> !arns.isEmpty())
        .map(this::describeContainerInstances)
        .collect(CompletableFutures.joinList())
        .thenApply(this::extractContainerInstancesFromResponses);
  }

  private CompletableFuture<DescribeContainerInstancesResponse> describeContainerInstances(
      List<String> arns) {
    return ecs.describeContainerInstances(describeRequest.copy().containerInstances(arns).build());
  }

  private List<ContainerInstance> extractContainerInstancesFromResponses(
      List<DescribeContainerInstancesResponse> r) {
    return r.stream()
        .flatMap(this::extractContainerInstancesFromResponse)
        .collect(Collectors.toList());
  }

  private Stream<ContainerInstance> extractContainerInstancesFromResponse(
      DescribeContainerInstancesResponse r) {
    return r.containerInstances().stream().map(ContainerInstance::from);
  }
}
