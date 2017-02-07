#!/bin/bash
# Copyright 2016-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may not use this file except in compliance with the License. A copy of the License is located at
#
#     http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

# Usage instructions - From one level above the scripts directory, running the following command generates the binary into <path-to-binary-destination>
#      . ./scripts/v1/generate_swagger_models.sh

# Normalize to working directory being build root (up one level from ./scripts)
ROOT=$( cd "$( dirname "${BASH_SOURCE[0]}" )/../.." && pwd )

GENERATED_DIR="${ROOT}/swagger/v1/generated"
cd "${GENERATED_DIR}"

# Remove models directory if it already exists
REMOVE_MODELS="rm -rf ./models ||:"
${REMOVE_MODELS}

# Remove client directory if it already exists
REMOVE_CLIENT="rm -rf ./client ||:"
${REMOVE_CLIENT}

# Generate models and client
SWAGGER_GENERATE="swagger generate client -f ../swagger.json -A blox_daemon_scheduler"
${SWAGGER_GENERATE}

cd ${ROOT}
