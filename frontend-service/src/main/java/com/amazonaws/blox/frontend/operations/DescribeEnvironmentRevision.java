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

import com.amazonaws.blox.frontend.models.EnvironmentRevision;
import io.swagger.annotations.ApiOperation;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class DescribeEnvironmentRevision extends EnvironmentController {
  @RequestMapping(path = "/{environmentName}/revisions/{revisionId}", method = RequestMethod.GET)
  @ApiOperation(value = "Describe Environment revisions")
  public DescribeEnvironmentRevisionResponse describeEnvironmentRevision(
      @PathVariable("cluster") String cluster,
      @PathVariable("environmentName") String environmentName,
      @PathVariable("revisionId") String revisionId) {

    return DescribeEnvironmentRevisionResponse.builder()
        .environmentRevision(
            EnvironmentRevision.builder()
                .environmentRevisionId(revisionId)
                .taskDefinition(null)
                .build())
        .build();
  }

  @Data
  @Builder
  // required for builder
  @AllArgsConstructor
  // required for mapstruct
  @NoArgsConstructor
  public static class DescribeEnvironmentRevisionResponse {
    private EnvironmentRevision environmentRevision;
  }
}
