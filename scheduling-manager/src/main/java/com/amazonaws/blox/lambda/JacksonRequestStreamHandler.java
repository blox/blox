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
package com.amazonaws.blox.lambda;

import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import com.amazonaws.services.lambda.runtime.RequestStreamHandler;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectReader;
import com.fasterxml.jackson.databind.ObjectWriter;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.lang.reflect.Type;
import java.lang.reflect.TypeVariable;
import java.util.Map;
import lombok.extern.log4j.Log4j2;
import org.apache.commons.lang3.reflect.TypeUtils;
import org.springframework.beans.factory.annotation.Autowired;

/**
 * Lambda RequestStreamHandler that allows for custom Jackson serialization of input/output types.
 *
 * <p>The default Lambda sandbox does not allow customization of the way that the function's inputs
 * and outputs are serialized. This class provides a way to inject a Jackson {@link ObjectMapper}
 * instance that can be configured as needed.
 */
@Log4j2
public class JacksonRequestStreamHandler<IN, OUT> implements RequestStreamHandler {

  private static final TypeVariable<Class<RequestHandler>>[] PARAMETERS =
      RequestHandler.class.getTypeParameters();

  private final ObjectReader reader;
  private final ObjectWriter writer;
  private final RequestHandler<IN, OUT> innerHandler;

  @Autowired
  public JacksonRequestStreamHandler(ObjectMapper mapper, RequestHandler<IN, OUT> innerHandler) {
    this.innerHandler = innerHandler;

    Map<TypeVariable<?>, Type> arguments =
        TypeUtils.getTypeArguments(innerHandler.getClass(), RequestHandler.class);
    Type inputType = arguments.get(PARAMETERS[0]);
    Type outputType = arguments.get(PARAMETERS[1]);

    this.reader = mapper.readerFor(mapper.constructType(inputType));
    this.writer = mapper.writerFor(mapper.constructType(outputType));
  }

  @Override
  public void handleRequest(InputStream input, OutputStream output, Context context)
      throws IOException {
    IN request = reader.readValue(input);
    log.debug("Request: {}", request);

    OUT response = innerHandler.handleRequest(request, context);
    log.debug("Response: {}", response);

    writer.writeValue(output, response);
  }
}
