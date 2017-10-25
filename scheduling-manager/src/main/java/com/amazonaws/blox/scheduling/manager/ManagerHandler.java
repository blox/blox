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

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsResponse;
import com.amazonaws.blox.lambda.LambdaFunction;
import com.amazonaws.blox.scheduling.scheduler.SchedulerInput;
import com.amazonaws.blox.scheduling.scheduler.SchedulerOutput;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot;
import com.amazonaws.blox.scheduling.state.ECSState;
import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import com.spotify.futures.CompletableFutures;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Stream;
import lombok.RequiredArgsConstructor;
import lombok.SneakyThrows;
import lombok.extern.log4j.Log4j2;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
@Log4j2
public class ManagerHandler implements RequestHandler<ManagerInput, ManagerOutput> {
  private final DataService data;
  private final ECSState ecs;
  private final LambdaFunction<SchedulerInput, SchedulerOutput> scheduler;

  @Override
  @SneakyThrows // TODO add checked exception handling
  public ManagerOutput handleRequest(ManagerInput input, Context context) {
    ListEnvironmentsResponse r =
        data.listEnvironments(
            ListEnvironmentsRequest.builder().cluster(input.getClusterArn()).build());
    List<String> environments = r.getEnvironmentIds();

    ClusterSnapshot state = ecs.snapshotState(input.getClusterArn());

    Stream<CompletableFuture<SchedulerOutput>> pendingRequests =
        environments
            .stream()
            .map(environmentId -> scheduler.callAsync(new SchedulerInput(state, environmentId)));

    List<SchedulerOutput> outputs = pendingRequests.collect(CompletableFutures.joinList()).join();

    return new ManagerOutput(input.getClusterArn(), outputs);
  }
}
