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
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsResponse;
import com.amazonaws.blox.frontend.mappers.ListEnvironmentsMapper;
import com.amazonaws.blox.frontend.models.RequestPagination;
import com.amazonaws.serverless.proxy.internal.model.ApiGatewayRequestContext;
import com.amazonaws.serverless.proxy.internal.servlet.AwsProxyHttpServletRequestReader;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import javax.servlet.http.HttpServletRequest;

@RunWith(MockitoJUnitRunner.StrictStubs.class)
public class ListEnvironmentsTest {
  private static final String CLUSTER = "cluster";
  private static final String ENVIRONMENT_NAME_PREFIX = "environmentNamePrefix";

  @Mock private HttpServletRequest request;
  @Mock private DataService dataService;
  @Mock private ListEnvironmentsMapper mapper;
  @Mock private ApiGatewayRequestContext context;
  @Mock private ListEnvironmentsRequest dsRequest;
  @Mock private ListEnvironmentsResponse dsResponse;
  @Mock private ListEnvironments.ListEnvironmentsResponse feResponse;
  @Mock private RequestPagination pagination;

  @InjectMocks private ListEnvironments api;

  @Before
  public void setup() {
    when(request.getAttribute(AwsProxyHttpServletRequestReader.API_GATEWAY_CONTEXT_PROPERTY))
        .thenReturn(context);
  }

  @Test
  public void testListEnvironments() throws Exception {
    when(mapper.toListEnvironmentsRequest(context, CLUSTER, ENVIRONMENT_NAME_PREFIX))
        .thenReturn(dsRequest);
    when(dataService.listEnvironments(dsRequest)).thenReturn(dsResponse);
    when(mapper.fromDataServiceResponse(dsResponse)).thenReturn(feResponse);

    final ListEnvironments.ListEnvironmentsResponse response =
        api.listEnvironments(CLUSTER, ENVIRONMENT_NAME_PREFIX, pagination);

    assertThat(response).isEqualTo(feResponse);
    verify(mapper).toListEnvironmentsRequest(context, CLUSTER, ENVIRONMENT_NAME_PREFIX);
    verify(dataService).listEnvironments(dsRequest);
    verify(mapper).fromDataServiceResponse(dsResponse);
  }
}
