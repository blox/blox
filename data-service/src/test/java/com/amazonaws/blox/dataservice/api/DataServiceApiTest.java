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

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.mockito.ArgumentMatchers.isA;
import static org.mockito.Mockito.when;

import com.amazonaws.blox.dataservice.arn.EnvironmentArnGenerator;
import com.amazonaws.blox.dataservice.handler.EnvironmentHandler;
import com.amazonaws.blox.dataservice.mapper.ApiModelMapper;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentHealth;
import com.amazonaws.blox.dataservice.model.EnvironmentStatus;
import com.amazonaws.blox.dataservice.model.EnvironmentVersion;
import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentExistsException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ServiceException;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentType;
import com.amazonaws.blox.dataservicemodel.v1.model.InstanceGroup;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateTargetEnvironmentRevisionRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateTargetEnvironmentRevisionResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRequest;
import java.time.Instant;
import java.util.UUID;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mapstruct.factory.Mappers;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;

@RunWith(MockitoJUnitRunner.StrictStubs.class)
public class DataServiceApiTest {

  private static final String ENVIRONMENT_NAME = "test";
  private static final String ACCOUNT_ID = "12345678912";
  private static final String TASK_DEFINITION_ARN =
      "arn:aws:ecs:us-east-1:" + ACCOUNT_ID + ":task-definition/sleep";
  private static final String ROLE_ARN = "arn:aws:iam::" + ACCOUNT_ID + ":role/testRole";
  private static final String CLUSTER_ARN = "arn:aws:ecs:us-east-1:" + ACCOUNT_ID + ":cluster/test";

  @Mock private EnvironmentHandler environmentHandler;

  private ApiModelMapper apiModelMapper = Mappers.getMapper(ApiModelMapper.class);
  private DataService dataService;

  private CreateEnvironmentRequest createEnvironmentRequest;
  private Environment environment;
  private Environment createdEnvironment;
  private CreateTargetEnvironmentRevisionRequest createTargetEnvironmentRevisionRequest;
  private EnvironmentVersion environmentVersion;
  private DescribeEnvironmentRequest describeEnvironmentRequest;

  @Before
  public void setup() {
    dataService = new DataServiceApi(environmentHandler, apiModelMapper);
    createEnvironmentRequest =
        CreateEnvironmentRequest.builder()
            .environmentName(ENVIRONMENT_NAME)
            .accountId(ACCOUNT_ID)
            .environmentType(EnvironmentType.Daemon)
            .instanceGroup(InstanceGroup.builder().cluster(CLUSTER_ARN).build())
            .role(ROLE_ARN)
            .taskDefinition(TASK_DEFINITION_ARN)
            .build();

    Environment.EnvironmentBuilder environmentBuilder =
        Environment.builder()
            .environmentName(createEnvironmentRequest.getEnvironmentName())
            .environmentId(
                EnvironmentArnGenerator.generateEnvironmentArn(
                    createEnvironmentRequest.getEnvironmentName(),
                    createEnvironmentRequest.getAccountId()))
            .type(com.amazonaws.blox.dataservice.model.EnvironmentType.Daemon)
            .instanceGroup(
                com.amazonaws.blox.dataservice.model.InstanceGroup.builder()
                    .cluster(CLUSTER_ARN)
                    .build())
            .role(createEnvironmentRequest.getRole())
            .taskDefinition(createEnvironmentRequest.getTaskDefinition());

    environment = environmentBuilder.build();
    createdEnvironment =
        environmentBuilder
            .environmentVersion(UUID.randomUUID().toString())
            .health(EnvironmentHealth.Healthy)
            .status(EnvironmentStatus.Inactive)
            .createdTime(Instant.now())
            .build();

    createTargetEnvironmentRevisionRequest =
        CreateTargetEnvironmentRevisionRequest.builder()
            .environmentId(createdEnvironment.getEnvironmentId())
            .environmentVersion(createdEnvironment.getEnvironmentVersion())
            .build();

    environmentVersion =
        EnvironmentVersion.builder()
            .environmentId(createTargetEnvironmentRevisionRequest.getEnvironmentId())
            .environmentVersion(createTargetEnvironmentRevisionRequest.getEnvironmentVersion())
            .cluster(createdEnvironment.getInstanceGroup().getCluster())
            .build();

    describeEnvironmentRequest =
        DescribeEnvironmentRequest.builder()
            .environmentId(createTargetEnvironmentRevisionRequest.getEnvironmentId())
            .environmentVersion(createTargetEnvironmentRevisionRequest.getEnvironmentVersion())
            .build();
  }

  @Test(expected = NullPointerException.class)
  public void constructorEnvironmentHandlerNull() {
    new DataServiceApi(null, apiModelMapper);
  }

  @Test(expected = NullPointerException.class)
  public void constructorMapperNull() {
    new DataServiceApi(environmentHandler, null);
  }

  @Test(expected = NullPointerException.class)
  public void createEnvironmentNullRequest() throws Exception {
    dataService.createEnvironment(null);
  }

  @Test(expected = EnvironmentExistsException.class)
  public void createEnvironmentEnvironmentExistsException() throws Exception {
    when(environmentHandler.createEnvironment(isA(Environment.class)))
        .thenThrow(new EnvironmentExistsException(""));

    dataService.createEnvironment(createEnvironmentRequest);
  }

  @Test(expected = ServiceException.class)
  public void createEnvironmentServiceException() throws Exception {
    when(environmentHandler.createEnvironment(isA(Environment.class)))
        .thenThrow(new ServiceException(""));

    dataService.createEnvironment(createEnvironmentRequest);
  }

  @Test
  public void createEnvironment() throws Exception {
    when(environmentHandler.createEnvironment(environment)).thenReturn(createdEnvironment);

    CreateEnvironmentResponse createEnvironmentResponse =
        dataService.createEnvironment(createEnvironmentRequest);
    assertEquals(
        createdEnvironment.getEnvironmentVersion(),
        createEnvironmentResponse.getEnvironmentVersion());
    assertEquals(
        createdEnvironment.getEnvironmentId(), createEnvironmentResponse.getEnvironmentId());
    assertEquals(
        createdEnvironment.getEnvironmentName(), createEnvironmentResponse.getEnvironmentName());
    assertEquals(
        createdEnvironment.getTaskDefinition(), createEnvironmentResponse.getTaskDefinition());
    assertEquals(createdEnvironment.getRole(), createEnvironmentResponse.getRole());
    assertEquals(
        createdEnvironment.getInstanceGroup().getCluster(),
        createEnvironmentResponse.getInstanceGroup().getCluster());
    assertEquals(
        createdEnvironment.getInstanceGroup().getAttributes(),
        createEnvironmentResponse.getInstanceGroup().getAttributes());
    assertEquals(
        createdEnvironment.getType().name(), createEnvironmentResponse.getEnvironmentType().name());
    assertEquals(
        createdEnvironment.getHealth().name(), createEnvironmentResponse.getEnvironmentHealth());
    assertEquals(
        createdEnvironment.getStatus().name(), createEnvironmentResponse.getEnvironmentStatus());
    assertNotNull(createEnvironmentResponse.getCreatedTime());
  }

  @Test(expected = NullPointerException.class)
  public void createTargetEnvironmentRevisionNullRequest() throws Exception {
    dataService.createTargetEnvironmentRevision(null);
  }

  @Test(expected = EnvironmentNotFoundException.class)
  public void createTargetEnvironmentRevisionEnvironmentNotFoundException() throws Exception {
    when(environmentHandler.createEnvironmentTargetVersion(
            createTargetEnvironmentRevisionRequest.getEnvironmentId(),
            createTargetEnvironmentRevisionRequest.getEnvironmentVersion()))
        .thenThrow(new EnvironmentNotFoundException(""));
    dataService.createTargetEnvironmentRevision(createTargetEnvironmentRevisionRequest);
  }

  @Test(expected = EnvironmentExistsException.class)
  public void createTargetEnvironmentRevisionEnvironmentExistsException() throws Exception {
    when(environmentHandler.createEnvironmentTargetVersion(
            createTargetEnvironmentRevisionRequest.getEnvironmentId(),
            createTargetEnvironmentRevisionRequest.getEnvironmentVersion()))
        .thenThrow(new EnvironmentExistsException(""));
    dataService.createTargetEnvironmentRevision(createTargetEnvironmentRevisionRequest);
  }

  @Test(expected = ServiceException.class)
  public void createTargetEnvironmentRevisionServiceException() throws Exception {
    when(environmentHandler.createEnvironmentTargetVersion(
            createTargetEnvironmentRevisionRequest.getEnvironmentId(),
            createTargetEnvironmentRevisionRequest.getEnvironmentVersion()))
        .thenThrow(new ServiceException(""));
    dataService.createTargetEnvironmentRevision(createTargetEnvironmentRevisionRequest);
  }

  @Test
  public void createTargetEnvironmentRevision() throws Exception {
    when(environmentHandler.createEnvironmentTargetVersion(
            createTargetEnvironmentRevisionRequest.getEnvironmentId(),
            createTargetEnvironmentRevisionRequest.getEnvironmentVersion()))
        .thenReturn(environmentVersion);

    final CreateTargetEnvironmentRevisionResponse createTargetEnvironmentRevisionResponse =
        dataService.createTargetEnvironmentRevision(createTargetEnvironmentRevisionRequest);
    assertEquals(
        environmentVersion.getEnvironmentId(),
        createTargetEnvironmentRevisionResponse.getEnvironmentId());
    assertEquals(
        environmentVersion.getEnvironmentVersion(),
        createTargetEnvironmentRevisionResponse.getEnvironmentVersion());
    assertEquals(
        environmentVersion.getCluster(), createTargetEnvironmentRevisionResponse.getCluster());
  }

  @Test(expected = NullPointerException.class)
  public void describeEnvironmentNullRequest() throws Exception {
    dataService.describeEnvironment(null);
  }
}
