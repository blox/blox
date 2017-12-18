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

import com.google.common.collect.ArrayListMultimap;
import com.google.common.collect.ListMultimap;
import com.google.common.collect.Multimaps;
import java.util.List;
import java.util.function.Function;
import steps.helpers.Memoized;

public class MemoizedWrapper implements Memoized {

  private final ListMultimap<Class<?>, Object> memory =
      Multimaps.synchronizedListMultimap(ArrayListMultimap.create());

  @Override
  public <T> T getFromHistory(Class<T> type, Function<List<Object>, Object> fn) {
    return null;
  }

  @Override
  public <T> T getLastFromHistory(Class<T> type) {
    return null;
  }

  @Override
  public final <T> void addToHistory(final Class<T> type, final T value) {
    memory.put(type, type.cast(value));
  }
}
