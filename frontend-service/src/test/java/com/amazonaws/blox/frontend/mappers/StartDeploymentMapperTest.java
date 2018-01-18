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
package com.amazonaws.blox.frontend.mappers;

import com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.StartDeploymentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.StartDeploymentResponse;
import com.amazonaws.blox.frontend.operations.StartDeployment;
import com.amazonaws.serverless.proxy.internal.model.ApiGatewayRequestContext;
import org.junit.Before;
import org.junit.Test;
import org.mapstruct.factory.Mappers;

import static org.assertj.core.api.Assertions.assertThat;

public class StartDeploymentMapperTest {

  private static final String ACCOUNT_ID = "123456789012";
  private static final String CLUSTER = "mycluster";
  private static final String ENVIRONMENT_NAME = "mynev";
  private static final String ENVIRONMENT_REVISION_ID = "1";
  private static final String DEPLOYMENT_ID = "deployment-id";

  private final StartDeploymentMapper mapper = Mappers.getMapper(StartDeploymentMapper.class);

  private EnvironmentId environmentId;

  @Before
  public void setUp() {
    environmentId =
        EnvironmentId.builder()
            .accountId(ACCOUNT_ID)
            .cluster(CLUSTER)
            .environmentName(ENVIRONMENT_NAME)
            .build();
  }

  @Test
  public void toDataServiceRequest() {
    final ApiGatewayRequestContext context = new ApiGatewayRequestContext();
    context.setAccountId(ACCOUNT_ID);

    final StartDeploymentRequest result =
        mapper.toDataServiceRequest(context, CLUSTER, ENVIRONMENT_NAME, ENVIRONMENT_REVISION_ID);

    assertThat(result.getEnvironmentRevisionId()).isEqualTo(ENVIRONMENT_REVISION_ID);
    assertThat(result.getEnvironmentId()).isEqualToComparingFieldByField(environmentId);
  }

  @Test
  public void fromDataServiceResponse() {
    final StartDeploymentResponse response =
        StartDeploymentResponse.builder()
            .deploymentId(DEPLOYMENT_ID)
            .environmentId(environmentId)
            .environmentRevisionId(ENVIRONMENT_REVISION_ID)
            .build();

    final StartDeployment.StartDeploymentResponse result = mapper.fromDataServiceResponse(response);

    assertThat(result.getDeploymentId()).isEqualTo(DEPLOYMENT_ID);
  }
}
