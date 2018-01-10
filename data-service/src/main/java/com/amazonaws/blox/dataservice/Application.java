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
package com.amazonaws.blox.dataservice;

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.serialization.DataServiceMapperFactory;
import com.amazonaws.blox.jsonrpc.JsonRpcLambdaHandler;
import com.amazonaws.services.dynamodbv2.AmazonDynamoDB;
import com.amazonaws.services.dynamodbv2.AmazonDynamoDBClient;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBMapper;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBMapperConfig;
import java.util.StringJoiner;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;

@Configuration
@ComponentScan("com.amazonaws.blox.dataservice")
public class Application {

  private static final String REGION_STAGE_SEPARATOR = "-";

  @Bean
  public JsonRpcLambdaHandler<DataService> serviceHandler(DataService service) {
    return new JsonRpcLambdaHandler<>(
        DataServiceMapperFactory.newMapper(), DataService.class, service);
  }

  @Bean
  public DynamoDBMapper dynamoDBMapper() {
    return new DynamoDBMapper(dynamoDBClient(), dynamoDbMapperConfig());
  }

  @Bean
  public AmazonDynamoDB dynamoDBClient() {
    return AmazonDynamoDBClient.builder().build();
  }

  private DynamoDBMapperConfig dynamoDbMapperConfig() {
    return new DynamoDBMapperConfig.Builder()
        .withSaveBehavior(DynamoDBMapperConfig.SaveBehavior.UPDATE)
        .withConsistentReads(DynamoDBMapperConfig.ConsistentReads.CONSISTENT)
        .build();
  }

  private String tableNamePrefix() {
    return regionStageStackPrefix() + "data-service-";
  }

  //TODO: remove hardcoded region and stage
  private String regionStageStackPrefix() {
    return new StringJoiner(REGION_STAGE_SEPARATOR).add("us-east-1").add("beta").toString();
  }
}
