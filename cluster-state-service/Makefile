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
LOCAL_BINARY=out/cluster-state-service
ROOT := $(shell pwd)

.PHONY: all
all: clean generate build unit-tests

.PHONY: generate-models
generate-models:
	./scripts/v1/generate_swagger_models.sh
	@echo "Generated swagger models"

.PHONY: build
build:	$(LOCAL_BINARY)

$(LOCAL_BINARY): $(SOURCES)
	./scripts/build_binary.sh ./out
	@echo "Built Cluster State Service"

.PHONY: generate
generate: generate-models $(SOURCES)
	PATH="$(ROOT)/scripts:${PATH}" go generate ./licenses/... ./copyright_gen/...
	@echo "Generated licenses and copyrights"

.PHONY: generate build-in-docker
build-in-docker:
	@docker build -f scripts/dockerfiles/Dockerfile.build -t "bloxoss/cluster-state-service-build:make" .
	@docker run --net=none \
	    -v "$(shell pwd)/out:/out" \
	    -v "$(shell pwd):/go/src/github.com/blox/blox/cluster-state-service" \
	    -v "$(shell dirname ${ROOT}):/go/src/github.com/blox/blox" \
	    "bloxoss/cluster-state-service-build:make"

.PHONY: certs
certs: misc/certs/ca-certificates.crt
misc/certs/ca-certificates.crt:
	docker build -t "bloxoss/cluster-state-service-cert-source:make" misc/certs/
	docker run "bloxoss/cluster-state-service-cert-source:make" cat /etc/ssl/certs/ca-certificates.crt > misc/certs/ca-certificates.crt

.PHONY: docker
docker: certs build-in-docker
	@cd scripts && ./create-cluster-state-service-scratch
	@docker build -f scripts/dockerfiles/Dockerfile.release -t "bloxoss/cluster-state-service:make" .
	@echo "Built Docker image \"bloxoss/cluster-state-service:make\""

.PHONY: docker-release
docker-release:
	@docker build -f scripts/dockerfiles/Dockerfile.cleanbuild -t "bloxoss/cluster-state-service-cleanbuild:make" .
	@echo "Built Docker image \"bloxoss/cluster-state-service-cleanbuild:make\""
	@docker run --net=none \
	    -v "$(shell pwd)/out:/out" \
	    -v "$(shell dirname ${ROOT}):/src/blox" \
	    "bloxoss/cluster-state-service-cleanbuild:make"

.PHONY: release
release: certs docker-release
	@cd scripts && ./create-cluster-state-service-scratch
	@docker build -f scripts/dockerfiles/Dockerfile.release -t "bloxoss/cluster-state-service:latest" .
	@echo "Built Docker image \"bloxoss/cluster-state-service:latest\""

.PHONY: get-deps
get-deps:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get github.com/tools/godep
	go get github.com/gucumber/gucumber/cmd/gucumber
	go get github.com/golang/mock/gomock
	go get github.com/golang/mock/mockgen

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
