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

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.exception.InternalServiceException;
import com.amazonaws.blox.dataservicemodel.v1.exception.InvalidParameterException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceExistsException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceInUseException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceNotFoundException;
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
import lombok.AllArgsConstructor;
import lombok.NonNull;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@AllArgsConstructor
public class DataServiceApi implements DataService {

  @NonNull private final CreateEnvironmentApi createEnvironmentApi;
  @NonNull private final DescribeEnvironmentApi describeEnvironmentApi;
  @NonNull private final StartDeploymentApi startDeploymentApi;
  @NonNull private final ListEnvironmentsApi listEnvironmentsApi;
  @NonNull private final DescribeEnvironmentRevisionApi describeEnvironmentRevisionApi;

  @Override
  public CreateEnvironmentResponse createEnvironment(
      @NonNull final CreateEnvironmentRequest request)
      throws ResourceExistsException, InvalidParameterException, InternalServiceException {

    return createEnvironmentApi.createEnvironment(request);
  }

  @Override
  public UpdateEnvironmentResponse updateEnvironment(UpdateEnvironmentRequest request)
      throws ResourceNotFoundException, InvalidParameterException, InternalServiceException {
    return null;
  }

  @Override
  public DescribeEnvironmentResponse describeEnvironment(
      @NonNull final DescribeEnvironmentRequest request)
      throws ResourceNotFoundException, InvalidParameterException, InternalServiceException {
    return describeEnvironmentApi.describeEnvironment(request);
  }

  @Override
  public ListEnvironmentsResponse listEnvironments(ListEnvironmentsRequest request)
      throws InvalidParameterException, InternalServiceException {
    return listEnvironmentsApi.listEnvironments(request);
  }

  @Override
  public DeleteEnvironmentResponse deleteEnvironment(DeleteEnvironmentRequest request)
      throws ResourceNotFoundException, ResourceInUseException, InvalidParameterException,
          InternalServiceException {
    return null;
  }

  @Override
  public DescribeEnvironmentRevisionResponse describeEnvironmentRevision(
      DescribeEnvironmentRevisionRequest request)
      throws ResourceNotFoundException, InvalidParameterException, InternalServiceException {
    return describeEnvironmentRevisionApi.describeEnvironmentRevision(request);
  }

  @Override
  public ListEnvironmentRevisionsResponse listEnvironmentRevisions(
      ListEnvironmentRevisionsRequest request)
      throws ResourceNotFoundException, InvalidParameterException, InternalServiceException {
    return null;
  }

  @Override
  public ListClustersResponse listClusters(ListClustersRequest request)
      throws InvalidParameterException, InternalServiceException {
    return null;
  }

  @Override
  public StartDeploymentResponse startDeployment(StartDeploymentRequest request)
      throws ResourceNotFoundException, InvalidParameterException, InternalServiceException {
    return startDeploymentApi.startDeployment(request);
  }
}
