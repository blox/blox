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

import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceExistsException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceInUseException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.model.Environment;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentHealth;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentStatus;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.*;
import cucumber.api.java8.En;
import cucumber.configuration.CucumberConfiguration;
import cucumber.steps.helpers.ExceptionContext;
import cucumber.steps.helpers.InputCreator;
import cucumber.steps.wrappers.DataServiceWrapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.annotation.DirtiesContext;
import org.springframework.test.context.ContextConfiguration;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;

@ContextConfiguration(classes = CucumberConfiguration.class)
@DirtiesContext
public class DataServiceSteps implements En {

  @Autowired private ExceptionContext exceptionContext;
  @Autowired private DataServiceWrapper dataServiceWrapper;
  @Autowired private InputCreator inputCreator;

  public DataServiceSteps() {
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

    When(
        "^I describe the created environment$",
        () -> {
          final EnvironmentId environmentId = getEnvironmentIdFromCreatedEnvironment();
          dataServiceWrapper.describeEnvironment(
              inputCreator.describeEnvironmentRequest(environmentId));
        });

    When(
        "^I describe the updated environment$",
        () -> {
          final EnvironmentId environmentId = getEnvironmentIdFromUpdatedEnvironment();
          dataServiceWrapper.describeEnvironment(
              inputCreator.describeEnvironmentRequest(environmentId));
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
                  environmentId.getEnvironmentName(), inputCreator.prefixName(newCluster)));
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
              inputCreator.deleteEnvironmentRequest(environmentId));
        });

    When(
        "^I try to describe the created environment$",
        () -> {
          final EnvironmentId environmentId = getEnvironmentIdFromCreatedEnvironment();
          dataServiceWrapper.tryDescribeEnvironment(
              inputCreator.describeEnvironmentRequest(environmentId));
        });

    Then(
        "^there should be an? \"?\'?(\\w*)\"?\'? thrown$",
        (final String exceptionName) -> {
          assertNotNull("Expecting an exception to be thrown", exceptionContext.getException());
          assertEquals(exceptionName, exceptionContext.getException().getClass().getSimpleName());
        });

    And(
        "^the resourceType is \"([^\"]*)\"$",
        (final String resourceType) -> {
          assertNotNull("Expecting an exception to be thrown", exceptionContext.getException());
          if (exceptionContext.getException() instanceof ResourceNotFoundException) {
            final ResourceNotFoundException exception =
                (ResourceNotFoundException) exceptionContext.getException();
            assertEquals(resourceType, exception.getResourceType());
          } else if (exceptionContext.getException() instanceof ResourceExistsException) {
            final ResourceExistsException exception =
                (ResourceExistsException) exceptionContext.getException();
            assertEquals(resourceType, exception.getResourceType());
          } else if (exceptionContext.getException() instanceof ResourceInUseException) {
            final ResourceInUseException exception =
                (ResourceInUseException) exceptionContext.getException();
            assertEquals(resourceType, exception.getResourceType());
          }
        });

    And(
        "^the resourceId contains \"([^\"]*)\"$",
        (final String resourceId) -> {
          assertNotNull("Expecting an exception to be thrown", exceptionContext.getException());
          if (exceptionContext.getException() instanceof ResourceNotFoundException) {
            final ResourceNotFoundException exception =
                (ResourceNotFoundException) exceptionContext.getException();
            assertEquals(resourceId, exception.getResourceType());
          } else if (exceptionContext.getException() instanceof ResourceExistsException) {
            final ResourceExistsException exception =
                (ResourceExistsException) exceptionContext.getException();
            assertEquals(resourceId, exception.getResourceType());
          } else if (exceptionContext.getException() instanceof ResourceInUseException) {
            final ResourceInUseException exception =
                (ResourceInUseException) exceptionContext.getException();
            assertEquals(resourceId, exception.getResourceType());
          }
        });
  }

  // TODO: Currently just check if the environment are equal to the other. May only need to just
  // compare the equality of some fields but not all
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
