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
package com.amazonaws.blox.lambda;

import static org.hamcrest.CoreMatchers.is;
import static org.junit.Assert.assertThat;

import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import com.amazonaws.services.lambda.runtime.RequestStreamHandler;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import lombok.Data;
import org.junit.Test;

public class JacksonRequestStreamHandlerTest {

  private final ObjectMapper mapper = new ObjectMapper();

  @Test
  public void roundtripsInputsAndOutputs() throws Exception {
    RequestStreamHandler handler = new JacksonRequestStreamHandler<>(mapper, new TestHandler());

    ByteArrayOutputStream outputStream = new ByteArrayOutputStream();
    handler.handleRequest(
        new ByteArrayInputStream("{\"input\":\"input\"}".getBytes()), outputStream, null);

    assertThat(outputStream.toString(), is("{\"output\":\"input\"}"));
  }

  abstract static class BaseHandler implements RequestHandler<FakeInput, FakeOutput> {}

  static class TestHandler extends BaseHandler {

    @Override
    public FakeOutput handleRequest(FakeInput input, Context context) {
      return new FakeOutput(input.getInput());
    }
  }

  @Data
  static class FakeInput {

    private final String input;
  }

  @Data
  static class FakeOutput {

    private final String output;
  }
}
