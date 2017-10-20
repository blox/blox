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

import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import java.util.concurrent.CompletableFuture;
import lombok.RequiredArgsConstructor;

@RequiredArgsConstructor
public class TestLambdaFunction<IN, OUT> implements LambdaFunction<IN, OUT> {

  final RequestHandler<IN, OUT> handler;

  @Override
  public CompletableFuture<OUT> callAsync(IN input) {
    return CompletableFuture.supplyAsync(() -> handler.handleRequest(input, fakeContext()));
  }

  @Override
  public CompletableFuture<Void> triggerAsync(IN input) {
    return CompletableFuture.runAsync(() -> handler.handleRequest(input, fakeContext()));
  }

  private Context fakeContext() {
    return null;
  }
}
