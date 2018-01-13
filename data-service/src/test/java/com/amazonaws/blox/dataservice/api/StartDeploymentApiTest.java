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
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentId;
import com.amazonaws.blox.dataservice.model.EnvironmentRevision;
import com.amazonaws.blox.dataservice.repository.EnvironmentRepository;
import com.amazonaws.blox.dataservicemodel.v1.exception.InternalServiceException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.StartDeploymentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.StartDeploymentResponse;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.InOrder;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;

import static org.assertj.core.api.Assertions.*;
import static org.mockito.Mockito.inOrder;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@RunWith(MockitoJUnitRunner.StrictStubs.class)
public class StartDeploymentApiTest {
  private static final String ENVIRONMENT_REVISION_ID = "revision-id";

  @Mock private ApiModelMapper apiModelMapper;
  @Mock private EnvironmentRepository environmentRepository;

  @Mock private com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId environmentIdWrapper;
  @Mock private EnvironmentId environmentId;
  @Mock private Environment environment;
  @Mock private EnvironmentRevision environmentRevision;

  @InjectMocks private StartDeploymentApi api;

  private StartDeploymentRequest request;

  @Before
  public void setUp() throws Exception {
    request =
        StartDeploymentRequest.builder()
            .environmentId(environmentIdWrapper)
            .environmentRevisionId(ENVIRONMENT_REVISION_ID)
            .build();

    when(apiModelMapper.toModelEnvironmentId(environmentIdWrapper)).thenReturn(environmentId);
    when(environmentRevision.getEnvironmentRevisionId()).thenReturn(ENVIRONMENT_REVISION_ID);
  }

  @Test
  public void startDeployment() throws Exception {
    // Given
    when(environmentRepository.getEnvironment(environmentId)).thenReturn(environment);
    when(environmentRepository.getEnvironmentRevision(environmentId, ENVIRONMENT_REVISION_ID))
        .thenReturn(environmentRevision);

    // When
    final StartDeploymentResponse response = api.startDeployment(request);

    // Then
    assertThat(response.getEnvironmentId()).isEqualTo(environmentIdWrapper);
    assertThat(response.getEnvironmentRevisionId()).isEqualTo(ENVIRONMENT_REVISION_ID);
    assertThat(response.getDeploymentId()).isNotEmpty();

    verify(apiModelMapper).toModelEnvironmentId(environmentIdWrapper);
    verify(environmentRepository).getEnvironment(environmentId);
    verify(environmentRepository).getEnvironmentRevision(environmentId, ENVIRONMENT_REVISION_ID);

    InOrder inOrder = inOrder(environment, environmentRepository);
    inOrder.verify(environment).setActiveEnvironmentRevisionId(ENVIRONMENT_REVISION_ID);
    inOrder.verify(environmentRepository).updateEnvironment(environment);
  }

  @Test(expected = ResourceNotFoundException.class)
  public void startDeploymentResourceNotFoundException() throws Exception {
    when(environmentRepository.getEnvironment(environmentId))
        .thenThrow(new ResourceNotFoundException(ResourceType.ENVIRONMENT, "env-name"));

    api.startDeployment(request);
  }

  @Test(expected = InternalServiceException.class)
  public void startDeploymentInternalServiceException() throws Exception {
    when(environmentRepository.getEnvironment(environmentId))
        .thenThrow(new InternalServiceException(""));

    api.startDeployment(request);
  }

  @Test(expected = InternalServiceException.class)
  public void startDeploymentInternalServiceExceptionWithUnknownException() throws Exception {
    when(environmentRepository.getEnvironment(environmentId))
        .thenThrow(new IllegalStateException(""));

    api.startDeployment(request);
  }
}
