package compress

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"github.com/pkg/errors"
)

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
