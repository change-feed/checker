package util

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// Compress takes a raw []byte and returns a compressed []byte
func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write(data)
	if err != nil {
		return nil, fmt.Errorf("compression failed: %v", err)
	}
	gz.Close()

	return buf.Bytes(), nil
}

// Decompress takes a compressed []byte and returns a decompressed []byte
func Decompress(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("gzip decompression failed: %v", err)
	}
	decoded, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("reading decompressed data failed: %v", err)
	}

	return decoded, nil
}
