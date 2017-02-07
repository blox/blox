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
	"os"

	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/client"
	"github.com/blox/blox/daemon-scheduler/pkg/facade"
	httptransport "github.com/go-openapi/runtime/client"
)

const (
	defaultCSSEndpoint = "localhost:3000"
)

func NewClusterState() (facade.ClusterState, error) {
	endpoint := os.Getenv("CSS_ENDPOINT")
	if len(endpoint) == 0 {
		endpoint = defaultCSSEndpoint
	}
	transport := httptransport.New(endpoint, "/v1", []string{"http"})
	httpclient := client.New(transport, nil)
	return facade.NewClusterState(httpclient)
}
