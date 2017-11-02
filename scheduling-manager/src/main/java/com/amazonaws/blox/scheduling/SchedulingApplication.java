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
package com.amazonaws.blox.scheduling;

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.serialization.DataServiceMapperFactory;
import com.amazonaws.blox.jsonrpc.JsonRpcLambdaClient;
import com.amazonaws.blox.lambda.JacksonRequestStreamHandler;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import com.amazonaws.services.lambda.runtime.RequestStreamHandler;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Profile;
import software.amazon.awssdk.config.ClientOverrideConfiguration;
import software.amazon.awssdk.services.ecs.ECSAsyncClient;
import software.amazon.awssdk.services.lambda.LambdaAsyncClient;

/** Common beans required by all Scheduling lambda functions. */
public abstract class SchedulingApplication {

  static {
    // HACK: Disable IPv6 for the current JVM, since the Netty `epoll` driver fails in a Lambda
    // function if IPv6 is enabled.
    // Remove after https://github.com/aws/aws-sdk-java-v2/issues/193 is fixed.
    System.setProperty("java.net.preferIPv4Stack", "true");
  }

  /** The X-Ray Trace ID provided by the Lambda runtime environment. */
  @Value("${_X_AMZN_TRACE_ID}")
  public String traceId;

  @Value("${data_service_function_name}")
  public String dataServiceFunctionName;

  @Bean
  public <IN, OUT> RequestStreamHandler streamHandler(RequestHandler<IN, OUT> innerHandler) {
    return new JacksonRequestStreamHandler<>(mapper(), innerHandler);
  }

  @Bean
  @Profile("!test")
  public LambdaAsyncClient lambdaClient() {
    // add trace ID to downstream calls, for X-Ray integration
    ClientOverrideConfiguration configuration =
        ClientOverrideConfiguration.builder()
            .addAdditionalHttpHeader("X-Amzn-Trace-Id", traceId)
            .build();

    return LambdaAsyncClient.builder().overrideConfiguration(configuration).build();
  }

  @Bean
  @Profile("!test")
  public ECSAsyncClient ecs() {
    return ECSAsyncClient.builder().build();
  }

  @Bean
  public ObjectMapper mapper() {
    return new ObjectMapper().findAndRegisterModules();
  }

  @Bean
  @Profile("!test")
  public DataService dataService(LambdaAsyncClient lambda) {
    return new JsonRpcLambdaClient(
            DataServiceMapperFactory.newMapper(), lambda, dataServiceFunctionName)
        .newProxy(DataService.class);
  }
}
