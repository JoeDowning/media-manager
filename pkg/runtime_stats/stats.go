package runtimestats

import (
	"fmt"

	"go.uber.org/zap"
)

type Stats struct {
	RawFilesChecked  int
	RawFilesFound    int
	RawFilesImported int

	LocalRawFilesChecked int
	LocalRawFilesFound   int
	LocalRawFilesMoved   int
	LocalRawFilesCopied  int

	LocalEditedFilesChecked int
	LocalEditedFilesFound   int
	LocalEditedFilesMoved   int
	LocalEditedFilesCopied  int

	ToUploadFilesChecked  int
	ToUploadFilesFound    int
	ToUploadFilesUploaded int
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) FinalStats(logger *zap.Logger) {
	logger.Info("Raw Statistics",
		zap.Int("raw_files_checked", s.RawFilesChecked),
		zap.Int("raw_files_found", s.RawFilesFound),
		zap.Int("raw_files_imported", s.RawFilesImported),

		zap.Int("local_raw_files_checked", s.LocalRawFilesChecked),
		zap.Int("local_raw_files_found", s.LocalRawFilesFound),
		zap.Int("local_raw_files_moved", s.LocalRawFilesMoved),
		zap.Int("local_raw_files_copied", s.LocalRawFilesCopied),

		zap.Int("local_edited_files_checked", s.LocalEditedFilesChecked),
		zap.Int("local_edited_files_found", s.LocalEditedFilesFound),
		zap.Int("local_edited_files_moved", s.LocalEditedFilesMoved),
		zap.Int("local_edited_files_copied", s.LocalEditedFilesCopied),

		zap.Int("to_upload_files_checked", s.ToUploadFilesChecked),
		zap.Int("to_upload_files_found", s.ToUploadFilesFound),
		zap.Int("to_upload_files_uploaded", s.ToUploadFilesUploaded),
	)

	logMsg := "Raw Files:          Checked: %d, Found: %d, Imported: %d"
	logger.Info(
		fmt.Sprintf(logMsg, s.RawFilesChecked, s.RawFilesFound, s.RawFilesImported),
	)

	logMsg = "Local Raw Files:    Checked: %d, Found: %d, Moved: %d, Copied: %d"
	logger.Info(
		fmt.Sprintf(logMsg, s.LocalRawFilesChecked, s.LocalRawFilesFound, s.LocalRawFilesMoved, s.LocalRawFilesCopied),
	)

	logMsg = "Local Edited Files: Checked: %d, Found: %d, Moved: %d, Copied: %d"
	logger.Info(
		fmt.Sprintf(logMsg, s.LocalEditedFilesChecked, s.LocalEditedFilesFound, s.LocalEditedFilesMoved, s.LocalEditedFilesCopied),
	)

	logMsg = "To Upload Files:    Checked: %d, Found: %d, Uploaded: %d"
	logger.Info(
		fmt.Sprintf(logMsg, s.ToUploadFilesChecked, s.ToUploadFilesFound, s.ToUploadFilesUploaded),
	)

	totalFilesChecked := s.RawFilesChecked + s.LocalRawFilesChecked + s.LocalEditedFilesChecked + s.ToUploadFilesChecked
	totalFilesFound := s.RawFilesFound + s.LocalRawFilesFound + s.LocalEditedFilesFound + s.ToUploadFilesFound
	totalFilesProcessed := s.RawFilesImported + s.LocalRawFilesMoved + s.LocalRawFilesCopied + s.LocalEditedFilesMoved + s.LocalEditedFilesCopied + s.ToUploadFilesUploaded

	logMsg = "Totals:             Checked: %d, Found: %d, Processed: %d"
	logger.Info(
		fmt.Sprintf(logMsg, totalFilesChecked, totalFilesFound, totalFilesProcessed),
	)
}
