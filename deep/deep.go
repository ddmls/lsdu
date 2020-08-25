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
		formattedSize := human.Format(fi.Size(), sizeFormatting)
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

// du.Fileinfo.Name() doesn't contain the full path, it is considered relative to current dir as returned by readdir(path); chdir(path)
// so someone has to call chdir to the directory that contains it before calling deepSize
// this is done by visitDir when deepSize calls itself
// note we do not include the dir size reported by stat (for its entries), just the size of files.
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
func ReadDirDeep(path string) ([]FileInfo, int64, error) {
	if fileInfo, err := du.Lstat(path); err != nil {
		return nil, 0, err
	} else if !fileInfo.IsDir() {
		return []FileInfo{{fileInfo, fileInfo.Size()}}, fileInfo.Size(), nil
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, 0, err
	}
	var dirEntries []FileInfo
	var totalSize int64
	err = visitDir(path, pwd, func(fileInfos []du.FileInfo) error {
		for _, fileInfo := range fileInfos {
			size, err := deepSize(fileInfo)
			if err != nil {
				return err
			}
			entry := FileInfo{fileInfo, size}
			dirEntries = append(dirEntries, entry)
			totalSize += size
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return dirEntries, totalSize, nil
}
