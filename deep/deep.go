package deep

import (
	"errors"
	"fmt"
	"os"
	"unicode"

	"github.com/ddmls/lsdu/du"
	"github.com/ddmls/lsdu/human"
)

// du -s -m *
// if initial path is a symbolic link to a dir it is followed
// however symbolinc links under it are not followed (so not checking for circular links)
// (ReadDir uses Lstat)
// hard links are not checked (they will be accounted more than once)

// DirEntry describes a file or a directory and all the files beneath it
type DirEntry struct {
	du.FileInfo
	path     string
	deepSize int64
}

// Size returns the size (apparent or allocated) of a file or all the files beneath a directory
func (fi DirEntry) Size() int64 {
	return fi.deepSize
}

func (fi DirEntry) String() string {
	return fmt.Sprintf("%v %d %s \"%s\"", fi.Mode(), fi.Size(), fi.ModTime().Format("2006-01-02"), fi.Name())
}

func containsSpace(s string) bool {
	for _, rune := range s {
		if unicode.IsSpace(rune) {
			return true
		}
	}
	return false
}

// MaybeQuote quotes a string with single quotes if it contains spaces
func MaybeQuote(s string) string {
	if containsSpace(s) {
		return fmt.Sprintf("'%s'", s)
	}
	return s
}

// Print displays directory entries with specified formatting and automatic padding
func Print(fis []DirEntry, sizeFormatting int) {
	var formattedSizes []string
	padding := 0
	for _, fi := range fis {
		formattedSize := human.Format(fi.Size(), sizeFormatting)
		formattedSizes = append(formattedSizes, formattedSize)
		if len(formattedSize) > padding {
			padding = len(formattedSize)
		}
	}

	for i, fi := range fis {
		fmt.Printf("%v %*s %s %s\n", fi.Mode(), padding, formattedSizes[i], fi.ModTime().Format("2006-01-02"), MaybeQuote(fi.Name()))
	}
}

func readDir(path string) ([]du.FileInfo, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fileInfos, err := du.Readdir(dir, 0)
	if err != nil {
		return nil, err
	}
	return fileInfos, err
}

// ReadDirDeep reads a directory and the deep size of each entry
func ReadDirDeep(path string) ([]DirEntry, int64, error) {
	// This test is only needed for the initial caller. We only call ourselves when the path is a directory.
	fileInfo, err := du.Lstat(path)
	if err != nil {
		return nil, 0, err
	}
	if !fileInfo.IsDir() {
		return []DirEntry{{fileInfo, path, fileInfo.Size()}}, fileInfo.Size(), nil
	}

	var totalSize int64
	var dirEntries []DirEntry
	fileInfos, err := readDir(path)
	if err != nil {
		if errors.Is(err, os.ErrPermission) || errors.Is(err, os.ErrNotExist) {
			fmt.Fprintln(os.Stderr, err)
			return []DirEntry{{fileInfo, path, fileInfo.Size()}}, fileInfo.Size(), nil
		}
		return nil, 0, err
	}
	for _, fileInfo := range fileInfos {
		var size int64
		entryPath := path + string(os.PathSeparator) + fileInfo.Name()
		if !fileInfo.IsDir() {
			size = fileInfo.Size()
		} else {
			_, size, err = ReadDirDeep(entryPath)
			if err != nil {
				return nil, 0, err
			}
			size += fileInfo.Size()
		}
		totalSize += size
		dirEntries = append(dirEntries, DirEntry{fileInfo, entryPath, size})
	}

	return dirEntries, totalSize, nil
}
