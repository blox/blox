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
import com.amazonaws.blox.frontend.mappers.UpdateEnvironmentMapper;
import com.amazonaws.blox.frontend.operations.UpdateEnvironment.UpdateEnvironmentRequest;
import com.amazonaws.blox.frontend.operations.UpdateEnvironment.UpdateEnvironmentResponse;
import org.junit.Test;
import org.springframework.beans.factory.annotation.Autowired;

public class UpdateEnvironmentTest extends EnvironmentControllerTestCase {

  private static final String NEW_TASK_DEFINITION = "new_task_definition";
  @Autowired UpdateEnvironment controller;
  @Autowired UpdateEnvironmentMapper mapper;

  @Test
  public void mapsInputsAndOutputsCorrectly() throws Exception {
    EnvironmentId id =
        EnvironmentId.builder()
            .accountId(ACCOUNT_ID)
            .cluster(TEST_CLUSTER)
            .environmentName(ENVIRONMENT_NAME)
            .build();

    when(dataService.updateEnvironment(any()))
        .thenReturn(
            com.amazonaws.blox.dataservicemodel.v1.model.wrappers.UpdateEnvironmentResponse
                .builder()
                .environmentRevisionId(ENVIRONMENT_REVISION_ID)
                .build());

    UpdateEnvironmentResponse response =
        controller.updateEnvironment(
            TEST_CLUSTER,
            ENVIRONMENT_NAME,
            UpdateEnvironmentRequest.builder()
                .environmentName(ENVIRONMENT_NAME)
                .taskDefinition(NEW_TASK_DEFINITION)
                .build());

    verify(dataService)
        .updateEnvironment(
            com.amazonaws.blox.dataservicemodel.v1.model.wrappers.UpdateEnvironmentRequest.builder()
                .environmentId(id)
                .taskDefinition(NEW_TASK_DEFINITION)
                .build());

    assertThat(response).isNotNull();

    assertThat(response.getEnvironmentRevisionId()).isEqualTo(ENVIRONMENT_REVISION_ID);
  }
}
