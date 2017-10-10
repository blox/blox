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
package com.amazonaws.blox.scheduling.handler;

import com.amazonaws.blox.scheduling.WorkflowApplication;
import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestStreamHandler;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.AnnotationConfigApplicationContext;

@Slf4j
public class MainLambdaHandler implements RequestStreamHandler {

  private static final ApplicationContext applicationContext =
      new AnnotationConfigApplicationContext(WorkflowApplication.class);

  @Override
  public void handleRequest(InputStream inputStream, OutputStream outputStream, Context context)
      throws IOException {

    final String lambdaAlias = getLambdaAlias(context.getInvokedFunctionArn());

    final Router router = applicationContext.getBean(Router.class);
    final StepHandler stepHandler = router.findStepHandler(lambdaAlias);
    log.debug("Invoking lambda {}", stepHandler);

    stepHandler.handleRequest(inputStream, outputStream, context);
  }

  /** Get the last part of the arn split by ":" which matches the lambda alias. */
  private String getLambdaAlias(final String lambdaArn) {
    final String[] parts = lambdaArn.split(":");
    return parts[parts.length - 1];
  }
}
