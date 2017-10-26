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

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentExistsException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentVersionNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentVersionOutdatedException;
import com.amazonaws.blox.dataservicemodel.v1.exception.InvalidParameterException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ServiceException;
import com.amazonaws.blox.dataservicemodel.v1.model.DeploymentConfiguration;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentHealth;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentStatus;
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
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.StartDeploymentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.StartDeploymentResponse;
import java.time.Instant;
import java.util.Collections;
import lombok.Builder;

/**
 * Temporary fake data service for testing without DataService in Lambda
 *
 * <p>TODO: Move to end to end/load tests, or replace with fixture data
 */
@Builder
public class FakeDataService implements DataService {
  @Builder.Default private String clusterArn = "default";
  @Builder.Default private String taskDefinition = "sleep:1";
  @Builder.Default private String environmentId = "TestEnvironment";

  @Override
  public CreateEnvironmentResponse createEnvironment(CreateEnvironmentRequest request)
      throws EnvironmentExistsException, InvalidParameterException, ServiceException {
    return null;
  }

  @Override
  public CreateTargetEnvironmentRevisionResponse createTargetEnvironmentRevision(
      CreateTargetEnvironmentRevisionRequest request)
      throws EnvironmentExistsException, EnvironmentNotFoundException, InvalidParameterException,
          ServiceException {
    return null;
  }

  @Override
  public StartDeploymentResponse startDeployment(StartDeploymentRequest request)
      throws EnvironmentNotFoundException, EnvironmentVersionNotFoundException,
          EnvironmentVersionOutdatedException, InvalidParameterException, ServiceException {
    return null;
  }

  @Override
  public ListEnvironmentsResponse listEnvironments(ListEnvironmentsRequest request) {
    return new ListEnvironmentsResponse(Collections.singletonList(environmentId));
  }

  @Override
  public ListClustersResponse listClusters(ListClustersRequest request) {
    return new ListClustersResponse(Collections.singletonList(clusterArn));
  }

  @Override
  public DescribeTargetEnvironmentRevisionResponse describeTargetEnvironmentRevision(
      DescribeTargetEnvironmentRevisionRequest request) {
    return DescribeTargetEnvironmentRevisionResponse.builder()
        .environmentId(request.getEnvironmentId())
        .environmentVersion("1")
        .cluster(clusterArn)
        .build();
  }

  @Override
  public DescribeEnvironmentResponse describeEnvironment(DescribeEnvironmentRequest request)
      throws EnvironmentNotFoundException, InvalidParameterException, ServiceException {
    return DescribeEnvironmentResponse.builder()
        .environmentId(request.getEnvironmentId())
        .environmentVersion("1")
        .environmentName(request.getEnvironmentId())
        .taskDefinition(taskDefinition)
        .role("")
        .status(EnvironmentStatus.ACTIVE)
        .health(EnvironmentHealth.HEALTHY)
        .createdTime(Instant.now())
        .instanceGroup(InstanceGroup.builder().cluster(clusterArn).build())
        .type(EnvironmentType.SingleTask)
        .deploymentConfiguration(DeploymentConfiguration.builder().build())
        .build();
  }
}
