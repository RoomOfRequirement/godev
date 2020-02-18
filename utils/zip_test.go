package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func TestZipUnZip(t *testing.T) {
	src := "test.txt"
	dst := "archive.zip"
	err := Zip(src, dst)
	assert.Error(t, err)
	f, err := os.Create(src)
	assert.NoError(t, err)
	assert.NotEmpty(t, f)
	err = Zip(src, dst)
	assert.NoError(t, err)

	defer func() {
		os.Remove(src)
		os.Remove(dst)
	}()

	// read-only
	_ = os.Chmod(dst, 0444)
	err = Zip(src, dst)
	assert.Error(t, err)

	src1 := "test"
	dst1 := "test.zip"
	os.Mkdir(src1, os.ModeDir|os.ModePerm)
	// recursive
	var filenames, dirNames []string
	for i := 0; i < 5; i++ {
		dirName := src1 + "/test" + strconv.Itoa(i)
		filename := dirName + "/file" + strconv.Itoa(i) + ".txt"
		filenames = append(filenames, filename)
		dirNames = append(dirNames, dirName)
		os.Mkdir(dirName, os.ModeDir|os.ModePerm)
		os.Create(filename)
	}
	err = Zip(src1, dst1)
	assert.NoError(t, err)

	defer func() {
		os.RemoveAll(src1)
		os.Remove(dst1)
	}()

	// unzip
	uz := "dst"
	err = UnZip(dst1, uz)
	assert.NoError(t, err)
	defer os.RemoveAll(uz)
}
