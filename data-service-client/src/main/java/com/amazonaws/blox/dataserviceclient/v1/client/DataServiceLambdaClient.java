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
package com.amazonaws.blox.dataserviceclient.v1.client;

import com.amazonaws.blox.dataservicemodel.v1.client.DataService;
import com.amazonaws.blox.dataservicemodel.v1.serialization.DataServiceMapperFactory;
import com.amazonaws.blox.jsonrpc.JsonRpcLambdaClient;
import software.amazon.awssdk.services.lambda.LambdaAsyncClient;

/** AWS DataService Lambda client using JSONRPC for routing. */
public class DataServiceLambdaClient {

  public static DataService dataService(
      final LambdaAsyncClient lambdaAsyncClient, final String dataServiceLambdaFunctionName) {

    //JsonRpcLambdaClient requires a lambda async client
    return new JsonRpcLambdaClient(
            DataServiceMapperFactory.newMapper(), lambdaAsyncClient, dataServiceLambdaFunctionName)
        .newProxy(DataService.class);
  }
}
