package file

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIs(t *testing.T) {
	dir, file := "./tmp/", "./tmp/tmp.txt"
	assert.False(t, IsExist(file))
	assert.False(t, IsDir(dir))
	assert.False(t, IsFile(file))
	err := os.MkdirAll(dir, os.ModePerm)
	assert.NoError(t, err)
	f, err := os.Create(file)
	assert.NoError(t, err)
	defer func() {
		f.Close()
		os.RemoveAll(dir)
	}()

	assert.True(t, IsExist(file))
	assert.True(t, IsDir(dir))
	assert.True(t, IsFile(file))
	assert.False(t, IsDir(file))
	assert.False(t, IsFile(dir))

	assert.NoError(t, Remove(file))
	assert.False(t, IsFile(file))
	assert.True(t, IsDir(dir))
	assert.NoError(t, Remove(dir))
	assert.False(t, IsDir(dir))
}

func TestCreate(t *testing.T) {
	dir, file := "./tmp/", "./tmp/tmp.txt"
	err := Create(file)
	assert.NoError(t, err)
	assert.True(t, IsExist(file))
	assert.True(t, IsDir(dir))
	assert.True(t, IsFile(file))
	err = CreateDir(dir)
	assert.Error(t, err, fmt.Errorf("%s already exists", dir))
	err = CreateDir(file)
	assert.Error(t, err, fmt.Errorf("%s already exist as a file", file))

	assert.NoError(t, Remove(dir))

	err = CreateDir(dir)
	assert.NoError(t, err)
	assert.True(t, IsDir(dir))
	err = Create(dir)
	assert.Error(t, err, fmt.Errorf("%s existed", dir))

	assert.NoError(t, Remove(dir))
}

func TestWalk(t *testing.T) {
	for i := 0; i < 5; i++ {
		_ = CreateDir(fmt.Sprintf("./tmp/%d/", i))
		for j := 0; j < 3; j++ {
			_ = Create(fmt.Sprintf("./tmp/%d/%d.txt", i, j))
		}
	}
	cnt := 0
	err := Walk("./tmp", func(path string) error {
		if IsExist(path) {
			if IsFile(path) {
				cnt++
			}
			return nil
		}
		return fmt.Errorf("%s not exist", path)
	}, "_")
	assert.NoError(t, err)
	assert.Equal(t, 15, cnt)

	err = Walk("./tmp/tmp", func(path string) error {
		return nil
	}, "_")
	assert.Error(t, err)

	err = Walk("./tmp", func(path string) error {
		return nil
	}, "")
	assert.NoError(t, err)

	err = Remove("./tmp")
	assert.NoError(t, err)
}
