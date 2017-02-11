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

package tests

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/blox/blox/cluster-state-service/canary/pkg/tests/util"
	"github.com/blox/blox/cluster-state-service/canary/pkg/wrappers"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

type canaryTest func() error
type metricName func() string

// CanaryRunner defines methods required to run a test
type CanaryRunner interface {
	Run(test canaryTest, metric metricName)
}

type canaryRunner struct {
	cloudwatchWrapper wrappers.CloudwatchWrapper
}

// NewCanaryRunner generates a new runner for the canary
func NewCanaryRunner(sess *session.Session) (CanaryRunner, error) {
	if sess == nil {
		return nil, errors.New("AWS session has to be initialized to initialize the canary runner. ")
	}
	cloudwatchWrapper, err := wrappers.NewCloudwatchWrapper(sess)
	if err != nil {
		return nil, err
	}
	return canaryRunner{
		cloudwatchWrapper: cloudwatchWrapper,
	}, nil
}

// Run runs the test 'test' and emits a success or a failure metric
// for 'metricName'
func (runner canaryRunner) Run(test canaryTest, metric metricName) {
	m := metric()
	err := test()
	var cwErr error
	if err != nil {
		log.Error(err.Error())
		if _, ok := errors.Cause(err).(util.CleanUpError); !ok {
			cwErr = runner.cloudwatchWrapper.EmitFailureMetric(&m)
		}
	} else {
		cwErr = runner.cloudwatchWrapper.EmitSuccessMetric(&m)
	}
	if cwErr != nil {
		log.Error(cwErr)
	}
}
