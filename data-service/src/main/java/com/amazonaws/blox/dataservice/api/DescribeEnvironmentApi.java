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
package com.amazonaws.blox.dataservice.api;

import com.amazonaws.blox.dataservice.mapper.ApiModelMapper;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentId;
import com.amazonaws.blox.dataservice.repository.EnvironmentRepository;
import com.amazonaws.blox.dataservicemodel.v1.exception.InternalServiceException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentResponse;

import lombok.NonNull;
import lombok.Value;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@Value
public class DescribeEnvironmentApi {
  @NonNull private ApiModelMapper apiModelMapper;
  @NonNull private EnvironmentRepository environmentRepository;

  public DescribeEnvironmentResponse describeEnvironment(
      @NonNull final DescribeEnvironmentRequest describeEnvironmentRequest)
      throws ResourceNotFoundException, InternalServiceException {
    final com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId environmentIdFromRequest =
        describeEnvironmentRequest.getEnvironmentId();
    final EnvironmentId environmentId =
        apiModelMapper.toModelEnvironmentId(environmentIdFromRequest);
    try {
      final Environment environment = environmentRepository.getEnvironment(environmentId);
      return DescribeEnvironmentResponse.builder()
          .environment(apiModelMapper.toWrapperEnvironment(environment))
          .build();
    } catch (final ResourceNotFoundException | InternalServiceException e) {
      log.error(e.getMessage(), e);
      throw e;
    } catch (final Exception e) {
      log.error(e.getMessage(), e);
      throw new InternalServiceException(e.getMessage(), e);
    }
  }
}
