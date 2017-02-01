#!/bin/bash
# Copyright 2016-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

# Usage:
# build_binary.sh  <relative-path-to-binary-destination-from-build-root>
set -e

# Normalize to working directory being build root (up one level from ./scripts)
ROOT=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )

cd "${ROOT}"
BINARY_DESTINATION_DIR=$1
mkdir -p ${BINARY_DESTINATION_DIR}

# Versioning: run the generator to setup the version and then always
# restore ourselves to a clean state
cp versioning/version.go versioning/_version.go
trap "mv versioning/_version.go versioning/version.go" EXIT SIGHUP SIGINT SIGTERM

cd versioning
go run gen/version-gen.go

cd "${ROOT}"

# Builds the handler binary from source in the specified destination paths
GOOS=$TARGET_GOOS CGO_ENABLED=0 go build -installsuffix cgo -a -ldflags '-s' -o ${BINARY_DESTINATION_DIR}/cluster-state-service ./
cp "${ROOT}/../LICENSE"  ${BINARY_DESTINATION_DIR}/
