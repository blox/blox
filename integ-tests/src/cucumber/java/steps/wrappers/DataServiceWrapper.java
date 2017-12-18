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
package steps.wrappers;

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.CreateEnvironmentResponse;
import java.util.function.Consumer;
import lombok.RequiredArgsConstructor;
import org.apache.commons.lang3.Validate;
import steps.helpers.ExceptionContext;
import steps.helpers.ThrowingFunction;

@RequiredArgsConstructor
public class DataServiceWrapper extends MemoizedWrapper {

  private final DataService dataService;
  private final ExceptionContext exceptionContext;

  public CreateEnvironmentResponse createEnvironment(
      CreateEnvironmentRequest createEnvironmentRequest) {
    return memoizeInputAndCall(createEnvironmentRequest, dataService::createEnvironment);
  }

  public void tryCreateEnvironment(CreateEnvironmentRequest createEnvironmentRequest) {
    captureException(createEnvironmentRequest, this::createEnvironment);
  }

  private <T, R> R memoizeInputAndCall(final T input, ThrowingFunction<T, R> fn) {
    Validate.notNull(input);
    addToHistory((Class<T>) input.getClass(), input);
    return fn.apply(input);
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
