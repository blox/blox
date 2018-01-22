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
package com.amazonaws.blox.scheduling;

import java.time.Instant;
import java.util.Collections;

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.exception.InternalServiceException;
import com.amazonaws.blox.dataservicemodel.v1.exception.InvalidParameterException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceExistsException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceInUseException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.model.Cluster;
import com.amazonaws.blox.dataservicemodel.v1.model.DeploymentConfiguration;
import com.amazonaws.blox.dataservicemodel.v1.model.Environment;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentHealth;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentRevision;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentStatus;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentType;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DeleteEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DeleteEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRevisionRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRevisionResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListClustersRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListClustersResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentRevisionsRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentRevisionsResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.StartDeploymentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.StartDeploymentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.UpdateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.UpdateEnvironmentResponse;
import lombok.Builder;

/**
 * Temporary fake data service for testing without DataService in Lambda
 *
 * <p>TODO: Move to end to end/load tests, or replace with fixture data
 */
@Builder
public class FakeDataService implements DataService {
  @Builder.Default private String accountId = "123456789012";
  @Builder.Default private String clusterName = "default";
  @Builder.Default private String taskDefinition = "sleep:1";
  @Builder.Default private String environmentName = "TestEnvironment";
  @Builder.Default private String deploymentMethod = "ReplaceAfterTerminate";
  @Builder.Default private String activeEnvironmentRevisionId = "1";

  @Override
  public CreateEnvironmentResponse createEnvironment(CreateEnvironmentRequest request)
      throws ResourceExistsException, InvalidParameterException, InternalServiceException {
    return null;
  }

  @Override
  public UpdateEnvironmentResponse updateEnvironment(UpdateEnvironmentRequest request)
      throws ResourceNotFoundException, InvalidParameterException, InternalServiceException {
    return null;
  }

  @Override
  public DescribeEnvironmentResponse describeEnvironment(DescribeEnvironmentRequest request)
      throws ResourceNotFoundException, InvalidParameterException, InternalServiceException {
    return DescribeEnvironmentResponse.builder()
        .environment(
            Environment.builder()
                .environmentId(request.getEnvironmentId())
                .role("")
                .environmentType(EnvironmentType.SingleTask)
                .createdTime(Instant.now())
                .lastUpdatedTime(Instant.now())
                .environmentHealth(EnvironmentHealth.HEALTHY)
                .environmentStatus(EnvironmentStatus.ACTIVE)
                .deploymentMethod(deploymentMethod)
                .deploymentConfiguration(DeploymentConfiguration.builder().build())
                .activeEnvironmentRevisionId(activeEnvironmentRevisionId)
                .build())
        .build();
  }

  @Override
  public ListEnvironmentsResponse listEnvironments(ListEnvironmentsRequest request) {
    return ListEnvironmentsResponse.builder()
        .environmentIds(
            Collections.singletonList(
                EnvironmentId.builder()
                    .accountId(accountId)
                    .cluster(clusterName)
                    .environmentName(environmentName)
                    .build()))
        .build();
  }

  @Override
  public DeleteEnvironmentResponse deleteEnvironment(DeleteEnvironmentRequest request)
      throws ResourceNotFoundException, ResourceInUseException, InvalidParameterException,
          InternalServiceException {
    throw new UnsupportedOperationException();
  }

  @Override
  public DescribeEnvironmentRevisionResponse describeEnvironmentRevision(
      DescribeEnvironmentRevisionRequest request)
      throws ResourceNotFoundException, InvalidParameterException, InternalServiceException {
    return DescribeEnvironmentRevisionResponse.builder()
        .environmentRevision(
            EnvironmentRevision.builder()
                .environmentId(request.getEnvironmentId())
                .environmentRevisionId("1")
                .taskDefinition(taskDefinition)
                .createdTime(Instant.now())
                .build())
        .build();
  }

  @Override
  public ListEnvironmentRevisionsResponse listEnvironmentRevisions(
      ListEnvironmentRevisionsRequest request)
      throws ResourceNotFoundException, InvalidParameterException, InternalServiceException {
    return null;
  }

  @Override
  public ListClustersResponse listClusters(ListClustersRequest request) {
    return ListClustersResponse.builder()
        .clusters(
            Collections.singletonList(
                Cluster.builder().accountId(accountId).clusterName(clusterName).build()))
        .build();
  }

  @Override
  public StartDeploymentResponse startDeployment(StartDeploymentRequest request)
      throws ResourceNotFoundException, InvalidParameterException, InternalServiceException {
    return null;
  }
}
