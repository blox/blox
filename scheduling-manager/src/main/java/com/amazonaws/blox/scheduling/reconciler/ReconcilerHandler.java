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
package com.amazonaws.blox.scheduling.reconciler;

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListClustersRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListClustersResponse;
import com.amazonaws.blox.lambda.LambdaFunction;
import com.amazonaws.blox.scheduling.manager.ManagerInput;
import com.amazonaws.blox.scheduling.manager.ManagerOutput;
import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import com.spotify.futures.CompletableFutures;
import java.util.List;
import java.util.Map;
import lombok.RequiredArgsConstructor;
import lombok.SneakyThrows;
import lombok.extern.log4j.Log4j2;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
@Log4j2
public class ReconcilerHandler implements RequestHandler<CloudWatchEvent<Map>, Void> {
  final DataService dataService;
  final LambdaFunction<ManagerInput, ManagerOutput> stateFunction;

  @Override
  @SneakyThrows // TODO add checked exception handling
  public Void handleRequest(CloudWatchEvent<Map> input, Context context) {
    ListClustersResponse r = dataService.listClusters(ListClustersRequest.builder().build());
    List<String> clusters = r.getClusters();

    clusters
        .stream()
        .map(c -> stateFunction.triggerAsync(new ManagerInput(c)))
        .collect(CompletableFutures.joinList())
        .join();

    return null;
  }
}
