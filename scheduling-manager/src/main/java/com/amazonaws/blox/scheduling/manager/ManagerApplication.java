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
package com.amazonaws.blox.scheduling.manager;

import com.amazonaws.blox.lambda.LambdaFunction;
import com.amazonaws.blox.lambda.AwsSdkV2LambdaFunction;
import com.amazonaws.blox.scheduling.SchedulingApplication;
import com.amazonaws.blox.scheduling.scheduler.SchedulerInput;
import com.amazonaws.blox.scheduling.scheduler.SchedulerOutput;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;
import software.amazon.awssdk.services.lambda.LambdaAsyncClient;

@Configuration
@ComponentScan({
  "com.amazonaws.blox.scheduling.manager",
  "com.amazonaws.blox.scheduling.state",
})
public class ManagerApplication extends SchedulingApplication {

  // Wired in through environment variable in CloudFormation template
  @Value("${scheduler_function_name}")
  String schedulerFunctionName;

  @Bean
  public LambdaFunction<SchedulerInput, SchedulerOutput> scheduler(
      LambdaAsyncClient lambda, ObjectMapper mapper) {
    return new AwsSdkV2LambdaFunction<>(
        lambda, mapper, SchedulerOutput.class, schedulerFunctionName);
  }
}
