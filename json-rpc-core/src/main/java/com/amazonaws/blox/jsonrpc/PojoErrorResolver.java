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
package com.amazonaws.blox.jsonrpc;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.googlecode.jsonrpc4j.ErrorResolver;
import java.lang.reflect.Method;
import java.util.Arrays;
import java.util.List;
import lombok.extern.slf4j.Slf4j;

@Slf4j
public class PojoErrorResolver implements ErrorResolver {
  // Since we encode the type name, just use one error code for all errors that were encoded by this
  // resolver. The number is arbitrary, and just chosen to be high enough to not accidentally
  // conflict with another error code.
  public static final int ERROR_CODE = 9001;

  private final ObjectMapper mapper;

  public PojoErrorResolver(final ObjectMapper mapper) {
    this.mapper = mapper.copy().addMixIn(Throwable.class, ThrowableSerializationMixin.class);
  }

  @Override
  public JsonError resolveError(
      final Throwable t, final Method method, final List<JsonNode> arguments) {

    final boolean isModeledException =
        Arrays.stream(method.getExceptionTypes()).anyMatch(aClass -> aClass.isInstance(t));

    if (isModeledException) {
      // We have to explicitly serialize the exception type here using the mapper, rather than
      // relying on the mapper serializing the entire JsonError object. This is needed because
      // Jackson looks up the mixins based on the type of the *variable* not the object instance.
      // Since the type of the data field on JsonError is just Object, it then doesn't correctly
      // apply the annotations from ThrowableSerializationMixin.
      final JsonNode node = mapper.valueToTree(t);
      return new JsonError(ERROR_CODE, t.getMessage(), node);
    }

    log.warn(
        "Exception type {} is not modeled in signature of {}, returning generic error",
        t.getClass(),
        method.getName());

    return new JsonError(JsonError.INTERNAL_ERROR.code, t.getMessage(), null);
  }
}
