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
package com.amazonaws.blox.dataservicemodel.v1.old.client;

import com.amazonaws.blox.dataservicemodel.v1.old.exception.EnvironmentActiveException;
import com.amazonaws.blox.dataservicemodel.v1.old.exception.EnvironmentExistsException;
import com.amazonaws.blox.dataservicemodel.v1.old.exception.EnvironmentNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.old.exception.EnvironmentTargetRevisionNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.old.exception.EnvironmentTargetRevisionExistsException;
import com.amazonaws.blox.dataservicemodel.v1.old.exception.InvalidParameterException;
import com.amazonaws.blox.dataservicemodel.v1.old.exception.ServiceException;
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

public interface DataService {

  /** Creates an environment record. */
  CreateEnvironmentResponse createEnvironment(CreateEnvironmentRequest request)
      throws EnvironmentExistsException, InvalidParameterException, ServiceException;

  /** Creates an environment target revision record. */
  CreateTargetEnvironmentRevisionResponse createTargetEnvironmentRevision(
      CreateTargetEnvironmentRevisionRequest request)
      throws EnvironmentTargetRevisionExistsException, EnvironmentNotFoundException,
          InvalidParameterException, ServiceException;

  /** Returns the environment record. */
  DescribeEnvironmentResponse describeEnvironment(DescribeEnvironmentRequest request)
      throws EnvironmentNotFoundException, InvalidParameterException, ServiceException;

  /** Returns the environment target revision record. */
  DescribeTargetEnvironmentRevisionResponse describeTargetEnvironmentRevision(
      DescribeTargetEnvironmentRevisionRequest request)
      throws EnvironmentTargetRevisionNotFoundException, InvalidParameterException,
          ServiceException;

  /** Lists all environments with a filter. */
  ListEnvironmentsResponse listEnvironments(ListEnvironmentsRequest request)
      throws InvalidParameterException, ServiceException;

  /** Lists all clusters the have environments running on them. */
  ListClustersResponse listClusters(ListClustersRequest request)
      throws InvalidParameterException, ServiceException;

  /** Deletes the provided environment if inactive or if forceDelete is true */
  DeleteEnvironmentResponse deleteEnvironment(DeleteEnvironmentRequest request)
      throws EnvironmentNotFoundException, EnvironmentActiveException, InvalidParameterException,
          ServiceException;

  /** Creates a deployment record which asynchronously starts a deployment. */
  StartDeploymentResponse startDeployment(StartDeploymentRequest request)
      throws EnvironmentNotFoundException, EnvironmentTargetRevisionNotFoundException,
          EnvironmentTargetRevisionExistsException, InvalidParameterException, ServiceException;
}
