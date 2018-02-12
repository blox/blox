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

import com.amazonaws.blox.dataservicemodel.v1.model.Environment;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRevisionRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRevisionResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.UpdateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.UpdateEnvironmentResponse;
import org.junit.Test;

public class UpdateEnvironmentIntegrationTest extends DataServiceIntegrationTestBase {

  private static final String NEW_TASK_DEFINITION = "newtaskdefinition";
  private final DataServiceModelBuilder models = DataServiceModelBuilder.builder().build();
  private final EnvironmentId id =
      models.environmentId().environmentName("EnvironmentToUpdate").build();

  @Test
  public void updateEnvironmentCreatesNewRevisionAndUpdatesLatest() throws Exception {
    dataService.createEnvironment(models.createEnvironmentRequest().environmentId(id).build());

    UpdateEnvironmentResponse response =
        dataService.updateEnvironment(
            UpdateEnvironmentRequest.builder()
                .environmentId(id)
                .taskDefinition(NEW_TASK_DEFINITION)
                .build());

    String newRevisionId = response.getEnvironmentRevisionId();

    DescribeEnvironmentRevisionResponse newRevision =
        dataService.describeEnvironmentRevision(
            DescribeEnvironmentRevisionRequest.builder()
                .environmentId(id)
                .environmentRevisionId(newRevisionId)
                .build());

    assertThat(newRevision.getEnvironmentRevision().getTaskDefinition())
        .isEqualTo(NEW_TASK_DEFINITION);

    Environment environment =
        dataService
            .describeEnvironment(DescribeEnvironmentRequest.builder().environmentId(id).build())
            .getEnvironment();

    assertThat(environment.getLatestEnvironmentRevisionId()).isEqualTo(newRevisionId);
  }
}
