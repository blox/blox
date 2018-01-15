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
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentStatus;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentResponse;
import cucumber.configuration.CucumberConfiguration;
import cucumber.api.PendingException;
import cucumber.api.java8.En;
import cucumber.steps.helpers.InputCreator;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.context.ContextConfiguration;

import cucumber.steps.wrappers.DataServiceWrapper;

@ContextConfiguration(classes = CucumberConfiguration.class)
public class CreateEnvironmentSteps implements En {
  @Autowired private DataServiceWrapper dataServiceWrapper;
  @Autowired private InputCreator inputCreator;

  public CreateEnvironmentSteps() {
    When(
        "^I create an environment$",
        () -> {
          dataServiceWrapper.createEnvironment(inputCreator.createEnvironmentRequest());
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
          assertNotNull(response.getEnvironment().getDeploymentMethod());
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
          dataServiceWrapper.createEnvironment(
              inputCreator.createEnvironmentRequest(environmentName));
        });

    Given(
        "^I create an environment named \"([^\"]*)\" in the cluster \"([^\"]*)\"$",
        (final String environmentName, final String clusterName) -> {
          dataServiceWrapper.createEnvironment(
              inputCreator.createEnvironmentRequest(environmentName, clusterName));
        });

    When(
        "^I try to create another environment with the name \"([^\"]*)\" in the cluster \"([^\"]*)\"$",
        (final String environmentName, final String clusterName) -> {
          dataServiceWrapper.tryCreateEnvironment(
              inputCreator.createEnvironmentRequest(environmentName, clusterName));
        });
  }
}
