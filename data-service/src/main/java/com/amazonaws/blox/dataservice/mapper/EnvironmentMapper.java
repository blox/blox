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
import com.amazonaws.blox.dataservice.repository.model.EnvironmentDDBRecord;
import com.amazonaws.blox.dataservice.repository.model.EnvironmentTargetVersionDDBRecord;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper
public interface EnvironmentMapper {

  @Mapping(source = "instanceGroup.cluster", target = "cluster")
  // TODO: Map attributes
  @Mapping(target = "attributes", ignore = true)
  @Mapping(target = "recordVersion", ignore = true)
  EnvironmentDDBRecord toEnvironmentDDBRecord(Environment environment);

  @Mapping(source = "cluster", target = "instanceGroup.cluster")
  //TODO: add to record
  @Mapping(target = "deploymentConfiguration", ignore = true)
  Environment toEnvironment(EnvironmentDDBRecord environmentDDBRecord);

  @Mapping(target = "recordVersion", ignore = true)
  EnvironmentTargetVersionDDBRecord toEnvironmentTargetVersionDDBRecord(
      EnvironmentVersion environmentVersion);

  EnvironmentVersion toEnvironmentVersion(
      EnvironmentTargetVersionDDBRecord environmentTargetVersionDDBRecord);
}
