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

import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import lombok.AllArgsConstructor;

@AllArgsConstructor
public class Encoder {

  private ObjectMapper mapper;

  public <T> T decode(InputStream input, Class<T> clazz) throws IOException {
    return mapper.readValue(input, clazz);
  }

  public <T> void encode(OutputStream output, T value) throws IOException {
    mapper.writeValue(output, value);
  }
}
