package config

import (
	"os"
	"path/filepath"
	"time"
)

const lastImportFilename = "last_import.txt"

func getLastImportPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".media-manager", lastImportFilename), nil
}

func GetLastImportDate() (time.Time, error) {
	filePath, err := getLastImportPath()
	if err != nil {
		return time.Time{}, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), nil
		}
		return time.Time{}, err
	}

	timestamp, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		return time.Time{}, err
	}

	return timestamp, nil
}

func SetLastImportDate(t time.Time) error {
	filePath, err := getLastImportPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data := t.Format(time.RFC3339)
	return os.WriteFile(filePath, []byte(data), 0644)
}
