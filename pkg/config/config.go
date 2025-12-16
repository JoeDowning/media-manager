package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	logLevel string `env:"log_level"`

	destinationPath string `env:"destination_path, required"`
	sourcePath      string `env:"source_path, required"`

	copyFiles bool `env:"copy_files"`
	moveFiles bool `env:"move_files"`

	importRaw    bool `env:"import_raw"`
	backupRaw    bool `env:"backup_raw"`
	backupEdited bool `env:"backup_edited"`
}

func GetConfig() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get env config: %w", err)
	}
	return cfg, nil
}

func (c Config) LogLevel() string {
	return c.logLevel
}

func (c Config) DestinationPath() string {
	return c.destinationPath
}

func (c Config) SourcePath() string {
	return c.sourcePath
}

func (c Config) CopyFiles() bool {
	return c.copyFiles
}

func (c Config) MoveFiles() bool {
	return c.moveFiles
}

func (c Config) ImportRaw() bool {
	return c.importRaw
}

func (c Config) BackupRaw() bool {
	return c.backupRaw
}

func (c Config) BackupEdited() bool {
	return c.backupEdited
}
