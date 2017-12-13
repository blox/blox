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

public interface DataService {

  /** Creates an environment record and a new revision record. */
  CreateEnvironmentResponse createEnvironment(CreateEnvironmentRequest request)
      throws ResourceExistsException, InvalidParameterException, InternalServiceException;

  /** Creates a new environment revision for the environment. */
  UpdateEnvironmentResponse updateEnvironment(UpdateEnvironmentRequest request)
      throws ResourceNotFoundException, InvalidParameterException, InternalServiceException;

  /** Returns the environment record. */
  DescribeEnvironmentResponse describeEnvironment(DescribeEnvironmentRequest request)
      throws ResourceNotFoundException, InvalidParameterException, InternalServiceException;

  /** Lists all environments with a filter. */
  ListEnvironmentsResponse listEnvironments(ListEnvironmentsRequest request)
      throws InvalidParameterException, InternalServiceException;

  /** Deletes the provided environment if inactive or if forceDelete is true */
  DeleteEnvironmentResponse deleteEnvironment(DeleteEnvironmentRequest request)
      throws ResourceNotFoundException, ResourceInUseException, InvalidParameterException,
          InternalServiceException;

  /** Returns the requested environment revision. */
  DescribeEnvironmentRevisionResponse describeEnvironmentRevision(
      DescribeEnvironmentRevisionRequest request)
      throws ResourceNotFoundException, InvalidParameterException, InternalServiceException;

  /** Lists all environment revisions for the environment. */
  ListEnvironmentRevisionsResponse listEnvironmentRevisions(ListEnvironmentRevisionsRequest request)
      throws ResourceNotFoundException, InvalidParameterException, InternalServiceException;

  /** Lists all clusters the have environments running on them. */
  ListClustersResponse listClusters(ListClustersRequest request)
      throws InvalidParameterException, InternalServiceException;

  /** Creates a deployment record which starts a deployment. */
  StartDeploymentResponse startDeployment(StartDeploymentRequest request)
      throws ResourceNotFoundException, InvalidParameterException, InternalServiceException;
}
