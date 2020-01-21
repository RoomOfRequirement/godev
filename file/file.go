// Package file provides some utils based on "os" and "path/filepath" packages
package file

import (
	"fmt"
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
