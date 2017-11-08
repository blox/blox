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
package steps.dataservice;

import configuration.CucumberConfiguration;
import cucumber.api.java8.En;
import org.apache.commons.lang3.NotImplementedException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.context.ContextConfiguration;
import steps.wrappers.DataServiceWrapper;

@ContextConfiguration(classes = CucumberConfiguration.class)
public class CreateEnvironmentSteps implements En {

  @Autowired private DataServiceWrapper dataServiceWrapper;

  public CreateEnvironmentSteps() {
    When(
        "^I create an environment$",
        () -> {
          throw new NotImplementedException("");
        });

    Then(
        "^the created environment response is valid$",
        () -> {
          throw new NotImplementedException("");
        });
  }
}
