// logger_test.go
package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		logLevel string
		expected bool
	}{
		{"debug", true},
		{"info", true},
		{"warn", true},
		{"error", true},
	}

	for _, tt := range tests {
		t.Run(tt.logLevel, func(t *testing.T) {
			err := Initialize(tt.logLevel)

			if tt.expected {
				assert.NoError(t, err, "Expected no error for log level: %s", tt.logLevel)
				assert.NotNil(t, Log, "Expected logger to be initialized for log level: %s", tt.logLevel)
			} else {
				assert.Error(t, err, "Expected error for invalid log level: %s", tt.logLevel)
				assert.Nil(t, Log, "Expected logger to be nil for invalid log level: %s", tt.logLevel)
			}
		})
	}
}

func TestLoggerOutput(t *testing.T) {
	// Initialize logger with a valid level
	err := Initialize("info")
	assert.NoError(t, err)

	// Check if the logger is set up correctly
	assert.NotNil(t, Log)

	// Log a message and check if it outputs correctly (this is more of an integration test)
	Log.Info("This is an info message")
	Log.Debug("This debug message should not appear if the level is info")
}
