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

import com.amazonaws.blox.scheduling.scheduler.engine.EnvironmentDescription;
import com.amazonaws.blox.scheduling.scheduler.engine.Scheduler;
import com.amazonaws.blox.scheduling.scheduler.engine.SchedulingAction;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot;
import java.util.List;

public abstract class DaemonScheduler implements Scheduler {

  @Override
  public List<SchedulingAction> schedule(
      ClusterSnapshot snapshot, EnvironmentDescription description) {
    DaemonEnvironment env = new DaemonEnvironment(description);
    ClusterSummary summary = new ClusterSummary(snapshot);

    return schedule(env, summary);
  }

  protected abstract List<SchedulingAction> schedule(DaemonEnvironment env, ClusterSummary summary);
}
