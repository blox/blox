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

import lombok.Value;

@Value
public class ResourceExistsException extends ClientException {

  private String resourceType;
  private String resourceId;

  public ResourceExistsException(String resourceType, String resourceId) {
    super(String.format("%s with id %s already exists", resourceType, resourceId));

    this.resourceType = resourceType;
    this.resourceId = resourceId;
  }
}
