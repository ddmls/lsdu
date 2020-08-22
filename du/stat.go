package du

import "golang.org/x/sys/unix"

// func Lstat

// https://stackoverflow.com/questions/20108520/get-amount-of-free-disk-space-using-go

// FreeSpace reports the disk space available to unprivileged user in bytes and the total disk space
func FreeSpace(path string) (uint64, uint64, error) {
	var statfs unix.Statfs_t

	if err := unix.Statfs(path, &statfs); err != nil {
		return 0, 0, err
	}
	return statfs.Bavail * uint64(statfs.Bsize),
		statfs.Blocks * uint64(statfs.Bsize),
		nil
}

// man 3 fstatat for an example, man 2 stat

// BlockSize reports the apparent size and allocated size of a file in bytes
func BlockSize(path string) (int64, int64, error) {
	var stat unix.Stat_t

	if err := unix.Lstat(path, &stat); err != nil {
		return 0, 0, err
	}
	// blocks are always 512-byte units (man 2 stat)
	return stat.Size,
		stat.Blocks * 512,
		nil
}
