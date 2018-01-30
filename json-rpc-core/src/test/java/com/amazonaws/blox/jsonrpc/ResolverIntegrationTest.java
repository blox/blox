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

import static org.assertj.core.api.Assertions.assertThat;

import com.amazonaws.blox.jsonrpc.fixtures.PojoService;
import com.amazonaws.blox.jsonrpc.fixtures.PojoTestException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ObjectNode;
import com.googlecode.jsonrpc4j.ErrorResolver;
import com.googlecode.jsonrpc4j.ErrorResolver.JsonError;
import com.googlecode.jsonrpc4j.ExceptionResolver;
import com.googlecode.jsonrpc4j.JsonRpcBasicServer;
import java.util.Collections;
import org.junit.Test;

public class ResolverIntegrationTest {
  ObjectMapper mapper = new ObjectMapper();

  @Test
  public void clientSideResolverResolvesServerSideErrors() throws Exception {
    ErrorResolver server = new PojoErrorResolver(mapper);
    ExceptionResolver client = new PojoExceptionResolver(mapper);

    PojoTestException originalException = new PojoTestException("alpha", "beta");

    JsonError originalError =
        server.resolveError(
            originalException,
            PojoService.class.getMethod("pojoCall", String.class),
            Collections.emptyList());

    ObjectNode response = wrapErrorAsResponse(originalError);

    Throwable deserializedThrowable = client.resolveException(response);

    assertThat(deserializedThrowable)
        .isInstanceOf(PojoTestException.class)
        .isEqualToComparingFieldByField(originalException);
  }

  private ObjectNode wrapErrorAsResponse(final JsonError originalError) {
    ObjectNode response = mapper.createObjectNode();
    response.set(JsonRpcBasicServer.ERROR, mapper.valueToTree(originalError));
    return response;
  }
}
