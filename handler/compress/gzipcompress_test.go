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
