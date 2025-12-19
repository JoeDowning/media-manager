package main

import (
	"fmt"

	"github.com/downing/media-manager/domain/sorting"
	"github.com/downing/media-manager/pkg/config"
	"go.uber.org/zap"
)

func importRawFiles(logger *zap.Logger, sortingService *sorting.Service) error {
	logger.Info("Starting import of raw files")

	lastImportDate, err := config.GetLastImportDate(logger)
	if err != nil {
		return fmt.Errorf("failed to get last import date: %w", err)
	}
	logger.Info("Last import date retrieved", zap.Time("last_import_date", lastImportDate))

	lastImportTime, err := sortingService.ImportRawFiles(lastImportDate)
	if err != nil {
		return fmt.Errorf("failed to import raw files: %w", err)
	}

	err = config.SetLastImportDate(lastImportTime)
	if err != nil {
		return fmt.Errorf("failed to set last import date: %w", err)
	}
	logger.Info("Last import date updated", zap.Time("new_last_import_date", lastImportTime))

	logger.Info("Import of raw files completed")
	return nil
}
