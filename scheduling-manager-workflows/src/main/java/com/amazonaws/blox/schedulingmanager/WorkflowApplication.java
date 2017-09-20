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
package com.amazonaws.blox.schedulingmanager;

import com.amazonaws.blox.schedulingmanager.deployment.steps.GetDeploymentData;
import com.amazonaws.blox.schedulingmanager.deployment.steps.GetStateData;
import com.amazonaws.blox.schedulingmanager.deployment.steps.StartDeployment;
import com.amazonaws.blox.schedulingmanager.handler.Encoder;
import com.amazonaws.blox.schedulingmanager.handler.Router;
import com.amazonaws.blox.schedulingmanager.handler.StepHandler;
import com.amazonaws.blox.schedulingmanager.task.steps.CheckTaskState;
import com.amazonaws.blox.schedulingmanager.task.steps.StartTask;
import com.amazonaws.blox.schedulingmanager.wrapper.ECSWrapperFactory;
import com.amazonaws.blox.schedulingmanager.wrapper.StepFunctionsWrapper;
import com.amazonaws.services.securitytoken.AWSSecurityTokenService;
import com.amazonaws.services.securitytoken.AWSSecurityTokenServiceClient;
import com.amazonaws.services.stepfunctions.AWSStepFunctions;
import com.amazonaws.services.stepfunctions.AWSStepFunctionsClient;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.util.HashMap;
import java.util.Map;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class WorkflowApplication {

  @Bean
  public AWSStepFunctions stepFunctionsClient() {
    return AWSStepFunctionsClient.builder().build();
  }

  @Bean
  public AWSSecurityTokenService stsClient() {
    return AWSSecurityTokenServiceClient.builder().build();
  }

  @Bean
  public ECSWrapperFactory ecsWrapperFactory() {
    return new ECSWrapperFactory(stsClient());
  }

  @Bean
  public StepFunctionsWrapper stepFunctionsWrapper() {
    return new StepFunctionsWrapper(stepFunctionsClient());
  }

  @Bean
  public Router router() {
    return new Router(handlers());
  }

  @Bean
  public Encoder encoder() {
    return new Encoder(objectMapper());
  }

  @Bean
  public ObjectMapper objectMapper() {
    return new ObjectMapper();
  }

  /**
   * This map is used for routing in the main workflow lambda handler. The lambda function aliases
   * used in the workflow definition need to match the map keys.
   */
  @Bean
  public Map<String, StepHandler> handlers() {
    final Map<String, StepHandler> handlers = new HashMap<String, StepHandler>();
    handlers.put("GetDeploymentData", getDeploymentData());
    handlers.put("GetStateData", getStateData());
    handlers.put("StartDeployment", startDeployment());
    handlers.put("StartTask", startTask());
    handlers.put("CheckTaskState", checkTaskState());
    return handlers;
  }

  @Bean
  public GetDeploymentData getDeploymentData() {
    return new GetDeploymentData(encoder());
  }

  @Bean
  public GetStateData getStateData() {
    return new GetStateData(encoder(), ecsWrapperFactory());
  }

  @Bean
  public StartDeployment startDeployment() {
    return new StartDeployment(encoder(), stepFunctionsWrapper());
  }

  @Bean
  public StartTask startTask() {
    return new StartTask(encoder(), ecsWrapperFactory());
  }

  @Bean
  public CheckTaskState checkTaskState() {
    return new CheckTaskState(encoder(), ecsWrapperFactory());
  }
}
