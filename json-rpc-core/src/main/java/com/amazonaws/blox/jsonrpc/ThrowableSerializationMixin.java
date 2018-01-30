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

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import com.fasterxml.jackson.annotation.JsonTypeInfo.As;
import com.fasterxml.jackson.annotation.JsonTypeInfo.Id;
import com.googlecode.jsonrpc4j.JsonRpcBasicServer;

/**
 * Jackson Mixin to control serialization of Throwables in error responses.
 *
 * <p>It suppress detailed fields of exceptions being serialized in JSON-RPC error responses, and
 * adds a property to the the serialized JSON that contains the Java class name of the exception, so
 * that the client can deserialize the correct exception class.
 */
@JsonIgnoreProperties(value = {"stackTrace", "suppressed", "cause", "localisedMessage"})
@JsonTypeInfo(
  use = Id.CLASS,
  include = As.PROPERTY,
  property = JsonRpcBasicServer.EXCEPTION_TYPE_NAME
)
interface ThrowableSerializationMixin {
  int ERROR_CODE = 9001;
}
