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

import static com.amazonaws.blox.dataservice.repository.model.EnvironmentDDBRecord.ENVIRONMENT_NAME_RANGE_KEY;
import static com.amazonaws.blox.dataservice.repository.model.EnvironmentDDBRecord.LATEST_ENVIRONMENT_REVISION_ID;

import com.amazonaws.AmazonServiceException;
import com.amazonaws.blox.dataservice.exception.ResourceType;
import com.amazonaws.blox.dataservice.mapper.EnvironmentMapper;
import com.amazonaws.blox.dataservice.model.Cluster;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentId;
import com.amazonaws.blox.dataservice.model.EnvironmentRevision;
import com.amazonaws.blox.dataservice.repository.model.EnvironmentDDBRecord;
import com.amazonaws.blox.dataservice.repository.model.EnvironmentRevisionDDBRecord;
import com.amazonaws.blox.dataservicemodel.v1.exception.InternalServiceException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceExistsException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceNotFoundException;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBMapper;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBMapperConfig;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBMapperConfig.SaveBehavior;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBQueryExpression;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBSaveExpression;
import com.amazonaws.services.dynamodbv2.model.AttributeValue;
import com.amazonaws.services.dynamodbv2.model.ComparisonOperator;
import com.amazonaws.services.dynamodbv2.model.Condition;
import com.amazonaws.services.dynamodbv2.model.ConditionalCheckFailedException;
import com.amazonaws.services.dynamodbv2.model.ExpectedAttributeValue;
import java.util.List;
import java.util.stream.Collectors;
import lombok.AllArgsConstructor;
import lombok.NonNull;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@AllArgsConstructor
public class EnvironmentRepositoryDDB implements EnvironmentRepository {

  @NonNull private DynamoDBMapper dynamoDBMapper;
  @NonNull private EnvironmentMapper environmentMapper;

  @Override
  public Environment createEnvironmentAndEnvironmentRevision(
      @NonNull final Environment environment,
      @NonNull final EnvironmentRevision environmentRevision)
      throws ResourceExistsException, InternalServiceException {

    createEnvironment(environment);
    createEnvironmentRevision(environmentRevision);

    return updateEnvironmentToValid(environment);
  }

  @Override
  public Environment updateEnvironment(Environment environment)
      throws ResourceNotFoundException, InternalServiceException {
    return updateEnvironment(environment, null);
  }

  @Override
  public EnvironmentRevision getEnvironmentRevision(
      @NonNull final EnvironmentId environmentId, @NonNull final String environmentRevisionId)
      throws ResourceNotFoundException, InternalServiceException {
    final String accountIdClusterEnvironmentName =
        environmentId.generateAccountIdClusterEnvironmentName();
    try {
      final EnvironmentRevisionDDBRecord environmentRevisionRecord =
          dynamoDBMapper.load(
              EnvironmentRevisionDDBRecord.class,
              accountIdClusterEnvironmentName,
              environmentRevisionId);
      if (environmentRevisionRecord == null) {
        throw new ResourceNotFoundException(
            ResourceType.ENVIRONMENT_REVISION, environmentRevisionId);
      }
      return environmentMapper.toEnvironmentRevision(environmentRevisionRecord);
    } catch (final AmazonServiceException e) {
      throw new InternalServiceException(
          String.format(
              "Could not load record with environment revision key %s and environment revision id %s",
              accountIdClusterEnvironmentName, environmentRevisionId),
          e);
    }
  }

  private Environment updateEnvironment(
      @NonNull final Environment environment, final DynamoDBSaveExpression dynamoDBSaveExpression)
      throws ResourceNotFoundException, InternalServiceException {

    try {
      EnvironmentDDBRecord environmentDDBRecord =
          dynamoDBMapper.load(
              EnvironmentDDBRecord.withKeys(
                  environment.getEnvironmentId().generateAccountIdCluster(),
                  environment.getEnvironmentId().getEnvironmentName()));

      if (environmentDDBRecord == null) {
        throw new ResourceNotFoundException(
            ResourceType.ENVIRONMENT, environment.getEnvironmentId().toString());
      }

      final EnvironmentDDBRecord environmentRecordUpdates =
          environmentMapper.toEnvironmentDDBRecord(environment);

      // optimistic locking will prevent an update if the record has been updated since this was loaded
      environmentRecordUpdates.setRecordVersion(environmentDDBRecord.getRecordVersion());

      dynamoDBMapper.save(environmentRecordUpdates, dynamoDBSaveExpression);
      return environmentMapper.toEnvironment(environmentRecordUpdates);

    } catch (ConditionalCheckFailedException e) {
      throw new InternalServiceException(
          String.format(
              "Environment %s has been modified and cannot be updated",
              environment.getEnvironmentId()),
          e);
    } catch (final AmazonServiceException e) {
      throw new InternalServiceException(
          String.format(
              "Could not save record with environment id %s", environment.getEnvironmentId()),
          e);
    }
  }

  private Environment updateEnvironmentToValid(@NonNull final Environment environment)
      throws InternalServiceException {

    environment.setValidEnvironment(true);

    try {
      // update succeeds if environment revision was created in this call
      return updateEnvironment(
          environment,
          new DynamoDBSaveExpression()
              .withExpectedEntry(
                  LATEST_ENVIRONMENT_REVISION_ID,
                  new ExpectedAttributeValue(
                      new AttributeValue(environment.getLatestEnvironmentRevisionId()))));
    } catch (final ResourceNotFoundException e) {
      throw new InternalServiceException(
          String.format(
              "Trying to update the environment %s to valid but the environment does not exist",
              environment.getEnvironmentId()),
          e);
    }
  }

  private void createEnvironment(@NonNull final Environment environment)
      throws ResourceExistsException, InternalServiceException {

    try {
      createEnvironmentDDB(environment);

      // if environment exists
    } catch (final ResourceExistsException e) {
      log.debug(
          String.format(
              "Environment %s exists. Checking if it's a valid environment.",
              environment.getEnvironmentId().toString()));

      // if the environment is valid, throw resource exists exception
      if (isEnvironmentValid(environment)) {
        throw e;
      }

      log.debug(
          String.format("Environment %s is not valid", environment.getEnvironmentId().toString()));

      // if it's not valid, that means creation failed sometimes before it was successfully set to valid.
      // clean up potential existing revisions and recreate the environment
      cleanupEnvironmentRevisions(environment.getEnvironmentId());
      // clears and replaces all attributes. Versioned field constraints are also disregarded.
      createEnvironmentDDB(environment, SaveBehavior.CLOBBER.config());
    }
  }

  /**
   * Create respecting optimistic locking. If an environment with specified keys exists or the item
   * has been updated after being retrieved, a conditional check exception will be thrown.
   */
  private Environment createEnvironmentDDB(@NonNull final Environment environment)
      throws ResourceExistsException, InternalServiceException {

    return createEnvironmentDDB(environment, SaveBehavior.UPDATE.config());
  }

  private Environment createEnvironmentDDB(
      @NonNull final Environment environment, DynamoDBMapperConfig dynamoDBMapperConfig)
      throws ResourceExistsException, InternalServiceException {

    final EnvironmentDDBRecord environmentDDBRecord =
        environmentMapper.toEnvironmentDDBRecord(environment);

    try {
      dynamoDBMapper.save(environmentDDBRecord, dynamoDBMapperConfig);
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

  private boolean isEnvironmentValid(@NonNull final Environment environment)
      throws InternalServiceException {
    try {
      final Environment retrievedEnvironment = getEnvironment(environment.getEnvironmentId());
      if (retrievedEnvironment.isValidEnvironment()) {
        log.info(
            String.format(
                "Environment %s exists and is valid. Cannot recreate.",
                environment.getEnvironmentId().toString()));
        return true;
      }
    } catch (final ResourceNotFoundException e) {
      log.info(
          String.format(
              "Environment %s does not exist. Skipping valid environment check.",
              environment.getEnvironmentId().toString()),
          e);
    }

    return false;
  }

  /** Check if environment revisions exist and clean them up. */
  private void cleanupEnvironmentRevisions(@NonNull final EnvironmentId environmentId)
      throws InternalServiceException {

    List<EnvironmentRevision> environmentRevisions = listEnvironmentRevisions(environmentId);
    for (EnvironmentRevision environmentRevision : environmentRevisions) {
      try {
        deleteEnvironmentRevision(environmentRevision);
      } catch (final InternalServiceException e) {
        log.info(
            "Skipping deleting environment revision with environment id %s and environment revision id %s because delete was unsuccessful.",
            environmentRevision.getEnvironmentId(),
            environmentRevision.getEnvironmentRevisionId(),
            e);
      }
    }
  }

  private EnvironmentRevision createEnvironmentRevision(
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

  @Override
  public Environment getEnvironment(@NonNull final EnvironmentId environmentId)
      throws ResourceNotFoundException, InternalServiceException {

    final String accountIdCluster = environmentId.generateAccountIdCluster();
    try {
      final EnvironmentDDBRecord environmentDDBRecord =
          dynamoDBMapper.load(
              EnvironmentDDBRecord.withKeys(accountIdCluster, environmentId.getEnvironmentName()));
      if (environmentDDBRecord == null) {
        throw new ResourceNotFoundException(ResourceType.ENVIRONMENT, environmentId.toString());
      }
      return environmentMapper.toEnvironment(environmentDDBRecord);
    } catch (final AmazonServiceException e) {
      throw new InternalServiceException(
          String.format(
              "Could not load record with environment id %s and environment name %s",
              environmentId.generateAccountIdCluster(), environmentId.getEnvironmentName()),
          e);
    }
  }

  @Override
  public List<Environment> listEnvironments(@NonNull final Cluster cluster)
      throws InternalServiceException {
    return listEnvironments(cluster, null);
  }

  @Override
  public List<Environment> listEnvironments(
      @NonNull final Cluster cluster, final String environmentNamePrefix)
      throws InternalServiceException {
    try {
      final DynamoDBQueryExpression<EnvironmentDDBRecord> queryExpression =
          new DynamoDBQueryExpression<EnvironmentDDBRecord>()
              .withHashKeyValues(
                  EnvironmentDDBRecord.withHashKeys(cluster.generateAccountIdCluster()));
      if (environmentNamePrefix != null) {
        queryExpression.withRangeKeyCondition(
            ENVIRONMENT_NAME_RANGE_KEY,
            new Condition()
                .withComparisonOperator(ComparisonOperator.BEGINS_WITH)
                .withAttributeValueList(new AttributeValue().withS(environmentNamePrefix)));
      }
      return dynamoDBMapper
          .query(EnvironmentDDBRecord.class, queryExpression)
          .stream()
          .map(environmentMapper::toEnvironment)
          .collect(Collectors.toList());
    } catch (final AmazonServiceException e) {
      throw new InternalServiceException(
          String.format("Could not query environments for cluster %s", cluster.toString()), e);
    }
  }

  @Override
  public List<EnvironmentRevision> listEnvironmentRevisions(
      @NonNull final EnvironmentId environmentId) throws InternalServiceException {
    try {
      return dynamoDBMapper
          .query(
              EnvironmentRevisionDDBRecord.class,
              new DynamoDBQueryExpression<EnvironmentRevisionDDBRecord>()
                  .withHashKeyValues(
                      EnvironmentRevisionDDBRecord.withHashKey(
                          environmentId.generateAccountIdClusterEnvironmentName())))
          .stream()
          .map(environmentMapper::toEnvironmentRevision)
          .collect(Collectors.toList());
    } catch (final AmazonServiceException e) {
      throw new InternalServiceException(
          String.format(
              "Could not query environment revisions for environment %s", environmentId.toString()),
          e);
    }
  }

  @Override
  public void deleteEnvironmentRevision(@NonNull final EnvironmentRevision environmentRevision)
      throws InternalServiceException {
    try {
      dynamoDBMapper.delete(
          EnvironmentRevisionDDBRecord.withKeys(
              environmentRevision.getEnvironmentId().generateAccountIdClusterEnvironmentName(),
              environmentRevision.getEnvironmentRevisionId()),
          SaveBehavior.CLOBBER.config());
    } catch (final AmazonServiceException e) {
      throw new InternalServiceException(
          String.format(
              "Could not delete environment revision with accountIdClusterEnvironmentName %s and environment revision %s",
              environmentRevision.getEnvironmentId().toString(),
              environmentRevision.getEnvironmentRevisionId()),
          e);
    }
  }
}
