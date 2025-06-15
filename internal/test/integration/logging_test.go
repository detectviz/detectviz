package integration

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"detectviz/pkg/shared/log"
)

// TestLoggerInitialization tests logger initialization with various configurations.
// zh: TestLoggerInitialization 測試各種配置下的日誌記錄器初始化。
func TestLoggerInitialization(t *testing.T) {
	t.Run("DefaultConfiguration", func(t *testing.T) {
		// Test default configuration
		config := log.DefaultLoggerConfig()
		logger, err := log.NewLogger(config)
		if err != nil {
			t.Fatalf("Failed to create logger with default config: %v", err)
		}

		if logger == nil {
			t.Fatal("Logger should not be nil")
		}

		logger.Info("Test message with default configuration")
		t.Log("Default configuration test passed")
	})

	t.Run("ConsoleConfiguration", func(t *testing.T) {
		// Test console configuration
		config := &log.LoggerConfig{
			Type:   "console",
			Level:  "debug",
			Format: "text",
			Output: "stdout",
		}

		logger, err := log.NewLogger(config)
		if err != nil {
			t.Fatalf("Failed to create console logger: %v", err)
		}

		logger.Debug("Debug message for console")
		logger.Info("Info message for console")
		logger.Warn("Warning message for console")
		logger.Error("Error message for console")
		t.Log("Console configuration test passed")
	})

	t.Run("FileConfiguration", func(t *testing.T) {
		// Create temporary directory for log files
		tempDir := t.TempDir()
		logFile := filepath.Join(tempDir, "test.log")

		config := &log.LoggerConfig{
			Type:   "file",
			Level:  "info",
			Format: "text",
			FileConfig: &log.FileConfig{
				Filename:   logFile,
				MaxSize:    1, // 1MB
				MaxBackups: 3,
				MaxAge:     7,
				Compress:   false,
			},
		}

		logger, err := log.NewLogger(config)
		if err != nil {
			t.Fatalf("Failed to create file logger: %v", err)
		}

		logger.Info("Test message to file")
		logger.Warn("Warning message to file")

		// Give some time for file write
		time.Sleep(100 * time.Millisecond)

		// Check if log file was created and contains expected content
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Fatalf("Log file %s was not created", logFile)
		}

		content, err := os.ReadFile(logFile)
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}

		logContent := string(content)
		if !strings.Contains(logContent, "Test message to file") {
			t.Errorf("Log file does not contain expected message. Content: %s", logContent)
		}

		if !strings.Contains(logContent, "INFO") {
			t.Errorf("Log file does not contain INFO level. Content: %s", logContent)
		}

		t.Log("File configuration test passed")
	})

	t.Run("BothConfiguration", func(t *testing.T) {
		// Create temporary directory for log files
		tempDir := t.TempDir()
		logFile := filepath.Join(tempDir, "both.log")

		config := &log.LoggerConfig{
			Type:   "both",
			Level:  "debug",
			Format: "text",
			Output: "stdout",
			FileConfig: &log.FileConfig{
				Filename:   logFile,
				MaxSize:    1,
				MaxBackups: 2,
				MaxAge:     1,
				Compress:   true,
			},
		}

		logger, err := log.NewLogger(config)
		if err != nil {
			t.Fatalf("Failed to create both logger: %v", err)
		}

		logger.Info("Message to both console and file")
		logger.Debug("Debug message to both outputs")

		// Give some time for file write
		time.Sleep(100 * time.Millisecond)

		// Check if log file was created
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Fatalf("Log file %s was not created", logFile)
		}

		content, err := os.ReadFile(logFile)
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}

		if !strings.Contains(string(content), "Message to both console and file") {
			t.Errorf("Log file does not contain expected message")
		}

		t.Log("Both configuration test passed")
	})
}

// TestGlobalLogger tests global logger functionality.
// zh: TestGlobalLogger 測試全域日誌記錄器功能。
func TestGlobalLogger(t *testing.T) {
	t.Run("GlobalLoggerUsage", func(t *testing.T) {
		// Test global logger functions
		log.Info("Global info message")
		log.Debug("Global debug message")
		log.Warn("Global warning message")
		log.Error("Global error message")

		// Test Printf compatibility
		log.Printf("Printf style message: %s, %d", "test", 42)

		t.Log("Global logger usage test passed")
	})

	t.Run("SetCustomGlobalLogger", func(t *testing.T) {
		// Create custom logger
		config := &log.LoggerConfig{
			Type:   "console",
			Level:  "debug",
			Format: "text",
			Output: "stderr",
		}

		customLogger, err := log.NewLogger(config)
		if err != nil {
			t.Fatalf("Failed to create custom logger: %v", err)
		}

		// Set as global logger
		log.SetGlobalLogger(customLogger)

		// Test global functions with custom logger
		log.Info("Message with custom global logger")

		// Verify we can get the global logger
		globalLogger := log.GetGlobalLogger()
		if globalLogger == nil {
			t.Fatal("Global logger should not be nil")
		}

		globalLogger.Info("Direct global logger usage")
		t.Log("Custom global logger test passed")
	})
}

// TestContextLogger tests context-aware logging.
// zh: TestContextLogger 測試上下文感知的日誌記錄。
func TestContextLogger(t *testing.T) {
	t.Run("ContextLogger", func(t *testing.T) {
		ctx := context.Background()

		// Test context logger (currently just returns global logger)
		logger := log.L(ctx)
		if logger == nil {
			t.Fatal("Context logger should not be nil")
		}

		logger.Info("Message from context logger")
		t.Log("Context logger test passed")
	})

	t.Run("NilContext", func(t *testing.T) {
		// Test with nil context
		logger := log.L(nil)
		if logger == nil {
			t.Fatal("Logger should not be nil even with nil context")
		}

		logger.Info("Message with nil context")
		t.Log("Nil context test passed")
	})
}

// TestLogLevels tests different log levels.
// zh: TestLogLevels 測試不同的日誌等級。
func TestLogLevels(t *testing.T) {
	t.Run("LogLevelFiltering", func(t *testing.T) {
		// Create temporary directory for log files
		tempDir := t.TempDir()
		logFile := filepath.Join(tempDir, "level_test.log")

		// Create logger with WARN level
		config := &log.LoggerConfig{
			Type:   "file",
			Level:  "warn",
			Format: "text",
			FileConfig: &log.FileConfig{
				Filename: logFile,
			},
		}

		logger, err := log.NewLogger(config)
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}

		// Log messages at different levels
		logger.Debug("This debug message should not appear")
		logger.Info("This info message should not appear")
		logger.Warn("This warning message should appear")
		logger.Error("This error message should appear")

		// Give some time for file write
		time.Sleep(100 * time.Millisecond)

		content, err := os.ReadFile(logFile)
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}

		logContent := string(content)

		// Check that only WARN and ERROR messages appear
		if strings.Contains(logContent, "debug message") {
			t.Errorf("Debug message should not appear in WARN level log")
		}
		if strings.Contains(logContent, "info message") {
			t.Errorf("Info message should not appear in WARN level log")
		}
		if !strings.Contains(logContent, "warning message") {
			t.Errorf("Warning message should appear in WARN level log")
		}
		if !strings.Contains(logContent, "error message") {
			t.Errorf("Error message should appear in WARN level log")
		}

		t.Log("Log level filtering test passed")
	})

	t.Run("LogLevelParsing", func(t *testing.T) {
		// Test log level parsing
		testCases := []struct {
			input    string
			expected log.LogLevel
		}{
			{"debug", log.DebugLevel},
			{"DEBUG", log.DebugLevel},
			{"info", log.InfoLevel},
			{"INFO", log.InfoLevel},
			{"warn", log.WarnLevel},
			{"WARN", log.WarnLevel},
			{"warning", log.WarnLevel},
			{"WARNING", log.WarnLevel},
			{"error", log.ErrorLevel},
			{"ERROR", log.ErrorLevel},
			{"fatal", log.FatalLevel},
			{"FATAL", log.FatalLevel},
			{"unknown", log.InfoLevel}, // Default to info
		}

		for _, tc := range testCases {
			result := log.ParseLogLevel(tc.input)
			if result != tc.expected {
				t.Errorf("ParseLogLevel(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		}

		t.Log("Log level parsing test passed")
	})
}

// TestLoggerConfiguration tests various logger configurations.
// zh: TestLoggerConfiguration 測試各種日誌記錄器配置。
func TestLoggerConfiguration(t *testing.T) {
	t.Run("InvalidConfiguration", func(t *testing.T) {
		// Test with invalid configuration
		config := &log.LoggerConfig{
			Type: "file",
			// Missing FileConfig should cause error
		}

		_, err := log.NewLogger(config)
		if err == nil {
			t.Fatal("Expected error for file type without FileConfig")
		}

		t.Logf("Got expected error: %v", err)
	})

	t.Run("DirectoryCreation", func(t *testing.T) {
		// Test automatic directory creation
		tempDir := t.TempDir()
		logDir := filepath.Join(tempDir, "logs", "subdir")
		logFile := filepath.Join(logDir, "test.log")

		config := &log.LoggerConfig{
			Type:   "file",
			Level:  "info",
			Format: "text",
			FileConfig: &log.FileConfig{
				Filename: logFile,
			},
		}

		logger, err := log.NewLogger(config)
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}

		logger.Info("Test message in nested directory")

		// Give some time for file write
		time.Sleep(100 * time.Millisecond)

		// Check if directory and file were created
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			t.Fatalf("Log directory %s was not created", logDir)
		}

		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Fatalf("Log file %s was not created", logFile)
		}

		t.Log("Directory creation test passed")
	})
}

// TestLoggerConcurrency tests concurrent logger usage.
// zh: TestLoggerConcurrency 測試並發日誌記錄器使用。
func TestLoggerConcurrency(t *testing.T) {
	t.Run("ConcurrentLogging", func(t *testing.T) {
		// Create temporary directory for log files
		tempDir := t.TempDir()
		logFile := filepath.Join(tempDir, "concurrent.log")

		config := &log.LoggerConfig{
			Type:   "file",
			Level:  "info",
			Format: "text",
			FileConfig: &log.FileConfig{
				Filename: logFile,
			},
		}

		logger, err := log.NewLogger(config)
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}

		// Launch multiple goroutines for concurrent logging
		const numGoroutines = 10
		const messagesPerGoroutine = 10
		done := make(chan bool, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				for j := 0; j < messagesPerGoroutine; j++ {
					logger.Info("Concurrent message", "goroutine", id, "message", j)
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		// Give some time for all writes to complete
		time.Sleep(200 * time.Millisecond)

		// Check if log file contains expected number of messages
		content, err := os.ReadFile(logFile)
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}

		logContent := string(content)
		messageCount := strings.Count(logContent, "Concurrent message")
		expectedCount := numGoroutines * messagesPerGoroutine

		if messageCount != expectedCount {
			t.Errorf("Expected %d messages, got %d", expectedCount, messageCount)
		}

		t.Log("Concurrent logging test passed")
	})
}

// TestLoggerCleanup tests logger cleanup functionality.
// zh: TestLoggerCleanup 測試日誌記錄器清理功能。
func TestLoggerCleanup(t *testing.T) {
	t.Run("SyncAndClose", func(t *testing.T) {
		// Test sync and close functions
		err := log.Sync()
		if err != nil {
			t.Errorf("Sync should not return error: %v", err)
		}

		err = log.Close()
		if err != nil {
			t.Errorf("Close should not return error: %v", err)
		}

		t.Log("Sync and close test passed")
	})
}
