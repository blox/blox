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
package com.amazonaws.blox.scheduling.deployment.steps;

import com.amazonaws.blox.scheduling.deployment.steps.types.StateData;
import com.amazonaws.blox.scheduling.handler.Encoder;
import com.amazonaws.blox.scheduling.handler.StepHandler;
import com.amazonaws.blox.scheduling.task.steps.types.TaskWorkflowInput;
import com.amazonaws.blox.scheduling.wrapper.StepFunctionsWrapper;
import com.amazonaws.services.lambda.runtime.Context;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.util.UUID;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@AllArgsConstructor
public class StartDeployment implements StepHandler {

  private static final String START_TASK_WF_ARN_ENV_VAR = "START_TASK_WF_ARN";
  private static final String START_TASK_WF_PREFIX = "StartTaskWorkflow";

  private Encoder encoder;
  private StepFunctionsWrapper stepFunctionsWrapper;

  @Override
  public void handleRequest(InputStream input, OutputStream output, Context context)
      throws IOException {

    final StateData stateData = encoder.decode(input, StateData.class);

    for (final String containerInstance : stateData.getInstances()) {

      final TaskWorkflowInput taskWorkflowInput =
          TaskWorkflowInput.builder()
              .taskDefinition(stateData.getTask())
              .cluster(stateData.getCluster())
              .containerInstance(containerInstance)
              .ecsRole(stateData.getEcsRole())
              .build();

      final String taskWorkflowInputJson = encoder.encode(taskWorkflowInput);
      final String stateMachineArn = System.getenv(START_TASK_WF_ARN_ENV_VAR);
      final String workflowName = START_TASK_WF_PREFIX + UUID.randomUUID().toString();

      stepFunctionsWrapper.startExecution(stateMachineArn, workflowName, taskWorkflowInputJson);
    }
  }
}
