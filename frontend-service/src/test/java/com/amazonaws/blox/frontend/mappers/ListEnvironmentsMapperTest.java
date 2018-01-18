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
package com.amazonaws.blox.frontend.mappers;

import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsResponse;
import com.amazonaws.blox.frontend.operations.ListEnvironments;
import com.amazonaws.serverless.proxy.internal.model.ApiGatewayRequestContext;
import org.junit.Before;
import org.junit.Test;
import org.mapstruct.factory.Mappers;
import java.util.Collections;

import static org.assertj.core.api.Assertions.assertThat;

public class ListEnvironmentsMapperTest {
  private static final String ACCOUNT_ID = "123456789012";
  private static final String CLUSTER = "cluster";
  private static final String ENVIRONMENT_NAME_PREFIX = "environmentNamePrefix";
  private static final String ENVIRONMENT_NAME = "environmentName";

  private final ListEnvironmentsMapper mapper = Mappers.getMapper(ListEnvironmentsMapper.class);

  @Test
  public void testToListEnvironmentsRequest() {
    final ApiGatewayRequestContext context = new ApiGatewayRequestContext();
    context.setAccountId(ACCOUNT_ID);

    final ListEnvironmentsRequest request =
        mapper.toListEnvironmentsRequest(context, CLUSTER, ENVIRONMENT_NAME_PREFIX);

    assertThat(request.getCluster().getAccountId()).isEqualTo(ACCOUNT_ID);
    assertThat(request.getCluster().getClusterName()).isEqualTo(CLUSTER);
    assertThat(request.getEnvironmentNamePrefix()).isEqualTo(ENVIRONMENT_NAME_PREFIX);
  }

  @Test
  public void testFromDataServiceResponse() {
    final ListEnvironmentsResponse dsResponse =
        ListEnvironmentsResponse.builder()
            .environmentIds(
                Collections.singletonList(
                    EnvironmentId.builder()
                        .accountId(ACCOUNT_ID)
                        .cluster(CLUSTER)
                        .environmentName(ENVIRONMENT_NAME)
                        .build()))
            .build();

    final ListEnvironments.ListEnvironmentsResponse response =
        mapper.fromDataServiceResponse(dsResponse);

    assertThat(response.getEnvironmentNames().size()).isEqualTo(1);
    assertThat(response.getEnvironmentNames().get(0)).isEqualTo(ENVIRONMENT_NAME);
    assertThat(response.getNextToken()).isNull();
  }
}
