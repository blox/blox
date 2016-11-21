#!/bin/bash
# Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

# Remove models directory if it already exists
RELATE_PATH_TO_GENERATED_MODELS="/handler/api/v1/models"
REMOVE_MODELS="rm -rf ${ROOT}${RELATE_PATH_TO_GENERATED_MODELS} ||:"
${REMOVE_MODELS}

# Generate models and client
RELATIVE_PATH_TO_SWAGGER_DIR="/handler/api/v1"
SWAGGER_GENERATE="swagger generate client -f swagger/swagger.json -A amazon_ecs_esh"
cd "${ROOT}/${RELATIVE_PATH_TO_SWAGGER_DIR}"
${SWAGGER_GENERATE}

# Remove the client (we only need the models for the server)
RELATE_PATH_TO_GENERATED_CLIENT="/handler/api/v1/client"
REMOVE_CLIENT="rm -rf ${ROOT}${RELATE_PATH_TO_GENERATED_CLIENT}"
${REMOVE_CLIENT}

cd ${ROOT}
