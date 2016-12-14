// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package event

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/blox/blox/cluster-state-service/handler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"testing"
)

const (
	streamName          = "test"
	kinesisMessageBody1 = "messageBody"
	kinesisMessageBody2 = "messageBody2"
)

type consumerMockKinesisContext struct {
	mockCtrl                      *gomock.Controller
	kinesisClient                 *mocks.MockKinesisAPI
	processor                     *mocks.MockProcessor
	getShardIteratorInput         *kinesis.GetShardIteratorInput
	getShardIteratorOutput        *kinesis.GetShardIteratorOutput
	getRecordsInput               *kinesis.GetRecordsInput
	getRecordsSecondInput         *kinesis.GetRecordsInput
	getRecordsNoMessagesOutput    *kinesis.GetRecordsOutput
	getRecordsFirstMessageOutput  *kinesis.GetRecordsOutput
	getRecordsSecondMessageOutput *kinesis.GetRecordsOutput
	getRecordsTwoMessagesOutput   *kinesis.GetRecordsOutput
	record1                       *kinesis.Record
	record2                       *kinesis.Record
	shardIteratorFromGetRecords   *string
}

func NewConsumerMockKinesisContext(t *testing.T) *consumerMockKinesisContext {
	context := consumerMockKinesisContext{}
	context.mockCtrl = gomock.NewController(t)
	context.kinesisClient = mocks.NewMockKinesisAPI(context.mockCtrl)
	context.processor = mocks.NewMockProcessor(context.mockCtrl)
	context.shardIteratorFromGetRecords = aws.String("getRecordsIterator")

	context.record1 = &kinesis.Record{
		Data: []byte(kinesisMessageBody1),
	}

	context.record2 = &kinesis.Record{
		Data: []byte(kinesisMessageBody2),
	}

	context.getShardIteratorInput = &kinesis.GetShardIteratorInput{
		ShardId:           aws.String("shardId-000000000000"),
		ShardIteratorType: aws.String("TRIM_HORIZON"),
		StreamName:        aws.String(streamName),
	}

	context.getShardIteratorOutput = &kinesis.GetShardIteratorOutput{
		ShardIterator: aws.String("shardId-123"),
	}

	context.getRecordsInput = &kinesis.GetRecordsInput{
		Limit:         aws.Int64(100),
		ShardIterator: aws.String("shardId-123"),
	}

	context.getRecordsSecondInput = &kinesis.GetRecordsInput{
		Limit:         aws.Int64(100),
		ShardIterator: context.shardIteratorFromGetRecords,
	}

	context.getRecordsNoMessagesOutput = &kinesis.GetRecordsOutput{
		Records:           []*kinesis.Record{},
		NextShardIterator: context.shardIteratorFromGetRecords,
	}

	context.getRecordsFirstMessageOutput = &kinesis.GetRecordsOutput{
		Records:           []*kinesis.Record{context.record1},
		NextShardIterator: context.shardIteratorFromGetRecords,
	}

	context.getRecordsSecondMessageOutput = &kinesis.GetRecordsOutput{
		Records:           []*kinesis.Record{context.record2},
		NextShardIterator: context.shardIteratorFromGetRecords,
	}

	context.getRecordsTwoMessagesOutput = &kinesis.GetRecordsOutput{
		Records:           []*kinesis.Record{context.record1, context.record2},
		NextShardIterator: context.shardIteratorFromGetRecords,
	}

	return &context
}

func TestNewConsumerNilKinesis(t *testing.T) {
	context := NewConsumerMockKinesisContext(t)
	defer context.mockCtrl.Finish()

	_, err := NewKinesisConsumer(nil, context.processor, streamName)
	if err == nil {
		t.Error("Expected an error when kinesis is nil")
	}
}

func TestNewConsumerKinesisNilProcessor(t *testing.T) {
	context := NewConsumerMockKinesisContext(t)
	defer context.mockCtrl.Finish()

	_, err := NewKinesisConsumer(context.kinesisClient, nil, streamName)
	if err == nil {
		t.Error("Expected an error when processor is nil")
	}
}

func TestNewConsumerKinesisEmptyQueueName(t *testing.T) {
	context := NewConsumerMockKinesisContext(t)
	defer context.mockCtrl.Finish()

	_, err := NewKinesisConsumer(context.kinesisClient, context.processor, "")
	if err == nil {
		t.Error("Expected an error when stream name is empty")
	}
}

func TestPollForKinesisEventsSingleMessages(t *testing.T) {
	mockContext := NewConsumerMockKinesisContext(t)
	defer mockContext.mockCtrl.Finish()

	mockContext.kinesisClient.EXPECT().GetShardIterator(gomock.Eq(mockContext.getShardIteratorInput)).Return(mockContext.getShardIteratorOutput, nil)

	c, err := NewKinesisConsumer(mockContext.kinesisClient, mockContext.processor, streamName)

	if err != nil {
		t.Errorf("Unexpected error when calling NewConsumer: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	mockContext.kinesisClient.EXPECT().GetRecords(mockContext.getRecordsInput).Return(mockContext.getRecordsFirstMessageOutput, nil)
	mockContext.processor.EXPECT().ProcessEvent(kinesisMessageBody1).Return(nil)
	mockContext.kinesisClient.EXPECT().GetRecords(mockContext.getRecordsSecondInput).Return(mockContext.getRecordsSecondMessageOutput, nil)
	mockContext.processor.EXPECT().ProcessEvent(kinesisMessageBody2).Return(nil).Do(func(x interface{}) {
		cancel()
	})

	c.PollForEvents(ctx)
}

func TestPollForKinesisEventsShardIteratorFailsOnce(t *testing.T) {
	mockContext := NewConsumerMockKinesisContext(t)
	defer mockContext.mockCtrl.Finish()

	mockContext.kinesisClient.EXPECT().GetShardIterator(gomock.Eq(mockContext.getShardIteratorInput)).Return(nil, errors.New("Shard iterator call failed."))
	mockContext.kinesisClient.EXPECT().GetShardIterator(gomock.Eq(mockContext.getShardIteratorInput)).Return(mockContext.getShardIteratorOutput, nil)

	c, err := NewKinesisConsumer(mockContext.kinesisClient, mockContext.processor, streamName)

	if err != nil {
		t.Errorf("Unexpected error when calling NewConsumer: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	mockContext.kinesisClient.EXPECT().GetRecords(mockContext.getRecordsInput).Return(mockContext.getRecordsFirstMessageOutput, nil)
	mockContext.processor.EXPECT().ProcessEvent(kinesisMessageBody1).Return(nil).Do(func(x interface{}) {
		cancel()
	})

	c.PollForEvents(ctx)
}

func TestPollForKinesisEventsGetRecordsFailsOnce(t *testing.T) {
	mockContext := NewConsumerMockKinesisContext(t)
	defer mockContext.mockCtrl.Finish()

	mockContext.kinesisClient.EXPECT().GetShardIterator(gomock.Eq(mockContext.getShardIteratorInput)).Return(mockContext.getShardIteratorOutput, nil)

	c, err := NewKinesisConsumer(mockContext.kinesisClient, mockContext.processor, streamName)

	if err != nil {
		t.Errorf("Unexpected error when calling NewConsumer: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	mockContext.kinesisClient.EXPECT().GetRecords(mockContext.getRecordsInput).Return(nil, errors.New("GetRecords call failed."))
	mockContext.kinesisClient.EXPECT().GetShardIterator(gomock.Eq(mockContext.getShardIteratorInput)).Return(mockContext.getShardIteratorOutput, nil)
	mockContext.kinesisClient.EXPECT().GetRecords(mockContext.getRecordsInput).Return(mockContext.getRecordsFirstMessageOutput, nil)
	mockContext.processor.EXPECT().ProcessEvent(kinesisMessageBody1).Return(nil).Do(func(x interface{}) {
		cancel()
	})

	c.PollForEvents(ctx)
}

func TestPollForKinesisEventsReceiveTwoMessages(t *testing.T) {
	mockContext := NewConsumerMockKinesisContext(t)
	defer mockContext.mockCtrl.Finish()

	mockContext.kinesisClient.EXPECT().GetShardIterator(gomock.Eq(mockContext.getShardIteratorInput)).Return(mockContext.getShardIteratorOutput, nil)

	c, err := NewKinesisConsumer(mockContext.kinesisClient, mockContext.processor, streamName)

	if err != nil {
		t.Errorf("Unexpected error when calling NewConsumer: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	pollCount := 0

	mockContext.kinesisClient.EXPECT().GetRecords(mockContext.getRecordsInput).Return(mockContext.getRecordsTwoMessagesOutput, nil)
	mockContext.processor.EXPECT().ProcessEvent(kinesisMessageBody1).Return(nil).Do(func(x interface{}) {
		pollCount++
		if pollCount == 2 {
			cancel()
		}
	})
	mockContext.processor.EXPECT().ProcessEvent(kinesisMessageBody2).Return(nil).Do(func(x interface{}) {
		pollCount++
		if pollCount == 2 {
			cancel()
		}
	})

	c.PollForEvents(ctx)
}
