package sorting

import (
	"fmt"
	"strings"
	"time"

	"github.com/downing/media-manager/domain/images"
	"go.uber.org/zap"
)

type SortCriteria struct {
	fileTypes []string

	rawPath    string
	localPath  string
	backupPath string
}

type Service struct {
	logger   *zap.Logger
	criteria SortCriteria
	files    fileManager
}

type fileManager interface {
	GetFilesInPath(path string) ([]string, error)
	GetFilesRecursivelyInPath(path string) ([]string, error)
	DoesFileExist(path string) (bool, error)
	DoesPathExist(path string) (bool, error)
	MoveFile(sourcePath, destinationPath string) error
	CopyFile(sourcePath, destinationPath string) error
}

func NewService(
	logging *zap.Logger,
	fileTypes []string,
	files fileManager,
	rawPath, localPath, backupPath string,
) *Service {
	return &Service{
		logger: logging,
		files:  files,
		criteria: SortCriteria{
			fileTypes:  fileTypes,
			rawPath:    rawPath,
			localPath:  localPath,
			backupPath: backupPath,
		},
	}
}

// ImportRawFiles imports raw files from the raw path to the local path based on the last import date.
func (s *Service) ImportRawFiles(lastImportDate time.Time) (time.Time, error) {
	files, err := s.files.GetFilesRecursivelyInPath(s.criteria.rawPath)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get files recursively in path [%s]: %w", s.criteria.rawPath, err)
	}
	s.logger.Info("Found files for import", zap.Int("file_count", len(files)))

	imageTypes := images.GetImageTypes()

	imageFiles := []string{}

	for _, file := range files {
		if !fileTypeIsInList(file, imageTypes) {
			s.logger.Debug("Skipping non-image file", zap.String("file", file))
			continue
		}
		imageFiles = append(imageFiles, file)
	}
	s.logger.Info("Filtered image files for import", zap.Int("image_file_count", len(imageFiles)))

	var importCount int
	var newestTime time.Time
	for _, file := range imageFiles {
		imgData, err := images.GetPhoto(nil, file)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to get photo data for file [%s]: %w", file, err)
		}

		imgTime := imgData.GetTimestamp()
		if imgTime.After(lastImportDate) {
			destPath := generateRawImportDestinationPath(s.criteria.localPath, imgData.GetFileName(), imgTime)
			err = s.files.CopyFile(file, destPath)
			if err != nil {
				return time.Time{}, fmt.Errorf("failed to copy file [%s] to [%s]: %w", file, destPath, err)
			}
			importCount++
			s.logger.Debug("Imported raw file", zap.String("source", file), zap.String("destination", destPath))
		}
		if imgTime.After(newestTime) {
			newestTime = imgTime
		}
	}

	s.logger.Info("Import raw files completed", zap.Int("imported_files_count", importCount))
	return newestTime, nil
}

func (s *Service) BackupLocalRawFiles() error {
	// Implementation for copying raw files to local based on criteria
	return nil
}

func (s *Service) BackupEditedFiles() error {
	// Implementation for copying edited files based on criteria
	return nil
}

func fileTypeIsInList(filePath string, fileTypes []string) bool {
	for _, fileType := range fileTypes {
		if len(filePath) >= len(fileType)+1 && strings.EqualFold(filePath[len(filePath)-len(fileType)-1:], "."+fileType) {
			return true
		}
	}
	return false
}

func generateRawImportDestinationPath(basePath, fileName string, timestamp time.Time) string {
	return fmt.Sprintf("%s/%d-%02d-%02d/%s", basePath, timestamp.Year(), timestamp.Month(), timestamp.Day(), fileName)
}
