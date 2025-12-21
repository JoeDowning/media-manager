package sorting

import (
	"fmt"
	"strings"
	"time"

	"github.com/downing/media-manager/domain/images"
	runtimestats "github.com/downing/media-manager/pkg/runtime_stats"
	"go.uber.org/zap"
)

type SortCriteria struct {
	FileTypes []string

	RawPath         string
	LocalRawPath    string
	LocalEditedPath string
	BackupPath      string

	MoveFiles bool
	CopyFiles bool
}

type Service struct {
	logger   *zap.Logger
	criteria SortCriteria
	files    fileManager
	stats    *runtimestats.Stats
}

type fileManager interface {
	GetFilesInPath(path string) ([]string, error)
	GetFilesRecursivelyInPath(path string) ([]string, error)
	DoesFileExist(path string) (bool, error)
	DoesPathExist(path string) (bool, error)
	MoveFile(sourcePath, destinationPath string) error
	CopyFile(sourcePath, destinationPath string) error
}

type statsManager interface {
	IncrementCounter(name string)
}

func NewService(
	logging *zap.Logger,
	files fileManager,
	sortingCriteria SortCriteria,
	stats *runtimestats.Stats,
) *Service {
	return &Service{
		logger:   logging,
		files:    files,
		criteria: sortingCriteria,
		stats:    stats,
	}
}

// ImportRawFiles imports raw files from the raw path to the local path based on the last import date.
func (s *Service) ImportRawFiles(lastImportDate time.Time) (time.Time, error) {
	files, err := s.files.GetFilesRecursivelyInPath(s.criteria.RawPath)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get files recursively in path [%s]: %w", s.criteria.RawPath, err)
	}
	s.logger.Info("Found files for import", zap.Int("file_count", len(files)))
	s.stats.RawFilesChecked += len(files)

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
	s.stats.RawFilesFound += len(imageFiles)

	var filesChecked int
	var newestTime time.Time
	for _, file := range imageFiles {
		filesChecked++
		logMsg := fmt.Sprintf("%d files remaining", len(imageFiles)-filesChecked)

		// get photo data
		imgData, err := images.GetPhoto(nil, file)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to get photo data for file [%s]: %w", file, err)
		}

		// get timestamp and if before last import date, skip
		imgTime := imgData.GetTimestamp()
		if imgTime.Before(lastImportDate) || imgTime.Equal(lastImportDate) {
			s.logger.Debug(logMsg, zap.String("file", file), zap.Bool("imported", false))
			continue
		}

		// create the new path of format <localRawPath>/<year>-<month>-<day>/<filename>
		destPath := generateRawImportDestinationPath(s.criteria.LocalRawPath, imgData.GetFileName(), imgTime)

		// copy the file to the new location
		err = s.files.CopyFile(file, destPath)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to copy file [%s] to [%s]: %w", file, destPath, err)
		}

		// update newest time
		if imgTime.After(newestTime) {
			newestTime = imgTime
		}
		s.logger.Debug(logMsg, zap.String("file", file), zap.Bool("imported", true))
		s.stats.RawFilesImported++
	}

	s.logger.Info("Import raw files completed", zap.Int("imported_files_count", s.stats.RawFilesImported), zap.Int("file_count", len(files)))
	return newestTime, nil
}

func (s *Service) BackupLocalRawFiles() error {
	files, err := s.files.GetFilesRecursivelyInPath(s.criteria.LocalRawPath)
	if err != nil {
		return fmt.Errorf("failed to get files recursively in path [%s]: %w", s.criteria.LocalRawPath, err)
	}
	s.logger.Info("Found files for import", zap.Int("file_count", len(files)))
	s.stats.LocalRawFilesChecked += len(files)

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
	s.stats.LocalRawFilesFound += len(imageFiles)

	var filesChecked int
	for _, file := range imageFiles {
		filesChecked++
		logMsg := fmt.Sprintf("%d files remaining", len(imageFiles)-filesChecked)

		// get photo data
		imgData, err := images.GetPhoto(nil, file)
		if err != nil {
			return fmt.Errorf("failed to get photo data for file [%s]: %w", file, err)
		}

		// create the new path of format <backupPath>/<year>/<month>/<day><hour><minute><second>_<filename>
		imgTime := imgData.GetTimestamp()
		destPath := generateRawBackupDestinationPath(s.criteria.BackupPath, imgData.GetFileName(), imgTime)

		// check if the file with the new name already exists at the destination
		exists, err := s.files.DoesFileExist(destPath)
		if err != nil {
			return fmt.Errorf("failed to check if file exists at destination [%s]: %w", destPath, err)
		}
		if exists {
			s.logger.Debug(logMsg, zap.String("file", file), zap.Bool("backed_up", false), zap.String("reason", "file already exists at destination"))
			continue
		}

		if s.criteria.CopyFiles {
			// copy the file to the new location
			err = s.files.CopyFile(file, destPath)
			if err != nil {
				return fmt.Errorf("failed to copy file [%s] to [%s]: %w", file, destPath, err)
			}
			s.stats.LocalRawFilesCopied++
		} else if s.criteria.MoveFiles {
			// move the file to the new location
			err = s.files.MoveFile(file, destPath)
			if err != nil {
				return fmt.Errorf("failed to move file [%s] to [%s]: %w", file, destPath, err)
			}
			s.stats.LocalRawFilesMoved++
		} else {
			s.logger.Warn("No file operation specified (neither move nor copy)", zap.String("file", file))
			continue
		}

		s.logger.Debug(logMsg, zap.String("file", file), zap.Bool("backed_up", true))
	}

	s.logger.Info("Backup of local raw files completed", zap.Int("file_count", s.stats.LocalRawFilesCopied+s.stats.LocalRawFilesMoved))
	return nil
}

func (s *Service) BackupEditedFiles() error {
	files, err := s.files.GetFilesRecursivelyInPath(s.criteria.LocalEditedPath)
	if err != nil {
		return fmt.Errorf("failed to get files recursively in path [%s]: %w", s.criteria.LocalEditedPath, err)
	}
	s.logger.Info("Found files for import", zap.Int("file_count", len(files)))
	s.stats.LocalEditedFilesChecked += len(files)

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
	s.stats.LocalEditedFilesFound += len(imageFiles)

	var filesChecked int
	for _, file := range imageFiles {
		filesChecked++
		logMsg := fmt.Sprintf("%d files remaining", len(imageFiles)-filesChecked)

		// get photo data
		imgData, err := images.GetPhoto(nil, file)
		if err != nil {
			return fmt.Errorf("failed to get photo data for file [%s]: %w", file, err)
		}

		// create the new path of format <backupPath>/<year>/<month>/<day><hour><minute><second>_<filename>
		imgTime := imgData.GetTimestamp()
		destPath := generateEditedBackupDestinationPath(s.criteria.BackupPath, imgData.GetFileName(), imgTime)

		// check if the file with the new name already exists at the destination
		exists, err := s.files.DoesFileExist(destPath)
		if err != nil {
			return fmt.Errorf("failed to check if file exists at destination [%s]: %w", destPath, err)
		}
		if exists {
			s.logger.Debug(logMsg, zap.String("file", file), zap.Bool("backed_up", false), zap.String("reason", "file already exists at destination"))
			continue
		}

		if s.criteria.CopyFiles {
			// copy the file to the new location
			err = s.files.CopyFile(file, destPath)
			if err != nil {
				return fmt.Errorf("failed to copy file [%s] to [%s]: %w", file, destPath, err)
			}
			s.stats.LocalEditedFilesCopied++
		} else if s.criteria.MoveFiles {
			// move the file to the new location
			err = s.files.MoveFile(file, destPath)
			if err != nil {
				return fmt.Errorf("failed to move file [%s] to [%s]: %w", file, destPath, err)
			}
			s.stats.LocalEditedFilesMoved++
		} else {
			s.logger.Warn("No file operation specified (neither move nor copy)", zap.String("file", file))
			continue
		}

		s.logger.Debug(logMsg, zap.String("file", file), zap.Bool("backed_up", true))
	}

	s.logger.Info("Backup of local edited files completed", zap.Int("file_count", s.stats.LocalEditedFilesCopied+s.stats.LocalEditedFilesMoved))
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

func generateRawBackupDestinationPath(basePath, fileName string, timestamp time.Time) string {
	return fmt.Sprintf("%s/raw/%d/%02d/%02d/%02d%02d%02d_%s", basePath,
		timestamp.Year(), timestamp.Month(), timestamp.Day(),
		timestamp.Hour(), timestamp.Minute(), timestamp.Second(),
		fileName)
}

func generateEditedBackupDestinationPath(basePath, fileName string, timestamp time.Time) string {
	return fmt.Sprintf("%s/edited/%d/%02d/%02d/%02d%02d%02d_%s", basePath,
		timestamp.Year(), timestamp.Month(), timestamp.Day(),
		timestamp.Hour(), timestamp.Minute(), timestamp.Second(),
		fileName)
}
