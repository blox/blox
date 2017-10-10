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
package com.amazonaws.blox.scheduling.wrapper;

import com.amazonaws.services.ecs.AmazonECSClient;
import com.amazonaws.services.securitytoken.AWSSecurityTokenService;
import lombok.AllArgsConstructor;
import lombok.NonNull;

@AllArgsConstructor
public class ECSWrapperFactory implements WrapperFactory<ECSWrapper> {

  @NonNull private AWSSecurityTokenService stsClient;

  public ECSWrapper getWrapperForRole(final String roleArn, final String roleSessionName) {
    return new ECSWrapper(
        AmazonECSClient.builder()
            .withCredentials(getCredentialsProvider(stsClient, roleArn, roleSessionName))
            .build());
  }
}
