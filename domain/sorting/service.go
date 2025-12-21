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

	moveFiles bool
	copyFiles bool
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
	moveFiles, copyFiles bool,
) *Service {
	return &Service{
		logger: logging,
		files:  files,
		criteria: SortCriteria{
			fileTypes:  fileTypes,
			rawPath:    rawPath,
			localPath:  localPath,
			backupPath: backupPath,
			moveFiles:  moveFiles,
			copyFiles:  copyFiles,
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

	var importCount, filesChecked int
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

		// create the new path of format <localPath>/<year>-<month>-<day>/<filename>
		destPath := generateRawImportDestinationPath(s.criteria.localPath, imgData.GetFileName(), imgTime)

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
		importCount++
	}

	s.logger.Info("Import raw files completed", zap.Int("imported_files_count", importCount), zap.Int("file_count", len(files)))
	return newestTime, nil
}

func (s *Service) BackupLocalRawFiles() error {
	files, err := s.files.GetFilesRecursivelyInPath(s.criteria.rawPath)
	if err != nil {
		return fmt.Errorf("failed to get files recursively in path [%s]: %w", s.criteria.rawPath, err)
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

	var backupCount, filesChecked int
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
		destPath := generateRawBackupDestinationPath(s.criteria.backupPath, imgData.GetFileName(), imgTime)

		// check if the file with the new name already exists at the destination
		exists, err := s.files.DoesFileExist(destPath)
		if err != nil {
			return fmt.Errorf("failed to check if file exists at destination [%s]: %w", destPath, err)
		}
		if exists {
			s.logger.Debug(logMsg, zap.String("file", file), zap.Bool("backed_up", false), zap.String("reason", "file already exists at destination"))
			continue
		}

		if s.criteria.copyFiles {
			// copy the file to the new location
			err = s.files.CopyFile(file, destPath)
			if err != nil {
				return fmt.Errorf("failed to copy file [%s] to [%s]: %w", file, destPath, err)
			}
		} else if s.criteria.moveFiles {
			// move the file to the new location
			err = s.files.MoveFile(file, destPath)
			if err != nil {
				return fmt.Errorf("failed to move file [%s] to [%s]: %w", file, destPath, err)
			}
		} else {
			s.logger.Warn("No file operation specified (neither move nor copy)", zap.String("file", file))
			continue
		}

		s.logger.Debug(logMsg, zap.String("file", file), zap.Bool("backed_up", true))
		backupCount++
	}

	s.logger.Info("Backup of local raw files completed", zap.Int("file_count", backupCount))
	return nil
}

func (s *Service) BackupEditedFiles() error {
	files, err := s.files.GetFilesRecursivelyInPath(s.criteria.rawPath)
	if err != nil {
		return fmt.Errorf("failed to get files recursively in path [%s]: %w", s.criteria.rawPath, err)
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

	var backupCount, filesChecked int
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
		destPath := generateEditedBackupDestinationPath(s.criteria.backupPath, imgData.GetFileName(), imgTime)

		// check if the file with the new name already exists at the destination
		exists, err := s.files.DoesFileExist(destPath)
		if err != nil {
			return fmt.Errorf("failed to check if file exists at destination [%s]: %w", destPath, err)
		}
		if exists {
			s.logger.Debug(logMsg, zap.String("file", file), zap.Bool("backed_up", false), zap.String("reason", "file already exists at destination"))
			continue
		}

		if s.criteria.copyFiles {
			// copy the file to the new location
			err = s.files.CopyFile(file, destPath)
			if err != nil {
				return fmt.Errorf("failed to copy file [%s] to [%s]: %w", file, destPath, err)
			}
		} else if s.criteria.moveFiles {
			// move the file to the new location
			err = s.files.MoveFile(file, destPath)
			if err != nil {
				return fmt.Errorf("failed to move file [%s] to [%s]: %w", file, destPath, err)
			}
		} else {
			s.logger.Warn("No file operation specified (neither move nor copy)", zap.String("file", file))
			continue
		}

		s.logger.Debug(logMsg, zap.String("file", file), zap.Bool("backed_up", true))
		backupCount++
	}

	s.logger.Info("Backup of local edited files completed", zap.Int("file_count", backupCount))
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
