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
package com.amazonaws.blox.jsonrpc;

import static org.assertj.core.api.Assertions.assertThatThrownBy;

import com.amazonaws.blox.jsonrpc.fixtures.TestException;
import com.amazonaws.blox.jsonrpc.fixtures.TestService;
import com.amazonaws.blox.test.FakeLambdaClient;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.Test;

public class ErrorMappingIntegrationTest {
  ObjectMapper mapper = new ObjectMapper();

  @Test
  public void itMapsPojoExceptionsFromFailingMethods() {
    TestService server =
        (input) -> {
          throw new TestException(input, "beta");
        };

    JsonRpcLambdaHandler<TestService> handler =
        new JsonRpcLambdaHandler<>(mapper, TestService.class, server);
    JsonRpcLambdaClient lambdaClient =
        new JsonRpcLambdaClient(mapper, new FakeLambdaClient<>(handler), "foo");

    TestService client = lambdaClient.newProxy(TestService.class);

    assertThatThrownBy(() -> client.throwsException("foo"))
        .isInstanceOf(TestException.class)
        .hasMessage("foo:beta")
        .hasFieldOrPropertyWithValue("value", "foo")
        .hasFieldOrPropertyWithValue("otherValue", "beta");
  }
}
