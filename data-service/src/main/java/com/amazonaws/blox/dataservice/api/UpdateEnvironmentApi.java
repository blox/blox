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
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceExistsException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.UpdateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.UpdateEnvironmentResponse;
import java.time.Instant;
import java.util.UUID;
import lombok.AllArgsConstructor;
import lombok.NonNull;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@AllArgsConstructor
public class UpdateEnvironmentApi {

  @NonNull private final ApiModelMapper mapper;
  @NonNull private final EnvironmentRepository repository;

  public UpdateEnvironmentResponse updateEnvironment(
      @NonNull final UpdateEnvironmentRequest request)
      throws ResourceNotFoundException, InternalServiceException, ResourceExistsException {

    final EnvironmentId id = mapper.toModelEnvironmentId(request.getEnvironmentId());
    final String environmentRevisionId = UUID.randomUUID().toString();

    // TODO this is basically ignoring instancegroup, because we're probably going to move that to
    // the Environment model, and not require a deployment to change that.
    repository.createEnvironmentRevision(
        EnvironmentRevision.builder()
            .environmentId(id)
            .environmentRevisionId(environmentRevisionId)
            .taskDefinition(request.getTaskDefinition())
            .createdTime(Instant.now())
            .build());

    // TODO this doesn't do any sort of concurrency control to ensure that things are left in a
    // consistent state under partial failure.
    final Environment environment = repository.getEnvironment(id);
    environment.setLatestEnvironmentRevisionId(environmentRevisionId);
    repository.updateEnvironment(environment);

    return UpdateEnvironmentResponse.builder().environmentRevisionId(environmentRevisionId).build();
  }
}
