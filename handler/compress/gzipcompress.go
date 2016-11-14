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

package compress

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Compress gzip compresses a string
func Compress(s string) ([]byte, error) {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)

	_, err := gzWriter.Write([]byte(s))
	if err != nil {
		return nil, errors.Wrapf(err, "Could not compress string: %v", s)
	}

	err = gzWriter.Close()
	if err != nil {
		return nil, errors.Wrapf(err, "Unexpected error while closing gzip writer")
	}

	return buf.Bytes(), nil
}

// Uncompress uncompresses a gzipped byte representation of a string
func Uncompress(b []byte) (string, error) {
	gzReader, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return "", errors.Wrapf(err, "Unexpected error while creating gzip reader")
	}
	defer gzReader.Close()

	result, err := ioutil.ReadAll(gzReader)
	if err != nil {
		return "", errors.Wrapf(err, "Could not read compressed bytes: %v", b)
	}

	return string(result), nil
}
