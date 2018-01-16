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

import com.amazonaws.blox.dataservice.model.Attribute;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBTypeConverter;
import java.util.StringJoiner;

public class AttributeDDBConverter implements DynamoDBTypeConverter<String, Attribute> {

  private static final String DELIMITER = "/";

  @Override
  public String convert(final Attribute attribute) {
    return new StringJoiner(DELIMITER)
        .add(attribute.getName())
        .add(attribute.getValue())
        .toString();
  }

  @Override
  public Attribute unconvert(final String stringValue) {
    return Attribute.builder()
        .name(stringValue.split(DELIMITER)[0])
        .value(stringValue.split(DELIMITER, 2)[1])
        .build();
  }
}
