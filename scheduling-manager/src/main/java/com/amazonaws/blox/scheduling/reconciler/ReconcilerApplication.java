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
package com.amazonaws.blox.scheduling.reconciler;

import com.amazonaws.blox.lambda.AwsSdkV2LambdaFunction;
import com.amazonaws.blox.lambda.LambdaFunction;
import com.amazonaws.blox.scheduling.SchedulingApplication;
import com.amazonaws.blox.scheduling.manager.ManagerInput;
import com.amazonaws.blox.scheduling.manager.ManagerOutput;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;
import software.amazon.awssdk.services.lambda.LambdaAsyncClient;

@Configuration
@ComponentScan("com.amazonaws.blox.scheduling.reconciler")
public class ReconcilerApplication extends SchedulingApplication {

  // Wired in through environment variable in CloudFormation template
  @Value("${manager_function_name}")
  String managerFunctionName;

  @Bean
  public LambdaFunction<ManagerInput, ManagerOutput> manager(
      LambdaAsyncClient lambda, ObjectMapper mapper) {
    return new AwsSdkV2LambdaFunction<>(lambda, mapper, ManagerOutput.class, managerFunctionName);
  }
}
