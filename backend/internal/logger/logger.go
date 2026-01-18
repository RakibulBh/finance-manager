package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log     *zap.Logger
	HTTPLog *zap.Logger
)

// InitLogger initializes the global logger instances
func InitLogger() error {
	// Ensure logs directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	// 1. Server Logger (Error/Warn only) -> logs/server.log
	serverConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.WarnLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{filepath.Join("logs", "server.log")},
		ErrorOutputPaths: []string{filepath.Join("logs", "server.log")},
		DisableCaller:    true, // Requested: disable caller
		DisableStacktrace: true,
	}

	var err error
	Log, err = serverConfig.Build()
	if err != nil {
		return err
	}

	// 2. HTTP Logger (Info level for requests) -> logs/http.log
	httpConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			// CallerKey:      "caller", // No caller needed for HTTP access logs
			MessageKey:     "msg",
			// StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
		OutputPaths:      []string{filepath.Join("logs", "http.log")},
		ErrorOutputPaths: []string{filepath.Join("logs", "http.log")},
		DisableCaller:    true,
		DisableStacktrace: true,
	}

	HTTPLog, err = httpConfig.Build()
	if err != nil {
		return err
	}

	return nil
}

// Error logs a message at ErrorLevel to server.log
func Error(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Error(msg, fields...)
	}
}

// Warn logs a message at WarnLevel to server.log
func Warn(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Warn(msg, fields...)
	}
}

// InfoHTTP logs a message at InfoLevel to http.log (exposed for middleware)
func InfoHTTP(msg string, fields ...zap.Field) {
	if HTTPLog != nil {
		HTTPLog.Info(msg, fields...)
	}
}

// Sync flushes any buffered log entries
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
	if HTTPLog != nil {
		_ = HTTPLog.Sync()
	}
}
