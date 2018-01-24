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
package com.amazonaws.blox.dataservice.test.data;

import com.amazonaws.blox.dataservice.model.DeploymentConfiguration;
import com.amazonaws.blox.dataservice.model.DeploymentConfiguration.DeploymentConfigurationBuilder;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.Environment.EnvironmentBuilder;
import com.amazonaws.blox.dataservice.model.EnvironmentHealth;
import com.amazonaws.blox.dataservice.model.EnvironmentId;
import com.amazonaws.blox.dataservice.model.EnvironmentId.EnvironmentIdBuilder;
import com.amazonaws.blox.dataservice.model.EnvironmentRevision;
import com.amazonaws.blox.dataservice.model.EnvironmentRevision.EnvironmentRevisionBuilder;
import com.amazonaws.blox.dataservice.model.EnvironmentStatus;
import com.amazonaws.blox.dataservice.model.EnvironmentType;
import java.time.Instant;
import lombok.Builder;
import lombok.Value;
import lombok.experimental.Accessors;

@Accessors(fluent = true)
@Builder
@Value
public class ModelBuilders {
  private String accountId = "0123456789012";
  private String clusterName = "Cluster";
  private String environmentName = "Environment";
  private String environmentRevisionId = "EnvironmentRevisionId";
  private String deploymentMethod = "DeploymentMethod";
  private String environmentRole = "Role";
  private EnvironmentStatus environmentStatus = EnvironmentStatus.Active;
  private EnvironmentHealth environmentHealth = EnvironmentHealth.Healthy;
  private EnvironmentType environmentType = EnvironmentType.Daemon;
  private String taskDefinition = "TaskDefinition";
  private Instant now = Instant.now();

  public EnvironmentRevisionBuilder environmentRevision() {
    return EnvironmentRevision.builder()
        .environmentId(environmentId().build())
        .environmentRevisionId(environmentRevisionId())
        .taskDefinition(taskDefinition())
        .createdTime(now());
  }

  public EnvironmentBuilder environment() {
    return Environment.builder()
        .environmentId(environmentId().build())
        .environmentType(environmentType())
        .environmentHealth(environmentHealth())
        .environmentStatus(environmentStatus())
        .role(environmentRole())
        .createdTime(now())
        .lastUpdatedTime(now())
        .deploymentConfiguration(deploymentConfiguration().build())
        .deploymentMethod(deploymentMethod())
        .latestEnvironmentRevisionId(environmentRevisionId())
        .activeEnvironmentRevisionId(environmentRevisionId());
  }

  private DeploymentConfigurationBuilder deploymentConfiguration() {
    return DeploymentConfiguration.builder();
  }

  public EnvironmentIdBuilder environmentId() {
    return EnvironmentId.builder()
        .accountId(accountId)
        .cluster(clusterName)
        .environmentName(environmentName);
  }

  public EnvironmentId environmentId(
      final String accountId, final String cluster, final String environmentName) {
    return environmentId()
        .accountId(accountId)
        .cluster(cluster)
        .environmentName(environmentName)
        .build();
  }
}
