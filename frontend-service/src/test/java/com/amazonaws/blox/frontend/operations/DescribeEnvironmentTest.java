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
package com.amazonaws.blox.frontend.operations;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import com.amazonaws.blox.dataservicemodel.v1.model.Environment;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId;
import com.amazonaws.blox.frontend.mappers.DescribeEnvironmentMapper;
import com.amazonaws.blox.frontend.operations.DescribeEnvironment.DescribeEnvironmentResponse;
import java.time.Instant;
import org.junit.Before;
import org.junit.Test;
import org.springframework.beans.factory.annotation.Autowired;

public class DescribeEnvironmentTest extends EnvironmentControllerTestCase {
  DescribeEnvironment controller = new DescribeEnvironment();
  @Autowired DescribeEnvironmentMapper mapper;

  @Before
  public void setupController() {
    controller.request = servletRequest;
    controller.dataService = dataService;
    controller.mapper = mapper;
  }

  @Test
  public void mapsInputsAndOutputsCorrectly() throws Exception {
    requestContext.setAccountId(ACCOUNT_ID);

    EnvironmentId id =
        EnvironmentId.builder()
            .accountId(ACCOUNT_ID)
            .cluster(TEST_CLUSTER)
            .environmentName(ENVIRONMENT_NAME)
            .build();

    Environment environment =
        Environment.builder()
            .environmentId(id)
            .role(ROLE)
            .environmentType(ENVIRONMENT_TYPE)
            .environmentHealth(HEALTHY)
            .environmentStatus(STATUS)
            .deploymentMethod(DEPLOYMENT_METHOD)
            .createdTime(Instant.now())
            .lastUpdatedTime(Instant.now())
            .build();

    when(dataService.describeEnvironment(any()))
        .thenReturn(
            com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentResponse
                .builder()
                .environment(environment)
                .build());

    DescribeEnvironmentResponse response =
        controller.describeEnvironment(TEST_CLUSTER, ENVIRONMENT_NAME);

    verify(dataService)
        .describeEnvironment(
            com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRequest
                .builder()
                .environmentId(id)
                .build());

    assertThat(response).isNotNull();

    assertThat(response.getEnvironment())
        .isEqualToIgnoringGivenFields(environment, "environmentName", "cluster", "environmentType");
    assertThat(response.getEnvironment())
        .hasFieldOrPropertyWithValue("cluster", id.getCluster())
        .hasFieldOrPropertyWithValue("environmentName", id.getEnvironmentName())
        .hasFieldOrPropertyWithValue("environmentType", ENVIRONMENT_TYPE_STRING);
  }
}
