package main

import (
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

type cfg struct {
	TwitchClientID    string        `env:"TWITCH_CLIENT_ID,required,notEmpty"`
	TwitchChannel     string        `env:"TWITCH_CHANNEL,required,notEmpty"`
	TwitchVideoHeight uint          `env:"TWITCH_VIDEO_HEIGHT" envDefault:"4096"`
	DeleteAfter       time.Duration `env:"DELETE_VOD_AFTER" envDefault:"72h"`
	WebPort           uint          `env:"WEB_SERVER_PORT" envDefault:"8080"`
}

var Config cfg

func LoadConfig() {
	logger := log.NewWithOptions(
		os.Stderr,
		log.Options{
			Level:           log.InfoLevel,
			Prefix:          "CFG",
			ReportTimestamp: true,
		},
	)

	godotenv.Overload()

	err := env.Parse(&Config)
	if err != nil {
		logger.Fatal("Invalid config", "err", err)
	}

	logger.Info("Loaded config", "cfg", Config)
}
