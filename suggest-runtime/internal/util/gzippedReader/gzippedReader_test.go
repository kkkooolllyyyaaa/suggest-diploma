package gzippedReader

import (
	"compress/gzip"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testDataStruct struct {
	A int
	S string
}

func TestGzippedJsonReader(t *testing.T) {
	file, err := NewGzippedJsonReader("not_existed_file")
	assert.Nil(t, file)
	assert.Error(t, err)

	data := make([]testDataStruct, 0, 8)

	testData(t, data)

	data = append(data, testDataStruct{1, "a"})
	testData(t, data)

	data = append(data, testDataStruct{2, "b"})
	testData(t, data)

	data = append(data, testDataStruct{3, "c"})
	testData(t, data)
}

func testData(t *testing.T, data []testDataStruct) {
	const plainFileName = ".test_json.txt"
	const gzippedFileName = ".test_json.txt.gz"

	{
		jsn, err := json.Marshal(data)
		require.NoError(t, err)

		file, err := os.Create(plainFileName)
		require.NoError(t, err)
		n, err := file.Write(jsn)
		require.NoError(t, err)
		require.Equal(t, n, len(jsn))
		err = file.Close()
		require.NoError(t, err)

		file, err = os.Create(gzippedFileName)
		require.NoError(t, err)
		gz := gzip.NewWriter(file)
		n, err = gz.Write(jsn)
		require.NoError(t, err)
		require.Equal(t, n, len(jsn))
		err = gz.Close()
		require.NoError(t, err)
		err = file.Close()
		require.NoError(t, err)
	}

	file, err := NewGzippedJsonReader(plainFileName)
	require.NoError(t, err)
	var dataPlain []testDataStruct
	err = file.DecodeJson(&dataPlain)
	assert.NoError(t, err)
	assert.Equal(t, data, dataPlain)
	file.Close()

	file, err = NewGzippedJsonReader(gzippedFileName)
	require.NoError(t, err)
	var dataGzipped []testDataStruct
	err = file.DecodeJson(&dataGzipped)
	assert.NoError(t, err)
	assert.Equal(t, data, dataGzipped)
	file.Close()

	_ = os.Remove(plainFileName)
	_ = os.Remove(gzippedFileName)
}
