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
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.NonNull;

@Data
@Builder
// required for builder
@AllArgsConstructor
// required for mapstruct
@NoArgsConstructor
public class Environment {

  @NonNull private EnvironmentId environmentId;
  @NonNull private String role;
  @NonNull private EnvironmentType environmentType;
  @NonNull private Instant createdTime;
  @NonNull private Instant lastUpdatedTime;
  @NonNull private EnvironmentHealth environmentHealth;
  @NonNull private EnvironmentStatus environmentStatus;
  @NonNull private String latestEnvironmentRevisionId;
  // environment and revision were both successfully created
  @NonNull private boolean validEnvironment;
  private DeploymentConfiguration deploymentConfiguration;
  private String activeEnvironmentRevisionId;
  @NonNull private String deploymentMethod;
}
