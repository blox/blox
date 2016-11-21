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
SWAGGER_SPEC := generated/v1/swagger.json
SWAGGER_GEN := $(shell find ./generated/v1 -name '*.go')
LOCAL_BINARY=bin/local/ecs-daemon-scheduler
GOINSTALL_BINARY=$(GOPATH)/bin/amazon-ecs-daemon-scheduler
ROOT := $(shell pwd)

.PHONY: all
all: clean build unit-tests

.PHONY: build
build: generate $(LOCAL_BINARY)

.PHONY: build-incr
build-incr: generate
	go install ./...
	cp $(GOINSTALL_BINARY) $(LOCAL_BINARY)

$(LOCAL_BINARY): $(SOURCES)
	./scripts/build_binary.sh ./out
	@echo "Built ECS Daemon Scheduler"

.PHONY: unit-tests
unit-tests:
	go test -v ./pkg/...

.PHONY: clean
clean:
	rm -rf ./out $(GOINSTALL_BINARY) ||:
	touch $(SWAGGER_SPEC)

$(SWAGGER_GEN): $(SWAGGER_SPEC)
	./scripts/generate_swagger_artifacts.sh

.PHONY: generate
generate: $(SWAGGER_GEN) $(SOURCES)
	PATH="$(ROOT)/scripts:${PATH}" go generate ./licenses/... ./copyright_gen/...

.PHONY: docker-build
docker-build:
	docker run -v $(shell pwd):/usr/src/app/src/github.com/aws/amazon-ecs-daemon-scheduler \
		--workdir=/usr/src/app/src/github.com/aws/amazon-ecs-daemon-scheduler \
		--env GOPATH=/usr/src/app \
		golang:1.7 make $(LOCAL_BINARY)

.PHONY: build-in-docker
build-in-docker:
	@docker build -f scripts/dockerfiles/Dockerfile.build -t "amazon/amazon-ecs-daemon-scheduler-build:make" .
	@docker run --net=none -v "$(shell pwd)/out:/out" -v "$(shell pwd):/go/src/github.com/aws/amazon-ecs-daemon-scheduler" "amazon/amazon-ecs-daemon-scheduler-build:make"

.PHONY: certs
certs: misc/certs/ca-certificates.crt
misc/certs/ca-certificates.crt:
	docker build -t "amazon/amazon-ecs-daemon-scheduler-cert-source:make" misc/certs/
	docker run "amazon/amazon-ecs-daemon-scheduler-cert-source:make" cat /etc/ssl/certs/ca-certificates.crt > misc/certs/ca-certificates.crt

.PHONY: docker
docker: certs build-in-docker
	@cd scripts && ./create-ecs-daemon-scheduler-scratch
	@docker build -f scripts/dockerfiles/Dockerfile.release -t "amazon/amazon-ecs-daemon-scheduler:make" .
	@echo "Built Docker image \"amazon/amazon-ecs-daemon-scheduler:make\""

.PHONY: docker-release
docker-release:
	@docker build -f scripts/dockerfiles/Dockerfile.cleanbuild -t "amazon/amazon-ecs-daemon-scheduler-cleanbuild:make" .
	@echo "Built Docker image \"amazon/amazon-ecs-daemon-scheduler-cleanbuild:make\""
	@docker run --net=none -v "$(shell pwd)/out:/out" -v "$(shell pwd):/src/amazon-ecs-daemon-scheduler" "amazon/amazon-ecs-daemon-scheduler-cleanbuild:make"

.PHONY: release
release: certs docker-release
	@cd scripts && ./create-ecs-daemon-scheduler-scratch
	@docker build -f scripts/dockerfiles/Dockerfile.release -t "amazon/amazon-ecs-daemon-scheduler:latest" .
	@echo "Built Docker image \"amazon/amazon-ecs-daemon-scheduler:latest\""

.PHONY: get-deps
get-deps:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get github.com/tools/godep
	go get github.com/gucumber/gucumber/cmd/gucumber
	go get github.com/golang/mock/gomock
	go get github.com/golang/mock/mockgen

.PHONY: e2e-tests
e2e-tests:
	gucumber -tags=@e2e

.PHONY: integ-tests
integ-tests:
	gucumber -tags=@integ
