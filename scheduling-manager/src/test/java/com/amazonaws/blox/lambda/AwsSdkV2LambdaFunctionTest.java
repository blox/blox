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

import static org.hamcrest.CoreMatchers.equalTo;
import static org.junit.Assert.assertThat;
import static org.mockito.Matchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.util.ByteBufferBackedInputStream;
import java.nio.ByteBuffer;
import java.util.concurrent.CompletableFuture;
import lombok.Data;
import org.junit.Before;
import org.junit.Test;
import org.mockito.ArgumentCaptor;
import software.amazon.awssdk.services.lambda.LambdaAsyncClient;
import software.amazon.awssdk.services.lambda.model.InvocationType;
import software.amazon.awssdk.services.lambda.model.InvokeRequest;
import software.amazon.awssdk.services.lambda.model.InvokeResponse;

public class AwsSdkV2LambdaFunctionTest {

  private final ArgumentCaptor<InvokeRequest> requestArgument =
      ArgumentCaptor.forClass(InvokeRequest.class);
  private final ObjectMapper mapper = new ObjectMapper();
  private final LambdaAsyncClient client = mock(LambdaAsyncClient.class);
  private final AwsSdkV2LambdaFunction<TestInput, TestOutput> function =
      new AwsSdkV2LambdaFunction<>(client, mapper, TestOutput.class, "TestFunction");

  @Before
  public void defaultStubs() throws Exception {
    when(client.invoke(any()))
        .thenReturn(
            CompletableFuture.completedFuture(
                InvokeResponse.builder()
                    .payload(ByteBuffer.wrap(mapper.writeValueAsBytes(new TestOutput("value"))))
                    .build()));
  }

  @Test
  public void deserializesOutputsUsingMapper() throws Exception {
    TestOutput output = function.call(new TestInput("test"));
    assertThat(output.getOutput(), equalTo("value"));
  }

  @Test
  public void serializesInputsUsingMapper() throws Exception {
    TestOutput output = function.call(new TestInput("test"));
    verify(client).invoke(requestArgument.capture());

    ByteBuffer actualPayload = requestArgument.getValue().payload();
    assertThat(
        mapper
            .readValue(new ByteBufferBackedInputStream(actualPayload), TestInput.class)
            .getInput(),
        equalTo("test"));
  }

  @Test
  public void passesCorrectFunctionName() {
    function.call(new TestInput("test"));
    verify(client).invoke(requestArgument.capture());

    assertThat(requestArgument.getValue().functionName(), equalTo("TestFunction"));
  }

  @Test
  public void callUsesRequestReply() {
    function.callAsync(new TestInput("test")).join();

    verify(client).invoke(requestArgument.capture());
    assertThat(
        requestArgument.getValue().invocationType(),
        equalTo(InvocationType.RequestResponse.toString()));
  }

  @Test
  public void triggerUsesEventInvocation() {
    function.triggerAsync(new TestInput("test")).join();

    verify(client).invoke(requestArgument.capture());
    assertThat(
        requestArgument.getValue().invocationType(), equalTo(InvocationType.Event.toString()));
  }

  @Data
  static class TestInput {

    final String input;
  }

  @Data
  static class TestOutput {

    final String output;
  }
}
