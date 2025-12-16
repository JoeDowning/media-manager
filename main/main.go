package main

import (
	"github.com/JoeDowning/media-manager/pkg/config"
	"github.com/JoeDowning/media-manager/pkg/logging"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	logger := logging.NewLogger(cfg.LogLevel)
	defer logger.Sync()

	logger.Info("Media Manager started",
		logging.Field("log_level", cfg.LogLevel),
	)
}
