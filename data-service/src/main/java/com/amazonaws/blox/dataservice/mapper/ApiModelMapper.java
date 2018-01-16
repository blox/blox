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

import com.amazonaws.blox.dataservice.model.Attribute;
import com.amazonaws.blox.dataservice.model.Cluster;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentId;
import com.amazonaws.blox.dataservice.model.EnvironmentRevision;
import com.amazonaws.blox.dataservice.model.EnvironmentType;
import com.amazonaws.blox.dataservice.model.DeploymentConfiguration;
import com.amazonaws.blox.dataservice.model.InstanceGroup;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentRequest;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper
public interface ApiModelMapper {

  @Mapping(target = "createdTime", ignore = true)
  @Mapping(target = "lastUpdatedTime", ignore = true)
  @Mapping(target = "environmentStatus", ignore = true)
  @Mapping(target = "environmentHealth", ignore = true)
  @Mapping(target = "activeEnvironmentRevisionId", ignore = true)
  @Mapping(target = "latestEnvironmentRevisionId", ignore = true)
  @Mapping(target = "validEnvironment", ignore = true)
  Environment toEnvironment(CreateEnvironmentRequest createEnvironmentRequest);

  com.amazonaws.blox.dataservicemodel.v1.model.Environment toWrapperEnvironment(
      Environment environment);

  Cluster toModelCluster(com.amazonaws.blox.dataservicemodel.v1.model.Cluster cluster);

  EnvironmentId toModelEnvironmentId(
      com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId environmentId);

  com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId toWrapperEnvironmentId(
      EnvironmentId environmentId);

  com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentRevision toWrapperEnvironmentRevision(
      EnvironmentRevision environmentRevision);

  EnvironmentType toModelEnvironmentType(
      com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentType environmentType);

  com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentType toWrapperEnvironmentType(
      EnvironmentType environmentType);

  InstanceGroup toModelInstanceGroup(
      com.amazonaws.blox.dataservicemodel.v1.model.InstanceGroup instanceGroup);

  com.amazonaws.blox.dataservicemodel.v1.model.InstanceGroup toWrapperInstanceGroup(
      InstanceGroup instanceGroup);

  Attribute toModelAttribute(com.amazonaws.blox.dataservicemodel.v1.model.Attribute attribute);

  com.amazonaws.blox.dataservicemodel.v1.model.Attribute toWrapperAttribute(Attribute attribute);

  DeploymentConfiguration toModelDeploymentConfiguration(
      com.amazonaws.blox.dataservicemodel.v1.model.DeploymentConfiguration deploymentConfiguration);

  com.amazonaws.blox.dataservicemodel.v1.model.DeploymentConfiguration
      toWrapperDeploymentConfiguration(DeploymentConfiguration deploymentConfiguration);
}
