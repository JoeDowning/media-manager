package main

import (
	"github.com/downing/media-manager/domain/sorting"
	"github.com/downing/media-manager/domain/uploader"
	"go.uber.org/zap"
)

func uploadEditedFiles(logger *zap.Logger, sortingService *sorting.Service, uploaderService *uploader.Service) error {
	logger.Info(". Starting upload of edited files")
	editedFiles, err := sortingService.ListEditedFilesToUpload()
	if err != nil {
		return err
	}

	err = uploaderService.UploadEditedFiles(editedFiles)
	if err != nil {
		return err
	}
	logger.Info(": Upload of edited files completed")
	return nil
}
