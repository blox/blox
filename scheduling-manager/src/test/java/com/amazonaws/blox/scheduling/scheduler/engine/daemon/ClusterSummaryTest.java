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

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import com.amazonaws.blox.scheduling.state.ClusterSnapshot;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot.ContainerInstance;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot.Task;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import org.junit.Test;

public class ClusterSummaryTest {
  private static final String CLUSTER_ARN = "cluster";

  private ClusterSnapshot snapshot =
      new ClusterSnapshot(CLUSTER_ARN, new ArrayList<>(), new ArrayList<>());

  @Test
  public void tasksForInstanceReturnsNoTasksForEmptyInstance() {
    given(instance("i-1"));

    ClusterSummary summary = new ClusterSummary(snapshot);

    assertThat(summary.tasksForInstance(instance("i-1"))).isEmpty();
  }

  @Test
  public void tasksForInstanceReturnsCorrectTasksForGivenInstance() {
    given(instance("i-1"), instance("i-2"), instance("i-3"));
    given(task("i-1", "t-1"), task("i-1", "t-2"), task("i-2", "t-3"));

    ClusterSummary summary = new ClusterSummary(snapshot);

    assertThat(summary.tasksForInstance(instance("i-1")))
        .containsExactly(task("i-1", "t-1"), task("i-1", "t-2"));
  }

  @Test
  public void tasksForInstanceThrowsForNonexistentInstance() {
    ClusterSummary summary = new ClusterSummary(snapshot);

    assertThatThrownBy(() -> summary.tasksForInstance(instance("i-1"))).hasMessageContaining("i-1");
  }

  private void given(ContainerInstance... instances) {
    snapshot.getInstances().addAll(Arrays.asList(instances));
  }

  private void given(Task... tasks) {
    snapshot.getTasks().addAll(Arrays.asList(tasks));
  }

  private ContainerInstance instance(String arn) {
    return ContainerInstance.builder().arn(arn).build();
  }

  private Task task(String instanceArn, String taskArn) {
    return Task.builder().containerInstanceArn(instanceArn).arn(taskArn).build();
  }

  private ClusterSnapshot snapshot(List<ContainerInstance> instances, List<Task> tasks) {
    return new ClusterSnapshot(CLUSTER_ARN, tasks, instances);
  }
}
