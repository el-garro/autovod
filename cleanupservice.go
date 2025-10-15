package main

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func CleanupService() {
	logger := log.NewWithOptions(
		os.Stderr,
		log.Options{
			Level:           log.InfoLevel,
			Prefix:          "CLEANUP",
			ReportTimestamp: true,
		},
	)

	logger.Info("Service started")

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
				logger.Info("Could not remove file", "name", f.Name, "err", err)
			} else {
				logger.Info("Removed file", "name", f.Name)
			}
		}
	}
}
