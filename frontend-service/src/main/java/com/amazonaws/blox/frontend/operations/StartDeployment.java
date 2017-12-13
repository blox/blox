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
package com.amazonaws.blox.frontend.operations;

import io.swagger.annotations.ApiOperation;
import lombok.Builder;
import lombok.Value;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class StartDeployment extends EnvironmentController {
  @RequestMapping(path = "/{environmentName}/deployments", method = RequestMethod.POST)
  @ApiOperation(value = "Deploy Environment revision")
  public StartDeploymentResponse startDeployment(
      @PathVariable("cluster") String cluster,
      @PathVariable("environmentName") String environmentName,
      @RequestParam("revisionId") String revisionId) {

    return StartDeploymentResponse.builder()
        .deploymentId(environmentName + "_" + revisionId + "_deploymentId")
        .build();
  }

  @Value
  @Builder
  public static class StartDeploymentResponse {
    private final String deploymentId;
  }
}
