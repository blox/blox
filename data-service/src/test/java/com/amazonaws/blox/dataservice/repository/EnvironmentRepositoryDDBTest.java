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
package com.amazonaws.blox.dataservice.repository;

import static com.amazonaws.blox.dataservice.repository.model.EnvironmentDDBRecord.ENVIRONMENT_NAME_RANGE_KEY;
import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.Assert.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.doNothing;
import static org.mockito.Mockito.doThrow;
import static org.mockito.Mockito.eq;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import com.amazonaws.AmazonServiceException;
import com.amazonaws.blox.dataservice.exception.ResourceType;
import com.amazonaws.blox.dataservice.mapper.EnvironmentMapper;
import com.amazonaws.blox.dataservice.model.Cluster;
import com.amazonaws.blox.dataservice.model.Environment;
import com.amazonaws.blox.dataservice.model.EnvironmentId;
import com.amazonaws.blox.dataservice.model.EnvironmentRevision;
import com.amazonaws.blox.dataservice.repository.model.EnvironmentDDBRecord;
import com.amazonaws.blox.dataservice.repository.model.EnvironmentRevisionDDBRecord;
import com.amazonaws.blox.dataservicemodel.v1.exception.InternalServiceException;
import com.amazonaws.blox.dataservicemodel.v1.exception.ResourceNotFoundException;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBMapper;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBMapperConfig;
import com.amazonaws.services.dynamodbv2.datamodeling.DynamoDBQueryExpression;
import com.amazonaws.services.dynamodbv2.datamodeling.PaginatedQueryList;
import com.amazonaws.services.dynamodbv2.model.AttributeValue;
import com.amazonaws.services.dynamodbv2.model.ComparisonOperator;
import com.amazonaws.services.dynamodbv2.model.Condition;
import java.util.List;
import java.util.Map;
import java.util.stream.Stream;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.junit.runner.RunWith;
import org.mockito.ArgumentCaptor;
import org.mockito.Captor;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;

@RunWith(MockitoJUnitRunner.StrictStubs.class)
public class EnvironmentRepositoryDDBTest {
  private static final String ACCOUNT_ID = "123456789012";
  private static final String CLUSTER = "mycluster";
  private static final String ENVIRONMENT_NAME = "myenv";
  private static final String ENVIRONMENT_REVISION_ID = "revision-id";
  private static final DynamoDBMapperConfig CONFIG =
      DynamoDBMapperConfig.SaveBehavior.CLOBBER.config();

  @Rule public ExpectedException thrown = ExpectedException.none();

  @Mock private DynamoDBMapper dynamoDBMapper;
  @Mock private EnvironmentMapper environmentMapper;
  @Mock private EnvironmentDDBRecord environmentDDBRecord;
  @Mock private Environment environment;
  @Mock private EnvironmentRevisionDDBRecord environmentRevisionDDBRecord;
  @Mock private EnvironmentRevision environmentRevision;
  @Mock private PaginatedQueryList<EnvironmentDDBRecord> environmentDDBRecords;

  @InjectMocks private EnvironmentRepositoryDDB environmentRepositoryDDB;

  @Captor private ArgumentCaptor<DynamoDBQueryExpression> ddbQueryExpressionCaptor;
  @Captor private ArgumentCaptor<EnvironmentDDBRecord> environmentDDBRecordArgumentCaptor;

  private EnvironmentId environmentId;
  private Cluster cluster;

  @Before
  public void setUp() {
    environmentId =
        EnvironmentId.builder()
            .accountId(ACCOUNT_ID)
            .cluster(CLUSTER)
            .environmentName(ENVIRONMENT_NAME)
            .build();

    cluster = Cluster.builder().accountId(ACCOUNT_ID).clusterName(CLUSTER).build();
  }

  @Test
  public void testGetEnvironmentRevision() throws Exception {
    // Given
    when(dynamoDBMapper.load(
            EnvironmentRevisionDDBRecord.class,
            environmentId.generateAccountIdClusterEnvironmentName(),
            ENVIRONMENT_REVISION_ID))
        .thenReturn(environmentRevisionDDBRecord);
    when(environmentMapper.toEnvironmentRevision(environmentRevisionDDBRecord))
        .thenReturn(environmentRevision);

    // When
    final EnvironmentRevision result =
        environmentRepositoryDDB.getEnvironmentRevision(environmentId, ENVIRONMENT_REVISION_ID);

    // Then
    verify(dynamoDBMapper)
        .load(
            EnvironmentRevisionDDBRecord.class,
            environmentId.generateAccountIdClusterEnvironmentName(),
            ENVIRONMENT_REVISION_ID);
    verify(environmentMapper).toEnvironmentRevision(environmentRevisionDDBRecord);
    assertThat(result).isEqualTo(environmentRevision);
  }

  @Test
  public void testGetEnvironmentRevisionNotFound() throws Exception {
    // Given
    when(dynamoDBMapper.load(
            EnvironmentRevisionDDBRecord.class,
            environmentId.generateAccountIdClusterEnvironmentName(),
            ENVIRONMENT_REVISION_ID))
        .thenReturn(null);

    // Expected exception
    thrown.expect(ResourceNotFoundException.class);
    thrown.expectMessage(
        String.format(
            "%s with id %s could not be found",
            ResourceType.ENVIRONMENT_REVISION, ENVIRONMENT_REVISION_ID));

    // When
    environmentRepositoryDDB.getEnvironmentRevision(environmentId, ENVIRONMENT_REVISION_ID);
  }

  @Test
  public void testGetEnvironmentRevisionInternalError() throws Exception {
    // Given
    when(dynamoDBMapper.load(
            EnvironmentRevisionDDBRecord.class,
            environmentId.generateAccountIdClusterEnvironmentName(),
            ENVIRONMENT_REVISION_ID))
        .thenThrow(AmazonServiceException.class);

    // Expected exception
    thrown.expect(InternalServiceException.class);

    // When
    environmentRepositoryDDB.getEnvironmentRevision(environmentId, ENVIRONMENT_REVISION_ID);
  }

  @Test
  public void testListEnvironmentsWithEnvironmentNamePrefix() throws Exception {
    final String environmentNamePrefix = "environmentNamePrefix";

    when(environmentDDBRecords.stream()).thenReturn(Stream.of(environmentDDBRecord));
    when(dynamoDBMapper.query(eq(EnvironmentDDBRecord.class), any(DynamoDBQueryExpression.class)))
        .thenReturn(environmentDDBRecords);
    when(environmentMapper.toEnvironment(environmentDDBRecord)).thenReturn(environment);

    final List<Environment> result =
        environmentRepositoryDDB.listEnvironments(cluster, environmentNamePrefix);

    verify(dynamoDBMapper)
        .query(eq(EnvironmentDDBRecord.class), ddbQueryExpressionCaptor.capture());
    verify(environmentMapper).toEnvironment(environmentDDBRecord);

    final EnvironmentDDBRecord queriedEnvironmentDDBRecord =
        (EnvironmentDDBRecord) ddbQueryExpressionCaptor.getValue().getHashKeyValues();
    final Map<String, Condition> queryConditions =
        ddbQueryExpressionCaptor.getValue().getRangeKeyConditions();

    assertThat(queriedEnvironmentDDBRecord.getAccountIdCluster())
        .isEqualTo(cluster.generateAccountIdCluster());
    assertThat(queryConditions).isNotEmpty().hasSize(1);
    assertThat(queryConditions).containsKey(ENVIRONMENT_NAME_RANGE_KEY);
    assertThat(queryConditions)
        .containsValue(
            new Condition()
                .withComparisonOperator(ComparisonOperator.BEGINS_WITH)
                .withAttributeValueList(new AttributeValue().withS(environmentNamePrefix)));
    assertThat(result.size()).isEqualTo(1);
    assertThat(result.get(0)).isEqualTo(environment);
  }

  @Test
  public void testListEnvironmentsWithoutEnvironmentNamePrefix() throws Exception {
    when(environmentDDBRecords.stream()).thenReturn(Stream.of(environmentDDBRecord));
    when(dynamoDBMapper.query(eq(EnvironmentDDBRecord.class), any(DynamoDBQueryExpression.class)))
        .thenReturn(environmentDDBRecords);
    when(environmentMapper.toEnvironment(environmentDDBRecord)).thenReturn(environment);

    final List<Environment> result = environmentRepositoryDDB.listEnvironments(cluster, null);

    verify(dynamoDBMapper)
        .query(eq(EnvironmentDDBRecord.class), ddbQueryExpressionCaptor.capture());
    verify(environmentMapper).toEnvironment(environmentDDBRecord);

    final EnvironmentDDBRecord queriedEnvironmentDDBRecord =
        (EnvironmentDDBRecord) ddbQueryExpressionCaptor.getValue().getHashKeyValues();

    assertThat(ddbQueryExpressionCaptor.getValue().getRangeKeyConditions()).isNull();
    assertThat(queriedEnvironmentDDBRecord.getAccountIdCluster())
        .isEqualTo(cluster.generateAccountIdCluster());
    assertThat(result.size()).isEqualTo(1);
    assertThat(result.get(0)).isEqualTo(environment);
  }

  @Test
  public void testListEnvironmentsEmptyResult() throws Exception {
    when(environmentDDBRecords.stream()).thenReturn(Stream.empty());
    when(dynamoDBMapper.query(eq(EnvironmentDDBRecord.class), any(DynamoDBQueryExpression.class)))
        .thenReturn(environmentDDBRecords);

    final List<Environment> result = environmentRepositoryDDB.listEnvironments(cluster, null);

    verify(dynamoDBMapper).query(eq(EnvironmentDDBRecord.class), any());
    verify(environmentMapper, never()).toEnvironment(environmentDDBRecord);
    assertThat(result.size()).isEqualTo(0);
  }

  @Test
  public void testListEnvironmentsInternalError() throws Exception {
    when(dynamoDBMapper.query(eq(EnvironmentDDBRecord.class), any(DynamoDBQueryExpression.class)))
        .thenThrow(AmazonServiceException.class);

    thrown.expect(InternalServiceException.class);
    thrown.expectMessage(
        String.format("Could not query environments for cluster %s", cluster.toString()));

    environmentRepositoryDDB.listEnvironments(cluster, null);
  }

  @Test
  public void testDeleteEnvironmentSuccess() throws Exception {
    doNothing().when(dynamoDBMapper).delete(any(EnvironmentDDBRecord.class));

    environmentRepositoryDDB.deleteEnvironment(environmentId);
    verify(dynamoDBMapper, times(1)).delete(environmentDDBRecordArgumentCaptor.capture());
    assertEquals(
        environmentId.generateAccountIdCluster(),
        environmentDDBRecordArgumentCaptor.getValue().getAccountIdCluster());
    assertEquals(
        environmentId.getEnvironmentName(),
        environmentDDBRecordArgumentCaptor.getValue().getEnvironmentName());
  }

  @Test
  public void testDeleteEnvironmentInternalError() throws Exception {
    doThrow(AmazonServiceException.class)
        .when(dynamoDBMapper)
        .delete(any(EnvironmentDDBRecord.class));

    thrown.expect(InternalServiceException.class);
    thrown.expectMessage(
        "Fail to delete environment with accountIdCluster 123456789012/mycluster and environment name myenv");

    environmentRepositoryDDB.deleteEnvironment(environmentId);
  }
}
