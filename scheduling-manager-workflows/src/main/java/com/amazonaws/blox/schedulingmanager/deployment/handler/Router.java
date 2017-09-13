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
package com.amazonaws.blox.schedulingmanager.deployment.handler;

import com.amazonaws.blox.schedulingmanager.deployment.exception.HandlerNotFoundException;
import com.amazonaws.blox.schedulingmanager.deployment.steps.StepHandler;
import java.util.Map;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@AllArgsConstructor
public class Router {

  private Map<String, StepHandler> handlers;

  /** Find stepHandler matching the stepName. */
  public StepHandler findStepHandler(final String stepName) {
    if (!handlers.containsKey(stepName)) {
      log.error("StepHandler {} not found", stepName);
      throw new HandlerNotFoundException(String.format("StepHandler %s not found", stepName));
    } else {
      return handlers.get(stepName);
    }
  }
}
