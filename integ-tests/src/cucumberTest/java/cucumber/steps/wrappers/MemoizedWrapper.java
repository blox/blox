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

import com.google.common.collect.ArrayListMultimap;
import com.google.common.collect.ListMultimap;
import com.google.common.collect.Multimaps;
import cucumber.steps.helpers.Memoized;
import cucumber.steps.helpers.ThrowingFunction;
import org.apache.commons.lang3.Validate;

public class MemoizedWrapper implements Memoized {

  private final ListMultimap<Class<?>, Object> memory =
      Multimaps.synchronizedListMultimap(ArrayListMultimap.create());

  @Override
  @SuppressWarnings("unchecked")
  public <T> T getLastFromHistory(final Class<T> type) {
    return (T) memory.get(type).get(memory.get(type).size() - 1);
  }

  @Override
  public final <T> void addToHistory(final Class<T> type, final T value) {
    memory.put(type, type.cast(value));
  }

  @Override
  @SuppressWarnings("unchecked")
  public final <T, R> R memoizeFunction(final T input, final ThrowingFunction<T, R> fn) {
    Validate.notNull(input);

    addToHistory((Class<T>) input.getClass(), input);
    final R result = fn.apply(input);
    Validate.notNull(result);
    addToHistory((Class<R>) result.getClass(), result);
    return result;
  }
}
