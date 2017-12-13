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
package com.amazonaws.blox.frontend.models;

import io.swagger.annotations.ApiParam;
import lombok.Data;
import lombok.Getter;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.RequestParam;

@Data
public class RequestPagination {
  private String nextToken;
  private Long maxResults;

  @ApiParam(name = "nextToken")
  public void setNextToken(@RequestParam("nextToken") String value) {
    this.nextToken = value;
  }

  @ApiParam(name = "maxResults")
  public void setMaxResults(@RequestParam("maxResults") Long value) {
    this.maxResults = value;
  }
}
