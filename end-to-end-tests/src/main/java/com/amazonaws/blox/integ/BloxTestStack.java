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

import com.amazonaws.auth.DefaultAWSCredentialsProviderChain;
import com.amazonaws.blox.Blox;
import java.util.List;
import software.amazon.awssdk.services.cloudformation.CloudFormationClient;
import software.amazon.awssdk.services.ecs.ECSClient;
import software.amazon.awssdk.services.ecs.model.Task;

/**
 * Facade over everything in an end-to-end Blox stack that's needed in tests.
 *
 * <p>TODO: it makes more sense to wire the individual support classes into the unit tests instead
 * of this whole stack, so this will probably be removed soon.
 */
public class BloxTestStack {

  private final String bloxEndpoint;
  private final CloudFormationClient cloudFormationClient;
  private final ECSClient ecsClient;

  private final CloudFormationStacks stacks;
  private final ECSClusterWrapper ecs;

  private final Blox blox;

  public BloxTestStack(String bloxEndpoint) {
    this.bloxEndpoint = bloxEndpoint;

    this.cloudFormationClient = CloudFormationClient.create();
    this.ecsClient = ECSClient.create();
    this.stacks = new CloudFormationStacks(cloudFormationClient);

    this.ecs = new ECSClusterWrapper(ecsClient, stacks);

    this.blox =
        Blox.builder()
            .iamCredentials(new DefaultAWSCredentialsProviderChain())
            .endpoint(this.bloxEndpoint)
            .build();
  }

  public Blox getBlox() {
    return this.blox;
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
    ecs.reset();
  }
}
