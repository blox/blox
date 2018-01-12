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

import java.util.UUID;

public class InputCreator {
  private static final String ENVIRONMENT_NAME_PREFIX = "test";
  private static final int ACCOUNT_ID_SIZE = 12;
  private static final String ACCOUNT_ID = generateAccountId();
  private static final String TASK_DEFINITION_ARN =
      "arn:aws:ecs:us-east-1:" + ACCOUNT_ID + ":task-definition/sleep";
  private static final String ROLE_ARN = "arn:aws:iam::" + ACCOUNT_ID + ":role/testRole";
  private static final String CLUSTER_NAME_PREFIX = "cluster";

  public CreateEnvironmentRequest createEnvironmentRequest() {
    final String environmentName = ENVIRONMENT_NAME_PREFIX + UUID.randomUUID().toString();
    final String clusterName = CLUSTER_NAME_PREFIX + UUID.randomUUID().toString();
    return createEnvironmentRequest(environmentName, clusterName);
  }

  public CreateEnvironmentRequest createEnvironmentRequest(final String environmentNamePrefix) {
    final String environmentName = environmentNamePrefix + UUID.randomUUID().toString();
    final String clusterName = CLUSTER_NAME_PREFIX + UUID.randomUUID().toString();
    return createEnvironmentRequest(environmentName, clusterName);
  }

  public DescribeEnvironmentRequest describeEnvironmentRequest(final String environmentNamePrefix) {
    final String environmentName = environmentNamePrefix + UUID.randomUUID().toString();
    final String clusterName = CLUSTER_NAME_PREFIX + UUID.randomUUID().toString();
    return describeEnvironmentRequest(environmentName, clusterName);
  }

  public CreateEnvironmentRequest createEnvironmentRequest(
      final String environmentName, final String cluster) {
    return CreateEnvironmentRequest.builder()
        .environmentId(
            EnvironmentId.builder()
                .accountId(ACCOUNT_ID)
                .cluster(cluster)
                .environmentName(environmentName)
                .build())
        .role(ROLE_ARN)
        .taskDefinition(TASK_DEFINITION_ARN)
        .environmentType(EnvironmentType.Daemon)
        .deploymentMethod("ReplaceAfterTerminate")
        .build();
  }

  public DescribeEnvironmentRequest describeEnvironmentRequest(
      final String environmentName, final String cluster) {
    return DescribeEnvironmentRequest.builder()
        .environmentId(
            EnvironmentId.builder()
                .accountId(ACCOUNT_ID)
                .cluster(cluster)
                .environmentName(environmentName)
                .build())
        .build();
  }

  public UpdateEnvironmentRequest updateEnvironmentRequestWithNewCluster(
      final String environmentName, final String cluster) {
    return UpdateEnvironmentRequest.builder()
        .environmentId(
            EnvironmentId.builder()
                .accountId(ACCOUNT_ID)
                .cluster(cluster)
                .environmentName(environmentName)
                .build())
        .role(ROLE_ARN)
        .taskDefinition(TASK_DEFINITION_ARN)
        .instanceGroup(InstanceGroup.builder().attributes(null).build())
        .build();
  }

  public DeleteEnvironmentRequest deleteEnvironmentRequest(
      final String environmentName, final String cluster) {
    return DeleteEnvironmentRequest.builder()
        .environmentId(
            EnvironmentId.builder()
                .accountId(ACCOUNT_ID)
                .cluster(cluster)
                .environmentName(environmentName)
                .build())
        .forceDelete(false)
        .build();
  }

  //TODO: random
  private static String generateAccountId() {
    return "123456789012";
  }
}
