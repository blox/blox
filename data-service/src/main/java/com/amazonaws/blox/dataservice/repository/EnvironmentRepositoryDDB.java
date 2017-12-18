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
import com.amazonaws.blox.dataservice.exception.ResourceType;
import com.amazonaws.blox.dataservice.mapper.EnvironmentMapper;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentRevision;
import com.amazonaws.blox.dataservice.repository.model.EnvironmentDDBRecord;
import com.amazonaws.blox.dataservice.repository.model.EnvironmentRevisionDDBRecord;
import com.amazonaws.blox.dataservicemodel.v1.exception.InternalServiceException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceExistsException;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBMapper;
import com.amazonaws.services.dynamodbv2.model.ConditionalCheckFailedException;
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
      throws ResourceExistsException, InternalServiceException {
    final EnvironmentDDBRecord environmentDDBRecord =
        environmentMapper.toEnvironmentDDBRecord(environment);

    try {
      dynamoDBMapper.save(environmentDDBRecord);
    } catch (final ConditionalCheckFailedException e) {
      throw new ResourceExistsException(
          ResourceType.ENVIRONMENT, environment.getEnvironmentId().toString());

    } catch (final AmazonServiceException e) {
      throw new InternalServiceException(
          String.format(
              "Could not save record with environment id %s", environment.getEnvironmentId()),
          e);
    }

    return environmentMapper.toEnvironment(environmentDDBRecord);
  }

  @Override
  public EnvironmentRevision createEnvironmentRevision(
      @NonNull final EnvironmentRevision environmentRevision)
      throws ResourceExistsException, InternalServiceException {

    final EnvironmentRevisionDDBRecord environmentRevisionDDBRecord =
        environmentMapper.toEnvironmentRevisionDDBRecord(environmentRevision);

    try {
      dynamoDBMapper.save(environmentRevisionDDBRecord);
    } catch (final ConditionalCheckFailedException e) {
      throw new ResourceExistsException(
          ResourceType.ENVIRONMENT_REVISION,
          environmentRevisionDDBRecord.getEnvironmentRevisionId());

    } catch (final AmazonServiceException e) {
      throw new InternalServiceException(
          String.format(
              "Could not save record with environment revision id %s",
              environmentRevision.getEnvironmentRevisionId()),
          e);
    }

    return environmentMapper.toEnvironmentRevision(environmentRevisionDDBRecord);
  }
}
