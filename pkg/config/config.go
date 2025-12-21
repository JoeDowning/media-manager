package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
)

func GetConfig() (Config, error) {
	var envCfg EnvConfig
	err := env.Parse(&envCfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get env config: %w", err)
	}

	cfg := Config{
		logLevel: envCfg.LogLevel,

		importRaw:    envCfg.ImportRaw,
		backupRaw:    envCfg.BackupRaw,
		backupEdited: envCfg.BackupEdited,
		uploadEdited: envCfg.UploadEdited,
	}

	switch envCfg.FileOperation {
	case "copy":
		cfg.copyFiles = true
	case "move":
		cfg.moveFiles = true
	default:
		return Config{}, fmt.Errorf("invalid file operation: %s, choose from [copy, move]", envCfg.FileOperation)
	}

	pathCfg, err := parsePathConfig(envCfg.PathConfig)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse path config: %w", err)
	}

	cfg.rawPath = pathCfg.rawPath
	cfg.localRawPath = pathCfg.localRawPath
	cfg.localEditedPath = pathCfg.localEditedPath
	cfg.backupPath = pathCfg.backupPath

	return cfg, nil
}

func parsePathConfig(pathCfg string) (pathConfig, error) {
	switch pathCfg {
	case "test":
		return pathConfig{
			rawPath:         "/Users/downing/Pictures/testing/raw",
			localRawPath:    "/Users/downing/Pictures/testing/rawsorted",
			localEditedPath: "/Users/downing/Pictures/testing/edited",
			backupPath:      "/Users/downing/Pictures/testing/backup",
		}, nil
	case "default":
		return pathConfig{
			rawPath:         "/Volumes/EOS_DIGITAL/DCIM",
			localRawPath:    "/Users/downing/Pictures/raw",
			localEditedPath: "/Users/downing/Pictures/edited",
			backupPath:      "/",
		}, nil
	default:
		return pathConfig{}, fmt.Errorf("unknown path config: %s, choose from [default]", pathCfg)
	}
}

func (c Config) LogLevel() string {
	return c.logLevel
}

func (c Config) RawPath() string {
	return c.rawPath
}

func (c Config) LocalRawPath() string {
	return c.localRawPath
}

func (c Config) LocalEditedPath() string {
	return c.localEditedPath
}

func (c Config) BackupPath() string {
	return c.backupPath
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

func (c Config) UploadEdited() bool {
	return c.uploadEdited
}

func (c Config) LogConfig(logger *zap.Logger) {
	logger.Info("Config on startup",
		zap.String("log_level", c.LogLevel()),
		zap.String("raw_path", c.RawPath()),
		zap.String("local_raw_path", c.LocalRawPath()),
		zap.String("local_edited_path", c.LocalEditedPath()),
		zap.String("backup_path", c.BackupPath()),
		zap.Bool("copy_files", c.CopyFiles()),
		zap.Bool("move_files", c.MoveFiles()),
		zap.Bool("import_raw", c.ImportRaw()),
		zap.Bool("backup_raw", c.BackupRaw()),
		zap.Bool("backup_edited", c.BackupEdited()),
		zap.Bool("upload_edited", c.UploadEdited()),
	)
}
