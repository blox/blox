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

import static org.junit.Assert.assertEquals;
import static org.mockito.Mockito.doNothing;
import static org.mockito.Mockito.doThrow;
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
import com.amazonaws.blox.dataservice.model.InstanceGroup;
import com.amazonaws.blox.dataservice.repository.model.EnvironmentDDBRecord;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentExistsException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentNotFoundException;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBMapper;
import java.time.Instant;
import java.util.UUID;
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
  private static final String ENVIRONMENT_VERSION = "12345678912";
  private static final String TASK_DEF = "sleep";
  private static final String ROLE = "roleArn";
  private static final String CLUSTER = "clusterArn";

  @Mock private DynamoDBMapper dynamoDBMapper;
  @Captor private ArgumentCaptor<EnvironmentDDBRecord> environmentDDBRecordCaptor;

  private EnvironmentRepositoryDDB environmentRepositoryDDB;
  private Environment environment;
  private EnvironmentDDBRecord environmentDDBRecord;

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
            .taskDefinition(TASK_DEF)
            .role(ROLE)
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
  public void createEnvironmentNullEnvironment()
      throws StorageException, EnvironmentExistsException {
    environmentRepositoryDDB.createEnvironment(null);
  }

  @Test(expected = StorageException.class)
  public void createEnvironmentStorageException()
      throws StorageException, EnvironmentExistsException {
    doThrow(new AmazonServiceException(""))
        .when(dynamoDBMapper)
        .save(isA(EnvironmentDDBRecord.class));
    environmentRepositoryDDB.createEnvironment(environment);
  }

  @Test
  public void createEnvironment() throws StorageException, EnvironmentExistsException {
    doNothing().when(dynamoDBMapper).save(isA(EnvironmentDDBRecord.class));

    final Environment createdEnvironment = environmentRepositoryDDB.createEnvironment(environment);
    assertEquals(environment, createdEnvironment);

    verify(dynamoDBMapper).save(environmentDDBRecordCaptor.capture());
    assertEquals(
        environmentDDBRecord.getAttributes(),
        environmentDDBRecordCaptor.getValue().getAttributes());
    assertEquals(
        environmentDDBRecord.getCluster(), environmentDDBRecordCaptor.getValue().getCluster());
    assertEquals(
        environmentDDBRecord.getCreatedTime(),
        environmentDDBRecordCaptor.getValue().getCreatedTime());
    assertEquals(
        environmentDDBRecord.getEnvironmentId(),
        environmentDDBRecordCaptor.getValue().getEnvironmentId());
    assertEquals(
        environmentDDBRecord.getEnvironmentName(),
        environmentDDBRecordCaptor.getValue().getEnvironmentName());
    assertEquals(
        environmentDDBRecord.getEnvironmentVersion(),
        environmentDDBRecordCaptor.getValue().getEnvironmentVersion());
    assertEquals(
        environmentDDBRecord.getHealth(), environmentDDBRecordCaptor.getValue().getHealth());
    assertEquals(
        environmentDDBRecord.getLastUpdatedTime(),
        environmentDDBRecordCaptor.getValue().getLastUpdatedTime());
    assertEquals(environmentDDBRecord.getRole(), environmentDDBRecordCaptor.getValue().getRole());
    assertEquals(
        environmentDDBRecord.getStatus(), environmentDDBRecordCaptor.getValue().getStatus());
    assertEquals(
        environmentDDBRecord.getTaskDefinition(),
        environmentDDBRecordCaptor.getValue().getTaskDefinition());
    assertEquals(environmentDDBRecord.getType(), environmentDDBRecordCaptor.getValue().getType());
  }

  @Test(expected = NullPointerException.class)
  public void getEnvironmentNullEnvironmentId()
      throws StorageException, EnvironmentNotFoundException {
    environmentRepositoryDDB.getEnvironment(null, ENVIRONMENT_VERSION);
  }

  @Test(expected = NullPointerException.class)
  public void getEnvironmentNullEnvironmentVersion()
      throws StorageException, EnvironmentNotFoundException {
    environmentRepositoryDDB.getEnvironment(ENVIRONMENT_ID, null);
  }

  @Test(expected = StorageException.class)
  public void geEnvironmentStorageException()
      throws StorageException, EnvironmentNotFoundException {
    doThrow(new AmazonServiceException(""))
        .when(dynamoDBMapper)
        .load(isA(EnvironmentDDBRecord.class));
    environmentRepositoryDDB.getEnvironment(ENVIRONMENT_ID, ENVIRONMENT_VERSION);
  }

  @Test
  public void geEnvironment() throws StorageException, EnvironmentNotFoundException {
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
}
