package config

type EnvConfig struct {
	logLevel string `env:"log_level"`

	pathConfig string `env:"path_config, required"`

	fileOperation string `env:"file_op, required"`

	importRaw    bool `env:"import_raw"`
	backupRaw    bool `env:"backup_raw"`
	backupEdited bool `env:"backup_edited"`
	uploadEdited bool `env:"upload_edited"`
}

type Config struct {
	logLevel string

	rawPath    string
	localPath  string
	backupPath string

	copyFiles bool
	moveFiles bool

	importRaw    bool
	backupRaw    bool
	backupEdited bool
	uploadEdited bool
}

type pathConfig struct {
	rawPath    string
	localPath  string
	backupPath string
}
