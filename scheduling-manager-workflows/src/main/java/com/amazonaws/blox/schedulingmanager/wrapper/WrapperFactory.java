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
package com.amazonaws.blox.schedulingmanager.wrapper;

import com.amazonaws.auth.AWSCredentialsProvider;
import com.amazonaws.auth.STSAssumeRoleSessionCredentialsProvider;
import com.amazonaws.services.securitytoken.AWSSecurityTokenService;
import java.time.Duration;
import lombok.NonNull;

public interface WrapperFactory<T> {

  T getWrapper(AWSCredentialsProvider credentials);

  default AWSCredentialsProvider getCredentialsProvider(
      @NonNull final AWSSecurityTokenService stsClient,
      @NonNull final String roleArn,
      @NonNull final String sessionNamePrefix) {

    final String sessionName =
        String.format("%s-%s", sessionNamePrefix, String.valueOf(System.currentTimeMillis()));

    return new STSAssumeRoleSessionCredentialsProvider.Builder(roleArn, sessionName)
        .withRoleSessionDurationSeconds((int) Duration.ofMinutes(15).getSeconds())
        .withStsClient(stsClient)
        .build();
  }
}
