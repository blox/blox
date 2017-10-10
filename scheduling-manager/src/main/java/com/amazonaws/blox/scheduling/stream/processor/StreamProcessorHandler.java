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
package com.amazonaws.blox.scheduling.stream.processor;

import com.amazonaws.blox.scheduling.stream.StreamProcessorApplication;
import com.amazonaws.blox.scheduling.stream.exception.InvalidRecordStructureException;
import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import com.amazonaws.services.lambda.runtime.events.DynamodbEvent;
import com.amazonaws.services.lambda.runtime.events.DynamodbEvent.DynamodbStreamRecord;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.AnnotationConfigApplicationContext;

@Slf4j
public class StreamProcessorHandler implements RequestHandler<DynamodbEvent, String> {

  public static final ApplicationContext applicationContext =
      new AnnotationConfigApplicationContext(StreamProcessorApplication.class);

  @Override
  public String handleRequest(final DynamodbEvent dynamodbEvent, final Context context) {

    for (DynamodbStreamRecord record : dynamodbEvent.getRecords()) {
      log.debug(record.toString());

      final RecordProcessor recordProcessor = applicationContext.getBean(RecordProcessor.class);

      try {
        recordProcessor.process(record);
      } catch (final InvalidRecordStructureException e) {
        log.error(e.getMessage(), e);
      }
    }

    return null;
  }
}
