package human

import (
	"fmt"
)

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
// Use width 6+suffix length
func Humanize(size int64) string {
	suffixes := [...]string{"b", "K", "M", "G", "T", "P", "E"}
	var unit int64 = KiB
	var i int
	for i = range suffixes {
		if size < unit {
			break
		}
		unit <<= 10
	}
	unit >>= 10
	if size < 1024 {
		return fmt.Sprintf("%d%s", size, suffixes[i])
	}
	return fmt.Sprintf("%.1f%s", float64(size)/float64(unit), suffixes[i])
	// Fixed point arithmetic truncates instead of rounding, causing small errors which may or may not be important
	// return fmt.Sprintf("%d.%d%s", size/unit, (10*size/unit)%10, suffixes[i])
}
