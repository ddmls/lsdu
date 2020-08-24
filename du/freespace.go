// +build linux darwin

package du

import "syscall"

// https://stackoverflow.com/questions/20108520/get-amount-of-free-disk-space-using-go
// see /usr/lib/go/src/syscall/ztypes_linux_amd64.go
// freebsd uses different types and the fields are not available on openbsd, netbsd

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
