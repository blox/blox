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
package com.amazonaws.blox.swagger;

import io.swagger.models.Swagger;
import java.util.TreeMap;
import lombok.Getter;
import lombok.Setter;

/**
 * Swagger model uses LinkedHashMap for data like paths, definitions. Due to no order guarantee, it
 * causes problem that the swagger.yml generated from the model may var even no swagger model change
 */
@Getter
@Setter
public class SortedModelFilter implements SwaggerFilter {

  @Override
  public void apply(Swagger swagger) {
    if (swagger.getDefinitions() != null) {
      swagger.setDefinitions(new TreeMap<>(swagger.getDefinitions()));
    }

    if (swagger.getPaths() != null) {
      swagger.setPaths(new TreeMap<>(swagger.getPaths()));
    }
  }
}
