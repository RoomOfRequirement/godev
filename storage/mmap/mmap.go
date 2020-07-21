package mmap

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"syscall"
)

// PROT_EXEC                         = 0x4
// PROT_NONE                         = 0x0
// PROT_READ                         = 0x1
// PROT_WRITE                        = 0x2
//
// MAP_ANON                          = 0x1000
// MAP_COPY                          = 0x2
// MAP_PRIVATE                       = 0x2
// MAP_SHARED                        = 0x1

const (
	// RDONLY maps the memory read-only.
	// Notice: If write to this MMap object, it will result in undefined behavior.
	RDONLY = 0x0
	// RDWR maps the memory as read-write. Underlying file sync changes with memory.
	RDWR = 0x1
	// COPY maps the memory as copy-on-write.
	// Notice: Write to this MMap object will only change data in memory, the underlying file will remain unchanged.
	COPY = 0x2
	// EXEC maps the memory as executable.
	EXEC = 0x4

	// ANON flag sets the mapped memory not backed by a file.
	// Notice: this will ignore input fd in syscall
	// https://stackoverflow.com/questions/34042915/what-is-the-purpose-of-map-anonymous-flag-in-mmap-system-call
	ANON = 0x1
)

// MMap ...
// reference: https://man7.org/linux/man-pages/man2/mmap.2.html
type MMap struct {
	data []byte
}

// New ...
func New(filename string, size int, prot, flags int, offset int64) (*MMap, error) {
	if size == 0 {
		return nil, fmt.Errorf("mmap: can NOT map region size == 0")
	}
	if offset%int64(os.Getpagesize()) != 0 {
		return nil, fmt.Errorf("mmap: offset must be a multiple of os pagesize")
	}
	sysProt := syscall.PROT_READ
	sysFlags := syscall.MAP_SHARED
	switch {
	// COPY
	case prot&COPY != 0:
		sysProt |= syscall.PROT_WRITE
		sysFlags = syscall.MAP_PRIVATE
	// RDWR
	case prot&RDWR != 0:
		sysProt |= syscall.PROT_WRITE
	// EXEC
	case prot&EXEC != 0:
		sysProt |= syscall.PROT_EXEC
	}
	var data []byte
	var err error

	// anonymous mapping, no file required, so ignore filename
	if flags&ANON != 0 {
		sysFlags |= syscall.MAP_ANON
		data, err = syscall.Mmap(-1, offset, size, sysProt, sysFlags)
	} else {
		// open or create
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		fi, err := f.Stat()
		if err != nil {
			return nil, err
		}
		fSize := fi.Size()
		// golang int is 64-bit on x64, 32-bit on x86 (like size_t)
		// so on 32-bit system, file > 2GB is not able to be mmap-ed
		if fSize != int64(int(fSize)) {
			return nil, fmt.Errorf("mmap: file %q size too large", filename)
		}
		if fSize < int64(size) {
			err = f.Truncate(int64(size))
			if err != nil {
				return nil, err
			}
		} else {
			data, err = syscall.Mmap(int(f.Fd()), offset, size, sysProt, sysFlags)
		}
	}

	if err != nil {
		return nil, err
	}
	// register finalizer to avoid fd leak when forget to close mmap
	m := &MMap{
		data: data,
	}
	runtime.SetFinalizer(m, (*MMap).Close)
	return m, nil
}

// Close ...
// io.Closer interface: https://golang.org/pkg/io/#Closer
func (m *MMap) Close() error {
	if m.data == nil {
		return nil
	}
	data := m.data
	m.data = nil
	// no need for a finalizer anymore
	runtime.SetFinalizer(m, nil)
	return syscall.Munmap(data)
}

// Len ...
func (m *MMap) Len() int {
	return len(m.data)
}

// At ...
func (m *MMap) At(i int) byte {
	return m.data[i]
}

// ReadAt ...
// io.ReadAt interface: https://golang.org/pkg/io/#ReaderAt
func (m *MMap) ReadAt(p []byte, off int64) (n int, err error) {
	if m.data == nil {
		return 0, fmt.Errorf("mmap: closed")
	}
	if off < 0 || off > int64(len(m.data)) {
		return 0, fmt.Errorf("mmap: invalid ReadAt offset: %d", off)
	}
	n = copy(p, m.data[off:])
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

// WriteAt ...
// io.WriteAt interface: https://golang.org/pkg/io/#WriteAt
// notice: only available when mmap prot >= RDWR, otherwise it will panic
func (m *MMap) WriteAt(p []byte, off int64) (n int, err error) {
	if m.data == nil {
		return 0, fmt.Errorf("mmap: closed")
	}
	if off < 0 || off > int64(len(m.data)) {
		return 0, fmt.Errorf("mmap: invalid WriteAt offset: %d", off)
	}
	n = copy(m.data[off:], p)
	if n < len(p) {
		return n, io.ErrShortWrite
	}
	return n, nil
}
