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
package com.amazonaws.blox.frontend.integration;

import static org.mockito.Mockito.mock;

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.frontend.MapperConfiguration;
import com.amazonaws.blox.frontend.integration.SampleController.Checker;
import com.amazonaws.serverless.proxy.internal.model.AwsProxyRequest;
import com.amazonaws.serverless.proxy.internal.model.AwsProxyResponse;
import com.amazonaws.serverless.proxy.internal.testutils.MockLambdaContext;
import com.amazonaws.serverless.proxy.spring.SpringLambdaContainerHandler;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Import;
import org.springframework.web.context.ConfigurableWebApplicationContext;
import org.springframework.web.servlet.config.annotation.EnableWebMvc;

@Configuration
@Import(MapperConfiguration.class)
@ComponentScan("com.amazonaws.blox.frontend.operations")
@EnableWebMvc
public class IntegrationTestConfig {

  @Autowired private ConfigurableWebApplicationContext applicationContext;

  @Bean
  public SpringLambdaContainerHandler<AwsProxyRequest, AwsProxyResponse>
      sprintLambdaContainerHandler() throws Exception {
    SpringLambdaContainerHandler<AwsProxyRequest, AwsProxyResponse> handler =
        SpringLambdaContainerHandler.getAwsProxyHandler(applicationContext);
    handler.setRefreshContext(false);
    return handler;
  }

  @Bean
  public SampleController sampleController() {
    return new SampleController();
  }

  @Bean
  public Checker sampleControllerChecker() {
    return mock(Checker.class);
  }

  @Bean
  public DataService dataService() {
    return mock(DataService.class);
  }

  @Bean
  public MockLambdaContext lambdaContext() {
    return new MockLambdaContext();
  }

  @Bean
  public ObjectMapper mapper() {
    return new ObjectMapper();
  }
}
