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

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.assertj.core.groups.Tuple.tuple;

import com.amazonaws.blox.dataservice.model.EnvironmentId;
import com.amazonaws.blox.dataservice.model.EnvironmentId.EnvironmentIdBuilder;
import com.amazonaws.blox.dataservice.test.data.ModelBuilders;
import java.util.Arrays;
import java.util.List;
import org.junit.Test;

public class ListClustersTest extends EnvironmentRepositoryTestBase {

  private static final String ACCOUNT_ID_ONE = "111111111111";
  private static final String ACCOUNT_ID_TWO = "222222222222";
  private static final String ACCOUNT_ID_THREE = "333333333333";
  ModelBuilders models = ModelBuilders.builder().build();

  @Test
  public void itDoesntIncludeDuplicateClusters() throws Exception {
    EnvironmentIdBuilder sameCluster =
        models.environmentId().accountId(ACCOUNT_ID_ONE).cluster("TestCluster");

    givenEnvironments(
        sameCluster.environmentName("EnvironmentOne").build(),
        sameCluster.environmentName("EnvironmentTwo").build());

    assertThat(repo.listClusters(ACCOUNT_ID_ONE, "TestCluster"))
        .extracting("accountId", "clusterName")
        .containsExactlyInAnyOrder(tuple(ACCOUNT_ID_ONE, "TestCluster"));
  }

  @Test
  public void itFiltersClustersByAccountIdIfGiven() throws Exception {
    EnvironmentIdBuilder sameAccountId = models.environmentId().accountId(ACCOUNT_ID_ONE);
    EnvironmentIdBuilder otherAccountId = models.environmentId().accountId(ACCOUNT_ID_TWO);

    givenEnvironments(
        sameAccountId.cluster("ClusterOne").environmentName("EnvironmentOne").build(),
        sameAccountId.cluster("ClusterTwo").environmentName("EnvironmentTwo").build(),
        otherAccountId.cluster("ClusterOne").environmentName("EnvironmentThree").build());

    assertThat(repo.listClusters(ACCOUNT_ID_ONE, null))
        .extracting("accountId", "clusterName")
        .containsExactlyInAnyOrder(
            tuple(ACCOUNT_ID_ONE, "ClusterOne"), tuple(ACCOUNT_ID_ONE, "ClusterTwo"));
  }

  @Test
  public void itFiltersClustersByAccountIdAndClusterPrefixIfGiven() throws Exception {
    EnvironmentIdBuilder sameAccountId = models.environmentId().accountId(ACCOUNT_ID_ONE);
    EnvironmentIdBuilder otherAccountId = models.environmentId().accountId(ACCOUNT_ID_TWO);

    givenEnvironments(
        sameAccountId.cluster("CommonPrefixOne").environmentName("EnvironmentOne").build(),
        sameAccountId.cluster("CommonPrefixTwo").environmentName("EnvironmentTwo").build(),
        otherAccountId.cluster("CommonPrefixThree").environmentName("EnvironmentFour").build(),
        sameAccountId.cluster("UncommonPrefix").environmentName("EnvironmentThree").build());

    assertThat(repo.listClusters(ACCOUNT_ID_ONE, "CommonPrefix"))
        .extracting("accountId", "clusterName")
        .containsExactlyInAnyOrder(
            tuple(ACCOUNT_ID_ONE, "CommonPrefixOne"), tuple(ACCOUNT_ID_ONE, "CommonPrefixTwo"));
  }

  @Test
  public void itReturnsAllClustersFromAllEnvironments() throws Exception {
    givenEnvironments(
        models.environmentId(ACCOUNT_ID_ONE, "ClusterOne", "EnvironmentOne"),
        models.environmentId(ACCOUNT_ID_TWO, "ClusterTwo", "EnvironmentTwo"),
        models.environmentId(ACCOUNT_ID_THREE, "ClusterThree", "EnvironmentThree"));

    assertThat(repo.listClusters(null, null))
        .extracting("accountId", "clusterName")
        .containsExactlyInAnyOrder(
            tuple(ACCOUNT_ID_ONE, "ClusterOne"),
            tuple(ACCOUNT_ID_TWO, "ClusterTwo"),
            tuple(ACCOUNT_ID_THREE, "ClusterThree"));
  }

  @Test
  public void itFailsIfPrefixGivenWithNoAccountId() throws Exception {
    assertThatThrownBy(() -> repo.listClusters(null, "SomePrefix"))
        .isInstanceOf(NullPointerException.class)
        .hasMessageContaining("accountId must be specified");
  }

  private void givenEnvironments(List<EnvironmentId> ids) throws Exception {
    for (EnvironmentId id : ids) {
      repo.createEnvironmentAndEnvironmentRevision(
          models.environment().environmentId(id).build(),
          models.environmentRevision().environmentId(id).build());
    }
  }

  private void givenEnvironments(EnvironmentId... ids) throws Exception {
    givenEnvironments(Arrays.asList(ids));
  }
}
