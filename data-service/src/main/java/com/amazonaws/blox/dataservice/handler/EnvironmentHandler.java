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
package com.amazonaws.blox.dataservice.handler;

import com.amazonaws.blox.dataservice.exception.StorageException;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.repository.EnvironmentRepository;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentExistsException;
import com.amazonaws.blox.dataservicemodel.v1.exception.EnvironmentNotFoundException;
import java.util.List;
import lombok.AllArgsConstructor;
import lombok.NonNull;
import org.springframework.stereotype.Component;

@Component
@AllArgsConstructor
public class EnvironmentHandler {

  @NonNull private EnvironmentRepository environmentRepository;

  public Environment createEnvironment(final String environmentName, final String accountId)
      throws EnvironmentExistsException, StorageException {
    return null;
  }

  public Environment describeEnvironment(final String environmentId) {
    return null;
  }

  public Environment describeEnvironment(
      final String environmentId, final String environmentVersion) {
    return null;
  }

  public List<String> listClustersWithEnvironments() {
    return null;
  }

  public Environment getLatestEnvironmentVersion(final String environmentId)
      throws EnvironmentNotFoundException, StorageException {
    return null;
  }
}
