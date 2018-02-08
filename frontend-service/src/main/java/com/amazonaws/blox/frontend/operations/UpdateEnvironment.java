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

import com.amazonaws.blox.dataservicemodel.v1.exception.InvalidParameterException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ServiceException;
import com.amazonaws.blox.frontend.mappers.UpdateEnvironmentMapper;
import io.swagger.annotations.ApiOperation;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class UpdateEnvironment extends EnvironmentController {
  @Autowired UpdateEnvironmentMapper mapper;

  @RequestMapping(
    path = "/{environmentName}",
    method = RequestMethod.PUT,
    consumes = "application/json"
  )
  @ApiOperation(value = "Update an existing Environment")
  public UpdateEnvironmentResponse updateEnvironment(
      @PathVariable("cluster") String cluster,
      @PathVariable("environmentName") String environmentName,
      @RequestBody UpdateEnvironmentRequest request)
      throws InvalidParameterException, ServiceException, ResourceNotFoundException {

    return mapper.fromDataServiceResponse(
        dataService.updateEnvironment(
            mapper.toDataServiceRequest(
                getApiGatewayRequestContext(), cluster, environmentName, request)));
  }

  @Data
  @Builder
  @AllArgsConstructor
  @NoArgsConstructor
  public static class UpdateEnvironmentResponse {
    private String environmentRevisionId;
  }

  @Data
  @Builder
  @AllArgsConstructor
  @NoArgsConstructor
  public static class UpdateEnvironmentRequest {
    private String environmentName;
    private String taskDefinition;
  }
}
