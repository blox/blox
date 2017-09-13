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
package com.amazonaws.blox.schedulingmanager.deployment.handler;

import static org.junit.Assert.assertEquals;

import com.amazonaws.blox.schedulingmanager.deployment.steps.types.DeploymentInput;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import org.junit.Before;
import org.junit.Test;

public class EncoderTest {

  private Encoder encoder;
  private ObjectMapper mapper;

  @Before
  public void setup() {
    mapper = new ObjectMapper();
    encoder = new Encoder(mapper);
  }

  @Test
  public void testDecode() throws IOException {
    final DeploymentInput deploymentInput =
        DeploymentInput.builder().account("1234").name("name1").build();

    final String inputJson = mapper.writeValueAsString(deploymentInput);
    final InputStream inputStream = new ByteArrayInputStream(inputJson.getBytes());

    final DeploymentInput result = encoder.decode(inputStream, DeploymentInput.class);
    assertEquals(result.getAccount(), deploymentInput.getAccount());
    assertEquals(result.getName(), deploymentInput.getName());
  }
}
