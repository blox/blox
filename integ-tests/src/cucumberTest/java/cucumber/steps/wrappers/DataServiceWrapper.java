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
package cucumber.steps.wrappers;

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentResponse;
import cucumber.steps.helpers.ExceptionContext;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.UpdateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DeleteEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.UpdateEnvironmentResponse;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DeleteEnvironmentResponse;
import java.util.function.Consumer;

import lombok.RequiredArgsConstructor;

@RequiredArgsConstructor
public class DataServiceWrapper extends MemoizedWrapper {

  private final DataService dataService;
  private final ExceptionContext exceptionContext;

  public CreateEnvironmentResponse createEnvironment(
      CreateEnvironmentRequest createEnvironmentRequest) {
    return memoizeFunction(createEnvironmentRequest, dataService::createEnvironment);
  }

  public DescribeEnvironmentResponse describeEnvironment(
      DescribeEnvironmentRequest describeEnvironmentRequest) {
    return memoizeFunction(describeEnvironmentRequest, dataService::describeEnvironment);
  }

  public UpdateEnvironmentResponse updateEnvironment(
      UpdateEnvironmentRequest updateEnvironmentRequest) {
    return memoizeFunction(updateEnvironmentRequest, dataService::updateEnvironment);
  }

  public DeleteEnvironmentResponse deleteEnvironment(
      DeleteEnvironmentRequest deleteEnvironmentRequest) {
    return memoizeFunction(deleteEnvironmentRequest, dataService::deleteEnvironment);
  }

  public void tryCreateEnvironment(CreateEnvironmentRequest createEnvironmentRequest) {
    captureException(createEnvironmentRequest, this::createEnvironment);
  }

  public void tryDescribeEnvironment(DescribeEnvironmentRequest describeEnvironmentRequest) {
    captureException(describeEnvironmentRequest, this::describeEnvironment);
  }

  public void tryUpdateEnvironment(UpdateEnvironmentRequest updateEnvironmentRequest) {
    captureException(updateEnvironmentRequest, this::updateEnvironment);
  }

  public void tryDeleteEnvironment(DeleteEnvironmentRequest deleteEnvironmentRequest) {
    captureException(deleteEnvironmentRequest, this::deleteEnvironment);
  }

  private <T> void captureException(final T input, final Consumer<T> consumer) {
    try {
      consumer.accept(input);
      throw new RuntimeException("Expected an exception, but none was thrown");
    } catch (final Exception e) {
      exceptionContext.setException(e);
    }
  }
}
