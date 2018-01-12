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
package cucumber.configuration;

import com.amazonaws.blox.dataserviceclient.v1.client.DataServiceLambdaClient;
import cucumber.steps.helpers.CloudFormationStacks;
import cucumber.steps.helpers.ExceptionContext;
import cucumber.steps.helpers.InputCreator;
import cucumber.steps.wrappers.DataServiceWrapper;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;
import software.amazon.awssdk.services.cloudformation.CloudFormationClient;
import software.amazon.awssdk.services.lambda.LambdaAsyncClient;

@Configuration
@ComponentScan("cucumber.steps")
public class CucumberConfiguration {

  public static final String DATASERVICE_STACK = "data-service";
  public static final String DATASERVICE_LAMBDA_FUNCTION_KEY = "DataService";

  @Bean
  public DataServiceWrapper dataServiceWrapper() {
    return new DataServiceWrapper(
        DataServiceLambdaClient.dataService(
            lambdaAsyncClient(),
            cloudFormationStacks().get(DATASERVICE_STACK).output(DATASERVICE_LAMBDA_FUNCTION_KEY)),
        exceptionContext());
  }

  @Bean
  public ExceptionContext exceptionContext() {
    return new ExceptionContext();
  }

  @Bean
  public LambdaAsyncClient lambdaAsyncClient() {
    return LambdaAsyncClient.create();
  }

  @Bean
  public CloudFormationStacks cloudFormationStacks() {
    return new CloudFormationStacks(cloudFormationClient());
  }

  @Bean
  public CloudFormationClient cloudFormationClient() {
    return CloudFormationClient.create();
  }

  @Bean
  public InputCreator inputCreator() {
    return new InputCreator();
  }
}
