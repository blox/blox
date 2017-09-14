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
package com.amazonaws.blox.schedulingmanager.deployment;

import com.amazonaws.blox.schedulingmanager.deployment.handler.Encoder;
import com.amazonaws.blox.schedulingmanager.deployment.handler.Router;
import com.amazonaws.blox.schedulingmanager.deployment.steps.CheckTaskState;
import com.amazonaws.blox.schedulingmanager.deployment.steps.GetDeploymentData;
import com.amazonaws.blox.schedulingmanager.deployment.steps.GetStateData;
import com.amazonaws.blox.schedulingmanager.deployment.steps.StartDeployment;
import com.amazonaws.blox.schedulingmanager.deployment.steps.StartTask;
import com.amazonaws.blox.schedulingmanager.deployment.steps.StepHandler;
import com.amazonaws.services.ecs.AmazonECS;
import com.amazonaws.services.ecs.AmazonECSClient;
import com.amazonaws.services.stepfunctions.AWSStepFunctions;
import com.amazonaws.services.stepfunctions.AWSStepFunctionsClient;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.util.HashMap;
import java.util.Map;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class DeploymentWorkflowApplication {

  @Bean
  public AmazonECS ecsClient() {
    return AmazonECSClient.builder().build();
  }

  @Bean
  public AWSStepFunctions stepFunctionsClient() {
    return AWSStepFunctionsClient.builder().build();
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
    return new GetStateData(encoder());
  }

  @Bean
  public StartDeployment startDeployment() {
    return new StartDeployment(encoder(), stepFunctionsClient());
  }

  @Bean
  public StartTask startTask() {
    return new StartTask(encoder(), ecsClient());
  }

  @Bean
  public CheckTaskState checkTaskState() {
    return new CheckTaskState(encoder());
  }
}
