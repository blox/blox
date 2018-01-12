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

import cucumber.api.java8.En;
import cucumber.configuration.CucumberConfiguration;
import cucumber.steps.helpers.ExceptionContext;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.context.ContextConfiguration;

import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertEquals;

@ContextConfiguration(classes = CucumberConfiguration.class)
public class CommonSteps implements En {

  @Autowired private ExceptionContext exceptionContext;

  public CommonSteps() {
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
}
