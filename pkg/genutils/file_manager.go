package genutils

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileManager struct{}

func NewFileManager() *FileManager {
	return &FileManager{}
}

func (fm *FileManager) GetFilesInPath(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(path, entry.Name()))
		}
	}
	return files, nil
}

func (fm *FileManager) GetFilesRecursivelyInPath(path string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk directory %s: %w", path, err)
		}
		if !d.IsDir() {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", path, err)
	}
	return files, nil
}

func (fm *FileManager) DoesFileExist(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check file %s: %w", path, err)
	}
	return !info.IsDir(), nil
}

func (fm *FileManager) DoesPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check path %s: %w", path, err)
	}
	return true, nil
}

func (fm *FileManager) MoveFile(sourcePath, destinationPath string) error {
	err := os.Rename(sourcePath, destinationPath)
	if err != nil {
		return fmt.Errorf("failed to move file from %s to %s: %w", sourcePath, destinationPath, err)
	}
	return nil
}

func (fm *FileManager) CopyFile(sourcePath, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", sourcePath, err)
	}
	defer sourceFile.Close()

	destDir := filepath.Dir(destinationPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", destDir, err)
	}

	destFile, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", destinationPath, err)
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy data from source file %s to destination file %s: %w", sourcePath, destinationPath, err)
	}

	return nil
}
