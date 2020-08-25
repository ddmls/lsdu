package human

import (
	"fmt"
	"math"
)

// Constants for sizes
const (
	KiB int64 = 1 << (10 * (iota + 1))
	MiB
	GiB
	TiB
	PiB
)

var suffixes = [...]string{"b", "K", "M", "G", "T", "P", "E"}

// SizesHuman, SizesInK, SizesInK choose the formatting of sizes used by Format
const (
	SizesBytes = iota
	SizesInK
	SizesInM
	SizesHuman
)

// Humanize reports the size in power of 1024 units (KiB, MiB etc)
func Humanize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d%s", size, suffixes[0])
	}
	b := int(math.Log2(float64(size))) / 10
	base := 1 << (10 * b)
	return fmt.Sprintf("%.1f%s", float64(size)/float64(base), suffixes[b])
}

// Base converts size to the specified base e.g. Base(size, KiB), optionally with a suffix
// If using a suffix, the printed result is a rounded float othewise a truncated integer
func Base(size, base int64, suffix bool) string {
	b := int(math.Log2(float64(base))) / 10
	if suffix {
		return fmt.Sprintf("%.1f%s", float64(size)/float64(base), suffixes[b])
	}
	return fmt.Sprintf("%d", size/base)
}

// Format return the size formatted according to sizeFormatting
func Format(size int64, sizeFormatting int) string {
	switch sizeFormatting {
	case SizesHuman:
		return Humanize(size)
	case SizesInK:
		return Base(size, KiB, false)
	case SizesInM:
		return Base(size, MiB, false)
	default:
		return fmt.Sprintf("%d", size)
	}
}
