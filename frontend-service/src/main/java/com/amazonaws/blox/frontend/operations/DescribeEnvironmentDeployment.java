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

import com.amazonaws.blox.frontend.models.Deployment;
import io.swagger.annotations.ApiOperation;
import lombok.Builder;
import lombok.Value;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class DescribeEnvironmentDeployment extends EnvironmentController {
  @RequestMapping(
    path = "/{environmentName}/deployments/{deploymentId}",
    method = RequestMethod.GET
  )
  @ApiOperation(value = "Describe Environment deployment history")
  public DescribeEnvironmentDeploymentResponse describeEnvironmentDeployment(
      @PathVariable("cluster") String cluster,
      @PathVariable("environmentName") String environmentName,
      @PathVariable("deploymentId") String deploymentId) {

    return DescribeEnvironmentDeploymentResponse.builder()
        .deployment(
            Deployment.builder()
                .environmentName(environmentName)
                .deploymentId(deploymentId)
                .oldTargetRevisionId(null)
                .newTargetRevisionId(null)
                .build())
        .build();
  }

  @Value
  @Builder
  public static class DescribeEnvironmentDeploymentResponse {
    private final Deployment deployment;
  }
}
