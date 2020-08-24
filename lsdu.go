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
	var sortBySize, reportFreeSpace bool
	flag.BoolVar(&deep.HumanSizes, "human", false, "display size in KiB, MiB etc")
	flag.BoolVar(&sortBySize, "sort", true, "sort by size")
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

	for i, path := range paths {
		dirEntries, err := deep.ReadDirDeep(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if sortBySize {
			sort.Slice(dirEntries, func(i, j int) bool {
				return dirEntries[i].Size() < dirEntries[j].Size()
			})
		}
		for _, entry := range dirEntries {
			fmt.Println(entry)
		}
		if reportFreeSpace {
			freeSpace, totalSpace, err := du.FreeSpace(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Printf("Free space %s/%s\n", human.Humanize(int64(freeSpace)), human.Humanize(int64(totalSpace)))
		}
		if i < len(paths)-1 {
			fmt.Println()
		}
	}
}
