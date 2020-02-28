// Package file provides some utils based on "os" and "path/filepath" packages
package file

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

// IsExist returns true if path exists
func IsExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true // may have other error like permission etc.
}

// IsDir returns true if path is dir
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// IsFile returns true if path is file
func IsFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

// Create creates file according to input path (nested dir will be created if required)
func Create(path string) error {
	if IsExist(path) {
		return fmt.Errorf("%s existed", path)
	}
	_ = CreateDir(filepath.Dir(path))
	file, err := os.Create(path)
	defer func() {
		_ = file.Close()
	}()
	return err
}

// CreateDir creates dir according to input dirPath (nested dir will be created if required)
func CreateDir(dirPath string) error {
	info, err := os.Stat(dirPath)
	// not exist
	if os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	// exist as a file
	if info.Mode().IsRegular() {
		return fmt.Errorf("%s already exist as a file", dirPath)
	}
	// exist
	return fmt.Errorf("%s already exists", dirPath)
}

// Remove removes file if input path is file or removes dir with all files if input path is dir
func Remove(path string) error {
	if IsDir(path) {
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}

// Walk walks through all sub-dirs and files of input root and call walkFunc on those paths
//	skipPattern is used to indicate which path is skipped
//	please notice, if skipPatten == "" (empty string), all paths will be skipped
func Walk(root string, walkFunc func(path string) error, skipPattern string) error {
	// TODO: concurrent walk like "x/tools/internal/fastwalk"?
	re := regexp.MustCompile(skipPattern)
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("fail to access path %q: %v", path, err)
		}
		// "" always return true
		if re.MatchString(path) {
			return nil
		}
		return walkFunc(path)
	})
}

// CopyFile copies src file to dst and keeps the file mode
func CopyFile(srcPath, dstPath string) error {
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	if srcInfo.IsDir() {
		return errors.New("src is a dir, not a file")
	}
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return os.Chmod(dstPath, srcInfo.Mode())
}

// CopyDir copies src dir to dst and keeps file mode
func CopyDir(srcPath, dstPath string) error {
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	// not dir
	if !srcInfo.IsDir() {
		return errors.New("src is not a dir")
	}
	// create dst dir
	err = os.MkdirAll(dstPath, srcInfo.Mode())
	if err != nil {
		return err
	}

	srcDir, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcDir.Close()

	// get all FileInfo (-1)
	fis, err := srcDir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		fp := filepath.Join(srcPath, fi.Name())
		dfp := filepath.Join(dstPath, fi.Name())
		// sub-dirs, handled recursively
		if fi.IsDir() {
			err = CopyDir(fp, dfp)
		} else {
			err = CopyFile(fp, dfp)
		}
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

// IsWritable ...
func IsWritable(filepath string) (bool, error) {
	file, err := os.OpenFile(filepath, os.O_WRONLY, 0666)
	if err != nil {
		if os.IsPermission(err) {
			return false, err
		}
	}
	defer file.Close()
	return true, nil

}

// SetWritable ...
func SetWritable(filepath string) error {
	return os.Chmod(filepath, 0222)
}

// SetReadOnly ...
func SetReadOnly(filepath string) error {
	return os.Chmod(filepath, 0444)
}

// SizeDir in bytes
func SizeDir(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("fail to access path %q: %v", p, err)
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// SizeFile in bytes
func SizeFile(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("fail to access path %q: %v", path, err)
	}
	return fi.Size(), nil
}

// Size in bytes
func Size(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("fail to access path %q: %v", path, err)
	}
	if fi.IsDir() {
		return SizeDir(path)
	}
	return fi.Size(), nil
}
