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

package regex

const (
	clusterNameRegexWithoutStart = "[a-zA-Z][a-zA-Z0-9_-]{1,254}$"
	ClusterNameRegex             = "^" + clusterNameRegexWithoutStart
	ClusterARNRegex              = "^(arn:aws:ecs:)([\\-\\w]+):[0-9]{12}:(cluster)/" + clusterNameRegexWithoutStart
	ClusterNameAsARNSuffixRegex  = "/" + clusterNameRegexWithoutStart
	TaskARNRegex                 = "^(arn:aws:ecs):([\\-\\w]+):[0-9]{12}:(task)\\/[\\-\\w]+$"
	InstanceARNRegex             = "^(arn:aws:ecs:)([\\-\\w]+):[0-9]{12}:(container\\-instance)\\/[\\-\\w]+$"
)
