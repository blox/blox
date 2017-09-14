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
package com.amazonaws.blox.schedulingmanager.deployment.steps.types;

import com.fasterxml.jackson.databind.annotation.JsonDeserialize;
import com.fasterxml.jackson.databind.annotation.JsonPOJOBuilder;
import lombok.Builder;
import lombok.Value;

@JsonDeserialize(builder = TaskWorkflowInput.TaskWorkflowInputBuilder.class)
@Builder
@Value
public class TaskWorkflowInput {

  private final String cluster;
  private final String containerInstance;
  private final String taskDefinition;
  private final String group;
  private final String startedBy;
  private final String taskRole;

  @JsonPOJOBuilder(withPrefix = "")
  public static final class TaskWorkflowInputBuilder {}
}
