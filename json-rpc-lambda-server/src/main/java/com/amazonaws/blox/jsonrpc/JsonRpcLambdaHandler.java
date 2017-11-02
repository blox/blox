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

import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestStreamHandler;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.googlecode.jsonrpc4j.JsonRpcBasicServer;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;

/**
 * An AWS Lambda RequestStreamHandler that can dispatch incoming JSON-RPC formatted payloads to a
 * handler interface.
 *
 * @param <T> the type of the interface that implements the handler.
 */
public class JsonRpcLambdaHandler<T> implements RequestStreamHandler {
  private final T service;

  JsonRpcBasicServer server;

  public JsonRpcLambdaHandler(Class<T> serviceClass, T service) {
    this(defaultObjectMapper(), serviceClass, service);
  }

  public JsonRpcLambdaHandler(ObjectMapper mapper, Class<T> serviceClass, T service) {
    this.service = service;
    this.server = new JsonRpcBasicServer(mapper, this.service, serviceClass);
  }

  private static ObjectMapper defaultObjectMapper() {
    return new ObjectMapper();
  }

  @Override
  public void handleRequest(InputStream input, OutputStream output, Context context)
      throws IOException {
    server.handleRequest(input, output);
  }
}
