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

import com.amazonaws.blox.dataservice.arn.EnvironmentArnGenerator;
import com.amazonaws.blox.dataservice.handler.EnvironmentHandler;
import com.amazonaws.blox.dataservice.mapper.ApiModelMapper;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentVersion;
import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentExistsException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentVersionNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentVersionOutdatedException;
import com.amazonaws.blox.dataservicemodel.v1.exception.InvalidParameterException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ServiceException;
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
import java.util.List;
import lombok.AllArgsConstructor;
import lombok.NonNull;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@AllArgsConstructor
public class DataServiceApi implements DataService {

  @NonNull private final EnvironmentHandler environmentHandler;
  @NonNull private final ApiModelMapper apiModelMapper;

  @Override
  public CreateEnvironmentResponse createEnvironment(
      @NonNull final CreateEnvironmentRequest request)
      throws EnvironmentExistsException, InvalidParameterException, ServiceException {

    //TODO: validate parameters and throw InvalidParameterException

    try {
      final Environment environmentRequest = apiModelMapper.toEnvironment(request);
      //TODO:move to the mapper
      environmentRequest.setEnvironmentId(
          EnvironmentArnGenerator.generateEnvironmentArn(
              request.getEnvironmentName(), request.getAccountId()));

      final Environment environment = environmentHandler.createEnvironment(environmentRequest);
      return apiModelMapper.toCreateEnvironmentResponse(environment);

    } catch (final EnvironmentExistsException | ServiceException e) {
      log.error(e.getMessage(), e);
      throw e;
    }
  }

  @Override
  public CreateTargetEnvironmentRevisionResponse createTargetEnvironmentRevision(
      @NonNull final CreateTargetEnvironmentRevisionRequest request)
      throws EnvironmentExistsException, EnvironmentNotFoundException, InvalidParameterException,
          ServiceException {

    //TODO: validate parameters and throw InvalidParameterException

    try {
      final EnvironmentVersion environmentVersion =
          environmentHandler.createEnvironmentTargetVersion(
              request.getEnvironmentId(), request.getEnvironmentVersion());

      return apiModelMapper.toCreateTargetEnvironmentRevisionResponse(environmentVersion);

    } catch (final EnvironmentNotFoundException | EnvironmentExistsException | ServiceException e) {
      log.error(e.getMessage(), e);
      throw e;
    }
  }

  @Override
  public DescribeEnvironmentResponse describeEnvironment(
      @NonNull final DescribeEnvironmentRequest request)
      throws EnvironmentNotFoundException, InvalidParameterException, ServiceException {
    //TODO: validate parameters and throw InvalidParameterException

    try {
      final Environment environment =
          environmentHandler.describeEnvironment(
              request.getEnvironmentId(), request.getEnvironmentVersion());
      return apiModelMapper.toDescribeEnvironmentResponse(environment);

    } catch (final EnvironmentNotFoundException | ServiceException e) {
      log.error(e.getMessage(), e);
      throw e;
    }
  }

  @Override
  public DescribeTargetEnvironmentRevisionResponse describeTargetEnvironmentRevision(
      DescribeTargetEnvironmentRevisionRequest request)
      throws EnvironmentVersionNotFoundException, InvalidParameterException, ServiceException {
    //TODO: validate parameters and throw InvalidParameterException

    try {
      final EnvironmentVersion environmentVersion =
          environmentHandler.describeEnvironmentTargetVersion(request.getEnvironmentId());
      return apiModelMapper.toDescribeTargetEnvironmentRevisionResponse(environmentVersion);

    } catch (final EnvironmentVersionNotFoundException | ServiceException e) {
      log.error(e.getMessage(), e);
      throw e;
    }
  }

  @Override
  public ListEnvironmentsResponse listEnvironments(ListEnvironmentsRequest request)
      throws InvalidParameterException, ServiceException {
    //TODO: validate parameters and throw InvalidParameterException

    try {
      final List<String> environmentIds =
          environmentHandler.listEnvironmentsWithCluster(request.getCluster());
      return ListEnvironmentsResponse.builder().environmentIds(environmentIds).build();

    } catch (final ServiceException e) {
      log.error(e.getMessage(), e);
      throw e;
    }
  }

  @Override
  public ListClustersResponse listClusters(ListClustersRequest request)
      throws InvalidParameterException, ServiceException {
    //TODO: validate parameters and throw InvalidParameterException

    try {
      final List<String> clusters = environmentHandler.listClusters();
      return ListClustersResponse.builder().clusters(clusters).build();

    } catch (final ServiceException e) {
      log.error(e.getMessage(), e);
      throw e;
    }
  }

  @Override
  public StartDeploymentResponse startDeployment(StartDeploymentRequest request)
      throws EnvironmentNotFoundException, EnvironmentVersionNotFoundException,
          EnvironmentVersionOutdatedException, InvalidParameterException, ServiceException {
    //TODO: validate parameters and throw InvalidParameterException
    throw new UnsupportedOperationException();
  }
}
