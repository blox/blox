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
package com.amazonaws.blox.dataservice.integration;

import com.amazonaws.blox.dataservicemodel.v1.model.Cluster;
import com.amazonaws.blox.dataservicemodel.v1.model.Cluster.ClusterBuilder;
import com.amazonaws.blox.dataservicemodel.v1.model.DeploymentConfiguration;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId.EnvironmentIdBuilder;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentType;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.*;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentRequest.CreateEnvironmentRequestBuilder;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DeleteEnvironmentRequest.DeleteEnvironmentRequestBuilder;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRequest.DescribeEnvironmentRequestBuilder;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRevisionRequest.DescribeEnvironmentRevisionRequestBuilder;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListClustersRequest.ListClustersRequestBuilder;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsRequest.ListEnvironmentsRequestBuilder;
import java.time.Instant;
import lombok.Builder;
import lombok.Value;
import lombok.experimental.Accessors;

@Accessors(fluent = true)
@Value
@Builder
public class DataServiceModelBuilder {
  private String accountId = "123456789012";
  private String clusterName = "Cluster";
  private String environmentName = "Environment";
  private String environmentRevisionId = "EnvironmentRevisionId";
  private String deploymentMethod = "DeploymentMethod";
  private String environmentRole = "Role";
  private String taskDefinition = "TaskDefinition";
  private EnvironmentType environmentType = EnvironmentType.Daemon;
  private Instant now = Instant.now();

  public CreateEnvironmentRequestBuilder createEnvironmentRequest() {
    return CreateEnvironmentRequest.builder()
        .taskDefinition(taskDefinition)
        .role(environmentRole)
        .deploymentConfiguration(DeploymentConfiguration.builder().build())
        .deploymentMethod(deploymentMethod)
        .environmentType(environmentType)
        .environmentId(environmentId().build());
  }

  public DescribeEnvironmentRequestBuilder describeEnvironmentRequest() {
    return DescribeEnvironmentRequest.builder().environmentId(environmentId().build());
  }

  public DescribeEnvironmentRevisionRequestBuilder describeEnvironmentRevisionRequest() {
    return DescribeEnvironmentRevisionRequest.builder()
        .environmentId(environmentId().build())
        .environmentRevisionId(environmentRevisionId);
  }

  public ListClustersRequestBuilder listClustersRequest() {
    return ListClustersRequest.builder().accountId(accountId).clusterNamePrefix(clusterName);
  }

  public ListEnvironmentsRequestBuilder listEnvironmentsRequest() {
    return ListEnvironmentsRequest.builder()
        .cluster(cluster().build())
        .environmentNamePrefix(environmentName);
  }

  public DeleteEnvironmentRequestBuilder deleteEnvironmentRequest() {
    return DeleteEnvironmentRequest.builder()
        .forceDelete(false)
        .environmentId(environmentId().build());
  }

  public EnvironmentIdBuilder environmentId() {
    return EnvironmentId.builder()
        .accountId(accountId)
        .cluster(clusterName)
        .environmentName(environmentName);
  }

  public ClusterBuilder cluster() {
    return Cluster.builder().accountId(accountId).clusterName(clusterName);
  }
}
