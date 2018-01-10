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

import com.amazonaws.blox.dataservicemodel.v1.model.Environment;
import com.amazonaws.serverless.proxy.internal.model.ApiGatewayRequestContext;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(uses = DeploymentConfigurationMapper.class)
public interface EnvironmentMapper {
  @Mapping(target = "environmentId.accountId", source = "context.accountId")
  @Mapping(target = "environmentId.cluster", source = "environment.cluster")
  @Mapping(target = "environmentId.environmentName", source = "environment.environmentName")
  // TODO: Add timestamps and status to frontend
  @Mapping(target = "createdTime", ignore = true)
  @Mapping(target = "lastUpdatedTime", ignore = true)
  @Mapping(target = "environmentStatus", ignore = true)
  Environment toDataService(
      ApiGatewayRequestContext context, com.amazonaws.blox.frontend.models.Environment environment);

  @Mapping(target = "cluster", source = "environmentId.cluster")
  @Mapping(target = "environmentName", source = "environmentId.environmentName")
  com.amazonaws.blox.frontend.models.Environment fromDataService(Environment environment);
}
