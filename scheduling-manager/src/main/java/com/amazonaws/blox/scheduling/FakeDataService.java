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

import com.amazonaws.blox.dataservicemodel.v1.old.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.old.exception.EnvironmentActiveException;
import com.amazonaws.blox.dataservicemodel.v1.old.exception.EnvironmentExistsException;
import com.amazonaws.blox.dataservicemodel.v1.old.exception.EnvironmentNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.old.exception.EnvironmentTargetRevisionNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.old.exception.EnvironmentTargetRevisionExistsException;
import com.amazonaws.blox.dataservicemodel.v1.old.exception.InvalidParameterException;
import com.amazonaws.blox.dataservicemodel.v1.old.exception.ServiceException;
import com.amazonaws.blox.dataservicemodel.v1.old.model.DeploymentConfiguration;
import com.amazonaws.blox.dataservicemodel.v1.old.model.EnvironmentHealth;
import com.amazonaws.blox.dataservicemodel.v1.old.model.EnvironmentStatus;
import com.amazonaws.blox.dataservicemodel.v1.old.model.EnvironmentType;
import com.amazonaws.blox.dataservicemodel.v1.old.model.InstanceGroup;
import com.amazonaws.blox.dataservicemodel.v1.old.model.wrappers.CreateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.old.model.wrappers.CreateEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.old.model.wrappers.CreateTargetEnvironmentRevisionRequest;
import com.amazonaws.blox.dataservicemodel.v1.old.model.wrappers.CreateTargetEnvironmentRevisionResponse;
import com.amazonaws.blox.dataservicemodel.v1.old.model.wrappers.DeleteEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.old.model.wrappers.DeleteEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.old.model.wrappers.DescribeEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.old.model.wrappers.DescribeEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.old.model.wrappers.DescribeTargetEnvironmentRevisionRequest;
import com.amazonaws.blox.dataservicemodel.v1.old.model.wrappers.DescribeTargetEnvironmentRevisionResponse;
import com.amazonaws.blox.dataservicemodel.v1.old.model.wrappers.ListClustersRequest;
import com.amazonaws.blox.dataservicemodel.v1.old.model.wrappers.ListClustersResponse;
import com.amazonaws.blox.dataservicemodel.v1.old.model.wrappers.ListEnvironmentsRequest;
import com.amazonaws.blox.dataservicemodel.v1.old.model.wrappers.ListEnvironmentsResponse;
import com.amazonaws.blox.dataservicemodel.v1.old.model.wrappers.StartDeploymentRequest;
import com.amazonaws.blox.dataservicemodel.v1.old.model.wrappers.StartDeploymentResponse;
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
      throws EnvironmentTargetRevisionExistsException, EnvironmentNotFoundException,
          InvalidParameterException, ServiceException {
    return null;
  }

  @Override
  public StartDeploymentResponse startDeployment(StartDeploymentRequest request)
      throws EnvironmentNotFoundException, EnvironmentTargetRevisionNotFoundException,
          EnvironmentTargetRevisionExistsException, InvalidParameterException, ServiceException {
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
  public DeleteEnvironmentResponse deleteEnvironment(DeleteEnvironmentRequest request)
      throws EnvironmentNotFoundException, EnvironmentActiveException, InvalidParameterException,
          ServiceException {
    throw new UnsupportedOperationException();
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
