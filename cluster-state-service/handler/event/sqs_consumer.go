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
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	sqsVisibilityTimeout = 10
	sqsWaitTimeSeconds   = 10
)

type sqsEventConsumer struct {
	sqs       sqsiface.SQSAPI
	queueURL  string
	processor Processor
}

func NewSQSConsumer(sqs sqsiface.SQSAPI, processor Processor, queueName string) (Consumer, error) {
	if sqs == nil {
		return nil, errors.Errorf("The SQS API interface is not initialized")
	}
	if processor == nil {
		return nil, errors.Errorf("The event processor is not initialized")
	}
	if queueName == "" {
		return nil, errors.Errorf("The SQS queue name is empty")
	}

	sqsQueueURL, err := getQueueURL(sqs, queueName)
	if err != nil {
		return nil, err
	}

	return &sqsEventConsumer{
		sqs:       sqs,
		queueURL:  sqsQueueURL,
		processor: processor,
	}, nil
}

func getQueueURL(client sqsiface.SQSAPI, queueName string) (string, error) {
	if client == nil {
		return "", errors.Errorf("The SQS client is not initialized")
	}
	if len(queueName) == 0 {
		return "", errors.Errorf("The queue name cannot be empty")
	}

	input := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}
	output, err := client.GetQueueUrl(input)

	if err != nil {
		return "", errors.Wrapf(err, "Could not get the queue url from queue name: %s", queueName)
	} else if output.QueueUrl == nil {
		return "", errors.Errorf("Queue url is empty: %v", output)
	}

	return aws.StringValue(output.QueueUrl), nil
}

func (sqsConsumer sqsEventConsumer) PollForEvents(ctx context.Context) {
	log.Infof("Starting to poll for events from SQS")
	statsTicker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-ctx.Done():
			return
		case <-statsTicker.C:
			go func() {
				sqsConsumer.logQueueStats(ctx)
			}()
		default:
			sqsConsumer.pollForMessages()
		}
	}
}

func (sqsConsumer sqsEventConsumer) logQueueStats(ctx context.Context) error {
	params := &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(sqsConsumer.queueURL),
		AttributeNames: []*string{
			aws.String("ApproximateNumberOfMessages"),
			aws.String("ApproximateNumberOfMessagesDelayed"),
			aws.String("ApproximateNumberOfMessagesNotVisible"),
		},
	}
	resp, err := sqsConsumer.sqs.GetQueueAttributes(params)
	if err != nil {
		return errors.Wrapf(err, "Error getting queue attributes for queue %s", sqsConsumer.queueURL)
	}
	for attrib, _ := range resp.Attributes {
		prop := resp.Attributes[attrib]
		i, _ := strconv.Atoi(*prop)
		log.Infof("SQS attribute[%s] = %d", attrib, i)
	}
	return nil
}

func (sqsConsumer sqsEventConsumer) pollForMessages() {
	receiveMessageInput := &sqs.ReceiveMessageInput{
		QueueUrl:          aws.String(sqsConsumer.queueURL),
		VisibilityTimeout: aws.Int64(sqsVisibilityTimeout),
		WaitTimeSeconds:   aws.Int64(sqsWaitTimeSeconds),
	}

	output, err := sqsConsumer.sqs.ReceiveMessage(receiveMessageInput)
	if err != nil {
		// wrap to get stack trace
		err = errors.Wrap(err, "Could not poll sqs")
		log.Errorf("%+v", err)
		return
	}

	if output == nil || output.Messages == nil {
		return
	}

	sqsConsumer.processMessages(output.Messages)
}

func (sqsConsumer sqsEventConsumer) processMessages(messages []*sqs.Message) {
	for _, message := range messages {
		err := sqsConsumer.processEvent(message)
		if err != nil {
			log.Errorf("Could not process message: %v: %+v", message, err)
			continue
		}

		err = sqsConsumer.deleteEvent(message)
		if err != nil {
			log.Errorf("Could not delete message %v: %+v", message, err)
		}
	}
}

func (sqsConsumer sqsEventConsumer) processEvent(message *sqs.Message) error {
	if message == nil {
		return errors.Errorf("The sqs message cannot be nil")
	}
	if message.Body == nil {
		return errors.Errorf("The sqs message body cannot be empty")
	}
	return sqsConsumer.processor.ProcessEvent(*message.Body)
}

func (sqsConsumer sqsEventConsumer) deleteEvent(message *sqs.Message) error {
	if message == nil {
		return errors.Errorf("The sqs message cannot be nil")
	}
	if message.ReceiptHandle == nil {
		return errors.Errorf("The sqs message receipt handle cannot be empty")
	}

	deleteMessageInput := &sqs.DeleteMessageInput{
		ReceiptHandle: message.ReceiptHandle,
		QueueUrl:      aws.String(sqsConsumer.queueURL),
	}

	_, err := sqsConsumer.sqs.DeleteMessage(deleteMessageInput)
	if err != nil {
		return errors.Wrap(err, "Could not delete message")
	}

	return nil
}
