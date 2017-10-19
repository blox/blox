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
import com.amazonaws.blox.dataservice.model.EnvironmentVersion;
import com.amazonaws.blox.dataservice.repository.model.EnvironmentDDBRecord;
import com.amazonaws.blox.dataservice.repository.model.EnvironmentTargetVersionDDBRecord;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentExistsException;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBMapper;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBQueryExpression;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBScanExpression;
import com.amazonaws.services.dynamodbv2.model.ConditionalCheckFailedException;
import java.util.List;
import java.util.stream.Collectors;
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
      throws EnvironmentExistsException, StorageException {
    final EnvironmentDDBRecord environmentDDBRecord =
        environmentMapper.toEnvironmentDDBRecord(environment);

    try {
      dynamoDBMapper.save(environmentDDBRecord);
    } catch (final ConditionalCheckFailedException e) {
      throw new EnvironmentExistsException(
          String.format(
              "Environment with id %s and version %s already exists",
              environment.getEnvironmentId(), environment.getEnvironmentVersion()));

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
      throws StorageException {

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
    return environmentMapper.toEnvironment(environmentDDBRecord);
  }

  @Override
  public EnvironmentVersion createEnvironmentTargetVersion(
      @NonNull final EnvironmentVersion environmentVersion)
      throws StorageException, EnvironmentExistsException {

    final EnvironmentTargetVersionDDBRecord environmentTargetVersionDDBRecord =
        environmentMapper.toEnvironmentTargetVersionDDBRecord(environmentVersion);

    try {
      dynamoDBMapper.save(environmentTargetVersionDDBRecord);
    } catch (final ConditionalCheckFailedException e) {
      throw new EnvironmentExistsException(
          String.format(
              "Environment with id %s and version %s already exists",
              environmentVersion.getEnvironmentId(), environmentVersion.getEnvironmentVersion()));

    } catch (final AmazonServiceException e) {
      throw new StorageException(
          String.format(
              "Could not save record with environment id %s",
              environmentVersion.getEnvironmentId()),
          e);
    }
    return environmentMapper.toEnvironmentVersion(environmentTargetVersionDDBRecord);
  }

  @Override
  public EnvironmentVersion getEnvironmentTargetVersion(@NonNull final String environmentId)
      throws StorageException {

    final EnvironmentTargetVersionDDBRecord environmentTargetVersionDDBRecord;
    try {
      environmentTargetVersionDDBRecord =
          dynamoDBMapper.load(EnvironmentTargetVersionDDBRecord.withHashKey(environmentId));
    } catch (final AmazonServiceException e) {
      throw new StorageException(
          String.format("Could not load record with environment id %s", environmentId), e);
    }
    return environmentMapper.toEnvironmentVersion(environmentTargetVersionDDBRecord);
  }

  @Override
  public List<String> listClusters() throws StorageException {
    try {
      final DynamoDBScanExpression scanExpression =
          new DynamoDBScanExpression()
              .withIndexName(EnvironmentTargetVersionDDBRecord.ENVIRONMENT_CLUSTER_GSI_NAME)
              .withConsistentRead(false);

      List<EnvironmentTargetVersionDDBRecord> scanResult =
          dynamoDBMapper.scan(EnvironmentTargetVersionDDBRecord.class, scanExpression);
      //TODO: integration test that covers pagination.
      return scanResult.stream().map(e -> e.getCluster()).distinct().collect(Collectors.toList());
    } catch (final AmazonServiceException e) {
      throw new StorageException("Could not scan environment target versions");
    }
  }

  @Override
  public List<String> listEnvironmentIdsByCluster(@NonNull final String cluster)
      throws StorageException {
    final DynamoDBQueryExpression<EnvironmentTargetVersionDDBRecord> queryExpression =
        new DynamoDBQueryExpression<EnvironmentTargetVersionDDBRecord>()
            .withIndexName(EnvironmentTargetVersionDDBRecord.ENVIRONMENT_CLUSTER_GSI_NAME)
            .withConsistentRead(false)
            .withHashKeyValues(EnvironmentTargetVersionDDBRecord.withGSIHashKey(cluster));

    try {
      //TODO: integration test that covers pagination.
      List<EnvironmentTargetVersionDDBRecord> queryResult =
          dynamoDBMapper.query(EnvironmentTargetVersionDDBRecord.class, queryExpression);

      return queryResult.stream().map(e -> e.getEnvironmentId()).collect(Collectors.toList());
    } catch (final AmazonServiceException e) {
      throw new StorageException(
          String.format("Could not query environment target versions for cluster %s", cluster), e);
    }
  }
}
