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

import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;

import com.amazonaws.blox.lambda.LambdaFunction;
import com.amazonaws.blox.lambda.TestLambdaFunction;
import com.amazonaws.blox.scheduling.LambdaHandlerTestCase;
import com.amazonaws.blox.scheduling.manager.ManagerInput;
import com.amazonaws.blox.scheduling.manager.ManagerOutput;
import com.amazonaws.blox.scheduling.reconciler.ReconcilerEntrypointTest.TestConfig;
import java.util.Collections;
import org.junit.Test;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Import;
import org.springframework.test.context.ContextConfiguration;

@ContextConfiguration(classes = TestConfig.class)
public class ReconcilerEntrypointTest extends LambdaHandlerTestCase {

  @Test
  public void convertsInputsAndOutputsFromJson() throws Exception {
    String result = callHandler(fixture("handlers/Reconciler.input.json"));
    assertThat(result, is("null"));
  }

  @Configuration
  @Import(ReconcilerApplication.class)
  public static class TestConfig {

    @Bean
    public LambdaFunction<ManagerInput, ManagerOutput> manager() {
      return new TestLambdaFunction<>(
          (input, context) -> new ManagerOutput(input.getClusterArn(), Collections.emptyList()));
    }
  }
}
