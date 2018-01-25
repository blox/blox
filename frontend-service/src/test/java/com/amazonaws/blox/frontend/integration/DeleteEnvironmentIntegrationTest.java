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
package com.amazonaws.blox.frontend.integration;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import com.amazonaws.blox.dataservicemodel.v1.model.Environment;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentHealth;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentType;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DeleteEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DeleteEnvironmentResponse;
import com.amazonaws.blox.frontend.operations.DeleteEnvironment;
import com.amazonaws.serverless.proxy.internal.model.AwsProxyRequest;
import com.amazonaws.serverless.proxy.internal.model.AwsProxyResponse;
import com.amazonaws.serverless.proxy.internal.testutils.AwsProxyRequestBuilder;
import java.time.Instant;
import org.junit.Before;
import org.junit.Test;
import org.mockito.ArgumentCaptor;

public class DeleteEnvironmentIntegrationTest extends IntegrationTestBase {

  private Environment environment;

  @Before
  public void setup() {
    environment =
        Environment.builder()
            .environmentId(environmentId)
            .role("role")
            .environmentType(EnvironmentType.Daemon)
            .environmentHealth(EnvironmentHealth.HEALTHY)
            .environmentStatus("ok")
            .deploymentMethod("TerminateThenReplace")
            .createdTime(Instant.now())
            .lastUpdatedTime(Instant.now())
            .build();
  }

  @Test
  public void testDeleteEnvironment() throws Exception {
    when(dataService.deleteEnvironment(any()))
        .thenReturn(DeleteEnvironmentResponse.builder().environment(environment).build());

    final AwsProxyRequest request =
        new AwsProxyRequestBuilder("/v1/myCluster/environments/myEnv", "DELETE").build();
    request.getRequestContext().setAccountId(ACCOUNT_ID);
    final AwsProxyResponse response = handler.proxy(request, lambdaContext);

    assertThat(response.getStatusCode()).isEqualTo(200);
    DeleteEnvironment.DeleteEnvironmentResponse deleteEnvironmentResponse =
        objectMapper.readValue(
            response.getBody(), DeleteEnvironment.DeleteEnvironmentResponse.class);
    assertThat(deleteEnvironmentResponse.getEnvironment())
        .isEqualToIgnoringGivenFields(environment, "environmentName", "cluster", "environmentType");
    assertThat(deleteEnvironmentResponse.getEnvironment())
        .hasFieldOrPropertyWithValue("cluster", CLUSTER_NAME)
        .hasFieldOrPropertyWithValue("environmentName", ENVIRONMENT_NAME)
        .hasFieldOrPropertyWithValue("environmentType", EnvironmentType.Daemon.toString());
    ArgumentCaptor<DeleteEnvironmentRequest> captor =
        ArgumentCaptor.forClass(DeleteEnvironmentRequest.class);
    verify(dataService).deleteEnvironment(captor.capture());
    assertThat(captor.getValue().isForceDelete()).isFalse();
  }

  @Test
  public void testDeleteEnvironmentWithForceDelete() throws Exception {
    when(dataService.deleteEnvironment(any()))
        .thenReturn(DeleteEnvironmentResponse.builder().environment(environment).build());

    final AwsProxyRequest request =
        new AwsProxyRequestBuilder("/v1/myCluster/environments/myEnv", "DELETE")
            .queryString("forceDelete", "true")
            .build();
    request.getRequestContext().setAccountId(ACCOUNT_ID);
    final AwsProxyResponse response = handler.proxy(request, lambdaContext);

    assertThat(response.getStatusCode()).isEqualTo(200);
    ArgumentCaptor<DeleteEnvironmentRequest> captor =
        ArgumentCaptor.forClass(DeleteEnvironmentRequest.class);
    verify(dataService).deleteEnvironment(captor.capture());
    assertThat(captor.getValue().isForceDelete()).isTrue();
  }
}
