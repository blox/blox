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
package com.amazonaws.blox.scheduling.scheduler.engine;

import java.util.concurrent.CompletableFuture;
import software.amazon.awssdk.services.ecs.ECSAsyncClient;

public interface SchedulingAction {
  // TODO: We probably want to give this code an ECS facade that can only start/stop tasks in a
  //       single cluster, instead of the full ECS API.
  CompletableFuture<Boolean> execute(ECSAsyncClient ecs);
}
