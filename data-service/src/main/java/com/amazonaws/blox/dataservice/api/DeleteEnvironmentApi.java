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
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DeleteEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DeleteEnvironmentResponse;
import java.util.List;
import lombok.AllArgsConstructor;
import lombok.NonNull;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@AllArgsConstructor
public class DeleteEnvironmentApi {
  @NonNull private final ApiModelMapper apiModelMapper;
  @NonNull private final EnvironmentRepository environmentRepository;

  // TODO: This deleteEnvironment is actually a reaper that removes the record from the dynamoDB tables. The actual deleteEnvironment should just update the status of the environment without removing the records
  public DeleteEnvironmentResponse deleteEnvironment(@NonNull DeleteEnvironmentRequest request)
      throws InternalServiceException, ResourceNotFoundException {
    final EnvironmentId environmentId =
        apiModelMapper.toModelEnvironmentId(request.getEnvironmentId());
    try {
      // Get the environment to check it exists
      final Environment environment = environmentRepository.getEnvironment(environmentId);

      // List and delete all environment revisions
      final List<EnvironmentRevision> environmentRevisions =
          environmentRepository.listEnvironmentRevisions(environmentId);
      // TODO: Need to stop all the running tasks from the scheduler manager

      for (final EnvironmentRevision revision : environmentRevisions) {
        environmentRepository.deleteEnvironmentRevision(revision);
      }

      // Delete the environment
      environmentRepository.deleteEnvironment(environmentId);
      return DeleteEnvironmentResponse.builder()
          .environment(apiModelMapper.toWrapperEnvironment(environment))
          .build();
    } catch (InternalServiceException | ResourceNotFoundException e) {
      log.error(e.getMessage(), e);
      throw e;
    } catch (Exception e) {
      log.error(e.getMessage(), e);
      throw new InternalServiceException(e.getMessage(), e);
    }
  }
}
