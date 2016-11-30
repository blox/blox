// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package httpclient

import (
	"net/http"
	"testing"

	"github.com/blox/blox/cluster-state-service/handler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UATestSuite struct {
	suite.Suite
	roundtripper     *mocks.MockRoundTripper
	bloxRoundTripper *bloxRoundTripper
}

func (testSuite *UATestSuite) SetupTest() {
	mockCtrl := gomock.NewController(testSuite.T())
	testSuite.roundtripper = mocks.NewMockRoundTripper(mockCtrl)
	testSuite.bloxRoundTripper = &bloxRoundTripper{}
	testSuite.bloxRoundTripper.transport = testSuite.roundtripper
}

func TestUATestSuite(t *testing.T) {
	suite.Run(t, new(UATestSuite))
}

func (testSuite *UATestSuite) TestBloxUA() {

	req := &http.Request{
		Header: make(http.Header),
	}
	rsp := &http.Response{}

	testSuite.roundtripper.EXPECT().RoundTrip(req).Return(nil, nil)

	rsp, err := testSuite.bloxRoundTripper.RoundTrip(req)

	actualUserAgent := req.Header.Get(userAgentHeader)
	testSuite.T().Logf("UA : %s", actualUserAgent)
	assert.Equal(testSuite.T(), actualUserAgent, userAgent,
		"Blox CSS user agent with unexpected value. Expected: %s Actual: %s", userAgent, actualUserAgent)
	assert.Nil(testSuite.T(), err, "Unexpected error when calling RoundTrip")
	assert.Nil(testSuite.T(), rsp, "Unexpected response when calling RoundTrip")
}
