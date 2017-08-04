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
package com.amazonaws.blox;

import static org.junit.Assert.assertEquals;

import com.amazonaws.blox.model.DescribeEnvironmentRequest;
import com.amazonaws.blox.model.DescribeEnvironmentResult;
import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.File;
import java.net.URL;
import org.junit.Test;

public class EnvironmentTest {
  private Blox client = Blox.builder().endpoint(endpoint()).build();

  @Test
  public void describeEnvironmentReturnsFakeEnvironment() throws Exception {
    DescribeEnvironmentResult result =
        client.describeEnvironment(new DescribeEnvironmentRequest().name("foo"));
    assertEquals("foo", result.getEnvironment().getName());
  }

  private static String endpoint() {
    try {
      JsonNode tree =
          new ObjectMapper(new JsonFactory())
              .readTree(new File(System.getProperty("blox.tests.stackoutputs")));
      String urlString = tree.get("ApiUrl").asText();

      // The generated client doesn't like the stage name in the endpoint, so strip it out:
      URL url = new URL(urlString);
      return new URL(url.getProtocol(), url.getHost(), url.getPort(), "/").toString();
    } catch (Exception e) {
      throw new RuntimeException(e);
    }
  }
}
