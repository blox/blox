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
package com.amazonaws.blox.dataservice.mapper;

import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentVersion;
import com.amazonaws.blox.dataservicemodel.v1.model.Attribute;
import com.amazonaws.blox.dataservicemodel.v1.model.DeploymentConfiguration;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentType;
import com.amazonaws.blox.dataservicemodel.v1.model.InstanceGroup;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateTargetEnvironmentRevisionResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeTargetEnvironmentRevisionResponse;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper
public interface ApiModelMapper {

  @Mapping(target = "environmentVersion", ignore = true)
  @Mapping(target = "status", ignore = true)
  @Mapping(target = "health", ignore = true)
  @Mapping(target = "createdTime", ignore = true)
  @Mapping(target = "lastUpdatedTime", ignore = true)
  @Mapping(source = "environmentType", target = "type")
  Environment toEnvironment(CreateEnvironmentRequest createEnvironmentRequest);

  @Mapping(source = "type", target = "environmentType")
  @Mapping(source = "health", target = "environmentHealth")
  @Mapping(source = "status", target = "environmentStatus")
  CreateEnvironmentResponse toCreateEnvironmentResponse(Environment environment);

  CreateTargetEnvironmentRevisionResponse toCreateTargetEnvironmentRevisionResponse(
      EnvironmentVersion environmentVersion);

  DescribeEnvironmentResponse toDescribeEnvironmentResponse(Environment environment);

  DescribeTargetEnvironmentRevisionResponse toDescribeTargetEnvironmentRevisionResponse(
      EnvironmentVersion environmentVersion);

  EnvironmentType toWrapperEnvironmentType(
      com.amazonaws.blox.dataservice.model.EnvironmentType environmentType);

  com.amazonaws.blox.dataservice.model.EnvironmentType toModelEnvironmentType(
      EnvironmentType environmentType);

  InstanceGroup toWrapperInstanceGroup(
      com.amazonaws.blox.dataservice.model.InstanceGroup instanceGroup);

  com.amazonaws.blox.dataservice.model.InstanceGroup toModelInstanceGroup(
      InstanceGroup instanceGroup);

  Attribute toWrapperAttribute(com.amazonaws.blox.dataservice.model.Attribute attribute);

  com.amazonaws.blox.dataservice.model.Attribute toModelAttribute(Attribute attribute);

  DeploymentConfiguration toWrapperDeploymentConfiguration(
      com.amazonaws.blox.dataservice.model.DeploymentConfiguration deploymentConfiguration);

  com.amazonaws.blox.dataservice.model.DeploymentConfiguration toModelDeploymentConfiguration(
      DeploymentConfiguration deploymentConfiguration);
}
