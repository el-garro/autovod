package main

import (
	"time"

	"github.com/charmbracelet/log"
)

func DownloadService(channel string) {
	log.Info("DOWNLOADSERVICE started", "channel", channel)
	for ; true; <-time.Tick(time.Minute * 30) {
		online, err := IsOnline(channel)
		if err != nil {
			log.Warn("DOWNLOADSERVICE Could not get online status", "err", err)
			continue
		}

		if !online {
			continue
		}

		vodurl, err := GetLatestVODUrl(channel)
		if err != nil {
			log.Warn("DOWNLOADSERVICE Could not get VOD url", "err", err)
			continue
		}

		log.Info("DOWNLOADSERVICE Downloading...", "url", vodurl)
		err = DownloadVOD(vodurl, Config.TwitchVideoHeight)
		if err != nil {
			log.Warn("DOWNLOADSERVICE Could not download VOD", "url", vodurl, "err", err)
			continue
		}

		log.Info("DOWNLOADSERVICE Downloaded VOD", "url", vodurl)
	}
}
