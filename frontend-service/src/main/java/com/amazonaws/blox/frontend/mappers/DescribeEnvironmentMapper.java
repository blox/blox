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
package com.amazonaws.blox.frontend.mappers;

import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentResponse;
import com.amazonaws.blox.frontend.operations.DescribeEnvironment;
import com.amazonaws.serverless.proxy.internal.model.ApiGatewayRequestContext;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(uses = {EnvironmentMapper.class, EnvironmentRevisionMapper.class})
public interface DescribeEnvironmentMapper {
  @Mapping(target = "environmentId.accountId", source = "context.accountId")
  @Mapping(target = "environmentId.cluster", source = "cluster")
  @Mapping(target = "environmentId.environmentName", source = "environmentName")
  DescribeEnvironmentRequest toDataServiceRequest(
      ApiGatewayRequestContext context, String cluster, String environmentName);

  DescribeEnvironment.DescribeEnvironmentResponse fromDataServiceResponse(
      DescribeEnvironmentResponse response);
}
