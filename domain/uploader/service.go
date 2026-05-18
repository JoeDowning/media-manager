package uploader

import (
	"fmt"

	runtimestats "github.com/downing/media-manager/pkg/runtime_stats"
	"go.uber.org/zap"
)

type fileUploader interface {
	UploadFile(filePath string) error
}

type Service struct {
	logger   *zap.Logger
	uploader fileUploader
	stats    *runtimestats.Stats
}

func NewService(logger *zap.Logger, uploader fileUploader, stats *runtimestats.Stats) *Service {
	return &Service{
		logger:   logger,
		uploader: uploader,
		stats:    stats,
	}
}

func (s *Service) UploadEditedFiles(files []string) error {
	for _, file := range files {
		err := s.uploader.UploadFile(file)
		if err != nil {
			return fmt.Errorf("failed to upload file %s: %w", file, err)
		}
	}
	return nil
}
