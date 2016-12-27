#!/bin/bash
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
#
# This script generates a file in go with the license contents as a constant

# Set language to C to make sorting consistent among different environments.
export LANG=C

set -e
outputfile=${1?Must provide an output file}
inputfile="$(<../../LICENSE)"

appendRepoLicense() {
  repo=$1
  echo "Adding license for ${repo}"
  inputfile+=$'\n'"***"$'\n'"$repo"$'\n\n'
  # Copy LICENSE* files
  for licensefile in $repo/LICENSE*; do
    if [ -f $licensefile ]; then
      inputfile+="$(<$licensefile)"$'\n'
    fi;
  done;
  # Copy COPYING* file
  if [ -f $repo/COPYING* ]; then
    inputfile+="$(<$repo/COPYING*)"$'\n'
  fi;
}

for registry in github.com golang.org; do
  for user in ./../../vendor/$registry/*; do
    for repo in $user/*; do
      if [[ $repo == *"go-buffruneio"* ]]; then
        # nop, since we add this explicitliy
	:
      else
        appendRepoLicense $repo
      fi
    done;
  done;
done;

for repo in ./../../vendor/gopkg.in/* ./../../vendor/google.golang.org/*; do
  appendRepoLicense $repo
done;

inputfile+="
***
./../../vendor/github.com/pelletier/go-buffruneio

Copyright (c) 2016 Thomas Pelletier

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the \"Software\"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

"

tr -d '\r' > "${outputfile}" << EOF 
// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package licenses

const License = \`$inputfile\`
EOF
