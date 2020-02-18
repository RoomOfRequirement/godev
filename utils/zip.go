package utils

import (
	"archive/zip"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// for complex usage, you may try this lib: https://github.com/mholt/archiver

// Zip ...
func Zip(srcPath, dstPath string) error {
	// create dst file
	fw, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer fw.Close()

	// zip.Writer
	zw := zip.NewWriter(fw)
	defer zw.Close()

	// recursively handle files/dirs
	return filepath.Walk(srcPath, func(path string, fi os.FileInfo, errBack error) error {
		if errBack != nil {
			return errBack
		}

		// zip file info
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		// From zip doc:
		// Because os.FileInfo's Name method returns only the base name of
		// the file it describes, it may be necessary to modify the Name field
		// of the returned header to provide the full path name of the file.
		fh.Name = filepath.ToSlash(filepath.Clean(path))

		// filepath.Clean will remove dir slash
		if fi.IsDir() {
			fh.Name += "/"
		}

		// writer
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}

		// if not regular file, write header info only
		// for example: dir
		if !fh.Mode().IsRegular() {
			return nil
		}

		// reader
		r, err := os.Open(path)
		defer r.Close()
		if err != nil {
			return err
		}

		sz, err := io.Copy(w, r)
		if err != nil {
			return err
		}

		log.Printf("successfully zip file: %s with %d bytes data\n", path, sz)

		return nil
	})
}

// UnZip ...
func UnZip(srcPath, dstPath string) error {
	dstPath = filepath.Clean(dstPath) + string(os.PathSeparator)
	zr, err := zip.OpenReader(srcPath)
	if err != nil {
		return err
	}
	defer zr.Close()

	// dst path dir
	if err := os.MkdirAll(dstPath, os.ModeDir|os.ModePerm); err != nil {
		return err
	}

	// traversal
	for _, f := range zr.File {
		path := filepath.Join(dstPath, f.Name)
		// Check for ZipSlip: https://snyk.io/research/zip-slip-vulnerability
		if !strings.HasPrefix(path, dstPath) {
			return errors.New("illegal file path: " + path + ", dstPath: " + dstPath)
		}
		if err := handle(f, path); err != nil {
			return err
		}
	}
	return nil
}

func handle(f *zip.File, path string) error {
	// reader
	r, err := f.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	// handle dir
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(path, f.Mode()); err != nil {
			return err
		}
		return nil
	}

	// handle file
	// file path
	if err := os.MkdirAll(filepath.Dir(path), f.Mode()); err != nil {
		return err
	}

	// writer
	// os.Create -> zip not store permissions
	w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer w.Close()

	sz, err := io.Copy(w, r)
	if err != nil {
		return err
	}

	log.Printf("successfully unzip file: %s with %d bytes data", f.Name, sz)

	return nil
}
