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
package com.amazonaws.blox.dataservice.api;

import com.amazonaws.blox.dataservice.exception.ResourceType;
import com.amazonaws.blox.dataservice.mapper.ApiModelMapper;
import com.amazonaws.blox.dataservice.model.*;
import com.amazonaws.blox.dataservice.repository.EnvironmentRepository;
import com.amazonaws.blox.dataservicemodel.v1.exception.InternalServiceException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentResponse;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mapstruct.factory.Mappers;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;

import java.time.Instant;

import static org.junit.Assert.assertEquals;
import static org.mockito.ArgumentMatchers.isA;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@RunWith(MockitoJUnitRunner.StrictStubs.class)
public class DescribeEnvironmentApiTest {
  private static final String ACCOUNT_ID = "123456789012";
  private static final String CLUSTER = "cluster";
  private static final String ENVIRONMENT_NAME = "name";
  private static final String ROLE_ARN = "role";
  private static final String ACTIVE_ENVIRONMENT_REVISION_ID = "123456789012_cluster_name";
  private static final String DEPLOYMENT_METHOD = "deploymentMethod";
  private static final ApiModelMapper apiModelMapper = Mappers.getMapper(ApiModelMapper.class);

  @Mock private EnvironmentRepository environmentRepository;

  private Environment environment;
  private EnvironmentId environmentId;
  private com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId environmentIdWrapper;
  private DescribeEnvironmentRequest describeEnvironmentRequest;
  private DescribeEnvironmentApi describeEnvironmentApi;

  @Before
  public void setup() {
    describeEnvironmentApi = new DescribeEnvironmentApi(apiModelMapper, environmentRepository);
    environmentId =
        EnvironmentId.builder()
            .accountId(ACCOUNT_ID)
            .cluster(CLUSTER)
            .environmentName(ENVIRONMENT_NAME)
            .build();
    environmentIdWrapper =
        com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId.builder()
            .accountId(ACCOUNT_ID)
            .cluster(CLUSTER)
            .environmentName(ENVIRONMENT_NAME)
            .build();
    environment =
        Environment.builder()
            .environmentId(environmentId)
            .role(ROLE_ARN)
            .environmentType(EnvironmentType.Daemon)
            .environmentStatus(EnvironmentStatus.Inactive)
            .environmentHealth(EnvironmentHealth.Healthy)
            .createdTime(Instant.now())
            .lastUpdatedTime(Instant.now())
            .activeEnvironmentRevisionId(ACTIVE_ENVIRONMENT_REVISION_ID)
            .latestEnvironmentRevisionId(ACTIVE_ENVIRONMENT_REVISION_ID)
            .deploymentConfiguration(new DeploymentConfiguration())
            .deploymentMethod(DEPLOYMENT_METHOD)
            .validEnvironment(true)
            .build();
    describeEnvironmentRequest =
        DescribeEnvironmentRequest.builder().environmentId(environmentIdWrapper).build();
  }

  @Test
  public void describeEnvironmentSuccess() throws Exception {
    when(environmentRepository.getEnvironment(environmentId)).thenReturn(environment);
    final DescribeEnvironmentResponse describeEnvironmentResponse =
        describeEnvironmentApi.describeEnvironment(describeEnvironmentRequest);

    verify(environmentRepository).getEnvironment(environmentId);

    assertEquals(
        describeEnvironmentResponse.getEnvironment().getEnvironmentId(), environmentIdWrapper);
    assertEquals(
        describeEnvironmentResponse.getEnvironment().getEnvironmentHealth(),
        environment.getEnvironmentHealth().name());
    assertEquals(
        describeEnvironmentResponse.getEnvironment().getActiveEnvironmentRevisionId(),
        environment.getActiveEnvironmentRevisionId());
    assertEquals(
        describeEnvironmentResponse.getEnvironment().getCreatedTime(),
        environment.getCreatedTime());
    assertEquals(
        describeEnvironmentResponse.getEnvironment().getEnvironmentStatus(),
        environment.getEnvironmentStatus().name());
    assertEquals(describeEnvironmentResponse.getEnvironment().getRole(), environment.getRole());
    assertEquals(
        describeEnvironmentResponse.getEnvironment().getDeploymentConfiguration(),
        apiModelMapper.toWrapperDeploymentConfiguration(environment.getDeploymentConfiguration()));
    assertEquals(
        describeEnvironmentResponse.getEnvironment().getEnvironmentType().name(),
        environment.getEnvironmentType().name());
    assertEquals(
        describeEnvironmentResponse.getEnvironment().getLastUpdatedTime(),
        environment.getLastUpdatedTime());
  }

  @Test(expected = ResourceNotFoundException.class)
  public void describeEnvironmentResourceNotFoundException() throws Exception {
    when(environmentRepository.getEnvironment(isA(EnvironmentId.class)))
        .thenThrow(new ResourceNotFoundException(ResourceType.ENVIRONMENT, ENVIRONMENT_NAME));
    describeEnvironmentApi.describeEnvironment(describeEnvironmentRequest);
  }

  @Test(expected = InternalServiceException.class)
  public void describeEnvironmentInternalServiceException() throws Exception {
    when(environmentRepository.getEnvironment(isA(EnvironmentId.class)))
        .thenThrow(new InternalServiceException(""));
    describeEnvironmentApi.describeEnvironment(describeEnvironmentRequest);
  }

  @Test(expected = InternalServiceException.class)
  public void describeEnvironmentInternalServiceExceptionWithUnknownException() throws Exception {
    when(environmentRepository.getEnvironment(isA(EnvironmentId.class)))
        .thenThrow(new IllegalStateException(""));
    describeEnvironmentApi.describeEnvironment(describeEnvironmentRequest);
  }
}
