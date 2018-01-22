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
import com.amazonaws.blox.dataservicemodel.v1.model.Environment;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentRevision;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRevisionRequest;
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
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import software.amazon.awssdk.services.ecs.ECSAsyncClient;

@Component
@RequiredArgsConstructor
@Slf4j
public class SchedulerHandler implements RequestHandler<SchedulerInput, SchedulerOutput> {
  private final DataService data;
  private final ECSAsyncClient ecs;
  private final SchedulerFactory schedulerFactory;

  @SneakyThrows
  @Override
  public SchedulerOutput handleRequest(SchedulerInput input, Context context) {
    log.debug("Request: {}", input);

    EnvironmentId environmentId = input.getEnvironmentId();

    Environment environment =
        data.describeEnvironment(
                DescribeEnvironmentRequest.builder().environmentId(environmentId).build())
            .getEnvironment();

    String activeEnvironmentId = environment.getActiveEnvironmentRevisionId();

    EnvironmentRevision activeEnvironmentRevision =
        data.describeEnvironmentRevision(
                DescribeEnvironmentRevisionRequest.builder()
                    .environmentId(environmentId)
                    .environmentRevisionId(activeEnvironmentId)
                    .build())
            .getEnvironmentRevision();

    EnvironmentDescription environmentDescription =
        EnvironmentDescription.builder()
            .environmentName(environmentId.getEnvironmentName())
            .activeEnvironmentRevisionId(activeEnvironmentId)
            .environmentType(
                EnvironmentDescription.EnvironmentType.valueOf(
                    environment.getEnvironmentType().toString()))
            .taskDefinitionArn(activeEnvironmentRevision.getTaskDefinition())
            .build();

    Scheduler s = schedulerFactory.schedulerFor(environmentDescription);

    List<SchedulingAction> actions = s.schedule(input.getSnapshot(), environmentDescription);

    List<Boolean> outcomes =
        actions.stream().map(a -> a.execute(ecs)).collect(CompletableFutures.joinList()).join();

    Map<Boolean, Long> outcomeCounts =
        outcomes
            .stream()
            .collect(Collectors.groupingBy(Function.identity(), Collectors.counting()));

    return new SchedulerOutput(
        input.getSnapshot().getClusterName(),
        input.getEnvironmentId(),
        outcomeCounts.getOrDefault(false, 0L),
        outcomeCounts.getOrDefault(true, 0L));
    // TODO: handle exceptions. captured in the lambda exception handling issue
  }
}
