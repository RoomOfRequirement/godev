package strutils

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// Compress compresses input string and return compressed bytes
func Compress(str string) ([]byte, error) {
	return CompressBytes(StringToBytes(str))
}

// Decompress decompresses input bytes and returns decompressed string
func Decompress(bytes []byte) (string, error) {
	ret, err := DecompressBytes(bytes)
	if err != nil {
		return "", err
	}
	return BytesToString(ret), nil
}

// CompressBytes compresses input bytes
// reference: https://golang.org/pkg/compress/gzip/
func CompressBytes(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	// Setting the Header fields is optional.
	// zw.Name = "a-new-hope.txt"
	// zw.Comment = "an epic space opera by George Lucas"
	// zw.ModTime = time.Date(1977, time.May, 25, 0, 0, 0, 0, time.UTC)

	_, err := zw.Write(data)
	if err != nil {
		return nil, err
	}

	err = zw.Flush()
	if err != nil {
		return nil, err
	}

	err = zw.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// DecompressBytes decompresses input bytes
func DecompressBytes(data []byte) ([]byte, error) {
	buf := bytes.NewReader(data)

	zr, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("Name: %s\nComment: %s\nModTime: %s\n\n", zr.Name, zr.Comment, zr.ModTime.UTC())

	ret, err := ioutil.ReadAll(zr)
	if err != nil {
		return nil, err
	}

	err = zr.Close()
	if err != nil {
		return nil, err
	}

	return ret, nil
}
