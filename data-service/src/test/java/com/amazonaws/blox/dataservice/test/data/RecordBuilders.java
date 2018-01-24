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

import com.amazonaws.blox.dataservice.model.EnvironmentHealth;
import com.amazonaws.blox.dataservice.model.EnvironmentStatus;
import com.amazonaws.blox.dataservice.model.EnvironmentType;
import com.amazonaws.blox.dataservice.repository.model.EnvironmentDDBRecord;
import com.amazonaws.blox.dataservice.repository.model.EnvironmentDDBRecord.EnvironmentDDBRecordBuilder;
import java.time.Instant;
import lombok.Builder;
import lombok.Value;
import lombok.experimental.Accessors;

@Accessors(fluent = true)
@Builder
@Value
public class RecordBuilders {

  private String accountId = "accountId";
  private String cluster = "cluster";
  private String environmentName = "environmentName";
  private EnvironmentType environmentType = EnvironmentType.Daemon;
  private EnvironmentHealth environmentHealth = EnvironmentHealth.Healthy;
  private EnvironmentStatus environmentStatus = EnvironmentStatus.Active;
  private String role = "role";
  private Instant instant = Instant.now();
  private String deploymentMethod = "deploymentMethod";
  private String environmentRevisionId = "";

  public String accountIdCluster() {
    return String.join("/", accountId(), cluster());
  }

  public EnvironmentDDBRecordBuilder environment() {
    return EnvironmentDDBRecord.builder()
        .accountIdCluster(accountIdCluster())
        .accountId(accountId())
        .clusterName(cluster())
        .environmentName(environmentName())
        .type(environmentType())
        .health(environmentHealth())
        .status(environmentStatus())
        .role(role())
        .createdTime(instant())
        .lastUpdatedTime(instant())
        .deploymentMethod(deploymentMethod())
        .latestEnvironmentRevisionId(environmentRevisionId())
        .activeEnvironmentRevisionId(environmentRevisionId());
  }
}
