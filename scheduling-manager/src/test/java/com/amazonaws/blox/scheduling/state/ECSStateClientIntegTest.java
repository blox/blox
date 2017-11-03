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
package com.amazonaws.blox.scheduling.state;

import static org.junit.Assert.assertThat;

import com.amazonaws.blox.testcategories.IntegrationTest;
import java.util.Map;
import java.util.stream.Collectors;
import org.hamcrest.Matchers;
import org.junit.BeforeClass;
import org.junit.Test;
import org.junit.experimental.categories.Category;
import software.amazon.awssdk.services.cloudformation.CloudFormationClient;
import software.amazon.awssdk.services.cloudformation.model.DescribeStacksRequest;
import software.amazon.awssdk.services.cloudformation.model.DescribeStacksResponse;
import software.amazon.awssdk.services.cloudformation.model.Output;
import software.amazon.awssdk.services.cloudformation.model.Parameter;
import software.amazon.awssdk.services.cloudformation.model.Stack;
import software.amazon.awssdk.services.cloudformation.model.StackStatus;
import software.amazon.awssdk.services.ecs.ECSAsyncClient;

@Category(IntegrationTest.class)
public class ECSStateClientIntegTest {

  private static final String STACK_NAME = "blox-test-cluster";
  public static String cluster;
  public static int expectedInstances = 0;
  public static int expectedTasks = 0;

  @BeforeClass
  public static void getStackOutputs() throws Exception {
    CloudFormationClient client = CloudFormationClient.create();

    DescribeStacksResponse stacks =
        client.describeStacks(DescribeStacksRequest.builder().stackName(STACK_NAME).build());
    Stack stack = stacks.stacks().get(0);

    if (!stack.stackStatus().equals(StackStatus.CREATE_COMPLETE.toString())) {
      throw new RuntimeException(
          String.format("Stack %s failed to create, status is %s", stack, stack.stackStatus()));
    }

    Map<String, String> outputs =
        stack.outputs().stream().collect(Collectors.toMap(Output::outputKey, Output::outputValue));

    Map<String, String> parameters =
        stack
            .parameters()
            .stream()
            .collect(Collectors.toMap(Parameter::parameterKey, Parameter::parameterValue));

    ECSStateClientIntegTest.cluster = outputs.get("cluster");
    ECSStateClientIntegTest.expectedInstances =
        Integer.parseInt(parameters.get("DesiredInstances"));
    ECSStateClientIntegTest.expectedTasks = Integer.parseInt(parameters.get("DesiredTasks"));
  }

  @Test
  public void describesAllTasksAndInstances() {
    // TODO add better integ tests to test multi-page scenarios
    ECSAsyncClient client = ECSAsyncClient.builder().build();
    ECSStateClient state = new ECSStateClient(client);

    ClusterSnapshot snapshot = state.snapshotState(cluster);

    assertThat(snapshot.getInstances(), Matchers.hasSize(expectedInstances));
    assertThat(snapshot.getTasks(), Matchers.hasSize(expectedTasks));
  }
}
