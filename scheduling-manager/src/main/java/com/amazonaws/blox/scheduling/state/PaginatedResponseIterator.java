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
package com.amazonaws.blox.scheduling.state;

import java.util.Iterator;
import java.util.List;
import java.util.Spliterator;
import java.util.Spliterators;
import java.util.function.Function;
import java.util.stream.Stream;
import java.util.stream.StreamSupport;
import lombok.RequiredArgsConstructor;

/**
 * Iterator that exposes paginated list API calls as an Iterator/Stream.
 *
 * @param <T> The response type returned by the list API call
 */
@RequiredArgsConstructor
class PaginatedResponseIterator<T> implements Iterator<T> {

  /** Function that extracts the pagination token from a response T */
  protected final Function<T, String> nextToken;

  /** Function that actually makes the list API call */
  protected final Function<String, T> list;

  private T previous = null;

  private String nextToken() {
    return previous == null ? null : nextToken.apply(previous);
  }

  @Override
  public boolean hasNext() {
    return previous == null || nextToken() != null;
  }

  @Override
  public T next() {
    previous = list.apply(nextToken());
    return previous;
  }

  /**
   * Return an ordered, synchronous stream that yields every page of results from the List API call.
   *
   * @return
   */
  public Stream<T> stream() {
    return StreamSupport.stream(
        Spliterators.spliteratorUnknownSize(this, Spliterator.ORDERED | Spliterator.NONNULL),
        false);
  }
}
