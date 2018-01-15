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

import com.amazonaws.blox.dataservice.mapper.ApiModelMapper;
import com.amazonaws.blox.dataservice.model.Cluster;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentId;
import com.amazonaws.blox.dataservice.repository.EnvironmentRepository;
import com.amazonaws.blox.dataservicemodel.v1.exception.InternalServiceException;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsRequest;
import com.amazonaws.blox.dataservicemodel.v1.model.wrappers.ListEnvironmentsResponse;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;

import java.util.Collections;

import static org.hamcrest.Matchers.is;
import static org.junit.Assert.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@RunWith(MockitoJUnitRunner.StrictStubs.class)
public class ListEnvironmentsApiTest {
  @Mock private ApiModelMapper apiModelMapper;
  @Mock private EnvironmentRepository environmentRepository;
  @Mock private com.amazonaws.blox.dataservicemodel.v1.model.Cluster clusterWrapper;
  @Mock private Cluster cluster;
  @Mock private com.amazonaws.blox.dataservicemodel.v1.model.EnvironmentId environmentIdWrapper;
  @Mock private EnvironmentId environmentId;
  @Mock private Environment environment;

  @InjectMocks private ListEnvironmentsApi api;

  private ListEnvironmentsRequest request;

  @Before
  public void setup() {
    request = ListEnvironmentsRequest.builder().cluster(clusterWrapper).build();
    when(apiModelMapper.toModelCluster(clusterWrapper)).thenReturn(cluster);
    when(environment.getEnvironmentId()).thenReturn(environmentId);
    when(apiModelMapper.toWrapperEnvironmentId(environmentId)).thenReturn(environmentIdWrapper);
  }

  @Test
  public void testListEnvironments() throws Exception {
    when(environmentRepository.listEnvironments(cluster))
        .thenReturn(Collections.singletonList(environment));

    final ListEnvironmentsResponse response = api.listEnvironments(request);

    verify(apiModelMapper).toModelCluster(clusterWrapper);
    verify(environmentRepository).listEnvironments(cluster);
    verify(apiModelMapper).toWrapperEnvironmentId(environmentId);

    assertThat(response.getEnvironmentIds().size(), is(1));
    assertThat(response.getEnvironmentIds().get(0), is(environmentIdWrapper));
  }

  @Test
  public void testListEnvironmentsEmptyResult() throws Exception {
    when(environmentRepository.listEnvironments(cluster)).thenReturn(Collections.emptyList());

    final ListEnvironmentsResponse response = api.listEnvironments(request);

    verify(apiModelMapper).toModelCluster(clusterWrapper);
    verify(environmentRepository).listEnvironments(cluster);
    verify(apiModelMapper, never()).toWrapperEnvironmentId(any());

    assertThat(response.getEnvironmentIds().size(), is(0));
  }

  @Test(expected = InternalServiceException.class)
  public void testListEnvironmentsInternalServiceException() throws Exception {
    when(environmentRepository.listEnvironments(cluster))
        .thenThrow(new InternalServiceException(""));

    api.listEnvironments(request);
  }

  @Test(expected = InternalServiceException.class)
  public void testListEnvironmentsUnknownException() throws Exception {
    when(environmentRepository.listEnvironments(cluster)).thenThrow(new IllegalStateException());

    api.listEnvironments(request);
  }
}
