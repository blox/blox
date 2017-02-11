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

package wrappers

import "github.com/aws/aws-sdk-go/aws/session"

// NewAWSSession creates an AWS session.
func NewAWSSession() (*session.Session, error) {
	var sess *session.Session
	var err error
	sess, err = session.NewSession()
	if err != nil {
		return nil, err
	}
	return sess, nil
}
