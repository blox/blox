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
#
# This script generates a file in go with the license contents as a constant

set -e

if [ -z "$1" ]
  then
    echo "Must provide at least one input directory"
    exit 1
fi

appendLicense(){
  echo "Slapping copyright notice on ${outputfile}"
  outputfile=${1?Must provide an output file}

  echo "// Copyright 2016-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the \"License\"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the \"license\" file accompanying this file. This file is distributed
// on an \"AS IS\" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

$(cat ${outputfile})" > ${outputfile}
}

for input in "$@"
do
  if [ -d ${input} ] ; then
    for file in `find ${input} -name "*.go"`; do
      if ! grep -q "// Copyright" ${file} ; then
        appendLicense ${file}
      fi
    done
  else
    if ! grep -q "// Copyright" ${input} ; then
      appendLicense ${input}
    fi
  fi
done
