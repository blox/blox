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
package com.amazonaws.blox.frontend;

import com.amazonaws.serverless.exceptions.ContainerInitializationException;
import com.amazonaws.serverless.proxy.internal.LambdaContainerHandler;
import com.amazonaws.serverless.proxy.internal.model.ApiGatewayRequestContext;
import com.amazonaws.serverless.proxy.internal.model.AwsProxyRequest;
import com.amazonaws.serverless.proxy.internal.model.AwsProxyResponse;
import com.amazonaws.serverless.proxy.internal.servlet.AwsHttpServletResponse;
import com.amazonaws.serverless.proxy.internal.servlet.AwsProxyHttpServletRequest;
import com.amazonaws.serverless.proxy.spring.SpringLambdaContainerHandler;
import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import lombok.extern.log4j.Log4j;

/** Entrypoint for mapping incoming API GW requests to Lambda into Spring controllers */
@Log4j
public class LambdaHandler implements RequestHandler<AwsProxyRequest, AwsProxyResponse> {
  private SpringLambdaContainerHandler<AwsProxyRequest, AwsProxyResponse> handler;

  public AwsProxyResponse handleRequest(AwsProxyRequest request, Context context) {
    ApiGatewayRequestContext requestContext = request.getRequestContext();

    log.info("Handling request: " + requestContext.getRequestId());
    log.debug("Caller identity: " + requestContext.getIdentity().getCaller());
    AwsProxyResponse response = getHandler().proxy(request, context);

    log.info(
        "Completed request: "
            + requestContext.getRequestId()
            + ", with response: "
            + response.getStatusCode());
    return response;
  }

  private LambdaContainerHandler<
          AwsProxyRequest, AwsProxyResponse, AwsProxyHttpServletRequest, AwsHttpServletResponse>
      getHandler() {
    // TODO: lazy initialization was copied from the example code, figure out why we can't
    // statically initialize the handler instead.
    if (handler == null) {
      try {
        log.info("Initializing handler (probably a cold start)");
        handler = SpringLambdaContainerHandler.getAwsProxyHandler(Application.class);
      } catch (ContainerInitializationException e) {
        throw new RuntimeException(e);
      }
    }
    return handler;
  }
}
