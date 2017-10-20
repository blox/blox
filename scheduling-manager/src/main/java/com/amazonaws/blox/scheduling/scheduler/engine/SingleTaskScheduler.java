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

import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentVersion;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot.ContainerInstance;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.Map.Entry;

/**
 * Temporary scheduler implementation that just starts a single task on the first instance in the
 * snapshot.
 */
public class SingleTaskScheduler implements Scheduler {

  @Override
  public List<SchedulingAction> schedule(ClusterSnapshot snapshot, EnvironmentVersion environment) {
    // Only run a Task on the first ContainerInstance in the snapshot:
    for (Entry<String, ContainerInstance> entry : snapshot.getInstances().entrySet()) {
      return Arrays.asList(
          new StartTask(snapshot.getClusterArn(), entry.getKey(), environment.getTaskDefinition()));
    }
    return Collections.emptyList();
  }
}
