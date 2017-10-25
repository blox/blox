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
package com.amazonaws.blox.scheduling.scheduler;

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeTargetEnvironmentRevisionRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeTargetEnvironmentRevisionResponse;
import com.amazonaws.blox.scheduling.scheduler.engine.EnvironmentDescription;
import com.amazonaws.blox.scheduling.scheduler.engine.Scheduler;
import com.amazonaws.blox.scheduling.scheduler.engine.SchedulerFactory;
import com.amazonaws.blox.scheduling.scheduler.engine.SchedulingAction;
import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import com.spotify.futures.CompletableFutures;
import java.util.List;
import java.util.Map;
import java.util.function.Function;
import java.util.stream.Collectors;
import lombok.RequiredArgsConstructor;
import lombok.SneakyThrows;
import lombok.extern.log4j.Log4j2;
import org.springframework.stereotype.Component;
import software.amazon.awssdk.services.ecs.ECSAsyncClient;

@Component
@RequiredArgsConstructor
@Log4j2
public class SchedulerHandler implements RequestHandler<SchedulerInput, SchedulerOutput> {
  private final DataService data;
  private final ECSAsyncClient ecs;
  private final SchedulerFactory schedulerFactory;

  @SneakyThrows
  @Override
  public SchedulerOutput handleRequest(SchedulerInput input, Context context) {
    log.debug("Request: {}", input);

    DescribeTargetEnvironmentRevisionResponse targetEnvironmentRevision =
        data.describeTargetEnvironmentRevision(
            DescribeTargetEnvironmentRevisionRequest.builder()
                .environmentId(input.getEnvironmentId())
                .build());
    DescribeEnvironmentResponse environment =
        data.describeEnvironment(
            DescribeEnvironmentRequest.builder()
                .environmentId(targetEnvironmentRevision.getEnvironmentId())
                .environmentVersion(targetEnvironmentRevision.getEnvironmentVersion())
                .build());

    EnvironmentDescription environmentDescription =
        EnvironmentDescription.builder()
            .environmentId(environment.getEnvironmentId())
            .environmentVersion(environment.getEnvironmentVersion())
            .taskDefinitionArn(environment.getTaskDefinition())
            .build();

    Scheduler s = schedulerFactory.schedulerFor(environment.getType());

    List<SchedulingAction> actions = s.schedule(input.getSnapshot(), environmentDescription);

    List<Boolean> outcomes =
        actions.stream().map(a -> a.execute(ecs)).collect(CompletableFutures.joinList()).join();

    Map<Boolean, Long> outcomeCounts =
        outcomes
            .stream()
            .collect(Collectors.groupingBy(Function.identity(), Collectors.counting()));

    return new SchedulerOutput(
        input.getSnapshot().getClusterArn(),
        input.getEnvironmentId(),
        outcomeCounts.getOrDefault(false, 0L),
        outcomeCounts.getOrDefault(true, 0L));
    //TODO: handle exceptions. captured in the lambda exception handling issue
  }
}
