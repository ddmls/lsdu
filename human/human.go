package human

import "fmt"

const minInt64 = -9223372036854775808
const maxInt64 = 9223372036854775807

// Constants for sizes
const (
	KiB int64 = 1 << (10 * (iota + 1))
	MiB
	GiB
	TiB
	PiB
)

// Humanize reports the size in power of 1024 units (KiB, MiB etc)
func Humanize(size int64) string {
	suffixes := [...]string{"b", "K", "M", "G", "T", "P"}
	for base := 0; base < len(suffixes); base++ {
		if size < 1024 {
			return fmt.Sprintf("%d%s", size, suffixes[base])
		}
		size >>= 10
	}
	// Too big
	return fmt.Sprintf("%d%s", size, suffixes[len(suffixes)-1])
}
