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
package com.amazonaws.blox.dataservice.model;

import java.time.Instant;
import lombok.Builder;
import lombok.Data;
import lombok.NonNull;

@Data
@Builder
public class Environment {

  @NonNull private String accountId;
  @NonNull private String environmentName;
  private String environmentVersion;
  private String taskDefinitionArn;
  private String roleArn;
  private InstanceGroup instanceGroup;
  private EnvironmentType type;
  private EnvironmentStatus status;
  private EnvironmentHealth health;
  //TODO: add to ddb record
  private DeploymentConfiguration deploymentConfiguration;
  private Instant createdTime;
  private Instant lastUpdatedTime;
}
