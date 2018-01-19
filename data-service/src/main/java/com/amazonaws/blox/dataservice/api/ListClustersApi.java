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

import com.amazonaws.blox.dataservice.mapper.ApiModelMapper;
import com.amazonaws.blox.dataservice.model.Cluster;
import com.amazonaws.blox.dataservice.repository.EnvironmentRepository;
import com.amazonaws.blox.dataservicemodel.v1.exception.InternalServiceException;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListClustersRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListClustersResponse;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;
import lombok.AllArgsConstructor;
import lombok.NonNull;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@AllArgsConstructor
public class ListClustersApi {

  @NonNull private final ApiModelMapper apiModelMapper;
  @NonNull private final EnvironmentRepository environmentRepository;

  public ListClustersResponse listClusters(@NonNull final ListClustersRequest request)
      throws InternalServiceException {
    try {
      List<Cluster> clusters =
          environmentRepository.listClusters(
              request.getAccountId(), request.getClusterNamePrefix());

      return ListClustersResponse.builder()
          .clusters(clusters.stream().map(apiModelMapper::toCluster).collect(Collectors.toList()))
          .build();
    } catch (final Exception e) {
      log.error(e.getMessage(), e);
      throw new InternalServiceException(e.getMessage(), e);
    }
  }
}
