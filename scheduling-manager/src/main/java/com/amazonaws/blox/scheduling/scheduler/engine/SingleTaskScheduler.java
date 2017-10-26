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

import com.amazonaws.blox.scheduling.state.ClusterSnapshot;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot.ContainerInstance;
import java.util.Arrays;
import java.util.Collections;
import java.util.Comparator;
import java.util.List;

/**
 * Temporary scheduler implementation that just starts a single task on the first instance in the
 * snapshot.
 *
 * <p>TODO Remove this, it's only a demo scheduler
 */
public class SingleTaskScheduler implements Scheduler {

  @Override
  public List<SchedulingAction> schedule(
      ClusterSnapshot snapshot, EnvironmentDescription environment) {
    // Naively only run a Task on the first ContainerInstance in the snapshot:
    ContainerInstance instance =
        snapshot
            .getInstances()
            .stream()
            .sorted(Comparator.comparing(ContainerInstance::getArn))
            .findFirst()
            .get();

    boolean hasTaskAlready =
        snapshot
            .getTasks()
            .stream()
            .anyMatch(
                task ->
                    task.getContainerInstanceArn().equals(instance.getArn())
                        && task.getGroup().equals(environment.getEnvironmentId()));

    if (!hasTaskAlready) {
      return Arrays.asList(
          StartTask.builder()
              .clusterArn(snapshot.getClusterArn())
              .taskDefinitionArn(environment.getTaskDefinitionArn())
              .group(environment.getEnvironmentId())
              .containerInstanceArn(instance.getArn())
              .build());
    } else {
      return Collections.emptyList();
    }
  }
}
