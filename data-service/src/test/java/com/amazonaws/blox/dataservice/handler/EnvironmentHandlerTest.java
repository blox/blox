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

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.mockito.Mockito.isA;
import static org.mockito.Mockito.when;

import com.amazonaws.blox.dataservice.arn.EnvironmentArnGenerator;
import com.amazonaws.blox.dataservice.exception.StorageException;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentHealth;
import com.amazonaws.blox.dataservice.model.EnvironmentStatus;
import com.amazonaws.blox.dataservice.model.EnvironmentType;
import com.amazonaws.blox.dataservice.model.EnvironmentVersion;
import com.amazonaws.blox.dataservice.model.InstanceGroup;
import com.amazonaws.blox.dataservice.repository.EnvironmentRepository;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentExistsException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentVersionNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ServiceException;
import java.time.Instant;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.UUID;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;

@RunWith(MockitoJUnitRunner.StrictStubs.class)
public class EnvironmentHandlerTest {

  private static final String ENVIRONMENT_NAME = "test";
  private static final String ACCOUNT_ID = "12345678912";
  private static final String ENVIRONMENT_VERSION = UUID.randomUUID().toString();
  private static final String ENVIRONMENT_ID =
      EnvironmentArnGenerator.generateEnvironmentArn(ENVIRONMENT_NAME, ACCOUNT_ID);
  private static final String TASK_DEFINITION_ARN =
      "arn:aws:ecs:us-east-1:" + ACCOUNT_ID + ":task-definition/sleep";
  private static final String ROLE_ARN = "arn:aws:iam::" + ACCOUNT_ID + ":role/testRole";
  private static final String CLUSTER_ARN = "arn:aws:ecs:us-east-1:" + ACCOUNT_ID + ":cluster/test";

  private EnvironmentHandler environmentHandler;
  private Environment createEnvironment;
  private Environment getEnvironment;

  @Mock private EnvironmentRepository environmentRepository;

  @Before
  public void setup() {
    environmentHandler = new EnvironmentHandler(environmentRepository);
    Environment.EnvironmentBuilder createEnvironmentBuilder =
        Environment.builder()
            .environmentId(ENVIRONMENT_ID)
            .environmentName(ENVIRONMENT_NAME)
            .taskDefinition(TASK_DEFINITION_ARN)
            .role(ROLE_ARN)
            .instanceGroup(InstanceGroup.builder().cluster(CLUSTER_ARN).build())
            .type(EnvironmentType.Daemon);

    createEnvironment = createEnvironmentBuilder.build();

    getEnvironment =
        createEnvironmentBuilder
            .health(EnvironmentHealth.Healthy)
            .status(EnvironmentStatus.Inactive)
            .createdTime(Instant.now())
            .environmentVersion(ENVIRONMENT_VERSION)
            .build();
  }

  @Test(expected = NullPointerException.class)
  public void constructorNullEnvironmentRepository() {
    new EnvironmentHandler(null);
  }

  @Test(expected = NullPointerException.class)
  public void createEnvironmentNullEnvironment() throws Exception {
    environmentHandler.createEnvironment(null);
  }

  @Test(expected = ServiceException.class)
  public void createEnvironmentServiceException() throws Exception {
    when(environmentRepository.createEnvironment(isA(Environment.class)))
        .thenThrow(new StorageException(""));

    environmentHandler.createEnvironment(createEnvironment);
  }

  @Test
  public void createEnvironment() throws Exception {
    final Environment expectedEnvironment = createEnvironment;
    expectedEnvironment.setHealth(EnvironmentHealth.Healthy);
    expectedEnvironment.setStatus(EnvironmentStatus.Inactive);

    when(environmentRepository.createEnvironment(isA(Environment.class)))
        .thenReturn(expectedEnvironment);

    final Environment result = environmentHandler.createEnvironment(createEnvironment);
    assertNotNull(result.getEnvironmentVersion());
    assertNotNull(result.getCreatedTime());

    assertEquals(expectedEnvironment.getEnvironmentId(), createEnvironment.getEnvironmentId());
    assertEquals(expectedEnvironment.getEnvironmentName(), createEnvironment.getEnvironmentName());
    assertEquals(expectedEnvironment.getTaskDefinition(), createEnvironment.getTaskDefinition());
    assertEquals(expectedEnvironment.getRole(), createEnvironment.getRole());
    assertEquals(
        expectedEnvironment.getInstanceGroup().getCluster(),
        createEnvironment.getInstanceGroup().getCluster());
    assertEquals(expectedEnvironment.getType(), createEnvironment.getType());
    assertEquals(expectedEnvironment.getHealth(), createEnvironment.getHealth());
    assertEquals(expectedEnvironment.getStatus(), createEnvironment.getStatus());
  }

  @Test(expected = NullPointerException.class)
  public void createEnvironmentTargetVersionNullEnvironmentId() throws Exception {
    environmentHandler.createEnvironmentTargetVersion(null, ENVIRONMENT_VERSION);
  }

  @Test(expected = NullPointerException.class)
  public void createEnvironmentTargetVersionNullEnvironmentVersion()
      throws EnvironmentNotFoundException, EnvironmentExistsException, ServiceException {
    environmentHandler.createEnvironmentTargetVersion(ENVIRONMENT_ID, null);
  }

  @Test(expected = ServiceException.class)
  public void createEnvironmentTargetVersionGetFails() throws Exception {
    when(environmentRepository.getEnvironment(ENVIRONMENT_ID, ENVIRONMENT_VERSION))
        .thenThrow(new StorageException(""));
    environmentHandler.createEnvironmentTargetVersion(ENVIRONMENT_ID, ENVIRONMENT_VERSION);
  }

  @Test(expected = EnvironmentNotFoundException.class)
  public void createEnvironmentTargetVersionGetReturnsNull() throws Exception {
    when(environmentRepository.getEnvironment(ENVIRONMENT_ID, ENVIRONMENT_VERSION))
        .thenReturn(null);
    environmentHandler.createEnvironmentTargetVersion(ENVIRONMENT_ID, ENVIRONMENT_VERSION);
  }

  @Test(expected = ServiceException.class)
  public void createEnvironmentTargetVersionCreateFails() throws Exception {
    when(environmentRepository.getEnvironment(ENVIRONMENT_ID, ENVIRONMENT_VERSION))
        .thenReturn(getEnvironment);
    when(environmentRepository.createEnvironmentTargetVersion(
            EnvironmentVersion.builder()
                .environmentId(getEnvironment.getEnvironmentId())
                .environmentVersion(getEnvironment.getEnvironmentVersion())
                .cluster(getEnvironment.getInstanceGroup().getCluster())
                .build()))
        .thenThrow(new StorageException(""));
    environmentHandler.createEnvironmentTargetVersion(ENVIRONMENT_ID, ENVIRONMENT_VERSION);
  }

  @Test(expected = NullPointerException.class)
  public void describeEnvironmentTargetVersionNullEnvironmentId() throws Exception {
    environmentHandler.describeEnvironmentTargetVersion(null);
  }

  @Test(expected = ServiceException.class)
  public void describeEnvironmentTargetVersionGetEnvironmentTargetVersionFails() throws Exception {
    when(environmentRepository.getEnvironmentTargetVersion(ENVIRONMENT_ID))
        .thenThrow(new StorageException(""));
    environmentHandler.describeEnvironmentTargetVersion(ENVIRONMENT_ID);
  }

  @Test(expected = EnvironmentVersionNotFoundException.class)
  public void describeEnvironmentTargetVersionGetEnvironmentTargetVersionNull() throws Exception {
    when(environmentRepository.getEnvironmentTargetVersion(ENVIRONMENT_ID)).thenReturn(null);
    environmentHandler.describeEnvironmentTargetVersion(ENVIRONMENT_ID);
  }

  @Test
  public void describeEnvironmentTargetVersion() throws Exception {
    final EnvironmentVersion expectedEnvironmentVersion =
        EnvironmentVersion.builder()
            .environmentId(getEnvironment.getEnvironmentId())
            .environmentVersion(getEnvironment.getEnvironmentVersion())
            .cluster(getEnvironment.getInstanceGroup().getCluster())
            .build();
    when(environmentRepository.getEnvironmentTargetVersion(
            expectedEnvironmentVersion.getEnvironmentId()))
        .thenReturn(expectedEnvironmentVersion);

    final EnvironmentVersion result =
        environmentHandler.describeEnvironmentTargetVersion(
            expectedEnvironmentVersion.getEnvironmentId());
    assertEquals(expectedEnvironmentVersion, result);
  }

  @Test(expected = NullPointerException.class)
  public void describeEnvironmentEnvironmentIdNull() throws Exception {
    environmentHandler.describeEnvironment(null, ENVIRONMENT_VERSION);
  }

  @Test(expected = NullPointerException.class)
  public void describeEnvironmentEnvironmentVersionNull() throws Exception {
    environmentHandler.describeEnvironment(ENVIRONMENT_ID, null);
  }

  @Test(expected = ServiceException.class)
  public void describeEnvironmentGetEnvironmentFails() throws Exception {
    when(environmentRepository.getEnvironment(ENVIRONMENT_ID, ENVIRONMENT_VERSION))
        .thenThrow(new StorageException(""));
    environmentHandler.describeEnvironment(ENVIRONMENT_ID, ENVIRONMENT_VERSION);
  }

  @Test(expected = EnvironmentNotFoundException.class)
  public void describeEnvironmentGetEnvironmentNull() throws Exception {
    when(environmentRepository.getEnvironment(ENVIRONMENT_ID, ENVIRONMENT_VERSION))
        .thenReturn(null);
    environmentHandler.describeEnvironment(ENVIRONMENT_ID, ENVIRONMENT_VERSION);
  }

  @Test
  public void describeEnvironment() throws Exception {
    when(environmentRepository.getEnvironment(
            getEnvironment.getEnvironmentId(), getEnvironment.getEnvironmentVersion()))
        .thenReturn(getEnvironment);

    final Environment result =
        environmentHandler.describeEnvironment(
            getEnvironment.getEnvironmentId(), getEnvironment.getEnvironmentVersion());
    assertEquals(getEnvironment, result);
  }

  @Test(expected = NullPointerException.class)
  public void listEnvironmentsWithClusterNullCluster() throws Exception {
    environmentHandler.listEnvironmentsWithCluster(null);
  }

  @Test(expected = ServiceException.class)
  public void listEnvironmentsWithClusterServiceException() throws Exception {
    when(environmentRepository.listEnvironmentIdsByCluster(CLUSTER_ARN))
        .thenThrow(new StorageException(""));
    environmentHandler.listEnvironmentsWithCluster(CLUSTER_ARN);
  }

  @Test
  public void listEnvironmentsWithCluster() throws Exception {
    final List<String> expectedEnvironments = new ArrayList<>(Arrays.asList(ENVIRONMENT_NAME));
    when(environmentRepository.listEnvironmentIdsByCluster(CLUSTER_ARN))
        .thenReturn(expectedEnvironments);
    final List<String> result = environmentHandler.listEnvironmentsWithCluster(CLUSTER_ARN);
    assertEquals(expectedEnvironments, result);
  }

  @Test(expected = ServiceException.class)
  public void listClustersServiceException() throws Exception {
    when(environmentRepository.listClusters()).thenThrow(new StorageException(""));
    environmentHandler.listClusters();
  }

  @Test
  public void listClusters() throws Exception {
    final List<String> expectedClusters = new ArrayList<>(Arrays.asList(CLUSTER_ARN));
    when(environmentRepository.listClusters()).thenReturn(expectedClusters);
    final List<String> result = environmentHandler.listClusters();
    assertEquals(expectedClusters, result);
  }
}
