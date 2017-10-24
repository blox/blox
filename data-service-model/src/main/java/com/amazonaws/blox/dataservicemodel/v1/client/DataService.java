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
package com.amazonaws.blox.dataservicemodel.v1.client;

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

public interface DataService {

  /** Creates an environment record. */
  CreateEnvironmentResponse createEnvironment(CreateEnvironmentRequest request)
      throws EnvironmentExistsException, InvalidParameterException, ServiceException;

  /** Creates an environment version record. */
  CreateTargetEnvironmentRevisionResponse createTargetEnvironmentRevision(
      CreateTargetEnvironmentRevisionRequest request)
      throws EnvironmentExistsException, EnvironmentNotFoundException, InvalidParameterException,
          ServiceException;

  /** Returns the environment record. */
  DescribeEnvironmentResponse describeEnvironment(DescribeEnvironmentRequest request)
      throws EnvironmentNotFoundException, InvalidParameterException, ServiceException;

  /** Returns the environment target version record. */
  DescribeTargetEnvironmentRevisionResponse describeTargetEnvironmentRevision(
      DescribeTargetEnvironmentRevisionRequest request)
      throws EnvironmentVersionNotFoundException, InvalidParameterException, ServiceException;

  /** Lists all environments with a filter. */
  ListEnvironmentsResponse listEnvironments(ListEnvironmentsRequest request)
      throws InvalidParameterException, ServiceException;

  /** Lists all clusters the have environments running on them. */
  ListClustersResponse listClusters(ListClustersRequest request)
      throws InvalidParameterException, ServiceException;

  /** Creates a deployment record which asynchronously starts a deployment. */
  StartDeploymentResponse startDeployment(StartDeploymentRequest request)
      throws EnvironmentNotFoundException, EnvironmentVersionNotFoundException,
          EnvironmentVersionOutdatedException, InvalidParameterException, ServiceException;
}
