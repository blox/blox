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
package com.amazonaws.blox.schedulingmanager.wrapper;

import com.amazonaws.AmazonClientException;
import com.amazonaws.services.stepfunctions.AWSStepFunctions;
import com.amazonaws.services.stepfunctions.model.StartExecutionRequest;
import com.amazonaws.services.stepfunctions.model.StartExecutionResult;
import lombok.AllArgsConstructor;
import lombok.NonNull;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@AllArgsConstructor
public class StepFunctionsWrapper {

  @NonNull private AWSStepFunctions stepFunctions;

  public StartExecutionResult startExecution(
      final String stateMachineArn, final String workflowName, final String inputJson) {
    final StartExecutionRequest startExecutionRequest =
        new StartExecutionRequest()
            .withStateMachineArn(stateMachineArn)
            .withInput(inputJson)
            .withName(workflowName);
    try {
      return stepFunctions.startExecution(startExecutionRequest);
    } catch (final AmazonClientException e) {
      log.error(
          "StateMachine {} workflow {} with input {} failed to start",
          stateMachineArn,
          workflowName,
          inputJson,
          e);
      throw e;
    }
  }
}
