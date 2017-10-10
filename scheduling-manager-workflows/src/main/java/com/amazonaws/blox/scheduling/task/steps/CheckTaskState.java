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
package com.amazonaws.blox.scheduling.task.steps;

import com.amazonaws.blox.scheduling.handler.StepHandler;
import com.amazonaws.blox.scheduling.handler.Encoder;
import com.amazonaws.blox.scheduling.task.steps.types.TaskData;
import com.amazonaws.blox.scheduling.wrapper.ECSWrapper;
import com.amazonaws.blox.scheduling.wrapper.ECSWrapperFactory;
import com.amazonaws.services.ecs.model.DescribeTasksResult;
import com.amazonaws.services.lambda.runtime.Context;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import lombok.AllArgsConstructor;
import lombok.NonNull;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@AllArgsConstructor
public class CheckTaskState implements StepHandler {

  private static final String ROLE_SESSION_NAME_PREFIX = "checkTaskState";

  @NonNull private Encoder encoder;
  @NonNull private ECSWrapperFactory ecsWrapperFactory;

  @Override
  public void handleRequest(InputStream input, OutputStream output, Context context)
      throws IOException {

    final TaskData taskData = encoder.decode(input, TaskData.class);

    final ECSWrapper ecsWrapper =
        ecsWrapperFactory.getWrapperForRole(taskData.getEcsRole(), ROLE_SESSION_NAME_PREFIX);

    final DescribeTasksResult describeTasksResult =
        ecsWrapper.describeTasks(taskData.getTasks(), taskData.getCluster());

    //assume one task for now
    taskData.setStatus(describeTasksResult.getTasks().get(0).getLastStatus());

    encoder.encode(output, taskData);
  }
}
