package deep

import (
	"fmt"
	"os"

	"github.com/ddmls/lsdu/du"
	"github.com/ddmls/lsdu/human"
)

// du -s -m *
// if initial path is a symbolic link to a dir it is followed
// however symbolinc links under it are not followed (so not checking for circular links)
// (ReadDir uses Lstat)
// hard links are not checked (they will be accounted more than once)

// SizesHuman, SizesInK, SizesInK choose the formatting of sizes used by Print
const (
	SizesBytes = iota
	SizesInK
	SizesInM
	SizesHuman
)

// FileInfo describes a file or a directory and all the files beneath it
type FileInfo struct {
	du.FileInfo
	deepSize int64
}

// Size returns the size (apparent or allocated) of a file or all the files beneath a directory
func (fi FileInfo) Size() int64 {
	return fi.deepSize
}

func (fi FileInfo) String() string {
	return fmt.Sprintf("%v %d %s '%s'", fi.Mode(), fi.Size(), fi.ModTime().Format("2006-01-02"), fi.Name())
}

// Print displays directory entries with specified formatting and automatic padding
func Print(fis []FileInfo, sizeFormatting int) {
	var formattedSizes []string
	padding := 0
	for _, fi := range fis {
		var formattedSize string
		switch sizeFormatting {
		case SizesHuman:
			formattedSize = human.Humanize(fi.Size())
		case SizesInK:
			formattedSize = human.Base(fi.Size(), human.KiB, false)
		case SizesInM:
			formattedSize = human.Base(fi.Size(), human.MiB, false)
		default:
			formattedSize = fmt.Sprintf("%d", fi.Size())
		}
		formattedSizes = append(formattedSizes, formattedSize)
		if len(formattedSize) > padding {
			padding = len(formattedSize)
		}
	}

	for i, fi := range fis {
		fmt.Printf("%v %*s %s '%s'\n", fi.Mode(), padding, formattedSizes[i], fi.ModTime().Format("2006-01-02"), fi.Name())
	}

}

func visitDir(path string,
	prevDir string,
	f func([]du.FileInfo) error) error {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()

	dirnames, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}

	if err := dir.Chdir(); err != nil {
		return err
	}

	var fileInfos []du.FileInfo
	for _, name := range dirnames {
		fileInfo, err := du.Lstat(name)
		if err != nil {
			return err
		}
		fileInfos = append(fileInfos, fileInfo)
	}

	if err := f(fileInfos); err != nil {
		return err
	}
	if prevDir != "" {
		if err := os.Chdir(prevDir); err != nil {
			return err
		}
	}

	return nil
}

// paths of fileInfo are considered relative to current dir
// as returned by readdir(path); chdir(path)
// note we do not include the dir size reported by stat (for its entries), just the size of files.
// The size reported is apparent size, not adjusted for cluster waste or holes.
func deepSize(
	fileInfo du.FileInfo,
) (int64, error) {
	if !fileInfo.IsDir() {
		return fileInfo.Size(), nil
	}

	var totalSize int64
	err := visitDir(fileInfo.Name(), "..", func(fileInfos []du.FileInfo) error {
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

// ReadDirDeep reads a directory and all the files beneath it or a file
func ReadDirDeep(path string) ([]FileInfo, error) {
	if fileInfo, err := du.Lstat(path); err != nil {
		return nil, err
	} else if !fileInfo.IsDir() {
		return []FileInfo{{fileInfo, fileInfo.Size()}}, nil
	}

	var dirEntries []FileInfo
	err := visitDir(path, "", func(fileInfos []du.FileInfo) error {
		for _, fileInfo := range fileInfos {
			size, err := deepSize(fileInfo)
			if err != nil {
				return err
			}
			entry := FileInfo{fileInfo, size}
			dirEntries = append(dirEntries, entry)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return dirEntries, nil
}
