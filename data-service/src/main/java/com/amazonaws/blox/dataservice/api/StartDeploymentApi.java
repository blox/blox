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
import com.amazonaws.blox.dataservice.model.EnvironmentRevision;
import com.amazonaws.blox.dataservice.repository.EnvironmentRepository;
import com.amazonaws.blox.dataservicemodel.v1.exception.InternalServiceException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.StartDeploymentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.StartDeploymentResponse;
import java.util.UUID;
import lombok.AllArgsConstructor;
import lombok.NonNull;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@AllArgsConstructor
public class StartDeploymentApi {

  @NonNull private final ApiModelMapper apiModelMapper;
  @NonNull private final EnvironmentRepository environmentRepository;

  public StartDeploymentResponse startDeployment(@NonNull final StartDeploymentRequest request)
      throws ResourceNotFoundException, InternalServiceException {
    try {
      final EnvironmentId environmentId =
          apiModelMapper.toModelEnvironmentId(request.getEnvironmentId());
      final Environment environment = environmentRepository.getEnvironment(environmentId);

      final EnvironmentRevision environmentRevision =
          environmentRepository.getEnvironmentRevision(
              environmentId, request.getEnvironmentRevisionId());

      environment.setActiveEnvironmentRevisionId(environmentRevision.getEnvironmentRevisionId());

      environmentRepository.updateEnvironment(environment);

      //TODO: Create a new deployment
      final String deploymentId = UUID.randomUUID().toString();

      return StartDeploymentResponse.builder()
          .environmentId(request.getEnvironmentId())
          .environmentRevisionId(request.getEnvironmentRevisionId())
          .deploymentId(deploymentId)
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
