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
package com.amazonaws.blox.schedulingmanager.handler;

import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@AllArgsConstructor
public class Encoder {

  private ObjectMapper mapper;

  public <T> T decode(InputStream input, Class<T> clazz) throws IOException {
    try {
      return mapper.readValue(input, clazz);
    } catch (final IOException e) {
      log.error("Could not parse input into the expected class {}", clazz.getName());
      throw e;
    }
  }

  public <T> void encode(OutputStream output, T value) throws IOException {
    mapper.writeValue(output, value);
  }

  public <T> String encode(T value) throws IOException {
    return mapper.writeValueAsString(value);
  }
}
