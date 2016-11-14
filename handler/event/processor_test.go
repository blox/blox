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
	"encoding/json"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/mocks"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/store"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"testing"
)

const (
	unknownEventType = "unknown"
)

type event struct {
	DetailType string `json:"detail-type"`
}

type processorMockContext struct {
	mockCtrl      *gomock.Controller
	stores        store.Stores
	taskStore     *mocks.MockTaskStore
	instanceStore *mocks.MockContainerInstanceStore
}

func NewProcessorMockContext(t *testing.T) *processorMockContext {
	context := processorMockContext{}
	context.mockCtrl = gomock.NewController(t)

	context.taskStore = mocks.NewMockTaskStore(context.mockCtrl)
	context.instanceStore = mocks.NewMockContainerInstanceStore(context.mockCtrl)

	context.stores = store.Stores{
		TaskStore:              context.taskStore,
		ContainerInstanceStore: context.instanceStore,
	}

	return &context
}

func TestNewProcessor(t *testing.T) {
	context := NewProcessorMockContext(t)
	defer context.mockCtrl.Finish()

	p := NewProcessor(context.stores)
	if p == nil {
		t.Error("NewProcessor returns nil")
	}
}

func TestProcessEventEmptyString(t *testing.T) {
	context := NewProcessorMockContext(t)
	defer context.mockCtrl.Finish()

	p := NewProcessor(context.stores)
	err := p.ProcessEvent("")

	if err == nil {
		t.Error("Expected ProcessEvent to return an error when passed an empty string")
	}
}

func TestProcessEventInvalidJson(t *testing.T) {
	context := NewProcessorMockContext(t)
	defer context.mockCtrl.Finish()

	p := NewProcessor(context.stores)

	err := p.ProcessEvent("invalidJson")

	if err == nil {
		t.Error("Expected ProcessEvent to return an error when passed an event with an unknown event type")
	}
}

func TestProcessEventUnknownEventType(t *testing.T) {
	context := NewProcessorMockContext(t)
	defer context.mockCtrl.Finish()

	p := NewProcessor(context.stores)

	e := event{
		DetailType: unknownEventType,
	}
	eventjson, _ := json.Marshal(e)

	err := p.ProcessEvent(string(eventjson))

	if err == nil {
		t.Error("Expected ProcessEvent to return an error when passed an event with an unknown event type")
	}
}

func TestProcessEventTaskEventFails(t *testing.T) {
	context := NewProcessorMockContext(t)
	defer context.mockCtrl.Finish()

	p := NewProcessor(context.stores)

	e := event{
		DetailType: taskType,
	}
	eventjson, _ := json.Marshal(e)

	context.taskStore.EXPECT().AddTask(string(eventjson)).Return(errors.New("AddTask failed"))

	err := p.ProcessEvent(string(eventjson))

	if err == nil {
		t.Error("Expected ProcessEvent to return an error when AddTask fails")
	}
}

func TestProcessEventTaskEvent(t *testing.T) {
	context := NewProcessorMockContext(t)
	defer context.mockCtrl.Finish()

	p := NewProcessor(context.stores)

	e := event{
		DetailType: taskType,
	}
	eventjson, _ := json.Marshal(e)

	context.taskStore.EXPECT().AddTask(string(eventjson)).Return(nil)

	err := p.ProcessEvent(string(eventjson))

	if err != nil {
		t.Error("Unexpected error in ProcessEvent")
	}
}

func TestProcessEventInstanceEventFails(t *testing.T) {
	context := NewProcessorMockContext(t)
	defer context.mockCtrl.Finish()

	p := NewProcessor(context.stores)

	e := event{
		DetailType: containerInstanceType,
	}
	eventjson, _ := json.Marshal(e)

	context.instanceStore.EXPECT().AddContainerInstance(string(eventjson)).Return(errors.New("AddInstance failed"))

	err := p.ProcessEvent(string(eventjson))

	if err == nil {
		t.Error("Expected ProcessEvent to return an error when AddInstance fails")
	}
}

func TestProcessEventInstanceEvent(t *testing.T) {
	context := NewProcessorMockContext(t)
	defer context.mockCtrl.Finish()

	p := NewProcessor(context.stores)

	e := event{
		DetailType: containerInstanceType,
	}
	eventjson, _ := json.Marshal(e)

	context.instanceStore.EXPECT().AddContainerInstance(string(eventjson)).Return(nil)

	err := p.ProcessEvent(string(eventjson))

	if err != nil {
		t.Error("Unexpected error in ProcessEvent")
	}
}
