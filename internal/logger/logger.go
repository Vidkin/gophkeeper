// Package logger provides a simple logging utility using the Zap logging library.
//
// This package allows for the initialization of a logger with a specified log level.
// It uses the Zap library to provide structured and performant logging capabilities.
package logger

import (
	"go.uber.org/zap"
)

// Log is the global logger instance that can be used throughout the application.
var Log *zap.Logger = zap.NewNop()

// Initialize sets up the logger with the specified log level.
//
// Parameters:
//   - logLevel: A string representing the desired log level (e.g., "debug", "info", "warn", "error").
//
// Returns:
//   - An error if the log level is invalid or if there is an issue during logger initialization.
//   - nil if the logger is successfully initialized.
//
// The function parses the provided log level and configures the logger accordingly. It creates a
// production logger with the specified log level. If the logger is successfully built, it replaces
// the global Log variable with the new logger instance. The logger is set to sync its output before
// the function returns, ensuring that all log entries are flushed.
func Initialize(logLevel string) error {
	lvl, err := zap.ParseAtomicLevel(logLevel)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	logger, err := cfg.Build()
	if err != nil {
		return err
	}
	defer logger.Sync()

	Log = logger
	return nil
}
