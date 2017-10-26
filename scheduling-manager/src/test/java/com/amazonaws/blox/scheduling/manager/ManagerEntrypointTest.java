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

import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.mockito.Matchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsResponse;
import com.amazonaws.blox.lambda.LambdaFunction;
import com.amazonaws.blox.lambda.TestLambdaFunction;
import com.amazonaws.blox.scheduling.LambdaHandlerTestCase;
import com.amazonaws.blox.scheduling.manager.ManagerEntrypointTest.TestConfig;
import com.amazonaws.blox.scheduling.scheduler.SchedulerInput;
import com.amazonaws.blox.scheduling.scheduler.SchedulerOutput;
import com.amazonaws.blox.scheduling.state.ClusterSnapshot;
import com.amazonaws.blox.scheduling.state.ECSState;
import java.util.Collections;
import org.junit.Test;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Import;
import org.springframework.test.context.ContextConfiguration;

@ContextConfiguration(classes = TestConfig.class)
public class ManagerEntrypointTest extends LambdaHandlerTestCase {

  @Test
  public void convertsInputsAndOutputsFromJson() throws Exception {
    String result = callHandler(fixture("handlers/Manager.input.json"));
    assertThat(result, is(fixtureAsString("handlers/Manager.output.json")));
  }

  @Configuration
  @Import(ManagerApplication.class)
  public static class TestConfig {

    @Bean
    public DataService dataService() throws Exception {
      return when(mock(DataService.class).listEnvironments(any()))
          .thenReturn(new ListEnvironmentsResponse(Collections.singletonList("SomeEnvironment")))
          .getMock();
    }

    @Bean
    public ECSState ecsState() {
      return when(mock(ECSState.class).snapshotState(any()))
          .thenReturn(
              new ClusterSnapshot(
                  "arn:aws:ecs:us-east-1:1234:cluster/default",
                  Collections.emptyList(),
                  Collections.emptyList()))
          .getMock();
    }

    @Bean
    public LambdaFunction<SchedulerInput, SchedulerOutput> scheduler() {
      return new TestLambdaFunction<>(
          (input, context) ->
              new SchedulerOutput(
                  input.getSnapshot().getClusterArn(), input.getEnvironmentId(), 0L, 0L));
    }
  }
}
