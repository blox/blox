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

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.StartDeploymentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.StartDeploymentResponse;
import com.amazonaws.blox.frontend.mappers.StartDeploymentMapper;
import com.amazonaws.serverless.proxy.internal.model.ApiGatewayRequestContext;
import com.amazonaws.serverless.proxy.internal.servlet.AwsProxyHttpServletRequestReader;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;

import javax.servlet.http.HttpServletRequest;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@RunWith(MockitoJUnitRunner.class)
public class StartDeploymentTest {

  private static final String CLUSTER = "myCluster";
  private static final String ENVIRONMENT_NAME = "TestEnvironment";
  private static final String ENVIRONMENT_REVISION_ID = "1";
  @Mock private HttpServletRequest request;
  @Mock private DataService dataService;
  @Mock private StartDeploymentMapper mapper;
  @Mock private ApiGatewayRequestContext context;

  @InjectMocks private StartDeployment api;

  @Before
  public void setUp() {
    when(request.getAttribute(AwsProxyHttpServletRequestReader.API_GATEWAY_CONTEXT_PROPERTY))
        .thenReturn(context);
  }

  @Test
  public void testStartDeployment() throws Exception {
    // Given
    final StartDeploymentRequest dsRequest = mock(StartDeploymentRequest.class);
    final StartDeploymentResponse dsResponse = mock(StartDeploymentResponse.class);
    final StartDeployment.StartDeploymentResponse feResponse =
        mock(StartDeployment.StartDeploymentResponse.class);

    when(mapper.toDataServiceRequest(context, CLUSTER, ENVIRONMENT_NAME, ENVIRONMENT_REVISION_ID))
        .thenReturn(dsRequest);
    when(dataService.startDeployment(dsRequest)).thenReturn(dsResponse);
    when(mapper.fromDataServiceResponse(dsResponse)).thenReturn(feResponse);

    // When
    final StartDeployment.StartDeploymentResponse response =
        api.startDeployment(CLUSTER, ENVIRONMENT_NAME, ENVIRONMENT_REVISION_ID);

    // Then
    assertThat(response).isEqualTo(feResponse);
    verify(mapper)
        .toDataServiceRequest(context, CLUSTER, ENVIRONMENT_NAME, ENVIRONMENT_REVISION_ID);
    verify(dataService).startDeployment(dsRequest);
    verify(mapper).fromDataServiceResponse(dsResponse);
  }
}
