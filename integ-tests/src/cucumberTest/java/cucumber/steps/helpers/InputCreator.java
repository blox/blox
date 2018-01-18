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
package cucumber.steps.helpers;

import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentType;
import com.amazonaws.blox.dataservicemodel.v1.model.InstanceGroup;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DeleteEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.UpdateEnvironmentRequest;
import java.util.Collections;
import java.util.StringJoiner;
import java.util.UUID;
import lombok.Getter;
import lombok.Setter;

public class InputCreator {
  private static final String DEFAULT_ACCOUNT_ID = "123456789012";
  private static final String DEFAULT_CLUSTER_NAME = "Cluster";
  private static final String DEFAULT_ENVIRONMENT_NAME = "Environment";
  private static final String NAMING_PREFIX = "blox-integ-tests";

  private final String sharedId = UUID.randomUUID().toString();
  @Getter @Setter private String accountId = DEFAULT_ACCOUNT_ID;

  private String getTaskDefinitionArn() {
    return "arn:aws:ecs:us-east-1:" + getAccountId() + ":task-definition/sleep";
  }

  private String getRoleArn() {
    return "arn:aws:iam::" + getAccountId() + ":role/testRole";
  }

  public String prefixName(final String name) {
    return new StringJoiner("-").add(NAMING_PREFIX).add(sharedId).add(name).toString();
  }

  public CreateEnvironmentRequest createEnvironmentRequest() {
    return createEnvironmentRequest(DEFAULT_ENVIRONMENT_NAME, DEFAULT_CLUSTER_NAME);
  }

  public CreateEnvironmentRequest createEnvironmentRequest(final String environmentName) {
    return createEnvironmentRequest(environmentName, DEFAULT_CLUSTER_NAME);
  }

  public DescribeEnvironmentRequest describeEnvironmentRequest(final String environmentName) {
    return describeEnvironmentRequest(environmentName, DEFAULT_CLUSTER_NAME);
  }

  public CreateEnvironmentRequest createEnvironmentRequest(
      final String environmentName, final String cluster) {
    EnvironmentId id = environmentId(environmentName, cluster);
    return createEnvironmentRequest(id);
  }

  private CreateEnvironmentRequest createEnvironmentRequest(final EnvironmentId id) {
    return CreateEnvironmentRequest.builder()
        .environmentId(id)
        .role(getRoleArn())
        .taskDefinition(getTaskDefinitionArn())
        .environmentType(EnvironmentType.Daemon)
        .deploymentMethod("ReplaceAfterTerminate")
        .build();
  }

  private EnvironmentId environmentId(final String environmentName, final String cluster) {
    return EnvironmentId.builder()
        .accountId(getAccountId())
        .cluster(prefixName(cluster))
        .environmentName(prefixName(environmentName))
        .build();
  }

  private DescribeEnvironmentRequest describeEnvironmentRequest(
      final String environmentName, final String cluster) {
    EnvironmentId id = environmentId(environmentName, cluster);
    return describeEnvironmentRequest(id);
  }

  public DescribeEnvironmentRequest describeEnvironmentRequest(final EnvironmentId id) {
    return DescribeEnvironmentRequest.builder().environmentId(id).build();
  }

  public UpdateEnvironmentRequest updateEnvironmentRequestWithNewCluster(
      final String environmentName, final String cluster) {
    EnvironmentId id = environmentId(environmentName, cluster);
    return updateEnvironmentRequest(id);
  }

  private UpdateEnvironmentRequest updateEnvironmentRequest(final EnvironmentId id) {
    return UpdateEnvironmentRequest.builder()
        .environmentId(id)
        .role(getRoleArn())
        .taskDefinition(getTaskDefinitionArn())
        .instanceGroup(InstanceGroup.builder().attributes(Collections.emptySet()).build())
        .build();
  }

  public DeleteEnvironmentRequest deleteEnvironmentRequest(final EnvironmentId environmentId) {
    return DeleteEnvironmentRequest.builder()
        .environmentId(environmentId)
        .forceDelete(false)
        .build();
  }

  public DeleteEnvironmentRequest deleteEnvironmentRequest(
      final String environmentName, final String cluster) {
    return DeleteEnvironmentRequest.builder()
        .environmentId(environmentId(environmentName, cluster))
        .forceDelete(false)
        .build();
  }
}
