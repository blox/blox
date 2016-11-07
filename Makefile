# Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You
# may not use this file except in compliance with the License. A copy of
# the License is located at
#
#     http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is
# distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF
# ANY KIND, either express or implied. See the License for the specific
# language governing permissions and limitations under the License.

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
	./scripts/build_binary.sh ./out
	@echo "Built event-stream handler"

.PHONY: generate
generate: generate-models $(SOURCES)
	PATH="$(ROOT)/scripts:${PATH}" go generate ./licenses/... ./copyright_gen/...

.PHONY: docker-build
docker-build:
	docker run -v $(shell pwd):/usr/src/app/src/github.com/aws/amazon-ecs-event-stream-handler \
		--workdir=/usr/src/app/src/github.com/aws/amazon-ecs-event-stream-handler \
		--env GOPATH=/usr/src/app \
		golang:1.7 make $(LOCAL_BINARY)

.PHONY: build-in-docker
build-in-docker:
	@docker build -f scripts/dockerfiles/Dockerfile.build -t "amazon/amazon-ecs-event-stream-handler-build:make" .
	@docker run --net=none -v "$(shell pwd)/out:/out" -v "$(shell pwd):/go/src/github.com/aws/amazon-ecs-event-stream-handler" "amazon/amazon-ecs-event-stream-handler-build:make"

.PHONY: certs
certs: misc/certs/ca-certificates.crt
misc/certs/ca-certificates.crt:
	docker build -t "amazon/amazon-ecs-event-stream-handler-cert-source:make" misc/certs/
	docker run "amazon/amazon-ecs-event-stream-handler-cert-source:make" cat /etc/ssl/certs/ca-certificates.crt > misc/certs/ca-certificates.crt

.PHONY: docker
docker: certs build-in-docker
	@cd scripts && ./create-ecs-event-stream-handler-scratch
	@docker build -f scripts/dockerfiles/Dockerfile.release -t "amazon/amazon-ecs-event-stream-handler:make" .
	@echo "Built Docker image \"amazon/amazon-ecs-event-stream-handler:make\""

.PHONY: docker-release
docker-release:
	@docker build -f scripts/dockerfiles/Dockerfile.cleanbuild -t "amazon/amazon-ecs-event-stream-handler-cleanbuild:make" .
	@echo "Built Docker image \"amazon/amazon-ecs-event-stream-handler-cleanbuild:make\""
	@docker run --net=none -v "$(shell pwd)/out:/out" -v "$(shell pwd):/src/amazon-ecs-event-stream-handler" "amazon/amazon-ecs-event-stream-handler-cleanbuild:make"

.PHONY: release
release: certs docker-release
	@cd scripts && ./create-ecs-event-stream-handler-scratch
	@docker build -f scripts/dockerfiles/Dockerfile.release -t "amazon/amazon-ecs-event-stream-handler:latest" .
	@echo "Built Docker image \"amazon/amazon-ecs-event-stream-handler:latest\""

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
	rm -rf ./out ||:
