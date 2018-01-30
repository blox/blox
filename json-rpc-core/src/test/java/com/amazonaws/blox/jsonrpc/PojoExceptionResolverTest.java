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

import com.amazonaws.blox.jsonrpc.fixtures.PojoTestException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.IntNode;
import com.fasterxml.jackson.databind.node.ObjectNode;
import com.fasterxml.jackson.databind.node.TextNode;
import com.googlecode.jsonrpc4j.JsonRpcBasicServer;
import org.junit.Test;

public class PojoExceptionResolverTest {

  private static final Integer RANDOM_ERROR_CODE = 123456;
  private static final String VALID_MESSAGE = "Something went horribly wrong.";
  private static final String OTHER_MESSAGE = "Something completely different went horribly wrong.";
  private ObjectMapper mapper = new ObjectMapper();

  private PojoExceptionResolver resolver = new PojoExceptionResolver(mapper);

  @Test
  public void itDeserializesPojoExceptionsWithJsonCreatorConstructor() {
    ObjectNode response =
        response(
            error(
                ThrowableSerializationMixin.ERROR_CODE,
                VALID_MESSAGE,
                data(PojoTestException.class.getTypeName(), "someValue=foo", "anotherValue=bar")));

    Throwable throwable = resolver.resolveException(response);

    assertThat(throwable)
        .isInstanceOf(PojoTestException.class)
        .hasMessage("foo:bar")
        .hasFieldOrPropertyWithValue("someValue", "foo")
        .hasFieldOrPropertyWithValue("anotherValue", "bar");
  }

  @Test
  public void itReturnsNullForNonErrorResponse() {
    ObjectNode response = response(null);

    Throwable throwable = resolver.resolveException(response);

    assertThat(throwable).isNull();
  }

  @Test
  public void itReturnsNullForOtherErrorCodes() {
    ObjectNode response =
        response(
            error(RANDOM_ERROR_CODE, VALID_MESSAGE, data(PojoTestException.class.getTypeName())));

    Throwable throwable = resolver.resolveException(response);

    assertThat(throwable).isNull();
  }

  @Test
  public void itReturnsNullForErrorsWithMissingData() {
    ObjectNode response =
        response(error(ThrowableSerializationMixin.ERROR_CODE, VALID_MESSAGE, null));

    Throwable throwable = resolver.resolveException(response);

    assertThat(throwable).isNull();
  }

  @Test
  public void itReturnsNullWhenJacksonDeserializationFails() {
    ObjectNode response =
        response(
            error(
                ThrowableSerializationMixin.ERROR_CODE,
                VALID_MESSAGE,
                data(PojoTestException.class.getTypeName(), "invalidProperty=badValue")));

    Throwable throwable = resolver.resolveException(response);

    assertThat(throwable).isNull();
  }

  @Test
  public void itDeserializesExceptionsWithOnlyMessageConstructor() {
    ObjectNode response =
        response(
            error(
                ThrowableSerializationMixin.ERROR_CODE,
                VALID_MESSAGE,
                data(RuntimeException.class.getTypeName())));

    Throwable throwable = resolver.resolveException(response);

    assertThat(throwable).isInstanceOf(RuntimeException.class).hasMessage(VALID_MESSAGE);
  }

  @Test
  public void itPrefersMessageFromDataOverMessageFromError() {
    ObjectNode response =
        response(
            error(
                ThrowableSerializationMixin.ERROR_CODE,
                VALID_MESSAGE,
                data(RuntimeException.class.getTypeName(), "message=" + OTHER_MESSAGE)));

    Throwable throwable = resolver.resolveException(response);

    assertThat(throwable).isInstanceOf(RuntimeException.class).hasMessage(OTHER_MESSAGE);
  }

  private ObjectNode response(ObjectNode error) {
    ObjectNode response = mapper.createObjectNode();
    response.set(JsonRpcBasicServer.ERROR, error);

    return response;
  }

  private ObjectNode error(Integer code, String message, ObjectNode data) {
    ObjectNode error = mapper.createObjectNode();
    if (code != null) error.set(JsonRpcBasicServer.ERROR_CODE, new IntNode(code));
    if (message != null) error.set(JsonRpcBasicServer.ERROR_MESSAGE, new TextNode(message));
    if (data != null) error.set(JsonRpcBasicServer.DATA, data);

    return error;
  }

  private ObjectNode data(String exceptionTypeName, String... keyValuePairs) {
    ObjectNode data = mapper.createObjectNode();
    if (exceptionTypeName != null)
      data.set(JsonRpcBasicServer.EXCEPTION_TYPE_NAME, new TextNode(exceptionTypeName));
    for (String pair : keyValuePairs) {
      String[] parts = pair.split("=");
      data.set(parts[0], new TextNode(parts[1]));
    }

    return data;
  }
}
