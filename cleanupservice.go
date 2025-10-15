package main

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func CleanupService() {
	log.Info("CLEANUPSERVICE started")

	for ; true; <-time.Tick(time.Minute) {
		path := "./vods"
		files, err := GetFileList(path)
		if err != nil {
			continue
		}

		for _, f := range files {
			if time.Since(f.ModTime) < time.Hour*72 {
				continue
			}

			err := os.Remove(path + "/" + f.Name)
			if err != nil {
				log.Info("CLEANUPSERVICE Could not remove file", "name", f.Name, "err", err)
			} else {
				log.Info("CLEANUPSERVICE Removed file", "name", f.Name)
			}
		}
	}
}
