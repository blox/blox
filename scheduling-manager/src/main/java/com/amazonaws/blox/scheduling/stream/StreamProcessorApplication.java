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
package com.amazonaws.blox.scheduling.stream;

import com.amazonaws.auth.profile.ProfileCredentialsProvider;
import com.amazonaws.services.dynamodbv2.AmazonDynamoDB;
import com.amazonaws.services.dynamodbv2.AmazonDynamoDBClient;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBMapper;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;

@Configuration
@ComponentScan("com.amazonaws.blox.scheduling")
public class StreamProcessorApplication {

  //TODO: remove profile
  private static final String CREDENTIALS_PROFILE_NAME = "daemon_prototype";

  @Bean
  public ProfileCredentialsProvider profileCredentialsProvider() {
    return new ProfileCredentialsProvider(CREDENTIALS_PROFILE_NAME);
  }

  @Bean
  public DynamoDBMapper dynamoDBMapper() {
    return new DynamoDBMapper(dynamoDBClient());
  }

  @Bean
  public AmazonDynamoDB dynamoDBClient() {
    return AmazonDynamoDBClient.builder().withCredentials(profileCredentialsProvider()).build();
  }

  @Bean
  public ObjectMapper objectMapper() {
    return new ObjectMapper();
  }
}
