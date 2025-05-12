package gzippedReader

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"strings"
)

type GzippedJsonReader struct {
	file       *os.File
	gzipReader *gzip.Reader
	reader     io.Reader

	filePath string
}

func NewGzippedJsonReader(filePath string) (*GzippedJsonReader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	var reader io.Reader = file
	var unzip *gzip.Reader
	if strings.HasSuffix(filePath, ".gz") {
		unzip, err = gzip.NewReader(file)
		if err != nil {
			_ = file.Close()
			return nil, err
		}
		unzip.Multistream(false)
		reader = unzip
	}

	return &GzippedJsonReader{
		file:       file,
		gzipReader: unzip,
		reader:     reader,
		filePath:   filePath,
	}, nil
}

func (g *GzippedJsonReader) Read(p []byte) (n int, err error) {
	return g.reader.Read(p)
}

func (g *GzippedJsonReader) Close() {
	if g.gzipReader != nil {
		_ = g.gzipReader.Close()
	}
	if g.file != nil {
		_ = g.file.Close()
	}
}

func (g *GzippedJsonReader) DecodeJson(v interface{}) error {
	jsonStream := bufio.NewReader(g)
	dec := json.NewDecoder(jsonStream)

	err := dec.Decode(v)
	if err != nil {
		return err
	}

	return nil
}
