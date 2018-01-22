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
import static org.assertj.core.api.SoftAssertions.assertSoftly;

import com.amazonaws.blox.scheduling.scheduler.engine.EnvironmentDescription;
import com.amazonaws.blox.scheduling.scheduler.engine.EnvironmentDescription.EnvironmentType;
import com.amazonaws.blox.scheduling.scheduler.engine.StartTask;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot.ContainerInstance;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot.Task;
import org.junit.Test;

public class DaemonEnvironmentTest {

  EnvironmentDescription env =
      EnvironmentDescription.builder()
          .clusterName("TestCluster")
          .environmentName("TestEnvironment")
          .environmentType(EnvironmentType.Daemon)
          .deploymentMethod("TestDeploymentMethod")
          .taskDefinitionArn("test-taskdef")
          .build();

  DaemonEnvironment environment = new DaemonEnvironment(env);

  private Task.TaskBuilder defaultTask() {
    return Task.builder()
        .arn("test-task")
        .containerInstanceArn("test-instance")
        .group(env.getEnvironmentName())
        .startedBy("blox")
        .taskDefinitionArn(env.getTaskDefinitionArn())
        .status("RUNNING");
  }

  @Test
  public void matchesHealthyTasksWithSameTaskDefinition() {
    boolean matches =
        environment.matchesTask(
            defaultTask().taskDefinitionArn(env.getTaskDefinitionArn()).status("RUNNING").build());

    assertThat(matches).isTrue();
  }

  @Test
  public void doesntMatchTaskWithDifferentTaskDefinition() {
    boolean matches =
        environment.matchesTask(defaultTask().taskDefinitionArn("different-taskdef").build());

    assertThat(matches).isFalse();
  }

  @Test
  public void doesntMatchTaskWithDifferentGroup() {
    boolean matches = environment.matchesTask(defaultTask().group("different-group").build());

    assertThat(matches).isFalse();
  }

  @Test
  public void doesntMatchUnhealthyTask() {
    boolean matches = environment.matchesTask(defaultTask().status("STOPPED").build());

    assertThat(matches).isFalse();
  }

  @Test
  public void startsTasksThatMatchesEnvironmentAndInstance() {
    ContainerInstance instance = ContainerInstance.builder().arn("instance-1").build();

    StartTask task = environment.startTaskFor(instance);

    assertSoftly(
        s -> {
          s.assertThat(task.getClusterName()).isEqualTo(env.getClusterName());
          s.assertThat(task.getContainerInstanceArn()).isEqualTo(instance.getArn());
          s.assertThat(task.getGroup()).isEqualTo(env.getEnvironmentName());
          s.assertThat(task.getTaskDefinitionArn()).isEqualTo(env.getTaskDefinitionArn());
        });
  }
}
