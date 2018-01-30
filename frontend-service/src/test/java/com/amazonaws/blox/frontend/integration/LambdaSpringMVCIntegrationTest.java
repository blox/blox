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
package com.amazonaws.blox.frontend.integration;

import static org.assertj.core.api.AssertionsForClassTypes.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import com.amazonaws.blox.frontend.integration.SampleController.SampleResponse;
import com.amazonaws.serverless.proxy.internal.model.AwsProxyRequest;
import com.amazonaws.serverless.proxy.internal.model.AwsProxyResponse;
import com.amazonaws.serverless.proxy.internal.testutils.AwsProxyRequestBuilder;
import org.junit.Test;

public class LambdaSpringMVCIntegrationTest extends IntegrationTestBase {

  @Test
  public void testRequestMapping() throws Exception {
    when(checker.handle(any())).thenReturn(new SampleResponse("ok"));

    final AwsProxyRequest request =
        new AwsProxyRequestBuilder("/test/testPath/sample", "POST")
            .queryString("query", "1")
            .build();

    final AwsProxyResponse response = handler.proxy(request, lambdaContext);

    assertThat(response.getStatusCode()).isEqualTo(200);
    verify(checker).handle("testPath", "1");
    SampleResponse sampleResponse =
        objectMapper.readValue(response.getBody(), SampleResponse.class);
    assertThat(sampleResponse.getResult()).isEqualTo("ok");
  }
}
