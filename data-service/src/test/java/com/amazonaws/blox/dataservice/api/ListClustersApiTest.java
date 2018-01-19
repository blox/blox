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

import static java.util.Arrays.asList;
import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.tuple;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

import com.amazonaws.blox.dataservice.mapper.ApiModelMapper;
import com.amazonaws.blox.dataservice.model.Cluster;
import com.amazonaws.blox.dataservice.repository.EnvironmentRepository;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListClustersRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListClustersResponse;
import org.junit.Test;
import org.mapstruct.factory.Mappers;

public class ListClustersApiTest {
  EnvironmentRepository repository = mock(EnvironmentRepository.class);
  ApiModelMapper mapper = Mappers.getMapper(ApiModelMapper.class);

  ListClustersApi api = new ListClustersApi(mapper, repository);

  @Test
  public void itListsAllClusters() throws Exception {
    when(repository.listClusters(null, null))
        .thenReturn(asList(cluster("1", "alpha"), cluster("2", "beta"), cluster("3", "gamma")));

    ListClustersResponse response = api.listClusters(ListClustersRequest.builder().build());

    assertThat(response.getClusters())
        .extracting("accountId", "clusterName")
        .containsExactlyInAnyOrder(tuple("1", "alpha"), tuple("2", "beta"), tuple("3", "gamma"));
  }

  private Cluster cluster(final String accountId, final String name) {
    return Cluster.builder().accountId(accountId).clusterName(name).build();
  }

  @Test
  public void itFiltersByAccountIdIfGiven() throws Exception {
    when(repository.listClusters("1", null)).thenReturn(asList(cluster("1", "alpha")));

    ListClustersResponse response =
        api.listClusters(ListClustersRequest.builder().accountId("1").build());

    assertThat(response.getClusters())
        .extracting("accountId", "clusterName")
        .containsExactly(tuple("1", "alpha"));
  }

  @Test
  public void itFiltersByAccountIdAndPrefixIfGiven() throws Exception {
    when(repository.listClusters("1", "alpha"))
        .thenReturn(asList(cluster("1", "alpha-one"), cluster("1", "alpha-two")));

    ListClustersResponse response =
        api.listClusters(
            ListClustersRequest.builder().accountId("1").clusterNamePrefix("alpha").build());

    assertThat(response.getClusters())
        .extracting("accountId", "clusterName")
        .containsExactlyInAnyOrder(tuple("1", "alpha-one"), tuple("1", "alpha-two"));
  }
}
