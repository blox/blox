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

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.hasItems;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;

import com.amazonaws.blox.dataservice.Application;
import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentExistsException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentVersionNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.exception.InvalidParameterException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ServiceException;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentType;
import com.amazonaws.blox.dataservicemodel.v1.model.InstanceGroup;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateTargetEnvironmentRevisionRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateTargetEnvironmentRevisionResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeTargetEnvironmentRevisionRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeTargetEnvironmentRevisionResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListClustersRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListClustersResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsRequest;
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
  private static final String ACCOUNT_ID = "12345678912";
  private static final String TASK_DEFINITION_ARN =
      "arn:aws:ecs:us-east-1:" + ACCOUNT_ID + ":task-definition/sleep";
  private static final String ROLE_ARN = "arn:aws:iam::" + ACCOUNT_ID + ":role/testRole";
  private static final String CLUSTER_ARN_PREFIX =
      "arn:aws:ecs:us-east-1:" + ACCOUNT_ID + ":cluster/";

  @Autowired private DataService dataService;

  @Test
  public void createAndDescribeEnvironment()
      throws InvalidParameterException, ServiceException, EnvironmentExistsException,
          EnvironmentNotFoundException {

    CreateEnvironmentResponse createEnvironmentResponse = createEnvironment();

    DescribeEnvironmentRequest describeEnvironmentRequest =
        DescribeEnvironmentRequest.builder()
            .environmentId(createEnvironmentResponse.getEnvironmentId())
            .environmentVersion(createEnvironmentResponse.getEnvironmentVersion())
            .build();

    DescribeEnvironmentResponse describeEnvironmentResponse =
        dataService.describeEnvironment(describeEnvironmentRequest);

    assertEquals(
        createEnvironmentResponse.getEnvironmentId(),
        describeEnvironmentResponse.getEnvironmentId());
    assertEquals(
        createEnvironmentResponse.getEnvironmentName(),
        describeEnvironmentResponse.getEnvironmentName());
    assertEquals(
        createEnvironmentResponse.getEnvironmentVersion(),
        describeEnvironmentResponse.getEnvironmentVersion());
    assertEquals(
        createEnvironmentResponse.getInstanceGroup().getCluster(),
        describeEnvironmentResponse.getInstanceGroup().getCluster());
    assertEquals(createEnvironmentResponse.getRole(), describeEnvironmentResponse.getRole());
    assertEquals(
        createEnvironmentResponse.getTaskDefinition(),
        describeEnvironmentResponse.getTaskDefinition());
  }

  @Test
  public void createAndDescribeEnvironmentTargetVersion()
      throws InvalidParameterException, ServiceException, EnvironmentExistsException,
          EnvironmentVersionNotFoundException, EnvironmentNotFoundException {

    CreateEnvironmentResponse createEnvironmentResponse = createEnvironment();
    CreateTargetEnvironmentRevisionResponse createTargetEnvironmentRevisionResponse =
        dataService.createTargetEnvironmentRevision(
            CreateTargetEnvironmentRevisionRequest.builder()
                .environmentId(createEnvironmentResponse.getEnvironmentId())
                .environmentVersion(createEnvironmentResponse.getEnvironmentVersion())
                .build());

    DescribeTargetEnvironmentRevisionResponse describeTargetEnvironmentRevisionResponse =
        dataService.describeTargetEnvironmentRevision(
            DescribeTargetEnvironmentRevisionRequest.builder()
                .environmentId(createTargetEnvironmentRevisionResponse.getEnvironmentId())
                .build());

    assertEquals(
        createTargetEnvironmentRevisionResponse.getEnvironmentId(),
        describeTargetEnvironmentRevisionResponse.getEnvironmentId());
    assertEquals(
        createTargetEnvironmentRevisionResponse.getEnvironmentVersion(),
        describeTargetEnvironmentRevisionResponse.getEnvironmentVersion());
    assertEquals(
        createTargetEnvironmentRevisionResponse.getCluster(),
        describeTargetEnvironmentRevisionResponse.getCluster());
  }

  @Test
  public void listAllClusters()
      throws ServiceException, EnvironmentExistsException, EnvironmentVersionNotFoundException,
          EnvironmentNotFoundException, InvalidParameterException {

    CreateEnvironmentResponse createEnvironmentResponse = createEnvironment();
    dataService.createTargetEnvironmentRevision(
        CreateTargetEnvironmentRevisionRequest.builder()
            .environmentId(createEnvironmentResponse.getEnvironmentId())
            .environmentVersion(createEnvironmentResponse.getEnvironmentVersion())
            .build());

    // same cluster
    CreateEnvironmentResponse sameClusterCreateEnvironmentResponse =
        createEnvironment(createEnvironmentResponse.getInstanceGroup().getCluster());
    dataService.createTargetEnvironmentRevision(
        CreateTargetEnvironmentRevisionRequest.builder()
            .environmentId(sameClusterCreateEnvironmentResponse.getEnvironmentId())
            .environmentVersion(sameClusterCreateEnvironmentResponse.getEnvironmentVersion())
            .build());

    // different cluster
    CreateEnvironmentResponse differentClusterCreateEnvironmentResponse = createEnvironment();
    dataService.createTargetEnvironmentRevision(
        CreateTargetEnvironmentRevisionRequest.builder()
            .environmentId(differentClusterCreateEnvironmentResponse.getEnvironmentId())
            .environmentVersion(differentClusterCreateEnvironmentResponse.getEnvironmentVersion())
            .build());

    ListClustersResponse listClustersResponse =
        dataService.listClusters(ListClustersRequest.builder().build());

    assertTrue(listClustersResponse.getClusters().size() >= 2);

    assertThat(
        listClustersResponse.getClusters(),
        hasItems(
            sameClusterCreateEnvironmentResponse.getInstanceGroup().getCluster(),
            differentClusterCreateEnvironmentResponse.getInstanceGroup().getCluster()));
  }

  @Test
  public void listEnvironmentsWithCluster()
      throws ServiceException, EnvironmentExistsException, EnvironmentNotFoundException,
          EnvironmentVersionNotFoundException, InvalidParameterException {

    List<String> environmentIds =
        dataService
            .listEnvironments(
                ListEnvironmentsRequest.builder()
                    .cluster(CLUSTER_ARN_PREFIX + "nonexistent")
                    .build())
            .getEnvironmentIds();
    assertTrue(environmentIds.isEmpty());

    // one environment
    CreateEnvironmentResponse createEnvironmentResponse = createEnvironment();
    CreateTargetEnvironmentRevisionResponse createTargetEnvironmentRevisionResponse =
        dataService.createTargetEnvironmentRevision(
            CreateTargetEnvironmentRevisionRequest.builder()
                .environmentId(createEnvironmentResponse.getEnvironmentId())
                .environmentVersion(createEnvironmentResponse.getEnvironmentVersion())
                .build());

    environmentIds =
        dataService
            .listEnvironments(
                ListEnvironmentsRequest.builder()
                    .cluster(createTargetEnvironmentRevisionResponse.getCluster())
                    .build())
            .getEnvironmentIds();
    assertTrue(environmentIds.size() == 1);
    assertEquals(createTargetEnvironmentRevisionResponse.getEnvironmentId(), environmentIds.get(0));

    // two environments on the same cluster
    final String clusterWithTwoEnvironments = CLUSTER_ARN_PREFIX + UUID.randomUUID().toString();

    CreateEnvironmentResponse createEnvironmentOnTheSameClusterOneResponse =
        createEnvironment(clusterWithTwoEnvironments);
    dataService.createTargetEnvironmentRevision(
        CreateTargetEnvironmentRevisionRequest.builder()
            .environmentId(createEnvironmentOnTheSameClusterOneResponse.getEnvironmentId())
            .environmentVersion(
                createEnvironmentOnTheSameClusterOneResponse.getEnvironmentVersion())
            .build());

    CreateEnvironmentResponse createEnvironmentOnTheSameClusterTwoResponse =
        createEnvironment(clusterWithTwoEnvironments);
    dataService.createTargetEnvironmentRevision(
        CreateTargetEnvironmentRevisionRequest.builder()
            .environmentId(createEnvironmentOnTheSameClusterTwoResponse.getEnvironmentId())
            .environmentVersion(
                createEnvironmentOnTheSameClusterTwoResponse.getEnvironmentVersion())
            .build());

    environmentIds =
        dataService
            .listEnvironments(
                ListEnvironmentsRequest.builder().cluster(clusterWithTwoEnvironments).build())
            .getEnvironmentIds();
    assertTrue(environmentIds.size() == 2);
    Set<String> environmentIdSet = new HashSet<>(environmentIds);
    assertTrue(
        environmentIdSet.contains(createEnvironmentOnTheSameClusterOneResponse.getEnvironmentId()));
    assertTrue(
        environmentIdSet.contains(createEnvironmentOnTheSameClusterTwoResponse.getEnvironmentId()));
  }

  private CreateEnvironmentRequest createEnvironmentRequest(
      final String environmentName, final String cluster) {
    return CreateEnvironmentRequest.builder()
        .environmentName(environmentName)
        .accountId(ACCOUNT_ID)
        .environmentType(EnvironmentType.Daemon)
        .instanceGroup(InstanceGroup.builder().cluster(cluster).build())
        .role(ROLE_ARN)
        .taskDefinition(TASK_DEFINITION_ARN)
        .build();
  }

  private CreateEnvironmentResponse createEnvironment()
      throws EnvironmentExistsException, InvalidParameterException, ServiceException {
    return createEnvironment(CLUSTER_ARN_PREFIX + UUID.randomUUID().toString());
  }

  private CreateEnvironmentResponse createEnvironment(final String cluster)
      throws EnvironmentExistsException, InvalidParameterException, ServiceException {
    final String environmentName = ENVIRONMENT_NAME_PREFIX + UUID.randomUUID().toString();

    CreateEnvironmentRequest createEnvironmentRequest =
        createEnvironmentRequest(environmentName, cluster);
    return dataService.createEnvironment(createEnvironmentRequest);
  }
}
