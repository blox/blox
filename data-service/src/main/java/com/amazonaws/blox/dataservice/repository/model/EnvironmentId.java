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

import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBTypeConverter;
import lombok.Builder;
import lombok.Value;

@Builder
@Value
public class EnvironmentId {

  private final String accountId;
  private final String environmentName;

  public static class Converter implements DynamoDBTypeConverter<String, EnvironmentId> {
    @Override
    public String convert(EnvironmentId id) {
      return String.join("_", id.accountId, id.environmentName);
    }

    @Override
    public EnvironmentId unconvert(String id) {
      String[] parts = id.split("_");
      return EnvironmentId.builder().accountId(parts[0]).environmentName(parts[1]).build();
    }
  }
}
