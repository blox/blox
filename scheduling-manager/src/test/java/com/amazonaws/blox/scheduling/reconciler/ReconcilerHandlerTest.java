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

import static org.hamcrest.CoreMatchers.hasItems;
import static org.junit.Assert.assertThat;
import static org.mockito.Matchers.any;
import static org.mockito.Mockito.when;

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListClustersResponse;
import com.amazonaws.blox.lambda.LambdaFunction;
import com.amazonaws.blox.scheduling.manager.ManagerInput;
import com.amazonaws.blox.scheduling.manager.ManagerOutput;
import java.util.Arrays;
import java.util.concurrent.CompletableFuture;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.ArgumentCaptor;
import org.mockito.Mock;
import org.mockito.runners.MockitoJUnitRunner;

@RunWith(MockitoJUnitRunner.class)
public class ReconcilerHandlerTest {

  public static final String FIRST_CLUSTER_ARN = "arn:::::cluster1";
  public static final String SECOND_CLUSTER_ARN = "arn:::::cluster2";

  private ArgumentCaptor<ManagerInput> input = ArgumentCaptor.forClass(ManagerInput.class);
  @Mock private DataService data;
  @Mock private LambdaFunction<ManagerInput, ManagerOutput> manager;

  @Test
  public void invokesManagerAsynchronouslyForAllClusters() throws Exception {
    when(data.listClusters(any()))
        .thenReturn(new ListClustersResponse(Arrays.asList(FIRST_CLUSTER_ARN, SECOND_CLUSTER_ARN)));

    when(manager.triggerAsync(input.capture())).thenReturn(CompletableFuture.completedFuture(null));

    ReconcilerHandler handler = new ReconcilerHandler(data, manager);
    handler.handleRequest(new CloudWatchEvent<>(), null);

    assertThat(
        input.getAllValues(),
        hasItems(new ManagerInput(FIRST_CLUSTER_ARN), new ManagerInput(SECOND_CLUSTER_ARN)));
  }
}
