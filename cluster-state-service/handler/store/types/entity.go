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

package types

import (
	"fmt"
)

// Entity represents an object stored in Etcd.
type Entity struct {
	Key string // Etcd key
	Value string // Etcd value
	Version string // Etcd mod_revision
}

func (entity Entity) String() string {
	return fmt.Sprintf("Key: '%s'. Value: '%s'. Version: '%s'.", entity.Key, entity.Value, entity.Version)
}
