package main

import (
	"fmt"
	"os"
)

// du -s -m *
// if path is a symbolic link to a dir it is followed
// however symbolinc links under it are not followed (so not checking for circular links)
// (ReadDir uses Lstat)
// hard links are not checked (they will be accounted more than once)

// const path = "/data/dimosd/Downloads"
const path = "/tmp/ama"

func errorExit(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func main() {
	dir, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open %s\n", path)
		os.Exit(1)
	}
	defer dir.Close()

	if fileInfo, err := dir.Stat(); err != nil || !fileInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "%s is not a directory\n", path)
		os.Exit(1)
	}

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		errorExit("Error reading contents of %s\n", path)
	}

	for _, fileInfo := range fileInfos {
		fmt.Println(fileInfo.Name(), fileInfo.Size())
	}

}
