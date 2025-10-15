package main

import (
	"time"

	"github.com/charmbracelet/log"
)

func main() {
	log.SetLevel(log.InfoLevel)

	log.Info("Initializing...")

	LoadConfig()
	log.Info("Config", "TWITCH_CLIENT_ID", Config.TwitchClientID, "TWITCH_CHANNEL", Config.TwitchChannel, "TWITCH_VIDEO_HEIGHT", Config.TwitchVideoHeight)

	log.Info("Starting services...")

	go WebService()
	go CleanupService()
	go DownloadService(Config.TwitchChannel)

	for ; true; <-time.Tick(time.Hour) {
	}
}
