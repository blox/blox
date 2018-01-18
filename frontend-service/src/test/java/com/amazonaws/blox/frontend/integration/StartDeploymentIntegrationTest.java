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
import static org.mockito.Mockito.when;

import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.StartDeploymentResponse;
import com.amazonaws.blox.frontend.operations.StartDeployment;
import com.amazonaws.serverless.proxy.internal.model.AwsProxyRequest;
import com.amazonaws.serverless.proxy.internal.model.AwsProxyResponse;
import com.amazonaws.serverless.proxy.internal.testutils.AwsProxyRequestBuilder;
import org.junit.Test;

public class StartDeploymentIntegrationTest extends IntegrationTestBase {

  private static final String DEPLOYMENT_ID = "deployment-id";
  private static final String ENVIRONMENT_REVISION_ID = "1";

  @Test
  public void testStartDeployment() throws Exception {
    // Given
    when(dataService.startDeployment(any()))
        .thenReturn(
            StartDeploymentResponse.builder()
                .environmentId(environmentId)
                .deploymentId(DEPLOYMENT_ID)
                .environmentRevisionId(ENVIRONMENT_REVISION_ID)
                .build());

    // When
    final AwsProxyRequest request =
        new AwsProxyRequestBuilder("/v1/myClsuter/environments/myEnv/deployments", "POST")
            .queryString("revisionId", "1")
            .build();
    request.getRequestContext().setAccountId(ACCOUNT_ID);

    final AwsProxyResponse response = handler.proxy(request, lambdaContext);

    // Then
    assertThat(response.getStatusCode()).isEqualTo(200);
    StartDeployment.StartDeploymentResponse startDeploymentResponse =
        objectMapper.readValue(response.getBody(), StartDeployment.StartDeploymentResponse.class);
    assertThat(startDeploymentResponse.getDeploymentId()).isEqualTo(DEPLOYMENT_ID);
  }

  @Test
  public void testStartDeploymentWithResourceNotFound() throws Exception {
    // Given
    when(dataService.startDeployment(any())).thenThrow(ResourceNotFoundException.class);

    // When
    final AwsProxyRequest request =
        new AwsProxyRequestBuilder("/v1/myClsuter/environments/myEnv/deployments", "POST")
            .queryString("revisionId", "1")
            .build();
    request.getRequestContext().setAccountId(ACCOUNT_ID);

    final AwsProxyResponse response = handler.proxy(request, lambdaContext);

    // Then
    assertThat(response.getStatusCode()).isGreaterThan(399);
    // TODO: Exception handling
    // assertThat(response.getStatusCode()).isEqualTo(404);
  }
}
