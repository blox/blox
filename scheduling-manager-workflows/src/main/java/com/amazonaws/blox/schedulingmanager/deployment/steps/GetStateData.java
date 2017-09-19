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

import com.amazonaws.blox.schedulingmanager.handler.Encoder;
import com.amazonaws.blox.schedulingmanager.deployment.steps.types.DeploymentData;
import com.amazonaws.blox.schedulingmanager.deployment.steps.types.StateData;
import com.amazonaws.blox.schedulingmanager.handler.StepHandler;
import com.amazonaws.services.lambda.runtime.Context;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@AllArgsConstructor
public class GetStateData implements StepHandler {

  private Encoder encoder;

  @Override
  public void handleRequest(InputStream input, OutputStream output, Context context)
      throws IOException {
    log.debug("getStateData lambda");

    final DeploymentData deploymentData = encoder.decode(input, DeploymentData.class);

    log.debug(
        "deployment data deployment id {} and clustername {}",
        deploymentData.getDeploymentId(),
        deploymentData.getClusterName());

    final StateData stateData =
        StateData.builder()
            .clusterName(deploymentData.getClusterName())
            .ecsRole(deploymentData.getEcsRole())
            .build();
    encoder.encode(output, stateData);
  }
}
