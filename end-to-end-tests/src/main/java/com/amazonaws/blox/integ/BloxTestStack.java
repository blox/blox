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
package com.amazonaws.blox.integ;

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import java.util.List;
import software.amazon.awssdk.services.cloudformation.CloudFormationClient;
import software.amazon.awssdk.services.dynamodb.DynamoDBClient;
import software.amazon.awssdk.services.ecs.ECSClient;
import software.amazon.awssdk.services.ecs.model.Task;
import software.amazon.awssdk.services.lambda.LambdaAsyncClient;

/**
 * Facade over everything in an end-to-end Blox stack that's needed in tests.
 *
 * <p>TODO: it makes more sense to wire the individual support classes into the unit tests instead
 * of this whole stack, so this will probably be removed soon.
 */
public class BloxTestStack {
  private CloudFormationClient cloudFormationClient = CloudFormationClient.create();
  private DynamoDBClient dynamoDBClient = DynamoDBClient.create();
  private ECSClient ecsClient = ECSClient.create();
  private LambdaAsyncClient lambdaClient = LambdaAsyncClient.create();

  private CloudFormationStacks stacks = new CloudFormationStacks(cloudFormationClient);

  private final DataServiceWrapper dataService =
      new DataServiceWrapper(lambdaClient, stacks, dynamoDBClient);
  private final ECSClusterWrapper ecs = new ECSClusterWrapper(ecsClient, stacks);

  public DataService createDataService() {
    // TODO: Replace with frontend once we have it wired up
    return dataService.createDataService();
  }

  public List<Task> describeTasks() {
    return ecs.describeTasks();
  }

  public String getTaskDefinition() {
    return ecs.getTaskDefinition();
  }

  public String getCluster() {
    return ecs.getCluster();
  }

  public void reset() {
    dataService.reset();
    ecs.reset();
  }
}
