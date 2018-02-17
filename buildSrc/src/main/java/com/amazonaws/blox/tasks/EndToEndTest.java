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
package com.amazonaws.blox.tasks;

import org.gradle.api.internal.tasks.options.Option;
import org.gradle.api.tasks.testing.Test;
/** Gradle task for end-to-end test */
public class EndToEndTest extends Test {
  private static final String AWS_REGION = "aws.region";
  private static final String AWS_PROFILE = "aws.profile";
  private static final String BLOX_ENDPOINT = "blox.tests.apiUrl";

  @Option(option = "aws-region", description = "AWS region to use for test", order = 1)
  public void setRegion(String region) {
    super.systemProperty(AWS_REGION, region);
  }

  @Option(option = "aws-profile", description = "AWS credential profile to use for test", order = 2)
  public void setProfile(String profile) {
    super.systemProperty(AWS_PROFILE, profile);
  }

  @Option(option = "endpoint", description = "Blox service endpoint", order = 3)
  public void setEndpoint(String endpoint) {
    super.systemProperty(BLOX_ENDPOINT, endpoint);
  }

  public void setDefaultRegion(String region) {
    if (!getSystemProperties().containsKey(AWS_REGION)) {
      setRegion(region);
    }
  }

  public void setDefaultProfile(String profile) {
    if (!getSystemProperties().containsKey(AWS_PROFILE)) {
      setProfile(profile);
    }
  }

  public void setDefaultEndpoint(String endpoint) {
    if (!getSystemProperties().containsKey(BLOX_ENDPOINT)) {
      setEndpoint(endpoint);
    }
  }
}
