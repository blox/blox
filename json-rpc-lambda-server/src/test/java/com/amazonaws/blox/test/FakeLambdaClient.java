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
package com.amazonaws.blox.test;

import com.amazonaws.blox.jsonrpc.JsonRpcLambdaHandler;
import com.fasterxml.jackson.databind.util.ByteBufferBackedInputStream;
import java.io.ByteArrayOutputStream;
import java.nio.ByteBuffer;
import java.util.concurrent.CompletableFuture;
import lombok.SneakyThrows;
import software.amazon.awssdk.services.lambda.LambdaAsyncClient;
import software.amazon.awssdk.services.lambda.model.InvokeRequest;
import software.amazon.awssdk.services.lambda.model.InvokeResponse;

public class FakeLambdaClient<T> implements LambdaAsyncClient {
  private final JsonRpcLambdaHandler<T> handler;

  public FakeLambdaClient(final JsonRpcLambdaHandler<T> handler) {
    this.handler = handler;
  }

  @Override
  public void close() {}

  @Override
  @SneakyThrows
  public CompletableFuture<InvokeResponse> invoke(final InvokeRequest invokeRequest) {
    ByteBufferBackedInputStream inputStream =
        new ByteBufferBackedInputStream(invokeRequest.payload());
    ByteArrayOutputStream outputStream = new ByteArrayOutputStream();

    handler.handleRequest(inputStream, outputStream, null);

    return CompletableFuture.completedFuture(
        InvokeResponse.builder().payload(ByteBuffer.wrap(outputStream.toByteArray())).build());
  }
}
