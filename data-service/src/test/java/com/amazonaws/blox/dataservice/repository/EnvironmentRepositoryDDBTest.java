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

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.containsInAnyOrder;
import static org.junit.Assert.assertEquals;
import static org.mockito.Mockito.doNothing;
import static org.mockito.Mockito.doThrow;
import static org.mockito.Mockito.eq;
import static org.mockito.Mockito.isA;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import com.amazonaws.AmazonServiceException;
import com.amazonaws.blox.dataservice.arn.EnvironmentArnGenerator;
import com.amazonaws.blox.dataservice.exception.StorageException;
import com.amazonaws.blox.dataservice.mapper.EnvironmentMapper;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentHealth;
import com.amazonaws.blox.dataservice.model.EnvironmentStatus;
import com.amazonaws.blox.dataservice.model.EnvironmentType;
import com.amazonaws.blox.dataservice.model.EnvironmentVersion;
import com.amazonaws.blox.dataservice.model.InstanceGroup;
import com.amazonaws.blox.dataservice.repository.model.EnvironmentDDBRecord;
import com.amazonaws.blox.dataservice.repository.model.EnvironmentTargetVersionDDBRecord;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentExistsException;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBMapper;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBQueryExpression;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBScanExpression;
import com.amazonaws.services.dynamodbv2.datamodeling.PaginatedQueryList;
import com.amazonaws.services.dynamodbv2.datamodeling.PaginatedScanList;
import com.amazonaws.services.dynamodbv2.model.ConditionalCheckFailedException;
import java.time.Instant;
import java.util.ArrayList;
import java.util.List;
import java.util.Set;
import java.util.UUID;
import java.util.stream.Collectors;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mapstruct.factory.Mappers;
import org.mockito.ArgumentCaptor;
import org.mockito.Captor;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;

@RunWith(MockitoJUnitRunner.StrictStubs.class)
public class EnvironmentRepositoryDDBTest {

  private static final String ENVIRONMENT_NAME = "test";
  private static final String ACCOUNT_ID = "12345678912";
  private static final String ENVIRONMENT_ID =
      EnvironmentArnGenerator.generateEnvironmentArn(ENVIRONMENT_NAME, ACCOUNT_ID);
  private static final String SECOND_ENVIRONMENT_ID =
      EnvironmentArnGenerator.generateEnvironmentArn(
          ENVIRONMENT_NAME + UUID.randomUUID().toString(), ACCOUNT_ID);
  private static final String ENVIRONMENT_VERSION = UUID.randomUUID().toString();
  private static final String TASK_DEF_ARN =
      "arn:aws:ecs:us-east-1:" + ACCOUNT_ID + ":task-definition/sleep";
  private static final String ROLE_ARN = "arn:aws:iam::" + ACCOUNT_ID + ":role/testRole";
  private static final String CLUSTER_ARN_PREFIX =
      "arn:aws:ecs:us-east-1:" + ACCOUNT_ID + ":cluster/";
  private static final String CLUSTER = CLUSTER_ARN_PREFIX + UUID.randomUUID().toString();
  private static final String SECOND_CLUSTER = CLUSTER_ARN_PREFIX + UUID.randomUUID().toString();

  @Mock private DynamoDBMapper dynamoDBMapper;

  @Mock
  private PaginatedScanList<EnvironmentTargetVersionDDBRecord>
      environmentTargetVersionDDBRecordPaginatedScanList;

  @Mock
  private PaginatedQueryList<EnvironmentTargetVersionDDBRecord>
      environmentTargetVersionDDBRecordPaginatedQueryList;

  @Captor private ArgumentCaptor<EnvironmentDDBRecord> environmentDDBRecordCaptor;

  @Captor
  private ArgumentCaptor<EnvironmentTargetVersionDDBRecord> environmentTargetVersionDDBRecordCaptor;

  @Captor private ArgumentCaptor<DynamoDBScanExpression> dynamoDBScanExpressionCaptor;
  @Captor private ArgumentCaptor<DynamoDBQueryExpression> dynamoDBQueryExpressionCaptor;

  private EnvironmentRepositoryDDB environmentRepositoryDDB;
  private Environment environment;
  private EnvironmentDDBRecord environmentDDBRecord;
  private EnvironmentVersion environmentVersion;
  private EnvironmentTargetVersionDDBRecord environmentTargetVersionDDBRecord;

  @Before
  public void setup() {
    environmentRepositoryDDB =
        new EnvironmentRepositoryDDB(dynamoDBMapper, Mappers.getMapper(EnvironmentMapper.class));

    environment =
        Environment.builder()
            .environmentName(ENVIRONMENT_NAME)
            .environmentId(
                EnvironmentArnGenerator.generateEnvironmentArn(ACCOUNT_ID, ENVIRONMENT_NAME))
            .environmentVersion(UUID.randomUUID().toString())
            .taskDefinition(TASK_DEF_ARN)
            .role(ROLE_ARN)
            .instanceGroup(InstanceGroup.builder().cluster(CLUSTER).build())
            .type(EnvironmentType.Daemon)
            .status(EnvironmentStatus.Inactive)
            .health(EnvironmentHealth.Healthy)
            .createdTime(Instant.now())
            .build();

    environmentDDBRecord =
        EnvironmentDDBRecord.builder()
            .environmentId(environment.getEnvironmentId())
            .environmentVersion(environment.getEnvironmentVersion())
            .environmentName(environment.getEnvironmentName())
            .taskDefinition(environment.getTaskDefinition())
            .role(environment.getRole())
            .cluster(environment.getInstanceGroup().getCluster())
            .type(environment.getType())
            .status(environment.getStatus())
            .health(environment.getHealth())
            .createdTime(environment.getCreatedTime())
            .build();

    environmentVersion =
        EnvironmentVersion.builder()
            .environmentId(environment.getEnvironmentId())
            .environmentVersion(environment.getEnvironmentVersion())
            .cluster(environment.getInstanceGroup().getCluster())
            .build();

    environmentTargetVersionDDBRecord =
        EnvironmentTargetVersionDDBRecord.builder()
            .environmentId(environmentVersion.getEnvironmentId())
            .environmentVersion(environmentVersion.getEnvironmentVersion())
            .cluster(environmentVersion.getCluster())
            .build();
  }

  @Test(expected = NullPointerException.class)
  public void constructorDynamoDbMapperNull() {
    new EnvironmentRepositoryDDB(null, Mappers.getMapper(EnvironmentMapper.class));
  }

  @Test(expected = NullPointerException.class)
  public void constructorEnvironmentMapperNull() {
    new EnvironmentRepositoryDDB(dynamoDBMapper, null);
  }

  @Test(expected = NullPointerException.class)
  public void createEnvironmentNullEnvironment() throws Exception {
    environmentRepositoryDDB.createEnvironment(null);
  }

  @Test(expected = EnvironmentExistsException.class)
  public void createEnvironmentEnvironmentExistsException() throws Exception {
    doThrow(new ConditionalCheckFailedException(""))
        .when(dynamoDBMapper)
        .save(isA(EnvironmentDDBRecord.class));
    environmentRepositoryDDB.createEnvironment(environment);
  }

  @Test(expected = StorageException.class)
  public void createEnvironmentStorageException() throws Exception {
    doThrow(new AmazonServiceException(""))
        .when(dynamoDBMapper)
        .save(isA(EnvironmentDDBRecord.class));
    environmentRepositoryDDB.createEnvironment(environment);
  }

  @Test
  public void createEnvironment() throws Exception {
    doNothing().when(dynamoDBMapper).save(isA(EnvironmentDDBRecord.class));

    final Environment createdEnvironment = environmentRepositoryDDB.createEnvironment(environment);
    assertEquals(environment, createdEnvironment);

    verify(dynamoDBMapper).save(environmentDDBRecordCaptor.capture());
    //TODO: assert attributes
    assertEquals(
        environment.getInstanceGroup().getCluster(),
        environmentDDBRecordCaptor.getValue().getCluster());
    assertEquals(
        environmentDDBRecord.getCreatedTime(),
        environmentDDBRecordCaptor.getValue().getCreatedTime());
    assertEquals(
        environmentDDBRecord.getEnvironmentId(),
        environmentDDBRecordCaptor.getValue().getEnvironmentId());
    assertEquals(
        environment.getEnvironmentName(),
        environmentDDBRecordCaptor.getValue().getEnvironmentName());
    assertEquals(
        environment.getEnvironmentVersion(),
        environmentDDBRecordCaptor.getValue().getEnvironmentVersion());
    assertEquals(environment.getHealth(), environmentDDBRecordCaptor.getValue().getHealth());
    assertEquals(
        environmentDDBRecord.getLastUpdatedTime(),
        environmentDDBRecordCaptor.getValue().getLastUpdatedTime());
    assertEquals(environment.getRole(), environmentDDBRecordCaptor.getValue().getRole());
    assertEquals(environment.getStatus(), environmentDDBRecordCaptor.getValue().getStatus());
    assertEquals(
        environment.getTaskDefinition(), environmentDDBRecordCaptor.getValue().getTaskDefinition());
    assertEquals(environment.getType(), environmentDDBRecordCaptor.getValue().getType());
  }

  @Test(expected = NullPointerException.class)
  public void getEnvironmentNullEnvironmentId() throws Exception {
    environmentRepositoryDDB.getEnvironment(null, ENVIRONMENT_VERSION);
  }

  @Test(expected = NullPointerException.class)
  public void getEnvironmentNullEnvironmentVersion() throws Exception {
    environmentRepositoryDDB.getEnvironment(ENVIRONMENT_ID, null);
  }

  @Test(expected = StorageException.class)
  public void geEnvironmentStorageException() throws Exception {
    doThrow(new AmazonServiceException(""))
        .when(dynamoDBMapper)
        .load(isA(EnvironmentDDBRecord.class));
    environmentRepositoryDDB.getEnvironment(ENVIRONMENT_ID, ENVIRONMENT_VERSION);
  }

  @Test
  public void geEnvironment() throws Exception {
    final EnvironmentDDBRecord recordWithKeys =
        EnvironmentDDBRecord.withKeys(
            environment.getEnvironmentId(), environment.getEnvironmentVersion());
    when(dynamoDBMapper.load(recordWithKeys)).thenReturn(environmentDDBRecord);

    final Environment loadedEnvironment =
        environmentRepositoryDDB.getEnvironment(
            environment.getEnvironmentId(), environment.getEnvironmentVersion());

    assertEquals(environment, loadedEnvironment);
    verify(dynamoDBMapper).load(recordWithKeys);
  }

  @Test(expected = NullPointerException.class)
  public void createEnvironmentTargetVersionNullEnvironmentTargetVersion() throws Exception {
    environmentRepositoryDDB.createEnvironmentTargetVersion(null);
  }

  @Test(expected = EnvironmentExistsException.class)
  public void createEnvironmentTargetVersionEnvironmentExistsException() throws Exception {
    doThrow(new ConditionalCheckFailedException(""))
        .when(dynamoDBMapper)
        .save(isA(EnvironmentTargetVersionDDBRecord.class));
    environmentRepositoryDDB.createEnvironmentTargetVersion(environmentVersion);
  }

  @Test(expected = StorageException.class)
  public void createEnvironmentTargetVersionStorageException() throws Exception {
    doThrow(new AmazonServiceException(""))
        .when(dynamoDBMapper)
        .save(isA(EnvironmentTargetVersionDDBRecord.class));
    environmentRepositoryDDB.createEnvironmentTargetVersion(environmentVersion);
  }

  @Test
  public void createEnvironmentTargetVersion() throws Exception {
    doNothing().when(dynamoDBMapper).save(isA(EnvironmentTargetVersionDDBRecord.class));

    final EnvironmentVersion createdEnvironmentVersion =
        environmentRepositoryDDB.createEnvironmentTargetVersion(environmentVersion);
    assertEquals(environmentVersion, createdEnvironmentVersion);

    verify(dynamoDBMapper).save(environmentTargetVersionDDBRecordCaptor.capture());
    assertEquals(
        environmentVersion.getCluster(),
        environmentTargetVersionDDBRecordCaptor.getValue().getCluster());
    assertEquals(
        environmentVersion.getEnvironmentId(),
        environmentTargetVersionDDBRecordCaptor.getValue().getEnvironmentId());
    assertEquals(
        environmentVersion.getEnvironmentVersion(),
        environmentTargetVersionDDBRecordCaptor.getValue().getEnvironmentVersion());
  }

  @Test(expected = NullPointerException.class)
  public void getEnvironmentTargetVersionNullEnvironmentId() throws Exception {
    environmentRepositoryDDB.getEnvironmentTargetVersion(null);
  }

  @Test(expected = StorageException.class)
  public void getEnvironmentTargetVersionStorageException() throws Exception {
    doThrow(new AmazonServiceException(""))
        .when(dynamoDBMapper)
        .load(isA(EnvironmentTargetVersionDDBRecord.class));
    environmentRepositoryDDB.getEnvironmentTargetVersion(environmentVersion.getEnvironmentId());
  }

  @Test
  public void getEnvironmentTargetVersion() throws Exception {
    when(dynamoDBMapper.load(isA(EnvironmentTargetVersionDDBRecord.class)))
        .thenReturn(environmentTargetVersionDDBRecord);

    final EnvironmentVersion environmentVersionResult =
        environmentRepositoryDDB.getEnvironmentTargetVersion(environmentVersion.getEnvironmentId());
    assertEquals(environmentVersion, environmentVersionResult);

    verify(dynamoDBMapper)
        .load(EnvironmentTargetVersionDDBRecord.withHashKey(environmentVersion.getEnvironmentId()));
  }

  @Test(expected = StorageException.class)
  public void listClustersStorageException() throws Exception {
    when(dynamoDBMapper.scan(
            eq(EnvironmentTargetVersionDDBRecord.class), isA(DynamoDBScanExpression.class)))
        .thenThrow(new AmazonServiceException(""));

    environmentRepositoryDDB.listClusters();
  }

  @Test
  public void listClusters() throws Exception {
    List<EnvironmentTargetVersionDDBRecord> environmentTargetVersionDDBRecords = new ArrayList<>();
    environmentTargetVersionDDBRecords.add(environmentTargetVersionDDBRecord);
    final EnvironmentTargetVersionDDBRecord secondRecord =
        EnvironmentTargetVersionDDBRecord.builder()
            .environmentId(SECOND_ENVIRONMENT_ID)
            .environmentVersion(UUID.randomUUID().toString())
            .cluster(SECOND_CLUSTER)
            .build();
    environmentTargetVersionDDBRecords.add(secondRecord);

    final EnvironmentTargetVersionDDBRecord thirdRecordSameCluster =
        EnvironmentTargetVersionDDBRecord.builder()
            .environmentId(SECOND_ENVIRONMENT_ID)
            .environmentVersion(UUID.randomUUID().toString())
            .cluster(SECOND_CLUSTER)
            .build();
    environmentTargetVersionDDBRecords.add(thirdRecordSameCluster);

    when(environmentTargetVersionDDBRecordPaginatedScanList.stream())
        .thenReturn(environmentTargetVersionDDBRecords.stream());

    when(dynamoDBMapper.scan(
            eq(EnvironmentTargetVersionDDBRecord.class), isA(DynamoDBScanExpression.class)))
        .thenReturn(environmentTargetVersionDDBRecordPaginatedScanList);

    final List<String> clusters = environmentRepositoryDDB.listClusters();
    final Set<String> expectedClusters =
        environmentTargetVersionDDBRecords
            .stream()
            .map(r -> r.getCluster())
            .collect(Collectors.toSet());
    assertThat(clusters, containsInAnyOrder(expectedClusters.toArray()));

    verify(dynamoDBMapper)
        .scan(eq(EnvironmentTargetVersionDDBRecord.class), dynamoDBScanExpressionCaptor.capture());
    final DynamoDBScanExpression expectedDynamodbScanExpression =
        new DynamoDBScanExpression()
            .withIndexName(EnvironmentTargetVersionDDBRecord.ENVIRONMENT_CLUSTER_GSI_NAME)
            .withConsistentRead(false);

    assertEquals(
        expectedDynamodbScanExpression.getIndexName(),
        dynamoDBScanExpressionCaptor.getValue().getIndexName());
  }

  @Test(expected = NullPointerException.class)
  public void listEnvironmentIdsByClusterNullCluster() throws Exception {
    environmentRepositoryDDB.listEnvironmentIdsByCluster(null);
  }

  @SuppressWarnings("unchecked")
  @Test(expected = StorageException.class)
  public void listEnvironmentIdsByClusterStorageException() throws Exception {
    when(dynamoDBMapper.query(
            eq(EnvironmentTargetVersionDDBRecord.class), isA(DynamoDBQueryExpression.class)))
        .thenThrow(new AmazonServiceException(""));
    environmentRepositoryDDB.listEnvironmentIdsByCluster(CLUSTER);
  }

  @SuppressWarnings("unchecked")
  @Test
  public void listEnvironmentIdsByCluster() throws Exception {
    List<EnvironmentTargetVersionDDBRecord> environmentTargetVersionDDBRecords = new ArrayList<>();
    environmentTargetVersionDDBRecords.add(environmentTargetVersionDDBRecord);
    final EnvironmentTargetVersionDDBRecord secondRecord =
        EnvironmentTargetVersionDDBRecord.builder()
            .environmentId(SECOND_ENVIRONMENT_ID)
            .environmentVersion(UUID.randomUUID().toString())
            .cluster(environmentTargetVersionDDBRecord.getCluster())
            .build();
    environmentTargetVersionDDBRecords.add(secondRecord);
    when(environmentTargetVersionDDBRecordPaginatedQueryList.stream())
        .thenReturn(environmentTargetVersionDDBRecords.stream());

    when(dynamoDBMapper.query(
            eq(EnvironmentTargetVersionDDBRecord.class), isA(DynamoDBQueryExpression.class)))
        .thenReturn(environmentTargetVersionDDBRecordPaginatedQueryList);

    final List<String> environmentsByCluster =
        environmentRepositoryDDB.listEnvironmentIdsByCluster(CLUSTER);

    final Set<String> expectedEnvironmentIds =
        environmentTargetVersionDDBRecords
            .stream()
            .map(r -> r.getEnvironmentId())
            .collect(Collectors.toSet());
    assertThat(environmentsByCluster, containsInAnyOrder(expectedEnvironmentIds.toArray()));

    verify(dynamoDBMapper)
        .query(
            eq(EnvironmentTargetVersionDDBRecord.class), dynamoDBQueryExpressionCaptor.capture());
    final DynamoDBQueryExpression<EnvironmentTargetVersionDDBRecord>
        expectedDynamodbQueryExpression =
            new DynamoDBQueryExpression<EnvironmentTargetVersionDDBRecord>()
                .withIndexName(EnvironmentTargetVersionDDBRecord.ENVIRONMENT_CLUSTER_GSI_NAME)
                .withConsistentRead(false)
                .withHashKeyValues(
                    EnvironmentTargetVersionDDBRecord.withGSIHashKey(
                        environmentTargetVersionDDBRecord.getCluster()));

    assertEquals(
        expectedDynamodbQueryExpression.getIndexName(),
        dynamoDBQueryExpressionCaptor.getValue().getIndexName());
    assertEquals(
        expectedDynamodbQueryExpression.getHashKeyValues(),
        dynamoDBQueryExpressionCaptor.getValue().getHashKeyValues());
  }
}
