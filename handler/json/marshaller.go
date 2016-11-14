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

package json

import (
	"encoding/json"
	"github.com/pkg/errors"
)

// UnmarshalJSON unmarshals the given string into the given struct
func UnmarshalJSON(s string, t interface{}) error {
	if len(s) == 0 {
		return errors.New("Cannot unmarshal empty string")
	}

	if t == nil {
		return errors.New("UnmarshalJSON needs a non-nil interface to unmarshal into")
	}

	err := json.Unmarshal([]byte(s), &t)
	if err != nil {
		return errors.Wrapf(err, "Could not unmarshal string %v", s)
	}

	return nil
}
