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

import static com.amazonaws.blox.test.JsonNodeAssert.assertThat;

import com.amazonaws.blox.jsonrpc.fixtures.PojoService;
import com.amazonaws.blox.jsonrpc.fixtures.PojoTestException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.TextNode;
import com.googlecode.jsonrpc4j.ErrorResolver.JsonError;
import com.googlecode.jsonrpc4j.JsonRpcBasicServer;
import java.util.Arrays;
import org.junit.Test;

public class PojoErrorResolverTest {
  private ObjectMapper mapper = new ObjectMapper();
  private PojoErrorResolver resolver = new PojoErrorResolver(mapper);

  @Test
  public void itPreservesErrorInformationForModeledExceptions() throws Exception {
    PojoTestException e =
        new PojoTestException("some value", "another value", new RuntimeException("oh dear!"));

    JsonError error =
        resolver.resolveError(
            e,
            PojoService.class.getMethod("pojoCall", String.class),
            Arrays.asList(new TextNode("foo")));

    JsonNode node = mapper.valueToTree(error);

    assertThat(node)
        .hasField("code", PojoErrorResolver.ERROR_CODE)
        .hasField("message", "some value:another value");

    JsonNode data = node.get("data");

    assertThat(data)
        .hasField(JsonRpcBasicServer.EXCEPTION_TYPE_NAME, PojoTestException.class.getName())
        .hasField("someValue", "some value")
        .hasField("anotherValue", "another value")
        .hasNoField("cause")
        .hasNoField("stackTrace")
        .hasNoField("suppressed")
        .hasNoField("localisedMessage");
  }

  @Test
  public void itReturnsGenericErrorForUnmodeledExceptions() throws Exception {
    PojoTestException e =
        new PojoTestException("some value", "another value", new RuntimeException("oh dear!"));

    JsonError error =
        resolver.resolveError(
            e,
            PojoService.class.getMethod("nonPojoCall", String.class),
            Arrays.asList(new TextNode("foo")));

    JsonNode node = mapper.valueToTree(error);

    assertThat(node)
        .hasField("code", JsonError.INTERNAL_ERROR.code)
        .hasField("message", "some value:another value")
        .hasNullField("data");
  }
}
