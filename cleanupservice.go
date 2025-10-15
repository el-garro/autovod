package main

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func CleanupService() {
	log.Info("CLEANUP_SERVICE started")

	for ; true; <-time.Tick(time.Minute * 10) {
		path := "./vods"
		files, err := GetFileList(path)
		if err != nil {
			continue
		}

		for _, f := range files {
			if time.Since(f.ModTime) < Config.DeleteAfter {
				continue
			}

			err := os.Remove(path + "/" + f.Name)
			if err != nil {
				log.Info("CLEANUP_SERVICE Could not remove file", "name", f.Name, "err", err)
			} else {
				log.Info("CLEANUP_SERVICE Removed file", "name", f.Name)
			}
		}
	}
}
