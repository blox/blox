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
package com.amazonaws.blox.schedulingmanager.wrapper;

import com.amazonaws.AmazonClientException;
import com.amazonaws.services.ecs.AmazonECS;
import com.amazonaws.services.ecs.model.DescribeTasksRequest;
import com.amazonaws.services.ecs.model.DescribeTasksResult;
import com.amazonaws.services.ecs.model.ListContainerInstancesRequest;
import com.amazonaws.services.ecs.model.ListContainerInstancesResult;
import com.amazonaws.services.ecs.model.StartTaskRequest;
import com.amazonaws.services.ecs.model.StartTaskResult;
import java.util.List;
import lombok.AllArgsConstructor;
import lombok.NonNull;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@AllArgsConstructor
public class ECSWrapper {

  @NonNull private final AmazonECS ecs;

  public StartTaskResult startTask(
      final String taskDefinition, final String containerInstance, final String cluster) {
    final StartTaskRequest startTaskRequest =
        new StartTaskRequest()
            .withTaskDefinition(taskDefinition)
            .withContainerInstances(containerInstance)
            .withCluster(cluster);

    try {
      return ecs.startTask(startTaskRequest);
    } catch (final AmazonClientException e) {
      log.error(
          "Could not start task with task definition {} on instance {} in cluster {}",
          startTaskRequest.getTaskDefinition(),
          startTaskRequest.getContainerInstances(),
          startTaskRequest.getCluster(),
          e);
      throw e;
    }
  }

  public DescribeTasksResult describeTasks(final List<String> tasks, final String cluster) {
    final DescribeTasksRequest describeTasksRequest =
        new DescribeTasksRequest().withTasks(tasks).withCluster(cluster);

    try {
      return ecs.describeTasks(describeTasksRequest);
    } catch (final AmazonClientException e) {
      log.error(
          "Could not describe tasks {} in cluster {}",
          describeTasksRequest.getTasks(),
          describeTasksRequest.getCluster(),
          e);
      throw e;
    }
  }

  public ListContainerInstancesResult listInstances(final String cluster) {
    final ListContainerInstancesRequest listContainerInstancesRequest =
        new ListContainerInstancesRequest().withCluster(cluster);

    try {
      return ecs.listContainerInstances(listContainerInstancesRequest);
    } catch (final AmazonClientException e) {
      log.error("Could not list instances in cluster {}", cluster);
      throw e;
    }
  }
}
