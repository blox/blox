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

import static org.junit.Assert.assertEquals;
import static org.mockito.Mockito.doNothing;
import static org.mockito.Mockito.doThrow;
import static org.mockito.Mockito.when;

import com.amazonaws.blox.dataservice.mapper.ApiModelMapper;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentId;
import com.amazonaws.blox.dataservice.model.EnvironmentRevision;
import com.amazonaws.blox.dataservice.repository.EnvironmentRepository;
import com.amazonaws.blox.dataservicemodel.v1.exception.InternalServiceException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceNotFoundException;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DeleteEnvironmentRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.DeleteEnvironmentResponse;
import java.util.Collections;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;

@RunWith(MockitoJUnitRunner.StrictStubs.class)
public class DeleteEnvironmentApiTest {

  @Mock private EnvironmentRepository environmentRepository;
  @Mock private ApiModelMapper apiModelMapper;
  @Mock private EnvironmentId environmentId;
  @Mock private com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId wrapperEnvironmentId;
  @Mock private Environment environment;
  @Mock private com.amazonaws.blox.dataservicemodel.v1.model.Environment wrapperEnvironment;
  @Mock private EnvironmentRevision environmentRevision;
  @Mock private ResourceNotFoundException resourceNotFoundException;

  private DeleteEnvironmentRequest deleteEnvironmentRequest;
  private DeleteEnvironmentApi deleteEnvironmentApi;

  @Before
  public void setup() {
    deleteEnvironmentApi = new DeleteEnvironmentApi(apiModelMapper, environmentRepository);
    deleteEnvironmentRequest =
        DeleteEnvironmentRequest.builder()
            .environmentId(wrapperEnvironmentId)
            .forceDelete(true)
            .build();
    when(apiModelMapper.toModelEnvironmentId(wrapperEnvironmentId)).thenReturn(environmentId);
    when(apiModelMapper.toWrapperEnvironment(environment)).thenReturn(wrapperEnvironment);
  }

  @Test
  public void testDeleteEnvironmentSuccess() throws Exception {
    when(environmentRepository.listEnvironmentRevisions(environmentId))
        .thenReturn(Collections.singletonList(environmentRevision));
    doNothing().when(environmentRepository).deleteEnvironmentRevision(environmentRevision);
    when(environmentRepository.getEnvironment(environmentId)).thenReturn(environment);
    doNothing().when(environmentRepository).deleteEnvironment(environmentId);

    final DeleteEnvironmentResponse deleteEnvironmentResponse =
        deleteEnvironmentApi.deleteEnvironment(deleteEnvironmentRequest);
    assertEquals(wrapperEnvironment, deleteEnvironmentResponse.getEnvironment());
  }

  @Test(expected = ResourceNotFoundException.class)
  public void testDeleteEnvironmentResourceNotFoundException() throws Exception {
    when(environmentRepository.getEnvironment(environmentId)).thenThrow(resourceNotFoundException);

    deleteEnvironmentApi.deleteEnvironment(deleteEnvironmentRequest);
  }

  @Test(expected = InternalServiceException.class)
  public void testDeleteEnvironmentInternalServiceException() throws Exception {
    when(environmentRepository.listEnvironmentRevisions(environmentId))
        .thenReturn(Collections.singletonList(environmentRevision));
    doThrow(InternalServiceException.class)
        .when(environmentRepository)
        .deleteEnvironmentRevision(environmentRevision);

    deleteEnvironmentApi.deleteEnvironment(deleteEnvironmentRequest);
  }

  @Test(expected = InternalServiceException.class)
  public void testDeleteEnvironmentUnhandledException() throws Exception {
    when(environmentRepository.listEnvironmentRevisions(environmentId))
        .thenReturn(Collections.singletonList(environmentRevision));
    doThrow(NullPointerException.class)
        .when(environmentRepository)
        .deleteEnvironmentRevision(environmentRevision);

    deleteEnvironmentApi.deleteEnvironment(deleteEnvironmentRequest);
  }
}
