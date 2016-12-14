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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
)

const (
	kinesisWaitTimeSeconds   = 10
	kinesisStartingShardId   = "shardId-000000000000"
	kinesisShardIteratorType = kinesis.ShardIteratorTypeTrimHorizon
	kinesisGetRecordsSize    = 100
)

type kinesisEventConsumer struct {
	kinesis    kinesisiface.KinesisAPI
	streamName string
	processor  Processor
	iterator   *string
}

func NewKinesisConsumer(kinesis kinesisiface.KinesisAPI, processor Processor, streamName string) (Consumer, error) {
	if kinesis == nil {
		return nil, errors.Errorf("The Kinesis API interface is not initialized")
	}
	if processor == nil {
		return nil, errors.Errorf("The event processor is not initialized")
	}
	if streamName == "" {
		return nil, errors.Errorf("The Kinesis stream name is empty")
	}

	return &kinesisEventConsumer{
		kinesis:    kinesis,
		streamName: streamName,
		processor:  processor,
	}, nil
}

func (kinesisConsumer *kinesisEventConsumer) PollForEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			kinesisConsumer.pollForMessages()
		}
	}
}

func (kinesisConsumer *kinesisEventConsumer) pollForMessages() {
	if kinesisConsumer.iterator == nil {
		iteratorRequest := &kinesis.GetShardIteratorInput{
			ShardId:           aws.String(kinesisStartingShardId),
			ShardIteratorType: aws.String(kinesisShardIteratorType),
			StreamName:        aws.String(kinesisConsumer.streamName),
		}
		iteratorResponse, err := kinesisConsumer.kinesis.GetShardIterator(iteratorRequest)
		if err != nil {
			log.Errorf("%+v", errors.Wrapf(err, "Could not get shard iterator"))
			return
		}
		kinesisConsumer.iterator = iteratorResponse.ShardIterator
	}

	recordsRequest := &kinesis.GetRecordsInput{
		Limit:         aws.Int64(kinesisGetRecordsSize),
		ShardIterator: kinesisConsumer.iterator,
	}
	recordsResponse, err := kinesisConsumer.kinesis.GetRecords(recordsRequest)
	if err != nil {
		log.Errorf("%+v", errors.Wrapf(err, "Unable to get records from kinesis"))
		kinesisConsumer.iterator = nil
		return
	}

	for _, record := range recordsResponse.Records {
		kinesisConsumer.processor.ProcessEvent(string(record.Data[:]))
	}

	kinesisConsumer.iterator = recordsResponse.NextShardIterator
	if len(recordsResponse.Records) == 0 {
		time.Sleep(kinesisWaitTimeSeconds * time.Second)
	}
}
