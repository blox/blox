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

.PHONY: all
all: clean build unit-test

.PHONY: build
build:	$(LOCAL_BINARY)

$(LOCAL_BINARY): $(SOURCES)
	. ./scripts/build_binary.sh ./bin/local
	@echo "Built event-stream handler"

.PHONY: unit-test
unit-test:
	go test -v -timeout 1s ./handler/... -short

.PHONY: clean
clean:
	rm -rf ./bin ||:
