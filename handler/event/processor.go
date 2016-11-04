// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the License). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the license file accompanying this file. This file is distributed
// on an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package event

import (
	"github.com/aws/amazon-ecs-event-stream-handler/handler/json"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/store"
	"github.com/pkg/errors"
)

// Unmarshal event message json by type
type eventType struct {
	Type string `json:"detail-type"`
}

// Detail-type in the event stream message must match one of these strings
const (
	taskType              = "ECS Task State Change"
	containerInstanceType = "ECS Container Instance State Change"
)

// Processor defines methods to process events
type Processor interface {
	ProcessEvent(event string) error
}

type eventProcessor struct {
	stores store.Stores
}

func NewProcessor(stores store.Stores) Processor {
	return eventProcessor{
		stores: stores,
	}
}

// ProcessEvent takes an event JSON, unmarhsals and stores it in the datastore
func (processor eventProcessor) ProcessEvent(event string) error {
	if len(event) == 0 {
		return errors.New("Event cannot be empty")
	}

	// Determine the type of event based on the detail-type in the message
	var et eventType
	err := json.UnmarshalJSON(event, &et)
	if err != nil {
		return err
	}

	switch et.Type {
	case taskType:
		err = processor.stores.TaskStore.AddTask(event)
		if err != nil {
			return err
		}

	case containerInstanceType:
		err = processor.stores.ContainerInstanceStore.AddContainerInstance(event)
		if err != nil {
			return err
		}

	default:
		return errors.Errorf("Unrecognized task type: %v", et.Type)
	}

	return nil
}
