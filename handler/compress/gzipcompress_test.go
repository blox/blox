// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the License). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the license file accompanying this file. This file is distributed
// on an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package compress

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

const (
	uncompressedStr = "Hello world"
)

var (
	compressedBytes = []byte{0x1f, 0x8b, 0x8, 0x0, 0x0, 0x9, 0x6e, 0x88, 0x0,
		0xff, 0xf2, 0x48, 0xcd, 0xc9, 0xc9, 0x57, 0x28, 0xcf, 0x2f, 0xca, 0x49, 0x1,
		0x4, 0x0, 0x0, 0xff, 0xff, 0x52, 0x9e, 0xd6, 0x8b, 0xb, 0x0, 0x0, 0x0}
)

func TestCompress(t *testing.T) {
	b, err := Compress(uncompressedStr)
	assert.Nil(t, err)
	assert.Equal(t, compressedBytes, b)
}

func TestUncompress(t *testing.T) {
	s, err := Uncompress(compressedBytes)
	assert.Nil(t, err)
	assert.Equal(t, uncompressedStr, s)
}

func TestCompressAndUncompress(t *testing.T) {
	b, err := Compress(uncompressedStr)
	assert.Nil(t, err)
	s, err := Uncompress(b)
	assert.Nil(t, err)
	assert.Equal(t, uncompressedStr, s)
}
