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
package com.amazonaws.blox.dataservicemodel.v1.exception;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;

@Getter
public class ResourceInUseException extends ClientException {

  private String resourceType;
  private String resourceId;

  @JsonCreator
  public ResourceInUseException(
      @JsonProperty("resourceType") String resourceType,
      @JsonProperty("resourceId") String resourceId,
      @JsonProperty("message") String message) {
    super(message);

    this.resourceType = resourceType;
    this.resourceId = resourceId;
  }
}
