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

package v1

import "github.com/gorilla/mux"

const (
	nextToken       = "nextToken"
	deploymentToken = "deploymentToken"
	cluster         = "cluster"
)

var (
	// TODO - Use cluster ARN regex from common regex package
	clusterARNRegex = "(arn:aws:ecs:)([\\-\\w]+):[0-9]{12}:(cluster)/[a-zA-Z][a-zA-Z0-9_-]{1,254}"
	clusterARNVal   = "{" + cluster + ":" + clusterARNRegex + "}"
)

func NewRouter(api API) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	s := r.Path("/v1").Subrouter()

	// health

	s.Path("/ping").
		Methods("GET").
		HandlerFunc(api.Ping)

	// environment

	s.Path("/environments").
		Methods("POST").
		HandlerFunc(api.CreateEnvironment)

	s.Path("/environments/{name}").
		Methods("GET").
		HandlerFunc(api.GetEnvironment)

	s.Path("/environments").
		Queries(cluster, clusterARNVal).
		Methods("GET").
		HandlerFunc(api.FilterEnvironments)

	s.Path("/environments").
		Methods("GET").
		HandlerFunc(api.ListEnvironments)

	s.Path("/environments").
		Queries(nextToken, "").
		Methods("GET").
		HandlerFunc(api.ListEnvironments)

	s.Path("/environments/{name}").
		Methods("DELETE").
		HandlerFunc(api.DeleteEnvironment)

	// deployment

	s.Path("/environments/{name}/deployments").
		Queries(deploymentToken, "{deploymentToken:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}").
		Methods("POST").
		HandlerFunc(api.CreateDeployment)

	s.Path("/environments/{name}/deployments/{id}").
		Methods("GET").
		HandlerFunc(api.GetDeployment)

	s.Path("/environments/{name}/deployments").
		Methods("GET").
		HandlerFunc(api.ListDeployments)

	s.Path("/environments/{name}/deployments").
		Queries(nextToken, "").
		Methods("GET").
		HandlerFunc(api.ListDeployments)

	return s
}
