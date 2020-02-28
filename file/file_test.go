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

func TestCopy(t *testing.T) {
	dir, file := "./tmp/", "./tmp/tmp.txt"
	err := Create(file)
	assert.NoError(t, err)
	subFile := dir + "sub/sub.txt"
	err = Create(subFile)
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	copyDir := "./copy/"
	err = CreateDir(copyDir)
	assert.NoError(t, err)
	defer os.RemoveAll(copyDir)
	err = CopyFile(file, copyDir+"dup.txt")
	assert.NoError(t, err)
	err = CopyDir(dir, copyDir)
	assert.NoError(t, err)

	// err file
	err = CopyFile("", copyDir)
	assert.Error(t, err)
	err = CopyFile(file, copyDir)
	assert.Error(t, err)
	err = CopyFile(dir, copyDir)
	assert.Error(t, err)
	// read-only file
	rofp := dir + "ro.txt"
	f, err := os.Create(rofp)
	assert.NoError(t, err)
	err = f.Chmod(0444)
	assert.NoError(t, err)
	err = CopyFile(rofp, copyDir+"dupro.txt")
	assert.NoError(t, err)

	// err dir
	err = CopyDir("", copyDir)
	assert.Error(t, err)
	err = CopyDir(file, copyDir)
	assert.Error(t, err)
	// read-only file
	err = CopyDir(dir, copyDir)
	assert.NoError(t, err)
	// read-only dir
	rodir := "ro"
	err = os.Mkdir(rodir, 0444)
	assert.NoError(t, err)
	defer os.RemoveAll(rodir)
	err = CopyDir(dir, rodir)
	assert.Error(t, err)
}

func TestReadWrite(t *testing.T) {
	dir, file := "./tmp/", "./tmp/tmp.txt"
	err := Create(file)
	assert.NoError(t, err)
	defer os.RemoveAll(dir)
	isWritable, err := IsWritable(file)
	assert.True(t, isWritable)
	assert.NoError(t, err)
	err = SetReadOnly(file)
	assert.NoError(t, err)
	isWritable, err = IsWritable(file)
	assert.False(t, isWritable)
	assert.Error(t, err) // permission denied
	err = SetWritable(file)
	isWritable, err = IsWritable(file)
	assert.True(t, isWritable)
	assert.NoError(t, err)
}

func TestDirSize(t *testing.T) {
	ds, err := SizeDir(".")
	assert.NoError(t, err)
	assert.True(t, ds > 0)

	ds, err = SizeDir("./not-exist")
	assert.Error(t, err)
	assert.True(t, ds == 0)
}

func TestFileSize(t *testing.T) {
	fs, err := SizeFile("./file.go")
	assert.NoError(t, err)
	assert.True(t, fs > 0)

	fs, err = SizeFile("./not-exist")
	assert.Error(t, err)
	assert.True(t, fs == 0)
}

func TestSize(t *testing.T) {
	sz, err := Size("./file.go")
	assert.NoError(t, err)
	assert.True(t, sz > 0)

	sz, err = Size(".")
	assert.NoError(t, err)
	assert.True(t, sz > 0)

	sz, err = Size("./not-exist")
	assert.Error(t, err)
	assert.True(t, sz == 0)
}
