// +build !windows
// +build !plan9

package du

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// If Go 2 drops syscall we have to use golang.org/x/sys/unix

// ReportApparentSize affects FileInfo.Size()
var ReportApparentSize bool

// https://stackoverflow.com/questions/20108520/get-amount-of-free-disk-space-using-go
// see /usr/lib/go/src/syscall/ztypes_linux_amd64.go
// compiles for linux, darwin but freebsd differs in its types

// FreeSpace reports the disk space available to unprivileged user in bytes and the total disk space
func FreeSpace(path string) (uint64, uint64, error) {
	var statfs syscall.Statfs_t

	if err := syscall.Statfs(path, &statfs); err != nil {
		return 0, 0, err
	}
	return statfs.Bavail * uint64(statfs.Bsize),
		statfs.Blocks * uint64(statfs.Bsize),
		nil
}

// FileInfo contains the information of os.FileInfo and the allocated size of a file which may be different if it contains holes.
type FileInfo interface {
	os.FileInfo
	AllocatedSize() int64
	ApparentSize() int64
}

type fileInfo struct {
	osfi          os.FileInfo
	allocatedSize int64
}

// We need to override Size() so lots of boilerplate

func (fi fileInfo) AllocatedSize() int64 {
	return fi.allocatedSize
}

func (fi fileInfo) ApparentSize() int64 {
	return fi.osfi.Size()
}

func (fi fileInfo) Name() string {
	return fi.osfi.Name()
}

func (fi fileInfo) Size() int64 {
	if ReportApparentSize {
		return fi.ApparentSize()
	}
	return fi.AllocatedSize()
}

func (fi fileInfo) Mode() os.FileMode {
	return fi.osfi.Mode()
}

func (fi fileInfo) ModTime() time.Time {
	return fi.osfi.ModTime()
}

func (fi fileInfo) IsDir() bool {
	return fi.osfi.IsDir()
}

func (fi fileInfo) Sys() interface{} {
	return fi.osfi.Sys()
}

// We assume an implementation of /usr/lib/go/src/os/types_unix.go as of go 1.15
// man 3 fstatat for an example
// man 2 stat

func xstat(path string, statFunc func(string) (os.FileInfo, error)) (FileInfo, error) {
	fi, err := statFunc(path)
	if err != nil {
		return nil, err
	}
	sys, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("Type assertion of *syscall.Stat_t failed")
	}
	// blocks are always 512-byte units (man 2 stat)
	return fileInfo{fi, sys.Blocks * 512}, nil
}

// Stat returns the info of os.Stat and the real allocated size of the file
func Stat(path string) (FileInfo, error) {
	return xstat(path, os.Stat)
}

// Lstat returns the info of os.Lstat and the real allocated size of the file
func Lstat(path string) (FileInfo, error) {
	return xstat(path, os.Lstat)
}
