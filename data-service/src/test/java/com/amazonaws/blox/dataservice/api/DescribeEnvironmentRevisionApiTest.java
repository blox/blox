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
import com.amazonaws.blox.dataservice.model.EnvironmentRevision;
import com.amazonaws.blox.dataservice.model.EnvironmentId;
import com.amazonaws.blox.dataservice.model.InstanceGroup;
import com.amazonaws.blox.dataservice.repository.EnvironmentRepository;
import com.amazonaws.blox.dataservicemodel.v1.exception.InternalServiceException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRevisionRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DescribeEnvironmentRevisionResponse;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mapstruct.factory.Mappers;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;

import java.time.Instant;
import java.util.HashSet;

import static org.junit.Assert.assertEquals;
import static org.mockito.ArgumentMatchers.isA;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@RunWith(MockitoJUnitRunner.StrictStubs.class)
public class DescribeEnvironmentRevisionApiTest {
  private static final String ACCOUNT_ID = "123456789012";
  private static final String CLUSTER = "cluster";
  private static final String ENVIRONMENT_NAME = "name";
  private static final String ENVIRONMENT_REVISION_ID = "123456789012_cluster_name";
  private static final String TASK_DEFINITION = "taskDefinition";

  private static final ApiModelMapper apiModelMapper = Mappers.getMapper(ApiModelMapper.class);

  @Mock private EnvironmentRepository environmentRepository;

  private EnvironmentRevision environmentRevision;
  private EnvironmentId environmentId;
  private InstanceGroup instanceGroup;
  private com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId environmentIdWrapper;
  private DescribeEnvironmentRevisionRequest describeEnvironmentRevisionRequest;
  private DescribeEnvironmentRevisionApi describeEnvironmentRevisionApi;

  @Before
  public void setup() {
    describeEnvironmentRevisionApi =
        new DescribeEnvironmentRevisionApi(apiModelMapper, environmentRepository);
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

    instanceGroup = InstanceGroup.builder().attributes(new HashSet<>()).build();

    environmentRevision =
        EnvironmentRevision.builder()
            .environmentRevisionId(ENVIRONMENT_REVISION_ID)
            .environmentId(environmentId)
            .createdTime(Instant.now())
            .instanceGroup(instanceGroup)
            .taskDefinition(TASK_DEFINITION)
            .build();
    describeEnvironmentRevisionRequest =
        DescribeEnvironmentRevisionRequest.builder()
            .environmentId(environmentIdWrapper)
            .environmentRevisionId(ENVIRONMENT_REVISION_ID)
            .build();
  }

  @Test
  public void describeEnvironmentRevisionSuccess() throws Exception {
    when(environmentRepository.getEnvironmentRevision(environmentId, ENVIRONMENT_REVISION_ID))
        .thenReturn(environmentRevision);
    final DescribeEnvironmentRevisionResponse describeEnvironmentRevisionResponse =
        describeEnvironmentRevisionApi.describeEnvironmentRevision(
            describeEnvironmentRevisionRequest);

    verify(environmentRepository).getEnvironmentRevision(environmentId, ENVIRONMENT_REVISION_ID);

    assertEquals(
        environmentIdWrapper,
        describeEnvironmentRevisionResponse.getEnvironmentRevision().getEnvironmentId());
    assertEquals(
        ENVIRONMENT_REVISION_ID,
        describeEnvironmentRevisionResponse.getEnvironmentRevision().getEnvironmentRevisionId());
    assertEquals(
        environmentRevision.getCreatedTime(),
        describeEnvironmentRevisionResponse.getEnvironmentRevision().getCreatedTime());
    assertEquals(
        environmentRevision.getInstanceGroup().getAttributes(),
        describeEnvironmentRevisionResponse
            .getEnvironmentRevision()
            .getInstanceGroup()
            .getAttributes());
    assertEquals(
        environmentRevision.getTaskDefinition(),
        describeEnvironmentRevisionResponse.getEnvironmentRevision().getTaskDefinition());
  }

  @Test(expected = ResourceNotFoundException.class)
  public void describeEnvironmentRevisionResourceNotFoundException() throws Exception {
    when(environmentRepository.getEnvironmentRevision(isA(EnvironmentId.class), isA(String.class)))
        .thenThrow(
            new ResourceNotFoundException(
                ResourceType.ENVIRONMENT_REVISION, ENVIRONMENT_REVISION_ID));
    describeEnvironmentRevisionApi.describeEnvironmentRevision(describeEnvironmentRevisionRequest);
  }

  @Test(expected = InternalServiceException.class)
  public void describeEnvironmentRevisionInternalServiceException() throws Exception {
    when(environmentRepository.getEnvironmentRevision(isA(EnvironmentId.class), isA(String.class)))
        .thenThrow(new InternalServiceException(""));
    describeEnvironmentRevisionApi.describeEnvironmentRevision(describeEnvironmentRevisionRequest);
  }

  @Test(expected = InternalServiceException.class)
  public void describeEnvironmentRevisionUnknownException() throws Exception {
    when(environmentRepository.getEnvironmentRevision(isA(EnvironmentId.class), isA(String.class)))
        .thenThrow(new IllegalStateException(""));
    describeEnvironmentRevisionApi.describeEnvironmentRevision(describeEnvironmentRevisionRequest);
  }
}
