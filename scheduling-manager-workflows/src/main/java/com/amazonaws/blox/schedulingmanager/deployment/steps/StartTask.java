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
package com.amazonaws.blox.schedulingmanager.deployment.steps;

import com.amazonaws.blox.schedulingmanager.deployment.handler.Encoder;
import com.amazonaws.blox.schedulingmanager.deployment.steps.types.TaskData;
import com.amazonaws.blox.schedulingmanager.deployment.steps.types.TaskWorkflowInput;
import com.amazonaws.services.ecs.AmazonECS;
import com.amazonaws.services.ecs.model.StartTaskRequest;
import com.amazonaws.services.lambda.runtime.Context;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@AllArgsConstructor
public class StartTask implements StepHandler {

  private Encoder encoder;
  private AmazonECS ecs;

  @Override
  public void handleRequest(InputStream input, OutputStream output, Context context)
      throws IOException {

    log.debug("start task lambda");

    final TaskWorkflowInput taskWorkflowInput = encoder.decode(input, TaskWorkflowInput.class);

    final StartTaskRequest startTaskRequest =
        new StartTaskRequest().withTaskDefinition(taskWorkflowInput.getTaskDefinition());

    //TODO: actually start a task. for now just test that this step gets invoked
    log.debug("Starting ECS task");

    final TaskData taskData =
        TaskData.builder().taskDefinition(taskWorkflowInput.getTaskDefinition()).build();

    encoder.encode(output, taskData);
  }
}
