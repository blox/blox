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
package com.amazonaws.blox.scheduling.state;

import java.util.List;
import java.util.Map;
import java.util.Set;
import lombok.Builder;
import lombok.Data;
import lombok.Value;

@Data
public class ClusterSnapshot {
  private final String clusterArn;
  private final List<Task> tasks;
  private final List<ContainerInstance> instances;

  @Value
  @Builder
  public static class Task {
    private final String arn;
    private final String containerInstanceArn;
    private final String taskDefinitionArn;
    private final String status;
    private final String group;
    private final String startedBy;

    public static Task from(software.amazon.awssdk.services.ecs.model.Task t) {
      return builder()
          .arn(t.taskArn())
          .containerInstanceArn(t.containerInstanceArn())
          .taskDefinitionArn(t.taskDefinitionArn())
          .status(t.desiredStatus())
          .group(t.group())
          .startedBy(t.startedBy())
          .build();
    }
  }

  @Value
  @Builder
  public static class ContainerInstance {
    private final String arn;

    public static ContainerInstance from(
        software.amazon.awssdk.services.ecs.model.ContainerInstance i) {
      return builder().arn(i.containerInstanceArn()).build();
    }
  }
}
