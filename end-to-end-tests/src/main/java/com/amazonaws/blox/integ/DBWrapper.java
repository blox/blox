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
package com.amazonaws.blox.integ;

import java.util.Map;
import java.util.Map.Entry;
import java.util.Set;
import java.util.stream.Collectors;
import lombok.RequiredArgsConstructor;
import software.amazon.awssdk.services.dynamodb.DynamoDBClient;
import software.amazon.awssdk.services.dynamodb.model.AttributeValue;
import software.amazon.awssdk.services.dynamodb.model.DeleteItemRequest;
import software.amazon.awssdk.services.dynamodb.model.DescribeTableRequest;
import software.amazon.awssdk.services.dynamodb.model.DescribeTableResponse;
import software.amazon.awssdk.services.dynamodb.model.KeySchemaElement;
import software.amazon.awssdk.services.dynamodb.model.ScanRequest;
import software.amazon.awssdk.services.dynamodb.model.ScanResponse;

@RequiredArgsConstructor
public class DBWrapper {
  private final DynamoDBClient ddb;

  public void reset() {
    // TODO: This is a temporary measure to clean up until we have a working DeleteEnvironment API.
    truncateTable("Environments");
    truncateTable("EnvironmentRevisions");
  }

  void truncateTable(String tableName) {
    // TODO: This will only delete one page of records, we'll replace this with the DeleteEnvironment API later.
    DescribeTableResponse table =
        ddb.describeTable(DescribeTableRequest.builder().tableName(tableName).build());

    Set<String> keys =
        table
            .table()
            .keySchema()
            .stream()
            .map(KeySchemaElement::attributeName)
            .collect(Collectors.toSet());

    ScanResponse records = ddb.scan(ScanRequest.builder().tableName(tableName).build());

    for (Map<String, AttributeValue> record : records.items()) {
      Map<String, AttributeValue> key =
          record
              .entrySet()
              .stream()
              .filter(e -> keys.contains(e.getKey()))
              .collect(Collectors.toMap(Entry::getKey, Entry::getValue));

      ddb.deleteItem(DeleteItemRequest.builder().tableName(tableName).key(key).build());
    }
  }
}
