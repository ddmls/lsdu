package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"syscall"

	"github.com/ddmls/lsdu/deep"
	"github.com/ddmls/lsdu/du"
	"github.com/ddmls/lsdu/human"
)

func main() {
	var sizesHuman, sizesInK, sizesInM, sortBySize, reportTotalSize, reportFreeSpace bool
	var sizeFormatting int
	flag.BoolVar(&sizesHuman, "human", false, "display size in KiB, MiB etc")
	flag.BoolVar(&sizesInK, "k", false, "display size in KiB")
	flag.BoolVar(&sizesInM, "m", false, "display size in MiB")
	flag.BoolVar(&sortBySize, "sort", true, "sort by size")
	flag.BoolVar(&reportTotalSize, "total", false, "report a total for all files")
	flag.BoolVar(&reportFreeSpace, "free", false, "report free space")
	flag.BoolVar(&du.ReportApparentSize, "apparent", false, "show apparent size")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s [OPTION]... [FILE|DIRECTORY]...\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	paths := []string{"."}
	if flag.NArg() > 0 {
		paths = flag.Args()
	}

	switch {
	case sizesInK:
		sizeFormatting = human.SizesInK
	case sizesInM:
		sizeFormatting = human.SizesInM
	case sizesHuman:
		sizeFormatting = human.SizesHuman
	default:
		sizeFormatting = human.SizesBytes
	}

	var grandTotalSize int64
	for _, path := range paths {
		dirEntries, totalSize, err := deep.ReadDirDeep(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if sortBySize {
			sort.Slice(dirEntries, func(i, j int) bool {
				return dirEntries[i].Size() < dirEntries[j].Size()
			})
		}
		deep.Print(dirEntries, sizeFormatting)
		grandTotalSize += totalSize
	}
	if reportTotalSize {
		fmt.Printf("Total size %s\n", human.Format(grandTotalSize, sizeFormatting))
	}
	if reportFreeSpace {
		// Remove duplicate filesystems with same Fsid
		type df struct {
			path       string
			freeSpace  uint64
			totalSpace uint64
		}
		uniqueDf := make(map[syscall.Fsid]df)
		for _, path := range paths {
			freeSpace, totalSpace, fsid, err := du.FreeSpace(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Could not get free disk information for", path)
				continue
			}
			uniqueDf[fsid] = df{path, freeSpace, totalSpace}

		}
		for _, v := range uniqueDf {
			fmt.Printf("Free space %s/%s: %s\n", human.Format(int64(v.freeSpace), sizeFormatting), human.Format(int64(v.totalSpace), sizeFormatting), v.path)
		}
	}
}
