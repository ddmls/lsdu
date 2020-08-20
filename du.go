package main

import (
	"du/human"
	"flag"
	"fmt"
	"os"
)

// du -s -m *
// if initial path is a symbolic link to a dir it is followed
// however symbolinc links under it are not followed (so not checking for circular links)
// (ReadDir uses Lstat)
// hard links are not checked (they will be accounted more than once)

type deepFileInfo struct {
	os.FileInfo
	deepSize int64
}

var humanSizes bool

func (d deepFileInfo) String() string {
	if humanSizes {
		return fmt.Sprintf("%v %s %s '%s'", d.Mode(), human.Humanize(d.deepSize), d.ModTime().Format("2006-01-02"), d.Name())
	}
	return fmt.Sprintf("%v %d %s '%s'", d.Mode(), d.deepSize, d.ModTime().Format("2006-01-02"), d.Name())
}

func visitDir(path string,
	relative bool,
	f func([]os.FileInfo) error) error {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	prevDir := ".."
	if !relative {
		prevDir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	if err := dir.Chdir(); err != nil {
		return err
	}
	if err := f(fileInfos); err != nil {
		return err
	}
	if err := os.Chdir(prevDir); err != nil {
		return err
	}

	return nil
}

// paths of fileInfo are considered relative to current dir
// as returned by readdir(path); chdir(path)
func deepSize(
	fileInfo os.FileInfo,
) (int64, error) {
	if !fileInfo.IsDir() {
		return fileInfo.Size(), nil
	}

	var totalSize int64
	err := visitDir(fileInfo.Name(), true, func(fileInfos []os.FileInfo) error {
		for _, fileInfo := range fileInfos {
			size, err := deepSize(fileInfo)
			if err != nil {
				return err
			}
			totalSize += size
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return totalSize, nil
}

func readDirDeep(path string) ([]deepFileInfo, error) {
	if fileInfo, err := os.Stat(path); err != nil {
		return nil, err
	} else if !fileInfo.IsDir() {
		return []deepFileInfo{{fileInfo, fileInfo.Size()}}, nil
	}

	var dirEntries []deepFileInfo
	err := visitDir(path, false, func(fileInfos []os.FileInfo) error {
		for _, fileInfo := range fileInfos {
			size, err := deepSize(fileInfo)
			if err != nil {
				return err
			}
			entry := deepFileInfo{fileInfo, size}
			dirEntries = append(dirEntries, entry)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return dirEntries, nil
}

func main() {
	flag.BoolVar(&humanSizes, "human", false, "display size in KiB, MiB etc")
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
		dirEntries, err := readDirDeep(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, entry := range dirEntries {
			fmt.Println(entry)
		}
		if i < len(paths)-1 {
			fmt.Println()
		}
	}
}
