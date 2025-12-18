package main

import (
	"fmt"

	"github.com/downing/media-manager/domain/sorting"
	"github.com/downing/media-manager/pkg/config"
	"go.uber.org/zap"
)

func importRawFiles(logger *zap.Logger, sortingService *sorting.Service) error {
	logger.Info("..... Starting import of raw files")

	lastImportDate, err := config.GetLastImportDate()
	if err != nil {
		return fmt.Errorf("failed to get last import date: %w", err)
	}
	logger.Info("::::. Last import date retrieved", zap.Time("last_import_date", lastImportDate))

	err = sortingService.ImportRawFiles(lastImportDate)
	if err != nil {
		return fmt.Errorf("failed to import raw files: %w", err)
	}
	logger.Info("::::: Import of raw files completed")
}
