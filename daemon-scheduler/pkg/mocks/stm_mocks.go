// Copyright 2016-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package mocks

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

// MockSTM embeds the STM interface and provides custom interceptors for
// Get() and Put() methods. Since STM interface has private methods, a
// generated mock can't be used in its place. Hence, the intent is to
// use the mockSTM struct in its place.
type MockSTM struct {
	concurrency.STM
	// GetFunc is the interceptor for the Get() method in the STM interface
	GetFunc func(string) string
	// PutFunc is the interceptor for the Put() method in the STM interface
	PutFunc func(key string, val string, opts ...clientv3.OpOption)
}

// Get implements the STM.Get() method by invoking the custom interceptor
// method
func (stm *MockSTM) Get(key string) string {
	return stm.GetFunc(key)
}

// Put implements the STM.Put() method by invoking the custom interceptor
// method
func (stm *MockSTM) Put(key string, val string, opts ...clientv3.OpOption) {
	stm.PutFunc(key, val, opts...)
}
