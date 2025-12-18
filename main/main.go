package main

import (
	"github.com/JoeDowning/media-manager/domain/files"
	"github.com/JoeDowning/media-manager/domain/sorting"
	"github.com/JoeDowning/media-manager/pkg/config"
	"github.com/JoeDowning/media-manager/pkg/genutils"
	"github.com/JoeDowning/media-manager/pkg/logging"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	logger := logging.NewLogger(cfg.LogLevel)
	defer logger.Sync()

	cfg.LogConfig(logger)

	logger.Info("Media Manager started")

	fileManager := files.NewService(genutils.NewFileManager())
	sortingService := sorting.NewService(
		logger,
		[]string{".jpg", ".png", ".mp4", ".mov"},
		fileManager,
		cfg.RawPath(),
		cfg.LocalPath(),
		cfg.BackupPath(),
	)

	if cfg.ImportRaw() {
		err := importRawFiles(logger, sortingService)
		if err != nil {
			logger.Error("Failed to import raw files", zap.Error(err))
			return
		}
	}

	if cfg.BackupRaw() {
		err := backupRawFiles(logger, sortingService)
		if err != nil {
			logger.Error("Failed to backup raw files", zap.Error(err))
			return
		}
	}

	if cfg.BackupEdited() {
		err := backupEditedFiles(logger, sortingService)
		if err != nil {
			logger.Error("Failed to backup edited files", zap.Error(err))
			return
		}
	}

	if cfg.UploadEdited() {
		err := uploadEditedFiles(logger, sortingService)
		if err != nil {
			logger.Error("Failed to upload edited files", zap.Error(err))
			return
		}
	}

	logger.Info("Media Manager completed")
}
