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
package com.amazonaws.blox.dataservice.handler;

import com.amazonaws.blox.dataservice.exception.StorageException;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentHealth;
import com.amazonaws.blox.dataservice.model.EnvironmentStatus;
import com.amazonaws.blox.dataservice.model.EnvironmentVersion;
import com.amazonaws.blox.dataservice.repository.EnvironmentRepository;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentExistsException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentVersionNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ServiceException;
import java.time.Instant;
import java.util.List;
import java.util.UUID;
import lombok.AllArgsConstructor;
import lombok.NonNull;
import org.springframework.stereotype.Component;

@Component
@AllArgsConstructor
public class EnvironmentHandler {

  @NonNull private EnvironmentRepository environmentRepository;

  //TODO: createEnvRequest only contains some fields that are set on the env object
  public Environment createEnvironment(@NonNull final Environment environment)
      throws EnvironmentExistsException, ServiceException {

    try {
      environment.setEnvironmentVersion(UUID.randomUUID().toString());
      environment.setCreatedTime(Instant.now());
      environment.setHealth(EnvironmentHealth.Healthy);
      environment.setStatus(EnvironmentStatus.Inactive);

      return environmentRepository.createEnvironment(environment);
    } catch (final StorageException e) {
      throw new ServiceException(
          String.format(
              "Exception occurred when creating environment with id %s and version %s",
              environment.getEnvironmentId(), environment.getEnvironmentVersion()),
          e);
    }
  }

  public EnvironmentVersion createEnvironmentTargetVersion(
      @NonNull final String environmentId, @NonNull final String environmentVersion)
      throws EnvironmentNotFoundException, EnvironmentExistsException, ServiceException {

    try {
      final Environment environment =
          environmentRepository.getEnvironment(environmentId, environmentVersion);

      if (environment == null) {
        throw new EnvironmentNotFoundException(
            String.format(
                "Environment with id %s and version %s does not exist",
                environmentId, environmentVersion));
      }

      return environmentRepository.createEnvironmentTargetVersion(
          EnvironmentVersion.builder()
              .environmentId(environment.getEnvironmentId())
              .environmentVersion(environment.getEnvironmentVersion())
              .cluster(environment.getInstanceGroup().getCluster())
              .build());
    } catch (final StorageException e) {
      throw new ServiceException(
          String.format(
              "Exception occurred when creating environment target version with id %s and version %s",
              environmentId, environmentVersion),
          e);
    }
  }

  public EnvironmentVersion describeEnvironmentTargetVersion(@NonNull final String environmentId)
      throws EnvironmentVersionNotFoundException, ServiceException {
    try {
      final EnvironmentVersion environmentVersion =
          environmentRepository.getEnvironmentTargetVersion(environmentId);

      if (environmentVersion == null) {
        throw new EnvironmentVersionNotFoundException(
            String.format("Could not find environment with id %s", environmentId));
      }

      return environmentVersion;
    } catch (final StorageException e) {
      throw new ServiceException(
          String.format("Exception occurred when getting environment with id %s", environmentId),
          e);
    }
  }

  public Environment describeEnvironment(
      @NonNull final String environmentId, @NonNull final String environmentVersion)
      throws EnvironmentNotFoundException, ServiceException {
    try {
      final Environment environment =
          environmentRepository.getEnvironment(environmentId, environmentVersion);

      if (environment == null) {
        throw new EnvironmentNotFoundException(
            String.format("Could not find environment with id %s", environmentId));
      }

      return environment;
    } catch (final StorageException e) {
      throw new ServiceException(
          String.format(
              "Exception occurred when getting environment with id %s and version %s",
              environmentId, environmentVersion),
          e);
    }
  }

  public List<String> listEnvironmentsWithCluster(@NonNull final String cluster)
      throws ServiceException {
    try {
      return environmentRepository.listEnvironmentIdsByCluster(cluster);
    } catch (final StorageException e) {
      throw new ServiceException(
          String.format("Exception occurred when listing clusters with cluster %s", cluster), e);
    }
  }

  public List<String> listClusters() throws ServiceException {
    try {
      return environmentRepository.listClusters();
    } catch (final StorageException e) {
      throw new ServiceException(String.format("Exception occurred when listing all clusters"), e);
    }
  }
}
