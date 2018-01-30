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
package com.amazonaws.blox.jsonrpc.fixtures;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;

@Getter
public class PojoTestException extends Exception {
  private final String someValue;
  private final String anotherValue;

  @JsonCreator
  public PojoTestException(
      @JsonProperty("someValue") String someValue,
      @JsonProperty("anotherValue") String anotherValue) {
    this(someValue, anotherValue, null);
  }

  public PojoTestException(String someValue, String anotherValue, Throwable cause) {
    super(someValue + ":" + anotherValue, cause);
    this.someValue = someValue;
    this.anotherValue = anotherValue;
  }
}
