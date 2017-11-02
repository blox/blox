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
package com.amazonaws.blox.scheduling;

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.scheduling.LambdaHandlerTestCase.TestConfigOverrides;
import com.amazonaws.services.lambda.runtime.RequestStreamHandler;
import java.io.BufferedReader;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.nio.file.Paths;
import java.util.stream.Collectors;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.io.ClassPathResource;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.junit4.SpringRunner;

@RunWith(SpringRunner.class)
@ActiveProfiles("test")
@ContextConfiguration(classes = TestConfigOverrides.class)
public abstract class LambdaHandlerTestCase {

  @Autowired protected RequestStreamHandler handler;

  public String callHandler(String input) throws IOException {
    ByteArrayInputStream inputStream = new ByteArrayInputStream(input.getBytes());
    return callHandler(inputStream);
  }

  protected String callHandler(InputStream inputStream) throws IOException {
    ByteArrayOutputStream outputStream = new ByteArrayOutputStream();

    handler.handleRequest(inputStream, outputStream, null);

    return outputStream.toString();
  }

  protected InputStream fixture(String path) throws IOException {
    return new ClassPathResource(Paths.get("fixtures", path).toString()).getInputStream();
  }

  protected String fixtureAsString(String path) throws IOException {
    try (BufferedReader reader = new BufferedReader(new InputStreamReader(fixture(path)))) {
      return reader.lines().collect(Collectors.joining("\n"));
    }
  }

  @Configuration
  static class TestConfigOverrides {
    @Bean
    public DataService dataService() {
      return FakeDataService.builder().build();
    }
  }
}
