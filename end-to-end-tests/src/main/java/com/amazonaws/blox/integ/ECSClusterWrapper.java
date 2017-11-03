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
package com.amazonaws.blox.integ;

import com.amazonaws.blox.integ.CloudFormationStacks.CfnStack;
import java.util.Collections;
import java.util.List;
import lombok.RequiredArgsConstructor;
import software.amazon.awssdk.services.ecs.ECSClient;
import software.amazon.awssdk.services.ecs.model.DescribeTasksRequest;
import software.amazon.awssdk.services.ecs.model.DesiredStatus;
import software.amazon.awssdk.services.ecs.model.ListTasksRequest;
import software.amazon.awssdk.services.ecs.model.StopTaskRequest;
import software.amazon.awssdk.services.ecs.model.Task;

/** Wrapper for interacting with a test ECS cluster */
@RequiredArgsConstructor
public class ECSClusterWrapper {
  private final ECSClient ecs;

  // TODO: For now, act on all tasks that match startedBy, we should change this to filter by prefix
  private final String startedBy = "blox";

  private final CfnStack stack;

  public ECSClusterWrapper(ECSClient ecs, CloudFormationStacks stacks) {
    this(ecs, stacks.get("blox-test-cluster"));
  }

  public String getTaskDefinition() {
    return stack.output("taskdef");
  }

  public String getCluster() {
    return stack.output("cluster");
  }

  public List<Task> describeTasks() {
    List<String> taskArns = listTasks();
    if (taskArns.isEmpty()) {
      return Collections.emptyList();
    }

    return ecs.describeTasks(
            DescribeTasksRequest.builder().cluster(getCluster()).tasks(taskArns).build())
        .tasks();
  }

  private List<String> listTasks() {
    return ecs.listTasks(
            ListTasksRequest.builder()
                .cluster(getCluster())
                .startedBy(startedBy)
                .desiredStatus(DesiredStatus.RUNNING)
                .build())
        .taskArns();
  }

  public void reset() {
    for (String task : listTasks()) {
      ecs.stopTask(StopTaskRequest.builder().cluster(getCluster()).task(task).build());
    }
  }
}
