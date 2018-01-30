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
package com.amazonaws.blox.test;

import com.fasterxml.jackson.databind.JsonNode;
import java.util.function.Function;
import org.assertj.core.api.AbstractAssert;

public class JsonNodeAssert extends AbstractAssert<JsonNodeAssert, JsonNode> {
  public JsonNodeAssert(final JsonNode actual) {
    super(actual, JsonNodeAssert.class);
  }

  public static JsonNodeAssert assertThat(JsonNode actual) {
    return new JsonNodeAssert(actual);
  }

  public JsonNodeAssert isObject() {
    isNotNull();

    if (!actual.isObject()) {
      failWithMessage("Expected node to be an Object node, but was not");
    }

    return this;
  }

  public JsonNodeAssert hasNoField(String name) {
    if (actual.has(name)) {
      failWithMessage("Expected no field named <%s>, but found a field with that name.", name);
    }

    return this;
  }

  public JsonNodeAssert hasField(String name) {
    if (!actual.has(name)) {
      failWithMessage(
          "Expected a field named <%s>, but did not find a field with that name.", name);
    }

    return this;
  }

  public JsonNodeAssert hasField(String name, String expected) {
    return hasField(name, JsonNode::asText, expected);
  }

  public JsonNodeAssert hasField(String name, Integer expected) {
    return hasField(name, JsonNode::asInt, expected);
  }

  private <T> JsonNodeAssert hasField(String name, Function<JsonNode, T> extractor, T expected) {
    isObject();
    hasField(name);

    JsonNode node = actual.get(name);
    if (!extractor.apply(node).equals(expected)) {
      failWithMessage(
          "Expected field <%s> to have value <%s>, but was <%s>", name, expected, actual);
    }

    return this;
  }

  public JsonNodeAssert hasNullField(final String name) {
    isObject();
    hasField(name);

    JsonNode node = actual.get(name);
    if (!node.isNull()) {
      failWithMessage("Expected field <%s> to be null, but was <%s>", name, node);
    }

    return this;
  }
}
