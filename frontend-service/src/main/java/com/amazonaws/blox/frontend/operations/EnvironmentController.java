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
import com.amazonaws.serverless.proxy.internal.model.ApiGatewayRequestContext;
import com.amazonaws.serverless.proxy.internal.servlet.AwsProxyHttpServletRequestReader;
import io.swagger.annotations.Api;
import javax.servlet.http.HttpServletRequest;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RequestMapping;

@Api
@RequestMapping(path = "/v1/{cluster}/environments", produces = "application/json")
public abstract class EnvironmentController {
  // Spring will autowire the current request into this field, so that we can get a handle on the
  // request context from API gateway.
  @Autowired HttpServletRequest request;

  @Autowired DataService dataService;

  /**
   * Get the API gateway request context from the current request.
   *
   * @return
   */
  protected ApiGatewayRequestContext getApiGatewayRequestContext() {
    return (ApiGatewayRequestContext)
        request.getAttribute(AwsProxyHttpServletRequestReader.API_GATEWAY_CONTEXT_PROPERTY);
  }
}
