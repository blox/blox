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
import com.amazonaws.blox.dataservice.model.EnvironmentHealth;
import com.amazonaws.blox.dataservice.model.EnvironmentRevision;
import com.amazonaws.blox.dataservice.model.EnvironmentStatus;
import com.amazonaws.blox.dataservice.repository.EnvironmentRepository;
import com.amazonaws.blox.dataservicemodel.v1.exception.InternalServiceException;
import com.amazonaws.blox.dataservicemodel.v1.exception.InvalidParameterException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceExistsException;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentResponse;
import java.time.Instant;
import java.util.UUID;
import lombok.AllArgsConstructor;
import lombok.NonNull;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@AllArgsConstructor
public class CreateEnvironmentApi {

  @NonNull private final ApiModelMapper apiModelMapper;
  @NonNull private final EnvironmentRepository environmentRepository;

  public CreateEnvironmentResponse createEnvironment(
      @NonNull final CreateEnvironmentRequest request)
      throws ResourceExistsException, InvalidParameterException, InternalServiceException {

    //TODO: validate

    final Environment environment = apiModelMapper.toEnvironment(request);
    environment.setEnvironmentStatus(EnvironmentStatus.Inactive);
    environment.setEnvironmentHealth(EnvironmentHealth.Healthy);
    environment.setCreatedTime(Instant.now());
    environment.setLastUpdatedTime(Instant.now());
    environment.setValidEnvironment(false);

    final String environmentRevisionId = UUID.randomUUID().toString();
    environment.setLatestEnvironmentRevisionId(environmentRevisionId);

    final EnvironmentRevision environmentRevision =
        EnvironmentRevision.builder()
            .environmentId(environment.getEnvironmentId())
            .environmentRevisionId(environment.getLatestEnvironmentRevisionId())
            .taskDefinition(request.getTaskDefinition())
            .instanceGroup(apiModelMapper.toModelInstanceGroup(request.getInstanceGroup()))
            .createdTime(Instant.now())
            .build();

    try {
      final Environment createdEnvironment =
          environmentRepository.createEnvironmentAndEnvironmentRevision(
              environment, environmentRevision);

      return CreateEnvironmentResponse.builder()
          .environment(apiModelMapper.toWrapperEnvironment(createdEnvironment))
          .environmentRevision(apiModelMapper.toWrapperEnvironmentRevision(environmentRevision))
          .build();

    } catch (final ResourceExistsException | InternalServiceException e) {
      log.error(e.getMessage(), e);
      throw e;
    } catch (final Exception e) {
      log.error(e.getMessage(), e);
      throw new InternalServiceException(e.getMessage(), e);
    }
  }
}
