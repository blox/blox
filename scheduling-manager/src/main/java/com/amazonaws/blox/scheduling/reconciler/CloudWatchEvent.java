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
package com.amazonaws.blox.scheduling.reconciler;

import com.amazonaws.regions.Regions;
import com.fasterxml.jackson.annotation.JsonProperty;
import java.time.Instant;
import java.util.List;
import lombok.Data;

/**
 * Representation of an Event from Cloudwatch.
 *
 * <p>See http://docs.aws.amazon.com/AmazonCloudWatch/latest/events/EventTypes.html for details.
 *
 * @param <T> The type of the "detail" field.
 */
@Data
public class CloudWatchEvent<T> {

  private String version;
  private String id;
  private String account;
  private String source;
  private Instant time;
  private Regions region;
  private List<String> resources;
  private T detail;

  @JsonProperty("detail-type")
  private String detailType;

  /**
   * Setter for Jackson to deserialize Regions enum from string
   *
   * @param region
   */
  public void setRegion(String region) {
    this.region = Regions.fromName(region);
  }
}
