package main

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

const DOWNLOAD_DIR = "./vods"

func main() {
	logger := log.NewWithOptions(
		os.Stderr,
		log.Options{
			Level:           log.InfoLevel,
			Prefix:          "MAIN",
			ReportTimestamp: true,
		},
	)

	logger.Info("Initializing...")
	LoadConfig()

	logger.Info("Starting services...")

	go WebService()
	go CleanupService()
	go DownloadService(Config.TwitchChannel)

	for ; true; <-time.Tick(time.Hour) {
	}
}
