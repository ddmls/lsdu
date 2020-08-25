package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

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
	flag.BoolVar(&reportTotalSize, "total", false, "report free space")
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

	for i, path := range paths {
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
		if reportTotalSize {
			fmt.Printf("Total size %s\n", human.Format(totalSize, sizeFormatting))
		}
		if reportFreeSpace {
			freeSpace, totalSpace, err := du.FreeSpace(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Printf("Free space %s/%s\n", human.Format(int64(freeSpace), sizeFormatting), human.Format(int64(totalSpace), sizeFormatting))
		}
		if i < len(paths)-1 {
			fmt.Println()
		}
	}
}
