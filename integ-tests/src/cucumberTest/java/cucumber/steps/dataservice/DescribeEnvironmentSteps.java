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

import com.amazonaws.blox.dataservicemodel.v1.model.Environment;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.*;

import cucumber.api.java8.En;
import cucumber.configuration.CucumberConfiguration;
import cucumber.steps.wrappers.DataServiceWrapper;
import cucumber.steps.helpers.InputCreator;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.context.ContextConfiguration;

import static org.junit.Assert.assertTrue;

@ContextConfiguration(classes = CucumberConfiguration.class)
public class DescribeEnvironmentSteps implements En {

  @Autowired private DataServiceWrapper dataServiceWrapper;
  @Autowired private InputCreator inputCreator;

  public DescribeEnvironmentSteps() {

    When(
        "^I describe the created environment$",
        () -> {
          final EnvironmentId environmentId = getEnvironmentIdFromCreatedEnvironment();
          dataServiceWrapper.describeEnvironment(
              inputCreator.describeEnvironmentRequest(
                  environmentId.getEnvironmentName(), environmentId.getCluster()));
        });

    When(
        "^I describe the updated environment$",
        () -> {
          final EnvironmentId environmentId = getEnvironmentIdFromUpdatedEnvironment();
          dataServiceWrapper.describeEnvironment(
              inputCreator.describeEnvironmentRequest(
                  environmentId.getEnvironmentName(), environmentId.getCluster()));
        });

    Then(
        "^the created and described environments match$",
        () -> {
          final CreateEnvironmentResponse createEnvironmentResponse =
              dataServiceWrapper.getLastFromHistory(CreateEnvironmentResponse.class);
          final DescribeEnvironmentResponse describeEnvironmentResponse =
              dataServiceWrapper.getLastFromHistory(DescribeEnvironmentResponse.class);
          checkEnvironmentEquality(
              createEnvironmentResponse.getEnvironment(),
              describeEnvironmentResponse.getEnvironment());
        });

    Given(
        "^I update the created environment with cluster name \"([^\"]*)\"$",
        (final String newCluster) -> {
          final EnvironmentId environmentId = getEnvironmentIdFromCreatedEnvironment();
          dataServiceWrapper.updateEnvironment(
              inputCreator.updateEnvironmentRequestWithNewCluster(
                  environmentId.getEnvironmentName(), newCluster));
        });

    Then(
        "^the updated and described environments match$",
        () -> {
          final UpdateEnvironmentResponse updateEnvironmentResponse =
              dataServiceWrapper.getLastFromHistory(UpdateEnvironmentResponse.class);
          final DescribeEnvironmentResponse describeEnvironmentResponse =
              dataServiceWrapper.getLastFromHistory(DescribeEnvironmentResponse.class);
          checkEnvironmentEquality(
              updateEnvironmentResponse.getEnvironment(),
              describeEnvironmentResponse.getEnvironment());
        });

    When(
        "^I try to describe a non-existent environment named \"([^\"]*)\"$",
        (final String environmentNamePrefix) -> {
          dataServiceWrapper.tryDescribeEnvironment(
              inputCreator.describeEnvironmentRequest(environmentNamePrefix));
        });

    Given(
        "^I delete the created environment$",
        () -> {
          final EnvironmentId environmentId = getEnvironmentIdFromCreatedEnvironment();
          dataServiceWrapper.deleteEnvironment(
              inputCreator.deleteEnvironmentRequest(
                  environmentId.getEnvironmentName(), environmentId.getCluster()));
        });

    When(
        "^I try to describe the created environment$",
        () -> {
          final EnvironmentId environmentId = getEnvironmentIdFromCreatedEnvironment();
          dataServiceWrapper.tryDescribeEnvironment(
              inputCreator.describeEnvironmentRequest(
                  environmentId.getEnvironmentName(), environmentId.getCluster()));
        });
  }

  // TODO: Currently just check if the environment are equal to the other. May only need to just compare the equality of some fields but not all
  private void checkEnvironmentEquality(
      final Environment thisEnvironment, final Environment otherEnvironment) {
    assertTrue(thisEnvironment.equals(otherEnvironment));
  }

  private EnvironmentId getEnvironmentIdFromCreatedEnvironment() {
    final CreateEnvironmentRequest createEnvironmentRequest =
        dataServiceWrapper.getLastFromHistory(CreateEnvironmentRequest.class);
    return createEnvironmentRequest.getEnvironmentId();
  }

  private EnvironmentId getEnvironmentIdFromUpdatedEnvironment() {
    final UpdateEnvironmentRequest updateEnvironmentRequest =
        dataServiceWrapper.getLastFromHistory(UpdateEnvironmentRequest.class);
    return updateEnvironmentRequest.getEnvironmentId();
  }
}
