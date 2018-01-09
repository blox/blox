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
package cucumber.steps.dataservice;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentHealth;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentStatus;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentType;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentResponse;
import cucumber.configuration.CucumberConfiguration;
import cucumber.api.PendingException;
import cucumber.api.java8.En;
import java.util.UUID;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.context.ContextConfiguration;
import cucumber.steps.wrappers.DataServiceWrapper;

@ContextConfiguration(classes = CucumberConfiguration.class)
public class CreateEnvironmentSteps implements En {

  private static final String ENVIRONMENT_NAME_PREFIX = "test";
  private static final int ACCOUNT_ID_SIZE = 12;
  private static final String ACCOUNT_ID = generateAccountId();
  private static final String TASK_DEFINITION_ARN =
      "arn:aws:ecs:us-east-1:" + ACCOUNT_ID + ":task-definition/sleep";
  private static final String ROLE_ARN = "arn:aws:iam::" + ACCOUNT_ID + ":role/testRole";
  private static final String CLUSTER_NAME_PREFIX = "cluster";

  @Autowired private DataServiceWrapper dataServiceWrapper;

  public CreateEnvironmentSteps() {
    When(
        "^I create an environment$",
        () -> {
          dataServiceWrapper.createEnvironment(createEnvironmentRequest());
        });

    Then(
        "the created environment response is valid",
        () -> {
          final CreateEnvironmentRequest request =
              dataServiceWrapper.getLastFromHistory(CreateEnvironmentRequest.class);
          final CreateEnvironmentResponse response =
              dataServiceWrapper.getLastFromHistory(CreateEnvironmentResponse.class);

          assertNotNull(response);
          assertNotNull(response.getEnvironment());
          assertNotNull(response.getEnvironmentRevision());

          assertEquals(request.getEnvironmentId(), response.getEnvironment().getEnvironmentId());
          assertEquals(request.getRole(), response.getEnvironment().getRole());
          assertEquals(
              request.getEnvironmentType(), response.getEnvironment().getEnvironmentType());
          assertEquals(
              request.getDeploymentConfiguration(),
              response.getEnvironment().getDeploymentConfiguration());

          assertNotNull(response.getEnvironment().getCreatedTime());
          assertNotNull(response.getEnvironment().getLatestEnvironmentRevisionId());
          assertEquals(EnvironmentHealth.HEALTHY, response.getEnvironment().getEnvironmentHealth());
          assertEquals(
              EnvironmentStatus.INACTIVE, response.getEnvironment().getEnvironmentStatus());

          assertEquals(
              request.getEnvironmentId(), response.getEnvironmentRevision().getEnvironmentId());
          assertEquals(
              request.getTaskDefinition(), response.getEnvironmentRevision().getTaskDefinition());
          assertEquals(
              request.getInstanceGroup(), response.getEnvironmentRevision().getInstanceGroup());

          assertNotNull(response.getEnvironmentRevision().getCreatedTime());
          assertNotNull(response.getEnvironmentRevision().getEnvironmentRevisionId());
        });

    When(
        "^I create an environment named \"([^\"]*)\"$",
        (String environmentName) -> {
          dataServiceWrapper.createEnvironment(createEnvironmentRequest(environmentName));
        });

    Given(
        "^I create an environment named \"([^\"]*)\" in the cluster \"([^\"]*)\"$",
        (String arg1, String arg2) -> {
          throw new PendingException();
        });

    When(
        "^I try to create another environment with the name \"([^\"]*)\" in the cluster \"([^\"]*)\"$",
        (String arg1, String arg2) -> {
          throw new PendingException();
        });

    Then(
        "^there should be a ResourceExistsException thrown$",
        () -> {
          throw new PendingException();
        });

    Then(
        "^the resourceType is \"([^\"]*)\"$",
        (String arg1) -> {
          throw new PendingException();
        });

    Then(
        "^the resourceId contains \"([^\"]*)\"$",
        (String arg1) -> {
          throw new PendingException();
        });
  }

  private CreateEnvironmentRequest createEnvironmentRequest() {
    final String environmentName = ENVIRONMENT_NAME_PREFIX + UUID.randomUUID().toString();
    final String clusterName = CLUSTER_NAME_PREFIX + UUID.randomUUID().toString();
    return createEnvironmentRequest(environmentName, clusterName);
  }

  private CreateEnvironmentRequest createEnvironmentRequest(final String environmentNamePrefix) {
    final String environmentName = environmentNamePrefix + UUID.randomUUID().toString();
    final String clusterName = CLUSTER_NAME_PREFIX + UUID.randomUUID().toString();
    return createEnvironmentRequest(environmentName, clusterName);
  }

  private CreateEnvironmentRequest createEnvironmentRequest(
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
        .build();
  }

  //TODO: random
  private static String generateAccountId() {
    return "12345678912";
  }
}
