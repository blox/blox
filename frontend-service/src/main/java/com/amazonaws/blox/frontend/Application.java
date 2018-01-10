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
package com.amazonaws.blox.frontend;

import com.amazonaws.blox.dataserviceclient.v1.client.DataServiceLambdaClient;
import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.EnableWebMvc;
import software.amazon.awssdk.services.lambda.LambdaAsyncClient;

@Configuration
@EnableWebMvc
@ComponentScan("com.amazonaws.blox.frontend")
public class Application {

  @Value("${data_service_function_name}")
  public String dataServiceFunctionName;

  @Bean
  public DataService dataService() {
    return DataServiceLambdaClient.dataService(lambdaClient(), dataServiceFunctionName);
  }

  @Bean
  public LambdaAsyncClient lambdaClient() {
    return LambdaAsyncClient.create();
  }
}
