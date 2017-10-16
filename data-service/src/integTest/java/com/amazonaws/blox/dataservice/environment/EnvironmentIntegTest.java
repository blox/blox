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
package com.amazonaws.blox.dataservice.environment;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;

import com.amazonaws.blox.dataservice.Application;
import com.amazonaws.blox.dataservice.exception.StorageException;
import com.amazonaws.blox.dataservice.handler.EnvironmentHandler;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentType;
import com.amazonaws.blox.dataservice.model.EnvironmentVersion;
import com.amazonaws.blox.dataservice.model.InstanceGroup;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentExistsException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentVersionNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ServiceException;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import java.util.UUID;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.junit4.SpringRunner;

//TODO: move integ tests to cucumber and add non happy paths
@RunWith(SpringRunner.class)
@ContextConfiguration(classes = {Application.class})
public class EnvironmentIntegTest {

  private static final String ENVIRONMENT_NAME_PREFIX = "test";

  @Autowired private EnvironmentHandler environmentHandler;

  @Test
  public void createAndDescribeEnvironment()
      throws StorageException, ServiceException, EnvironmentExistsException,
          EnvironmentNotFoundException {
    final String environmentName = ENVIRONMENT_NAME_PREFIX + UUID.randomUUID().toString();

    Environment environment = environmentObject(environmentName);
    environmentHandler.createEnvironment(environment);
    Environment describeResult =
        environmentHandler.describeEnvironment(
            environment.getEnvironmentId(), environment.getEnvironmentVersion());

    assertEquals(environment, describeResult);
  }

  @Test
  public void createAndDescribeEnvironmentTargetVersion()
      throws StorageException, ServiceException, EnvironmentExistsException,
          EnvironmentVersionNotFoundException, EnvironmentNotFoundException,
          EnvironmentNotFoundException {

    final String environmentName = ENVIRONMENT_NAME_PREFIX + UUID.randomUUID().toString();
    Environment environment = environmentObject(environmentName);
    environmentHandler.createEnvironment(environment);
    environmentHandler.createEnvironmentTargetVersion(
        environment.getEnvironmentId(), environment.getEnvironmentVersion());
    EnvironmentVersion describeResult =
        environmentHandler.describeEnvironmentTargetVersion(environment.getEnvironmentId());

    assertEquals(environment.getEnvironmentId(), describeResult.getEnvironmentId());
    assertEquals(environment.getEnvironmentVersion(), describeResult.getEnvironmentVersion());
    assertEquals(environment.getInstanceGroup().getCluster(), describeResult.getCluster());
  }

  @Test
  public void listAllClusters()
      throws ServiceException, EnvironmentExistsException, EnvironmentVersionNotFoundException,
          EnvironmentNotFoundException {
    final String environmentName = ENVIRONMENT_NAME_PREFIX + UUID.randomUUID().toString();
    Environment environment = environmentObject(environmentName);
    environmentHandler.createEnvironment(environment);
    environmentHandler.createEnvironmentTargetVersion(
        environment.getEnvironmentId(), environment.getEnvironmentVersion());

    // same cluster
    final Environment environment1 = environment;
    environment1.setEnvironmentId(UUID.randomUUID().toString());
    environmentHandler.createEnvironment(environment1);
    environmentHandler.createEnvironmentTargetVersion(
        environment1.getEnvironmentId(), environment1.getEnvironmentVersion());

    // different cluster
    final Environment environment2 = environment;
    environment2.setEnvironmentId(UUID.randomUUID().toString());
    environment2.setInstanceGroup(
        InstanceGroup.builder().cluster("cluster" + UUID.randomUUID()).build());
    environmentHandler.createEnvironment(environment2);
    environmentHandler.createEnvironmentTargetVersion(
        environment2.getEnvironmentId(), environment2.getEnvironmentVersion());

    List<String> clusters = environmentHandler.listClusters();
    assertTrue(clusters.size() >= 2);
    Set<String> clusterSet = new HashSet<>(clusters);
    assertTrue(clusterSet.contains(environment1.getInstanceGroup().getCluster()));
    assertTrue(clusterSet.contains(environment2.getInstanceGroup().getCluster()));
  }

  @Test
  public void listEnvironmentsWithCluster()
      throws ServiceException, EnvironmentExistsException, EnvironmentNotFoundException,
          EnvironmentVersionNotFoundException {
    List<String> environmentIds = environmentHandler.listEnvironmentsWithCluster("nonexistent");
    assertTrue(environmentIds.isEmpty());

    // one environment
    Environment environment =
        environmentObject(ENVIRONMENT_NAME_PREFIX + UUID.randomUUID().toString());
    final String clusterWithOneEnvironment = "cluster" + UUID.randomUUID().toString();
    environment.getInstanceGroup().setCluster(clusterWithOneEnvironment);
    environmentHandler.createEnvironment(environment);
    environmentHandler.createEnvironmentTargetVersion(
        environment.getEnvironmentId(), environment.getEnvironmentVersion());

    environmentIds = environmentHandler.listEnvironmentsWithCluster(clusterWithOneEnvironment);
    assertTrue(environmentIds.size() == 1);
    assertEquals(environment.getEnvironmentId(), environmentIds.get(0));

    // two environments on same cluster
    final String clusterWithTwoEnvironments = "cluster" + UUID.randomUUID().toString();
    Environment environmentOne = environment;
    environmentOne.setEnvironmentId(UUID.randomUUID().toString());
    environmentOne.getInstanceGroup().setCluster(clusterWithTwoEnvironments);
    environmentHandler.createEnvironment(environmentOne);
    environmentHandler.createEnvironmentTargetVersion(
        environmentOne.getEnvironmentId(), environmentOne.getEnvironmentVersion());

    Environment environmentTwo = environment;
    environmentTwo.setEnvironmentId(UUID.randomUUID().toString());
    environmentTwo.getInstanceGroup().setCluster(clusterWithTwoEnvironments);
    environmentHandler.createEnvironment(environmentTwo);
    environmentHandler.createEnvironmentTargetVersion(
        environmentTwo.getEnvironmentId(), environmentTwo.getEnvironmentVersion());

    environmentIds = environmentHandler.listEnvironmentsWithCluster(clusterWithTwoEnvironments);
    assertTrue(environmentIds.size() == 2);
    Set<String> environmentIdSet = new HashSet<>(environmentIds);
    assertTrue(environmentIdSet.contains(environmentOne.getEnvironmentId()));
    assertTrue(environmentIdSet.contains(environmentTwo.getEnvironmentId()));
  }

  private Environment environmentObject(final String environmentName) {
    return Environment.builder()
        .environmentName(environmentName)
        .environmentId(environmentName + "1234")
        .environmentVersion(UUID.randomUUID().toString())
        .type(EnvironmentType.Daemon)
        .instanceGroup(
            InstanceGroup.builder().cluster("cluster" + UUID.randomUUID().toString()).build())
        .role("role1")
        .taskDefinition("task1")
        .build();
  }
}
