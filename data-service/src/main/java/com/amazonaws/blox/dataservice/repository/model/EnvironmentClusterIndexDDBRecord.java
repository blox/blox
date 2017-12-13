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
package com.amazonaws.blox.dataservice.repository.model;

import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBIndexHashKey;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBIndexRangeKey;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBTable;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@DynamoDBTable(tableName = "Environments")
@Data
@Builder
//Required for builder because an empty constructor exists
@AllArgsConstructor
//Required for dynamodbmapper
@NoArgsConstructor
public class EnvironmentClusterIndexDDBRecord {

  public static final String ENVIRONMENT_CLUSTER_GSI_NAME = "environmentClusterIndex";
  public static final String ENVIRONMENT_CLUSTER_INDEX_HASH_KEY = "accountId";
  public static final String ENVIRONMENT_CLUSTER_INDEX_RANGE_KEY = "cluster";

  @DynamoDBIndexHashKey(
    attributeName = ENVIRONMENT_CLUSTER_INDEX_HASH_KEY,
    globalSecondaryIndexName = ENVIRONMENT_CLUSTER_GSI_NAME
  )
  private String accountId;

  @DynamoDBIndexRangeKey(
    attributeName = ENVIRONMENT_CLUSTER_INDEX_RANGE_KEY,
    globalSecondaryIndexName = ENVIRONMENT_CLUSTER_GSI_NAME
  )
  private String clusterName;
}
