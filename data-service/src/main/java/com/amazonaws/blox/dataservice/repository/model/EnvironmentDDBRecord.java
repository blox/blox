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

import com.amazonaws.blox.dataservice.model.EnvironmentHealth;
import com.amazonaws.blox.dataservice.model.EnvironmentStatus;
import com.amazonaws.blox.dataservice.model.EnvironmentType;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBHashKey;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBIndexHashKey;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBIndexRangeKey;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBRangeKey;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBTable;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBTypeConverted;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBTypeConvertedEnum;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBVersionAttribute;
import java.time.Instant;
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
public class EnvironmentDDBRecord {

  public static final String ACCOUNT_ID_CLUSTER_HASH_KEY = "accountIdCluster";
  public static final String ENVIRONMENT_NAME_RANGE_KEY = "environmentName";
  public static final String LATEST_ENVIRONMENT_REVISION_ID = "latestEnvironmentRevisionId";

  public static final String ENVIRONMENT_CLUSTER_GSI_NAME = "environmentClusterIndex";
  public static final String ENVIRONMENT_CLUSTER_INDEX_HASH_KEY = "accountId";
  public static final String ENVIRONMENT_CLUSTER_INDEX_RANGE_KEY = "cluster";

  public static EnvironmentDDBRecord withHashKeys(final String accountIdCluster) {
    return EnvironmentDDBRecord.builder().accountIdCluster(accountIdCluster).build();
  }

  public static EnvironmentDDBRecord withKeys(
      final String accountIdCluster, final String environmentName) {

    return EnvironmentDDBRecord.builder()
        .accountIdCluster(accountIdCluster)
        .environmentName(environmentName)
        .build();
  }

  @DynamoDBHashKey(attributeName = ACCOUNT_ID_CLUSTER_HASH_KEY)
  private String accountIdCluster;

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

  @DynamoDBRangeKey(attributeName = ENVIRONMENT_NAME_RANGE_KEY)
  private String environmentName;

  @DynamoDBVersionAttribute private Long recordVersion;

  @DynamoDBTypeConverted(converter = InstantDDBConverter.class)
  private Instant createdTime;

  @DynamoDBTypeConverted(converter = InstantDDBConverter.class)
  private Instant lastUpdatedTime;

  //TODO: add deploymentConfiguration to ddb record
  private String role;
  private String activeEnvironmentRevisionId;
  private String latestEnvironmentRevisionId;
  private Boolean validEnvironment;
  private String deploymentMethod;

  @DynamoDBTypeConvertedEnum private EnvironmentType type;

  @DynamoDBTypeConvertedEnum private EnvironmentStatus status;

  @DynamoDBTypeConvertedEnum private EnvironmentHealth health;
}
