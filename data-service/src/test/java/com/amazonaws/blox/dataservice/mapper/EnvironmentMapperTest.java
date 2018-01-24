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
package com.amazonaws.blox.dataservice.mapper;

import static org.assertj.core.api.Assertions.assertThat;

import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentId;
import com.amazonaws.blox.dataservice.repository.model.EnvironmentDDBRecord;
import com.amazonaws.blox.dataservice.test.data.ModelBuilders;
import com.amazonaws.blox.dataservice.test.data.RecordBuilders;
import org.junit.Test;
import org.mapstruct.factory.Mappers;

public class EnvironmentMapperTest {

  private static final String ACCOUNT_ID = "accountId";
  private static final String CLUSTER = "cluster";
  private static final String ENVIRONMENT_NAME = "environmentName";
  private static final String ACCOUNT_ID_CLUSTER = "accountId/cluster";

  EnvironmentMapper mapper = Mappers.getMapper(EnvironmentMapper.class);

  ModelBuilders models = ModelBuilders.builder().build();
  RecordBuilders records = RecordBuilders.builder().build();

  @Test
  public void itUnpacksEnvironmentIdCorrectly() throws Exception {
    // The EnvironmentMapper has to both unpack the EnvironmentId into individual
    // accountId/clusterName attributes (so that we can index on those attributes), and store it as
    // a concatenated hash/range key.
    //
    // Since this is pretty easy to accidentally break in the mapper, we test it here.
    EnvironmentId id =
        EnvironmentId.builder()
            .accountId(ACCOUNT_ID)
            .cluster(CLUSTER)
            .environmentName(ENVIRONMENT_NAME)
            .build();
    EnvironmentDDBRecord record =
        mapper.toEnvironmentDDBRecord(models.environment().environmentId(id).build());

    assertThat(record.getAccountId()).isEqualTo(ACCOUNT_ID);
    assertThat(record.getClusterName()).isEqualTo(CLUSTER);
    assertThat(record.getEnvironmentName()).isEqualTo(ENVIRONMENT_NAME);
    assertThat(record.getAccountIdCluster()).isEqualTo(ACCOUNT_ID_CLUSTER);
  }

  @Test
  public void itMapsEnvironmentIdFromHashKeyOnly() throws Exception {
    // When the account ID and cluster in the accountIdCluster field differs from that in the
    // separate accountId/cluster fields
    String accountIdCluster = "accountIdFromId/clusterFromId";
    String accountIdFromId = "accountIdFromId";
    String clusterFromId = "clusterFromId";

    Environment environment =
        mapper.toEnvironment(
            records
                .environment()
                .accountIdCluster(accountIdCluster)
                .accountId(ACCOUNT_ID)
                .clusterName(CLUSTER)
                .environmentName(ENVIRONMENT_NAME)
                .build());

    // Then prefer the values in accountIdCluster:
    EnvironmentId id = environment.getEnvironmentId();
    assertThat(id.getAccountId()).isEqualTo(accountIdFromId);
    assertThat(id.getCluster()).isEqualTo(clusterFromId);
    assertThat(id.getEnvironmentName()).isEqualTo(ENVIRONMENT_NAME);
  }
}
