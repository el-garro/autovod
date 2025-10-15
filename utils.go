package main

import (
	"os"
	"time"

	"github.com/dustin/go-humanize"
)

type FileInfo struct {
	Name    string
	Size    string
	ModTime time.Time
}

func GetFileList(dir string) ([]FileInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		files = append(files, FileInfo{
			Name:    info.Name(),
			Size:    humanize.Bytes(uint64(info.Size())),
			ModTime: info.ModTime(),
		})
	}

	return files, nil
}
