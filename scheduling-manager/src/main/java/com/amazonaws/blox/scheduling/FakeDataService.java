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
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.DeploymentConfiguration;
import com.amazonaws.blox.dataservicemodel.v1.model.DeploymentType;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeTargetEnvironmentRevisionRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeTargetEnvironmentRevisionResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentVersion;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListClustersRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListClustersResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.StartDeploymentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.StartDeploymentResponse;
import java.util.Arrays;
import java.util.List;
import java.util.stream.Collectors;
import lombok.Builder;

@Builder
/**
 * Temporary fake data service for load-testing in Lambda
 *
 * <p>TODO: Move to end to end/load tests, or replace with fixture data
 */
public class FakeDataService implements DataService {
  public static final List<String> GREEK =
      Arrays.asList(
          "Alpha",
          "Beta",
          "Gamma",
          "Delta",
          "Epsilon",
          "Zeta",
          "Eta",
          "Theta",
          "Iota",
          "Kappa",
          "Lambda",
          "Mu",
          "Nu",
          "Xi",
          "Omicron",
          "Pi",
          "Rho",
          "Sigma",
          "Tau",
          "Upsilon",
          "Phi",
          "Chi",
          "Psi",
          "Omega",
          "Ultra" /* not greek, but need 25 things */);
  public static final List<String> THING = Arrays.asList("Blaster", "Boomer", "Slapper", "Zapper");
  @Builder.Default private int clusters = 100;
  @Builder.Default private int environmentsPerCluster = 10;

  @Override
  public CreateEnvironmentResponse createEnvironment(CreateEnvironmentRequest request)
      throws EnvironmentExistsException, InvalidParameterException, ServiceException {
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
    return new ListEnvironmentsResponse(
        names("EnvironmentFor" + request.getClusterArn())
            .stream()
            .limit(environmentsPerCluster)
            .collect(Collectors.toList()));
  }

  @Override
  public ListClustersResponse listClusters(ListClustersRequest request) {
    return new ListClustersResponse(
        names("Cluster").stream().limit(clusters).collect(Collectors.toList()));
  }

  @Override
  public DescribeTargetEnvironmentRevisionResponse describeTargetEnvironmentRevision(
      DescribeTargetEnvironmentRevisionRequest request) {
    return new DescribeTargetEnvironmentRevisionResponse(
        EnvironmentVersion.builder()
            .deploymentType(DeploymentType.SingleTask)
            .deploymentConfiguration(DeploymentConfiguration.builder().build())
            .build());
  }

  private List<String> names(String suffix) {
    return GREEK
        .stream()
        .flatMap(g -> THING.stream().map(t -> g + t + suffix))
        .collect(Collectors.toList());
  }
}
