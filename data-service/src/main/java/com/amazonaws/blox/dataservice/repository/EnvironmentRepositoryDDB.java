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
package com.amazonaws.blox.dataservice.repository;

import com.amazonaws.AmazonServiceException;
import com.amazonaws.blox.dataservice.exception.StorageException;
import com.amazonaws.blox.dataservice.mapper.EnvironmentMapper;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.repository.model.EnvironmentDDBRecord;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentNotFoundException;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBMapper;
import lombok.AllArgsConstructor;
import lombok.NonNull;
import org.springframework.stereotype.Component;

@Component
@AllArgsConstructor
public class EnvironmentRepositoryDDB implements EnvironmentRepository {

  @NonNull private DynamoDBMapper dynamoDBMapper;

  @NonNull private EnvironmentMapper environmentMapper;

  @Override
  public Environment createEnvironment(@NonNull final Environment environment)
      throws StorageException {
    final EnvironmentDDBRecord environmentDDBRecord =
        environmentMapper.toEnvironmentDDBRecord(environment);

    try {
      dynamoDBMapper.save(environmentDDBRecord);
    } catch (final AmazonServiceException e) {
      throw new StorageException(
          String.format(
              "Could not save record with environment id %s", environment.getEnvironmentId()),
          e);
    }
    return environmentMapper.toEnvironment(environmentDDBRecord);
  }

  @Override
  public Environment getEnvironment(
      @NonNull final String environmentId, @NonNull final String environmentVersion)
      throws StorageException, EnvironmentNotFoundException {

    final EnvironmentDDBRecord environmentDDBRecord;
    try {
      environmentDDBRecord =
          dynamoDBMapper.load(EnvironmentDDBRecord.withKeys(environmentId, environmentVersion));
    } catch (final AmazonServiceException e) {
      throw new StorageException(
          String.format(
              "Could not load record with environment id %s and version %s",
              environmentId, environmentVersion),
          e);
    }

    if (environmentDDBRecord == null) {
      throw new EnvironmentNotFoundException(
          String.format(
              "Could not find environment with id %s and version %s",
              environmentId, environmentVersion));
    }

    return environmentMapper.toEnvironment(environmentDDBRecord);
  }
}
