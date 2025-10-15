package main

import (
	"time"

	"github.com/charmbracelet/log"
)

func DownloadService(channel string) {
	log.Info("DOWNLOAD_SERVICE started", "channel", channel)
	for ; true; <-time.Tick(time.Minute * 10) {
		online, err := IsOnline(channel)
		if err != nil {
			log.Warn("DOWNLOAD_SERVICE Could not get online status", "err", err)
			continue
		}

		if !online {
			continue
		}

		vodurl, err := GetLatestVODUrl(channel)
		if err != nil {
			log.Warn("DOWNLOAD_SERVICE Could not get VOD url", "err", err)
			continue
		}

		log.Info("DOWNLOAD_SERVICE Downloading...", "url", vodurl)
		err = DownloadVOD(vodurl, Config.TwitchVideoHeight)
		if err != nil {
			log.Warn("DOWNLOAD_SERVICE Could not download VOD", "url", vodurl, "err", err)
			continue
		}

		log.Info("DOWNLOAD_SERVICE Downloaded VOD", "url", vodurl)
	}
}
