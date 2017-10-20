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

import java.util.Collections;
import lombok.SneakyThrows;

/**
 * Temporary fake implementation of ECS State client that just takes time and returns an empty
 * cluster.
 */
public class ECSStateClient implements ECSState {

  @Override
  @SneakyThrows
  public ClusterSnapshot snapshotState(String clusterArn) {
    Thread.sleep(5000); // this will take a while
    return new ClusterSnapshot(clusterArn, Collections.emptyMap(), Collections.emptyMap());
  }
}
