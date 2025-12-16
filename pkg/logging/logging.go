package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(level string) *zap.Logger {
	var lvl zapcore.Level
	switch level {
	case "debug":
		lvl = zapcore.DebugLevel
	case "info":
		lvl = zapcore.InfoLevel
	case "warn":
		lvl = zapcore.WarnLevel
	case "error":
		lvl = zapcore.ErrorLevel
	default:
		lvl = zapcore.DebugLevel
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // Use colors for levels
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // Human-readable time format
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Use a custom console encoder for pretty printing
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Output to console
	logOutput := zapcore.AddSync(os.Stdout)

	// Create a core to write logs to standard output
	core := zapcore.NewCore(encoder, logOutput, lvl)

	// Create the logger with the core
	return zap.New(core)
}
