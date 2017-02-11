# Copyright 2016-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
LOCAL_BINARY=out/css-canary
ROOT := $(shell pwd)
PARENT := $(shell dirname ${ROOT})

.PHONY: all
all: clean build

.PHONY: clean
clean:
	rm -rf ./out ||:

.PHONY: build
build:	$(LOCAL_BINARY)

$(LOCAL_BINARY): $(SOURCES)
	./scripts/build_binary.sh ./out
	@echo "Built Cluster State Service Canary"

.PHONY: build-in-docker
build-in-docker:
	@docker build -f scripts/dockerfiles/Dockerfile.build -t "bloxoss/cluster-state-service-canary-build:make" .
	@echo "Built Docker image \"bloxoss/cluster-state-service-canary-build:make\""
	@docker run --net=none \
			-v "$(shell pwd)/out:/out" \
			-v "$(shell dirname ${PARENT}):/src/blox" \
			"bloxoss/cluster-state-service-canary-build:make"

.PHONY: docker-certs
docker-certs: misc/certs/ca-certificates.crt
misc/certs/ca-certificates.crt:
	docker build -f scripts/dockerfiles/Dockerfile.certs -t "bloxoss/cluster-state-service-canary-cert-source:make" .
	docker run "bloxoss/cluster-state-service-canary-cert-source:make" cat /etc/ssl/certs/ca-certificates.crt > out/docker-ca-certificates.crt

.PHONY: docker-release
docker-release: docker-certs build-in-docker
	@docker build -f scripts/dockerfiles/Dockerfile.release -t "bloxoss/cluster-state-service-canary:latest" .
	@echo "Built Docker image \"bloxoss/cluster-state-service-canary:latest\""
