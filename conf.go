package main

import (
	"os"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

type cfg struct {
	TwitchClientID    string
	TwitchChannel     string
	TwitchVideoHeight string
}

var Config *cfg

func LoadConfig() {
	Config = &cfg{}
	godotenv.Overload()

	Config.TwitchClientID = os.Getenv("TWITCH_CLIENT_ID")
	Config.TwitchChannel = os.Getenv("TWITCH_CHANNEL")
	Config.TwitchVideoHeight = os.Getenv("TWITCH_VIDEO_HEIGHT")

	height, _ := strconv.Atoi(Config.TwitchVideoHeight)
	if height <= 0 {
		Config.TwitchVideoHeight = "4096"
	}

	if Config.TwitchChannel == "" || Config.TwitchClientID == "" {
		log.Fatal("Invalid config", "TWITCH_CLIENT_ID", Config.TwitchClientID, "TWITCH_CHANNEL", Config.TwitchChannel, "TWITCH_VIDEO_HEIGHT", Config.TwitchVideoHeight)
	}

}
