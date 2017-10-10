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

import com.amazonaws.blox.scheduling.stream.exception.InvalidRecordStructureException;
import com.amazonaws.services.dynamodbv2.model.AttributeValue;
import com.amazonaws.services.lambda.runtime.events.DynamodbEvent.DynamodbStreamRecord;
import java.util.Map;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;

@Slf4j
@Component
@AllArgsConstructor
public class DeploymentRecordProcessor implements RecordProcessor {

  //new item
  private static final String INSERT = "INSERT";
  private static final String DEPLOYMENTS_HASH_KEY = "deploymentId";

  @Override
  public void process(final DynamodbStreamRecord record) throws InvalidRecordStructureException {
    if (record.getEventName().equals(INSERT)) {
      final Map<String, AttributeValue> recordKeys = record.getDynamodb().getKeys();
      if (!recordKeys.containsKey(DEPLOYMENTS_HASH_KEY)) {
        throw new InvalidRecordStructureException(
            String.format("Record keys %s do not contain deploymentId", recordKeys));
      }

      final String deploymentId = recordKeys.get(DEPLOYMENTS_HASH_KEY).getS();
      log.debug(String.format("DeploymentId %s", deploymentId));

      //TODO: get and verify deployment from DDB
      //TODO: start deployment workflow

    } else {
      //TODO: add metrics on modify and remove?
    }
  }
}
