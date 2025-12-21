package main

import (
	"github.com/downing/media-manager/domain/sorting"
	"go.uber.org/zap"
)

func backupRawFiles(logger *zap.Logger, sortingService *sorting.Service) error {
	logger.Info("Starting backup of raw files")

	err := sortingService.BackupLocalRawFiles()
	if err != nil {
		return err
	}

	logger.Info("Backup of raw files completed")
	return nil
}
