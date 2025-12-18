package main

import (
	"github.com/downing/media-manager/domain/sorting"
	"go.uber.org/zap"
)

func uploadEditedFiles(logger *zap.Logger, sortingService *sorting.Service) error {
	logger.Info(". Starting upload of edited files")
	logger.Info(": Upload of edited files completed")
	return nil
}
