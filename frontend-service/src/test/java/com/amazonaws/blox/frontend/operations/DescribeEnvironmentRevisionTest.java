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

import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentRevision;
import com.amazonaws.blox.frontend.mappers.DescribeEnvironmentRevisionMapper;
import com.amazonaws.blox.frontend.operations.DescribeEnvironmentRevision.DescribeEnvironmentRevisionResponse;
import java.time.Instant;
import org.junit.Before;
import org.junit.Test;
import org.springframework.beans.factory.annotation.Autowired;

public class DescribeEnvironmentRevisionTest extends EnvironmentControllerTestCase {
  DescribeEnvironmentRevision controller = new DescribeEnvironmentRevision();
  @Autowired DescribeEnvironmentRevisionMapper mapper;

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

    EnvironmentRevision environmentRevision =
        EnvironmentRevision.builder()
            .environmentId(id)
            .environmentRevisionId(ENVIRONMENT_REVISION_ID)
            .instanceGroup(instanceGroupWithAttributeDS(ATTRIBUTE_NAME, ATTRIBUTE_VALUE))
            .taskDefinition(TASK_DEFINITION)
            .createdTime(Instant.now())
            .build();

    when(dataService.describeEnvironmentRevision(any()))
        .thenReturn(
            com.amazonaws.blox.dataservicemodel.v1.model.wrappers
                .DescribeEnvironmentRevisionResponse.builder()
                .environmentRevision(environmentRevision)
                .build());

    DescribeEnvironmentRevisionResponse response =
        controller.describeEnvironmentRevision(
            TEST_CLUSTER, ENVIRONMENT_NAME, ENVIRONMENT_REVISION_ID);

    verify(dataService)
        .describeEnvironmentRevision(
            com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRevisionRequest
                .builder()
                .environmentId(id)
                .environmentRevisionId(ENVIRONMENT_REVISION_ID)
                .build());

    assertThat(response).isNotNull();

    assertThat(response.getEnvironmentRevision());
  }
}
