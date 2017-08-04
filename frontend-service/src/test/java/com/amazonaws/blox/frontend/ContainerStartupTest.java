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

import static org.junit.Assert.assertEquals;

import com.amazonaws.serverless.proxy.internal.model.AwsProxyResponse;
import com.amazonaws.serverless.proxy.internal.testutils.AwsProxyRequestBuilder;
import com.amazonaws.serverless.proxy.internal.testutils.MockLambdaContext;
import org.junit.Test;

public final class ContainerStartupTest {
  private static final LambdaHandler handler = new LambdaHandler();

  @Test
  public final void handleRealisticRequestSuccessfully() {
    AwsProxyResponse response =
        handler.handleRequest(
            new AwsProxyRequestBuilder().method("GET").path("/environments/test-env").build(),
            new MockLambdaContext());

    assertEquals(200, response.getStatusCode());
    assertEquals("{\"name\":\"test-env\"}", response.getBody());
  }
}
