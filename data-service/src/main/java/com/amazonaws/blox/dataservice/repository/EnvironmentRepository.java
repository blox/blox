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
package com.amazonaws.blox.dataservice.repository;

import com.amazonaws.blox.dataservice.model.Cluster;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentId;
import com.amazonaws.blox.dataservice.model.EnvironmentRevision;
import com.amazonaws.blox.dataservicemodel.v1.exception.InternalServiceException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceExistsException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceNotFoundException;
import java.util.List;

/** Methods for interacting with environment and environment revision objects in the repository. */
public interface EnvironmentRepository {

  Environment createEnvironmentAndEnvironmentRevision(
      Environment environment, EnvironmentRevision environmentRevision)
      throws ResourceExistsException, InternalServiceException;

  Environment getEnvironment(EnvironmentId environmentId)
      throws ResourceNotFoundException, InternalServiceException;

  List<Environment> listEnvironments(Cluster cluster) throws InternalServiceException;

  List<Environment> listEnvironments(Cluster cluster, String environmentNamePrefix)
      throws InternalServiceException;

  Environment updateEnvironment(Environment environment)
      throws ResourceNotFoundException, InternalServiceException;

  EnvironmentRevision getEnvironmentRevision(
      EnvironmentId environmentId, String environmentRevisionId)
      throws ResourceNotFoundException, InternalServiceException;

  List<EnvironmentRevision> listEnvironmentRevisions(EnvironmentId environmentId)
      throws InternalServiceException;

  void deleteEnvironmentRevision(EnvironmentRevision environmentRevision)
      throws InternalServiceException;

  List<Cluster> listClusters(String accountId, String clusterNamePrefix);

  void deleteEnvironment(EnvironmentId environmentId) throws InternalServiceException;
}
