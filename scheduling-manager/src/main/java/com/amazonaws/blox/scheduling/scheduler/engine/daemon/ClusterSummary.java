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
package com.amazonaws.blox.scheduling.scheduler.engine.daemon;

import com.amazonaws.blox.scheduling.state.ClusterSnapshot;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot.ContainerInstance;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot.Task;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.function.Function;
import java.util.stream.Collectors;

/** Wrapper around {@link ClusterSnapshot} that indexes Tasks by ContainerInstance. */
public class ClusterSummary {
  private final ClusterSnapshot snapshot;
  private final Map<String, List<Task>> tasksByInstanceArn;
  private final Map<String, ContainerInstance> instancesByInstanceArn;

  public ClusterSummary(ClusterSnapshot snapshot) {
    this.snapshot = snapshot;

    this.tasksByInstanceArn =
        snapshot.getTasks().stream().collect(Collectors.groupingBy(Task::getContainerInstanceArn));

    this.instancesByInstanceArn =
        snapshot
            .getInstances()
            .stream()
            .collect(Collectors.toMap(ContainerInstance::getArn, Function.identity()));
  }

  public List<Task> tasksForInstance(ContainerInstance instance) {
    String arn = instance.getArn();

    if (instancesByInstanceArn.containsKey(arn)) {
      return tasksByInstanceArn.getOrDefault(arn, Collections.emptyList());
    } else {
      throw new IndexOutOfBoundsException(
          String.format("ContainerInstance with ARN %s not found in snapshot", arn));
    }
  }

  public List<ContainerInstance> getInstances() {
    return snapshot.getInstances();
  }
}
