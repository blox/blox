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

import com.googlecode.jsonrpc4j.IJsonRpcClient;
import java.lang.reflect.Type;
import java.util.Map;
import lombok.RequiredArgsConstructor;

/**
 * Adapter that allows using a {@link JsonRpcLambdaClient} as the transport for jsonrpc4j's dynamic
 * client proxy by implementing the {@link IJsonRpcClient} interface.
 *
 * <p>This class is package-local, because it should only be used through the {@link
 * JsonRpcLambdaClient#newProxy(Class)} method.
 */
@RequiredArgsConstructor
class JsonRpcLambdaProxyAdapter implements IJsonRpcClient {
  private final JsonRpcLambdaClient innerClient;

  @Override
  public void invoke(String methodName, Object argument) throws Throwable {
    innerClient.invoke(methodName, argument, null);
  }

  @Override
  public Object invoke(String methodName, Object argument, Type returnType) throws Throwable {
    return innerClient.invoke(methodName, argument, returnType);
  }

  @Override
  public Object invoke(
      String methodName, Object argument, Type returnType, Map<String, String> extraHeaders)
      throws Throwable {
    return innerClient.invoke(methodName, argument, returnType);
  }

  @Override
  public <T> T invoke(String methodName, Object argument, Class<T> returnType) throws Throwable {
    return innerClient.invoke(methodName, argument, returnType);
  }

  @Override
  public <T> T invoke(
      String methodName, Object argument, Class<T> returnType, Map<String, String> extraHeaders)
      throws Throwable {
    return innerClient.invoke(methodName, argument, returnType);
  }
}
