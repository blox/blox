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

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.util.ByteBufferBackedInputStream;
import com.googlecode.jsonrpc4j.IJsonRpcClient;
import com.googlecode.jsonrpc4j.JsonRpcClient;
import com.googlecode.jsonrpc4j.ProxyUtil;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.lang.reflect.Type;
import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CompletionException;
import lombok.extern.slf4j.Slf4j;
import software.amazon.awssdk.services.lambda.LambdaAsyncClient;
import software.amazon.awssdk.services.lambda.model.InvokeRequest;
import software.amazon.awssdk.services.lambda.model.InvokeResponse;

/**
 * A JSON-RPC client that translates calls into Lambda function invocations.
 *
 * <p>The typical way to use this class is to use the {@link #newProxy(Class)} method to create a
 * dynamic client proxy that corresponds to a Java interface. However, this interface only supports
 * synchronously invoking remote functions.
 *
 * <p>For asynchronous invocation, callers must use {@link #invokeAsync(String, Object, Class)} to
 * make a raw JSON-RPC service call.
 */
@Slf4j
public class JsonRpcLambdaClient {
  private final LambdaAsyncClient lambda;
  private final String functionName;

  /** The inner JSON-RPC client that only knows how to read/write to Streams */
  private final JsonRpcClient streamClient;

  public JsonRpcLambdaClient(ObjectMapper mapper, LambdaAsyncClient lambda, String functionName) {
    this.lambda = lambda;
    this.functionName = functionName;
    this.streamClient = new JsonRpcClient(mapper);
  }

  /**
   * Create a new dynamic JSON-RPC proxy for an interface that dispatches calls as Lambda function
   * invocations using this client.
   *
   * @return a dynamic proxy that implements serviceInterface
   */
  public <T> T newProxy(Class<T> serviceInterface) {
    IJsonRpcClient client = new JsonRpcLambdaProxyAdapter(this);
    return ProxyUtil.createClientProxy(
        Thread.currentThread().getContextClassLoader(), serviceInterface, client);
  }

  /**
   * Invoke a JSON-RPC method asynchronously
   *
   * @throws IOException if serializing the request payload fails
   */
  @SuppressWarnings("unchecked")
  public <T> CompletableFuture<T> invokeAsync(
      String methodName, Object argument, Class<T> returnType) throws IOException {
    return invokeAsync(methodName, argument, (Type) returnType).thenApply(r -> (T) r);
  }

  /**
   * Invoke a JSON-RPC method synchronously.
   *
   * <p>This method should not be used directly, prefer using {@link #newProxy(Class)} to create a
   * strongly-typed dynamic proxy for a Java interface instead.
   *
   * @throws IOException if serializing the request payload or deserializing the response payload
   *     fails
   * @throws Throwable if the response is a serialized exception from the server
   */
  @SuppressWarnings("unchecked")
  public <T> T invoke(String methodName, Object argument, Class<T> returnType) throws Throwable {
    return (T) invoke(methodName, argument, (Type) returnType);
  }

  /**
   * Invoke a JSON-RPC method asynchronously without compile-time type-checking.
   *
   * <p>This overload is private, because it should only be used by the {@link
   * JsonRpcLambdaProxyAdapter} (which requires a non-typesafe invoke method).
   */
  private CompletableFuture<Object> invokeAsync(String methodName, Object argument, Type returnType)
      throws IOException {

    ByteBuffer requestPayload = writeRequest(methodName, argument);

    InvokeRequest invokeRequest =
        InvokeRequest.builder().functionName(functionName).payload(requestPayload).build();

    CompletableFuture<InvokeResponse> pendingRequest = lambda.invoke(invokeRequest);
    return pendingRequest.thenApply(r -> readResponse(returnType, r));
  }

  /**
   * Invoke a JSON-RPC method synchronously without compile-time type-checking.
   *
   * <p>This overload is package-private, because it should only be used by the {@link
   * JsonRpcLambdaProxyAdapter} (which requires a non-typesafe invoke method).
   */
  Object invoke(String methodName, Object argument, Type returnType) throws Throwable {
    try {
      return invokeAsync(methodName, argument, returnType).join();
    } catch (CompletionException ex) {
      // Unwrap CompletionExceptions as normal exceptions
      throw ex.getCause();
    }
  }

  /**
   * Serialize a JSON-RPC request into a ByteBuffer for use as a Lambda request payload.
   *
   * @throws IOException if serializing the request payload fails
   */
  private ByteBuffer writeRequest(String methodName, Object argument) throws IOException {
    ByteArrayOutputStream requestStream = new ByteArrayOutputStream();

    streamClient.invoke(methodName, argument, requestStream);
    log.trace("Raw request payload: " + requestStream.toString());

    return ByteBuffer.wrap(requestStream.toByteArray());
  }

  /** Deserialize a JSON-RPC response from a ByteBuffer from a Lambda response payload. */
  private Object readResponse(Type returnType, InvokeResponse response) {
    try {
      log.trace(
          "Raw response payload: " + StandardCharsets.UTF_8.decode(response.payload()).toString());

      return streamClient.readResponse(
          returnType, new ByteBufferBackedInputStream(response.payload()));
    } catch (Throwable t) {
      if (t instanceof IOException) {
        log.warn("Could not read JSON-RPC response: ", t);
      } else {
        // This exception is raised when deserialized from the remote call, and must be
        // handled/logged in calling code; so only log at DEBUG level
        log.debug("Read an exception value from JSON-RPC response: ", t);
      }

      // wrap the raw Throwable from the JSON RPC client so that this method can be used in
      // CompletableFuture::thenApply
      throw new CompletionException(t);
    }
  }
}
