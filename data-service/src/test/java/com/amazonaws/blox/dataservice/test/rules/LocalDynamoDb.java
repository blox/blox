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
package com.amazonaws.blox.dataservice.test.rules;

import com.amazonaws.services.dynamodbv2.AmazonDynamoDB;
import com.amazonaws.services.dynamodbv2.local.embedded.DynamoDBEmbedded;
import com.amazonaws.services.dynamodbv2.local.shared.access.AmazonDynamoDBLocal;
import org.junit.rules.ExternalResource;

public class LocalDynamoDb extends ExternalResource {
  private AmazonDynamoDBLocal localDdb;

  @Override
  protected void before() throws Throwable {
    localDdb = DynamoDBEmbedded.create();
  }

  @Override
  protected void after() {
    localDdb.shutdown();
  }

  public AmazonDynamoDB client() {
    return localDdb.amazonDynamoDB();
  }
}
