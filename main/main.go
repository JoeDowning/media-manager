package main

import (
	"time"

	"github.com/downing/media-manager/domain/files"
	"github.com/downing/media-manager/domain/sorting"
	"github.com/downing/media-manager/domain/uploader"
	"github.com/downing/media-manager/pkg/config"
	"github.com/downing/media-manager/pkg/genutils"
	googlephotos "github.com/downing/media-manager/pkg/google_photos"
	"github.com/downing/media-manager/pkg/logging"
	runtimestats "github.com/downing/media-manager/pkg/runtime_stats"

	"go.uber.org/zap"
)

func main() {
	startTime := time.Now()
	var importDuration, backupRawDuration, backupEditedDuration, uploadEditedDuration time.Duration
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	logger := logging.NewLogger(cfg.LogLevel())
	defer logger.Sync()

	cfg.LogConfig(logger)

	logger.Info("Media Manager started")

	stats := runtimestats.NewStats()

	fileManager := files.NewService(genutils.NewFileManager())
	sortingService := sorting.NewService(
		logger,
		fileManager,
		toSortingCtiteria(cfg),
		stats,
	)
	photosUploader := googlephotos.NewUploader()
	uploaderService := uploader.NewService(
		logger,
		photosUploader,
		stats,
	)

	if cfg.ImportRaw() {
		importStart := time.Now()
		err := importRawFiles(logger, sortingService)
		if err != nil {
			logger.Error("Failed to import raw files", zap.Error(err))
			return
		}
		importDuration = time.Since(importStart)
	}

	if cfg.BackupRaw() {
		backupRawStart := time.Now()
		err := backupRawFiles(logger, sortingService)
		if err != nil {
			logger.Error("Failed to backup raw files", zap.Error(err))
			return
		}
		backupRawDuration = time.Since(backupRawStart)
	}

	if cfg.BackupEdited() {
		backupEditedStart := time.Now()
		err := backupEditedFiles(logger, sortingService)
		if err != nil {
			logger.Error("Failed to backup edited files", zap.Error(err))
			return
		}
		backupEditedDuration = time.Since(backupEditedStart)
	}

	if cfg.UploadEdited() {
		uploadEditedStart := time.Now()
		err := uploadEditedFiles(logger, sortingService, uploaderService)
		if err != nil {
			logger.Error("Failed to upload edited files", zap.Error(err))
			return
		}
		uploadEditedDuration = time.Since(uploadEditedStart)
	}

	stats.FinalStats(logger)
	logMsg := "Media Manager completed in " + time.Since(startTime).Round(500*time.Millisecond).String()
	if cfg.ImportRaw() {
		logMsg += ", Import Duration: " + importDuration.Round(500*time.Millisecond).String()
	}
	if cfg.BackupRaw() {
		logMsg += ", Backup Raw Duration: " + backupRawDuration.Round(500*time.Millisecond).String()
	}
	if cfg.BackupEdited() {
		logMsg += ", Backup Edited Duration: " + backupEditedDuration.Round(500*time.Millisecond).String()
	}
	if cfg.UploadEdited() {
		logMsg += ", Upload Edited Duration: " + uploadEditedDuration.Round(500*time.Millisecond).String()
	}
	logger.Info(logMsg)
}

func toSortingCtiteria(cfg config.Config) sorting.SortCriteria {
	return sorting.SortCriteria{
		FileTypes:       []string{".jpg", ".png", ".mp4", ".mov"},
		RawPath:         cfg.RawPath(),
		LocalRawPath:    cfg.LocalRawPath(),
		LocalEditedPath: cfg.LocalEditedPath(),
		BackupPath:      cfg.BackupPath(),
		MoveFiles:       cfg.MoveFiles(),
		CopyFiles:       cfg.CopyFiles(),
	}
}
