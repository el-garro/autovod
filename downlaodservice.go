package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/log"
)

func DownloadService(channel string) {
	logger := log.NewWithOptions(
		os.Stderr,
		log.Options{
			Level:           log.InfoLevel,
			Prefix:          "DOWNLOAD",
			ReportTimestamp: true,
		},
	)
	logger.Info("Service started", "channel", channel)

	for ; true; <-time.Tick(time.Minute * 10) {
		online, err := IsOnline(channel)
		if err != nil {
			logger.Warn("Could not get online status", "err", err)
			continue
		}

		if !online {
			continue
		}

		vodurl, err := GetLatestVODUrl(channel)
		if err != nil {
			logger.Warn("Could not get VOD url", "err", err)
			continue
		}

		logger.Info("Checking yt-dlp version...")
		err = UpdateYtDlp()
		if err != nil {
			logger.Warn("Could not update yt-dlp", "err", err)
			continue
		}
		logger.Info("Downloading...", "url", vodurl)
		err = DownloadVOD(vodurl, Config.TwitchVideoHeight)
		if err != nil {
			logger.Warn(" Could not download VOD", "url", vodurl, "err", err)
			continue
		}

		logger.Info("Downloaded VOD", "url", vodurl)
	}
}

func UpdateYtDlp() error {
	cmd := exec.Command("python3", "-m", "pip", "install", "-U", "yt-dlp")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yt-dlp update failed: %s", stderr.String())
	}

	return nil
}
