# Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

# This dockerfile generates a cert bundle that can be packaged with the 
# release container to invoke AWS APIs over TLS.

FROM debian:latest

RUN apt-get update &&  \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

RUN cp /etc/ca-certificates.conf /tmp/caconf && cat /tmp/caconf | \
    grep -v "mozilla/CNNIC_ROOT\.crt" > /etc/ca-certificates.conf && \
    update-ca-certificates --fresh
