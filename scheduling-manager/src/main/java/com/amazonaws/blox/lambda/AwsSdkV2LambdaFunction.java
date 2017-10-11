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
package com.amazonaws.blox.lambda;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectReader;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.fasterxml.jackson.databind.util.ByteBufferBackedInputStream;
import java.nio.ByteBuffer;
import java.util.concurrent.CompletableFuture;
import lombok.SneakyThrows;
import lombok.extern.log4j.Log4j2;
import software.amazon.awssdk.services.lambda.LambdaAsyncClient;
import software.amazon.awssdk.services.lambda.model.InvocationType;
import software.amazon.awssdk.services.lambda.model.InvokeRequest;
import software.amazon.awssdk.services.lambda.model.InvokeResponse;

@Log4j2
public class AwsSdkV2LambdaFunction<IN, OUT> implements LambdaFunction<IN, OUT> {

  private final LambdaAsyncClient lambda;
  private final ObjectReader reader;
  private final ObjectWriter writer;
  private final String functionName;

  public AwsSdkV2LambdaFunction(
      LambdaAsyncClient lambda, ObjectMapper mapper, Class<OUT> outputClass, String functionName) {

    this.lambda = lambda;

    this.reader = mapper.readerFor(outputClass);
    this.writer = mapper.writer();

    this.functionName = functionName;
  }

  @Override
  public CompletableFuture<OUT> callAsync(IN input) {
    CompletableFuture<OUT> response =
        callLambda(InvocationType.RequestResponse, input).thenApply(this::deserialize);

    response.thenAccept(r -> log.debug("response from {}: {}", functionName, r));

    return response;
  }

  @Override
  public CompletableFuture<Void> triggerAsync(IN input) {
    CompletableFuture<InvokeResponse> response = callLambda(InvocationType.Event, input);

    return response.thenAccept(
        r -> log.debug("response from {}: {}", functionName, r.statusCode()));
  }

  private CompletableFuture<InvokeResponse> callLambda(InvocationType type, IN input) {
    log.debug("calling '{}' as {} with payload: {}", functionName, type, input);

    ByteBuffer payload = serialize(input);
    InvokeRequest request =
        InvokeRequest.builder()
            .functionName(functionName)
            .invocationType(type)
            .payload(payload)
            .build();

    return lambda.invoke(request);
  }

  @SneakyThrows
  private ByteBuffer serialize(IN input) {
    return ByteBuffer.wrap(writer.writeValueAsBytes(input));
  }

  @SneakyThrows
  private OUT deserialize(InvokeResponse response) {
    return reader.readValue(new ByteBufferBackedInputStream(response.payload()));
  }
}
