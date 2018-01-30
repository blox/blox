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
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceExistsException;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentResponse;
import org.junit.Test;

public class CreateEnvironmentIntegrationTest extends DataServiceIntegrationTestBase {
  private static final String ACCOUNT_ID = "123456789012";
  private static final String ENVIRONMENT_NAME = "environmentName";
  private static final String CLUSTER_ONE = "cluster1";
  private static final String CLUSTER_TWO = "cluster2";
  private static final String TASK_DEFINITION = "taskDefinition";
  private static final DataServiceModelBuilder models = DataServiceModelBuilder.builder().build();
  private static final EnvironmentId createdEnvironmentId1 =
      models
          .environmentId()
          .accountId(ACCOUNT_ID)
          .environmentName(ENVIRONMENT_NAME)
          .cluster(CLUSTER_ONE)
          .build();
  private static final EnvironmentId createdEnvironmentId2 =
      models
          .environmentId()
          .accountId(ACCOUNT_ID)
          .environmentName(ENVIRONMENT_NAME)
          .cluster(CLUSTER_TWO)
          .build();

  @Test
  public void testCreateEnvironmentSuccessful() throws Exception {
    final CreateEnvironmentResponse createEnvironmentResponse =
        dataService.createEnvironment(
            models
                .createEnvironmentRequest()
                .taskDefinition(TASK_DEFINITION)
                .environmentId(createdEnvironmentId1)
                .build());
    assertThat(createEnvironmentResponse.getEnvironment().getEnvironmentId())
        .isEqualTo(createdEnvironmentId1);
  }

  @Test
  public void testCreateAnEnvironmentAlreadyExist() throws Exception {
    dataService.createEnvironment(
        models
            .createEnvironmentRequest()
            .taskDefinition(TASK_DEFINITION)
            .environmentId(createdEnvironmentId1)
            .build());

    assertThatThrownBy(
            () ->
                dataService.createEnvironment(
                    models
                        .createEnvironmentRequest()
                        .taskDefinition(TASK_DEFINITION)
                        .environmentId(createdEnvironmentId1)
                        .build()))
        .isInstanceOf(ResourceExistsException.class)
        .hasMessageContaining(
            String.format("environment with id %s already exists", createdEnvironmentId1));
  }

  @Test
  public void testCreateTwoEnvironmentsWithTheSameNameButDifferentClusters() throws Exception {
    final CreateEnvironmentResponse createEnvironmentResponse1 =
        dataService.createEnvironment(
            models
                .createEnvironmentRequest()
                .taskDefinition(TASK_DEFINITION)
                .environmentId(createdEnvironmentId1)
                .build());
    assertThat(createEnvironmentResponse1.getEnvironment().getEnvironmentId())
        .isEqualTo(createdEnvironmentId1);

    final CreateEnvironmentResponse createEnvironmentResponse2 =
        dataService.createEnvironment(
            models
                .createEnvironmentRequest()
                .taskDefinition(TASK_DEFINITION)
                .environmentId(createdEnvironmentId2)
                .build());
    assertThat(createEnvironmentResponse2.getEnvironment().getEnvironmentId())
        .isEqualTo(createdEnvironmentId2);
  }
}
