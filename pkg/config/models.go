package config

type EnvConfig struct {
	LogLevel string `env:"log_level"`

	PathConfig string `env:"path_config"`

	FileOperation string `env:"file_op"`

	ImportRaw    bool `env:"import_raw"`
	BackupRaw    bool `env:"backup_raw"`
	BackupEdited bool `env:"backup_edited"`
	UploadEdited bool `env:"upload_edited"`
}

type Config struct {
	logLevel string

	rawPath         string
	localRawPath    string
	localEditedPath string
	backupPath      string

	copyFiles bool
	moveFiles bool

	importRaw    bool
	backupRaw    bool
	backupEdited bool
	uploadEdited bool
}

type pathConfig struct {
	rawPath         string
	localRawPath    string
	localEditedPath string
	backupPath      string
}
