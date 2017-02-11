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

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/pkg/errors"
)

var (
	logNamespace = "CSSCanary"
)

// CloudwatchWrapper defines methods to access wrapper methods to call Cloudwatch APIs
type CloudwatchWrapper interface {
	EmitSuccessMetric(metricName *string) error
	EmitFailureMetric(metricName *string) error
}

type cloudwatchClientWrapper struct {
	client cloudwatchiface.CloudWatchAPI
}

// NewCloudwatchWrapper returns a new CloudwatchWrapper for the canary
func NewCloudwatchWrapper(sess *session.Session) (CloudwatchWrapper, error) {
	if sess == nil {
		return nil, errors.New("AWS session has to be initialized to initialize the Cloudwatch client. ")
	}
	cloudwatchClient := cloudwatch.New(sess)
	return cloudwatchClientWrapper{
		client: cloudwatchClient,
	}, nil
}

// EmitSuccessMetric emits a metric with value 1 for metric 'metricName'
func (wrapper cloudwatchClientWrapper) EmitSuccessMetric(metricName *string) error {
	metricValue := float64(1)
	return wrapper.emitMetric(metricName, &metricValue)
}

// EmitSuccessMetric emits a metric with value 0 for metric 'metricName'
func (wrapper cloudwatchClientWrapper) EmitFailureMetric(metricName *string) error {
	metricValue := float64(0)
	return wrapper.emitMetric(metricName, &metricValue)
}

func (wrapper cloudwatchClientWrapper) emitMetric(metricName *string, metricValue *float64) error {
	currentTime := time.Now()
	metricData := cloudwatch.MetricDatum{
		MetricName: metricName,
		Value:      metricValue,
		Timestamp:  &currentTime,
	}
	in := cloudwatch.PutMetricDataInput{
		Namespace:  &logNamespace,
		MetricData: []*cloudwatch.MetricDatum{&metricData},
	}

	_, err := wrapper.client.PutMetricData(&in)
	if err != nil {
		return errors.Wrapf(err, "Failed to emit Cloudwatch metric for '%s' wth value '%f'. ",
			*metricName, *metricValue)
	}
	return nil
}
