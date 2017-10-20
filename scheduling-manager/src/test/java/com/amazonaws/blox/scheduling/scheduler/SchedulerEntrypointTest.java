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

import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.mockito.Matchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

import com.amazonaws.blox.scheduling.LambdaHandlerTestCase;
import com.amazonaws.blox.scheduling.scheduler.SchedulerEntrypointTest.TestConfig;
import com.amazonaws.blox.scheduling.scheduler.engine.SchedulerFactory;
import java.util.Collections;
import org.junit.Test;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Import;
import org.springframework.test.context.ContextConfiguration;
import software.amazon.awssdk.services.ecs.ECSAsyncClient;

@ContextConfiguration(classes = TestConfig.class)
public class SchedulerEntrypointTest extends LambdaHandlerTestCase {

  @Test
  public void convertsInputsAndOutputsFromJson() throws Exception {
    String result = callHandler(fixture("handlers/Scheduler.input.json"));
    assertThat(result, is(fixtureAsString("handlers/Scheduler.output.json")));
  }

  @Configuration
  @Import(SchedulerApplication.class)
  public static class TestConfig {
    @Bean
    public ECSAsyncClient ecs() {
      return mock(ECSAsyncClient.class);
    }

    @Bean
    public SchedulerFactory schedulerFactory() {
      return when(mock(SchedulerFactory.class).schedulerFor(any()))
          .thenReturn((snapshot, deploymentConfiguration) -> Collections.emptyList())
          .getMock();
    }
  }
}
