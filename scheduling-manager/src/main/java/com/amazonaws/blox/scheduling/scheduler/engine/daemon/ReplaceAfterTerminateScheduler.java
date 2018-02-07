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

import com.amazonaws.blox.scheduling.scheduler.engine.SchedulingAction;
import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.Stream;

public class ReplaceAfterTerminateScheduler extends DaemonScheduler {
  public static final String ID = "ReplaceAfterTerminate";

  public List<SchedulingAction> schedule(DaemonEnvironment env, ClusterSummary summary) {
    Stream<SchedulingAction> stopTaskActions =
        summary
            .getInstances()
            .stream()
            .flatMap(i -> summary.tasksForInstance(i).stream().filter(env::isTaskStoppable))
            .map(env::stopTaskFor);

    Stream<SchedulingAction> startTaskActions =
        summary
            .getInstances()
            .stream()
            .filter(i -> env.isMissingHealthyTask(summary.tasksForInstance(i)))
            .map(env::startTaskFor);

    return Stream.concat(startTaskActions, stopTaskActions).collect(Collectors.toList());
  }
}
