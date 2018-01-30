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

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ObjectNode;
import com.fasterxml.jackson.databind.node.TextNode;
import com.googlecode.jsonrpc4j.ExceptionResolver;
import com.googlecode.jsonrpc4j.JsonRpcBasicServer;
import lombok.extern.slf4j.Slf4j;

@Slf4j
public class PojoExceptionResolver implements ExceptionResolver {
  private final ObjectMapper mapper;

  public PojoExceptionResolver(final ObjectMapper mapper) {
    this.mapper = mapper.copy().addMixIn(Throwable.class, ThrowableSerializationMixin.class);
  }

  /**
   * Resolves a given JSON-RPC error response into a Throwable.
   *
   * <p>Callers expect this method to return null (instead of raising an exception) if it cannot
   * deserialize the given response object.
   */
  @Override
  public Throwable resolveException(final ObjectNode response) {
    log.trace("Resolving exception from JSON response {}", response);

    final JsonNode error = response.get(JsonRpcBasicServer.ERROR);
    if (error == null || !error.isObject()) {
      log.warn("No error information found in JSON response {}", response);
      return null;
    }

    final JsonNode code = error.get(JsonRpcBasicServer.ERROR_CODE);
    if (code == null || !code.isInt() || code.asInt() != ThrowableSerializationMixin.ERROR_CODE) {
      log.warn("Not resolving exception for unsupported error code {}", code);
      return null;
    }

    final JsonNode data = error.get(JsonRpcBasicServer.DATA);
    if (data == null || !data.isObject()) {
      log.warn("No error details included in data field of JSON error {}", error);
      return null;
    }

    final ObjectNode dataObject = ((ObjectNode) data);
    if (!dataObject.has(JsonRpcBasicServer.ERROR_MESSAGE)) {
      dataObject.set(
          JsonRpcBasicServer.ERROR_MESSAGE,
          new TextNode(error.get(JsonRpcBasicServer.ERROR_MESSAGE).asText()));
    }

    try {
      return mapper.treeToValue(data, Throwable.class);
    } catch (final JsonProcessingException e) {
      log.warn("Unable to convert JSON response error '{}' into Throwable", data, e);
      return null;
    }
  }
}
