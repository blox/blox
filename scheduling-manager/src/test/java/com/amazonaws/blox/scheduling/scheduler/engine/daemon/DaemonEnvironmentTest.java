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
import com.amazonaws.blox.scheduling.scheduler.engine.StopTask;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot.ContainerInstance;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot.Task;
import java.util.Collections;
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
        environment.isMatchingTask(
            defaultTask().taskDefinitionArn(env.getTaskDefinitionArn()).status("RUNNING").build());

    assertThat(matches).isTrue();
  }

  @Test
  public void matchesTaskWithDifferentTaskDefinition() {
    boolean matches =
        environment.isMatchingTask(defaultTask().taskDefinitionArn("different-taskdef").build());

    assertThat(matches).isTrue();
  }

  @Test
  public void doesntMatchTaskWithDifferentGroup() {
    boolean matches = environment.isMatchingTask(defaultTask().group("different-group").build());

    assertThat(matches).isFalse();
  }

  @Test
  public void doesntMatchUnhealthyTask() {
    boolean matches = environment.isMatchingTask(defaultTask().status("STOPPED").build());

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

  @Test
  public void hasTaskToStop() {
    boolean stop =
        environment.isTaskStoppable(defaultTask().taskDefinitionArn("different-taskdef").build());

    assertThat(stop).isTrue();
  }

  @Test
  public void noTasksToStopForDesiredTasks() {
    boolean stop = environment.isTaskStoppable(defaultTask().build());

    assertThat(stop).isFalse();
  }

  @Test
  public void noTasksToStopForAlreadyStoppedTasks() {
    boolean stop =
        environment.isTaskStoppable(
            defaultTask().taskDefinitionArn("different-taskdef").status("STOPPED").build());

    assertThat(stop).isFalse();
  }

  @Test
  public void noTasksToStopForAnotherEnvironment() {
    boolean stop =
        environment.isTaskStoppable(
            defaultTask().taskDefinitionArn("different-taskdef").group("different-group").build());

    assertThat(stop).isFalse();
  }

  @Test
  public void noHealthyTasks() {
    boolean isMissingHealthyTask =
        environment.isMissingHealthyTask(
            Collections.singletonList(defaultTask().status("STOPPED").build()));

    assertThat(isMissingHealthyTask).isTrue();
  }

  @Test
  public void hasHealthyTasks() {
    boolean isMissingHealthyTask =
        environment.isMissingHealthyTask(Collections.singletonList(defaultTask().build()));

    assertThat(isMissingHealthyTask).isFalse();
  }

  @Test
  public void hasHealthyTasksWithAnotherVersion() {
    boolean isMissingHealthyTask =
        environment.isMissingHealthyTask(
            Collections.singletonList(
                defaultTask().taskDefinitionArn("different-taskdef").build()));

    assertThat(isMissingHealthyTask).isFalse();
  }

  @Test
  public void hasHealthyTasksOnAnotherEnvironment() {
    boolean isMissingHealthyTask =
        environment.isMissingHealthyTask(
            Collections.singletonList(defaultTask().group("different-environment").build()));

    assertThat(isMissingHealthyTask).isTrue();
  }

  @Test
  public void stopTasksOnEnvironmentWithDifferentVersion() {
    Task taskWithDifferentVersion = defaultTask().taskDefinitionArn("different-taskdef").build();

    StopTask stopTask = environment.stopTaskFor(taskWithDifferentVersion);

    assertSoftly(
        s -> {
          s.assertThat(stopTask.getClusterName()).isEqualTo(env.getClusterName());
          s.assertThat(stopTask.getTask())
              .isEqualTo(taskWithDifferentVersion.getTaskDefinitionArn());
          s.assertThat(stopTask.getReason())
              .isEqualTo(
                  String.format(
                      "Stopped by deployment to %s@%s",
                      env.getEnvironmentName(), env.getTaskDefinitionArn()));
        });
  }
}
