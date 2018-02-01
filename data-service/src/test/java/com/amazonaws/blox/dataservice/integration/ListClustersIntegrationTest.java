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
package com.amazonaws.blox.dataservice.integration;

import static org.assertj.core.api.Assertions.assertThat;

import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListClustersResponse;
import java.util.Arrays;
import org.junit.Test;

public class ListClustersIntegrationTest extends DataServiceIntegrationTestBase {
  private static final String CLUSTER_ONE = "cluster1";
  private static final String CLUSTER_TWO = "cluster2";
  private static final DataServiceModelBuilder models = DataServiceModelBuilder.builder().build();
  private static final EnvironmentId createdEnvironmentId1 =
      models.environmentId().cluster(CLUSTER_ONE).build();
  private static final EnvironmentId createdEnvironmentId2 =
      models.environmentId().cluster(CLUSTER_TWO).build();

  @Test
  public void testListClusterWithTwoEnvironments() throws Exception {
    dataService.createEnvironment(
        models.createEnvironmentRequest().environmentId(createdEnvironmentId1).build());
    dataService.createEnvironment(
        models.createEnvironmentRequest().environmentId(createdEnvironmentId2).build());
    final ListClustersResponse listClustersResponse =
        dataService.listClusters(models.listClustersRequest().clusterNamePrefix(null).build());
    assertThat(listClustersResponse.getClusters())
        .isEqualTo(
            Arrays.asList(
                models.cluster().clusterName(CLUSTER_ONE).build(),
                models.cluster().clusterName(CLUSTER_TWO).build()));
  }
}
