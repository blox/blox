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
package com.amazonaws.blox.frontend.operations;

import com.amazonaws.blox.frontend.models.PaginatedResponse;
import com.amazonaws.blox.frontend.models.RequestPagination;
import io.swagger.annotations.ApiOperation;
import java.util.Collections;
import java.util.List;
import lombok.Builder;
import lombok.Value;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class ListEnvironments extends EnvironmentController {

  @RequestMapping(method = RequestMethod.GET)
  @ApiOperation(value = "List all environments")
  public ListEnvironmentsResponse listEnvironments(
      @PathVariable("cluster") String cluster,
      @RequestParam(value = "environmentNamePrefix", required = false) String environmentNamePrefix,
      @ModelAttribute RequestPagination pagination) {

    return ListEnvironmentsResponse.builder().environmentNames(Collections.emptyList()).build();
  }

  @Value
  @Builder
  public static class ListEnvironmentsResponse implements PaginatedResponse {

    private final List<String> environmentNames;
    private final String nextToken;
  }
}
