# Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may not use this file except in compliance with the License. A copy of the License is located at
#
#     http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
LOCAL_BINARY=bin/local/ecs-event-stream-handler
ROOT := $(shell pwd)

.PHONY: all
all: clean build unit-tests

.PHONY: generate-models
generate-models:
	. ./scripts/v1/generate_swagger_models.sh
	@echo "Generated swagger models"

.PHONY: build
build:	$(LOCAL_BINARY)

$(LOCAL_BINARY): $(SOURCES)
	. ./scripts/build_binary.sh ./bin/local
	@echo "Built event-stream handler"

.PHONY: generate
generate: generate-models $(SOURCES)
	PATH="$(ROOT)/scripts:${PATH}" go generate ./licenses/... ./copyright_gen/...

.PHONY: get-deps
get-deps:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get github.com/tools/godep
	go get github.com/gucumber/gucumber/cmd/gucumber

.PHONY: unit-tests
unit-tests:
	go test -v -timeout 1s ./handler/... -short

# Start the server before running this target. More details in Readme under ./internal directory.
.PHONY: e2e-test
e2e-tests:
	gucumber -tags=@e2e

.PHONY: clean
clean:
	rm -rf ./bin ||:
