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
package cucumber.steps.scheduling.daemon;

import static org.assertj.core.api.Assertions.assertThat;

import com.amazonaws.blox.scheduling.scheduler.engine.EnvironmentDescription;
import com.amazonaws.blox.scheduling.scheduler.engine.EnvironmentDescription.EnvironmentDescriptionBuilder;
import com.amazonaws.blox.scheduling.scheduler.engine.Scheduler;
import com.amazonaws.blox.scheduling.scheduler.engine.SchedulingAction;
import com.amazonaws.blox.scheduling.scheduler.engine.StartTask;
import com.amazonaws.blox.scheduling.scheduler.engine.daemon.ReplaceAfterTerminateScheduler;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot.ContainerInstance;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot.Task;
import cucumber.api.DataTable;
import cucumber.api.java8.En;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

public class ReplaceAfterTerminateSteps implements En {

  private EnvironmentDescription environment;
  private ClusterSnapshot snapshot;
  private Scheduler scheduler = new ReplaceAfterTerminateScheduler();

  private Set<SchedulingAction> actions = new HashSet<>();

  public ReplaceAfterTerminateSteps() {
    Given(
        "^a Daemon environment named \"([^\"]*)\":$",
        (String name, DataTable properties) -> {
          EnvironmentDescriptionBuilder builder = environmentDescriptionFromTable(properties);

          environment = builder.environmentName(name).build();
        });

    Given(
        "^a cluster named \"([^\"]*)\"",
        (String name) -> {
          snapshot = new ClusterSnapshot(name, new ArrayList<>(), new ArrayList<>());
        });

    Given(
        "^the cluster has the following instances and tasks:$",
        (DataTable table) -> {
          updateSnapshotFromTable(table);
        });

    When(
        "^the scheduler runs$",
        () -> {
          actions.clear();
          actions.addAll(scheduler.schedule(snapshot, environment));
        });

    Then(
        "^it should start the following tasks:$",
        (DataTable startTasksTable) -> {
          List<StartTask> startTasks = startTaskActionsFromTable(startTasksTable);

          assertThat(actions).containsAll(startTasks);

          actions.removeAll(startTasks);
        });

    And(
        "^it should not take any further actions$",
        () -> {
          assertThat(actions).isEmpty();
        });
  }

  private EnvironmentDescriptionBuilder environmentDescriptionFromTable(DataTable properties) {
    return properties.asList(EnvironmentDescriptionBuilder.class).get(0);
  }

  private List<StartTask> startTaskActionsFromTable(DataTable startTasksTable) {
    return startTasksTable
        .asList(StartTask.StartTaskBuilder.class)
        .stream()
        .map(b -> b.clusterName(environment.getClusterName()).build())
        .collect(Collectors.toList());
  }

  private void updateSnapshotFromTable(DataTable table) {
    snapshot.getInstances().clear();
    snapshot.getTasks().clear();
    for (Map<String, String> row : table.asMaps(String.class, String.class)) {
      String instanceArn = row.get("instance");
      snapshot.getInstances().add(ContainerInstance.builder().arn(instanceArn).build());

      for (String taskDescription : row.get("tasks").split(",")) {
        if (!taskDescription.isEmpty()) {
          snapshot.getTasks().add(taskFromDescription(instanceArn, taskDescription));
        }
      }
    }
  }

  private Task taskFromDescription(String instanceArn, String description) {
    String[] parts = description.split(":");
    String group = parts.length > 3 ? parts[3] : environment.getEnvironmentName();
    return Task.builder()
        .containerInstanceArn(instanceArn)
        .arn(parts[0])
        .taskDefinitionArn(parts[1])
        .status(parts[2])
        .group(group)
        .startedBy("blox")
        .build();
  }
}
