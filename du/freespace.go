// +build linux darwin

package du

import "syscall"

// https://stackoverflow.com/questions/20108520/get-amount-of-free-disk-space-using-go
// see /usr/lib/go/src/syscall/ztypes_linux_amd64.go for Statfs_t.Block, Statfs_t.Bsize
// and /usr/lib/go/src/syscall/types_linux.go for Fsid
// freebsd uses different types (int64) and the fields are not available on openbsd, netbsd

// FreeSpace reports the disk space available to unprivileged user in bytes, the total disk space and fsid
func FreeSpace(path string) (uint64, uint64, syscall.Fsid, error) {
	var statfs syscall.Statfs_t
	var fsid syscall.Fsid // Zero value

	if err := syscall.Statfs(path, &statfs); err != nil {
		return 0, 0, fsid, err
	}
	return statfs.Bavail * uint64(statfs.Bsize),
		statfs.Blocks * uint64(statfs.Bsize),
		statfs.Fsid,
		nil
}
