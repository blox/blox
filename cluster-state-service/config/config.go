// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the License). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the license file accompanying this file. This file is distributed
// on an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package config

// EtcdEndpoints represents the etcd servers to connect to.
var EtcdEndpoints []string

// QueueName represents the queue name to listen to for ECS events. Formatted as
// a URI with the scheme determining the type.  For example sqs://name or kinesis://name
var QueueNameURI string

// CSSBindAddr represents the address CSS listens on.
var CSSBindAddr string

// PrintVersion represents the flag to set when printing version information.
var PrintVersion bool
