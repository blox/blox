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
package steps.helpers;

import java.util.function.Function;

@FunctionalInterface
public interface ThrowingFunction<T, R> extends Function<T, R> {

  @Override
  default R apply(T t) {
    try {
      return applyThrows(t);
    } catch (Exception e) {
      throw new RuntimeException(e);
    }
  }

  R applyThrows(T t) throws Exception;
}
