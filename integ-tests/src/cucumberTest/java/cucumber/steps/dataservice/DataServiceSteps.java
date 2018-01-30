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

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;

import com.amazonaws.blox.dataservicemodel.v1.model.Cluster;
import com.amazonaws.blox.dataservicemodel.v1.model.Environment;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentHealth;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId;
import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentStatus;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListClustersResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.UpdateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.UpdateEnvironmentResponse;
import cucumber.api.DataTable;
import cucumber.api.java8.En;
import cucumber.configuration.CucumberConfiguration;
import cucumber.steps.helpers.ExceptionContext;
import cucumber.steps.helpers.InputCreator;
import cucumber.steps.wrappers.DataServiceWrapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.annotation.DirtiesContext;
import org.springframework.test.context.ContextConfiguration;

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

    // TODO: Remove; changing the cluster of an existent environment is unsupported.
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
          assertThat(exceptionContext.getException())
              .isNotNull()
              .satisfies(t -> assertThat(t.getClass().getSimpleName()).isEqualTo(exceptionName));
        });

    And(
        "^its ([^ ]*) is \"([^\"]*)\"$",
        (final String field, final String value) -> {
          assertThat(exceptionContext.getException()).hasFieldOrPropertyWithValue(field, value);
        });

    And(
        "^its ([^ ]*) contains \"([^\"]*)\"$",
        (final String field, final String value) -> {
          assertThat(exceptionContext.getException())
              .hasFieldOrProperty(field)
              .extracting(field)
              .allSatisfy(s -> assertThat(s).asString().contains(value));
        });

    When(
        "^I list clusters$",
        () -> {
          dataServiceWrapper.listClusters(inputCreator.listClustersRequest());
        });

    Given(
        "^I am using account ID (\\d+)$",
        (String accountId) -> {
          inputCreator.setAccountId(accountId);
        });

    Then(
        "^no clusters are returned$",
        () -> {
          final ListClustersResponse response =
              dataServiceWrapper.getLastFromHistory(ListClustersResponse.class);
          assertThat(response.getClusters()).isEmpty();
        });

    Then(
        "^these clusters are returned:$",
        (DataTable table) -> {
          final ListClustersResponse response =
              dataServiceWrapper.getLastFromHistory(ListClustersResponse.class);

          Cluster[] expectedClusters =
              table
                  .asList(Cluster.class)
                  .stream()
                  .peek(c -> c.setClusterName(inputCreator.prefixName(c.getClusterName())))
                  .toArray(Cluster[]::new);

          assertThat(response.getClusters()).containsExactlyInAnyOrder(expectedClusters);
        });

    When(
        "^I list environments in cluster \"([^\"]*)\"$",
        (final String clusterName) -> {
          dataServiceWrapper.listEnvironments(
              inputCreator.listEnvironmentsRequest(clusterName, null));
        });

    When(
        "^I list environments in cluster \"([^\"]*)\" with name prefix \"([^\"]*)\"$",
        (final String clusterName, final String environmentNamePrefix) -> {
          dataServiceWrapper.listEnvironments(
              inputCreator.listEnvironmentsRequest(
                  clusterName, inputCreator.prefixName(environmentNamePrefix)));
        });

    Then(
        "^these environments are returned$",
        (DataTable table) -> {
          final ListEnvironmentsResponse response =
              dataServiceWrapper.getLastFromHistory(ListEnvironmentsResponse.class);

          EnvironmentId[] expectedEnvironmentIds =
              table
                  .asList(EnvironmentId.class)
                  .stream()
                  .peek(
                      e -> {
                        e.setCluster(inputCreator.prefixName(e.getCluster()));
                        e.setEnvironmentName(inputCreator.prefixName(e.getEnvironmentName()));
                      })
                  .toArray(EnvironmentId[]::new);

          assertThat(response.getEnvironmentIds())
              .containsExactlyInAnyOrder(expectedEnvironmentIds);
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
