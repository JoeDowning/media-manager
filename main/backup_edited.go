package main

import (
	"github.com/downing/media-manager/domain/sorting"
	"go.uber.org/zap"
)

func backupEditedFiles(logger *zap.Logger, sortingService *sorting.Service) error {
	logger.Info("Starting backup of edited files")

	err := sortingService.BackupEditedFiles()
	if err != nil {
		return err
	}

	logger.Info("Backup of edited files completed")
	return nil
}
