package fake

import (
	"bytes"
	"context"
	"testing"
	"time"
)

// TestFakeLogger tests the fake logger functionality.
// zh: TestFakeLogger 測試假日誌記錄器功能。
func TestFakeLogger(t *testing.T) {
	t.Run("NewFakeLogger", func(t *testing.T) {
		logger := NewFakeLogger()
		if logger == nil {
			t.Fatal("NewFakeLogger should not return nil")
		}

		if logger.level != InfoLevel {
			t.Errorf("Expected default level InfoLevel, got %v", logger.level)
		}

		if logger.CountEntries() != 0 {
			t.Errorf("Expected 0 entries initially, got %d", logger.CountEntries())
		}

		if logger.disabled {
			t.Error("Logger should not be disabled initially")
		}

		t.Log("NewFakeLogger test passed")
	})

	t.Run("NewFakeLoggerWithTracing", func(t *testing.T) {
		traceID := "trace-123"
		spanID := "span-456"
		logger := NewFakeLoggerWithTracing(traceID, spanID)

		if logger.traceID != traceID {
			t.Errorf("Expected trace ID '%s', got '%s'", traceID, logger.traceID)
		}

		if logger.spanID != spanID {
			t.Errorf("Expected span ID '%s', got '%s'", spanID, logger.spanID)
		}

		t.Log("NewFakeLoggerWithTracing test passed")
	})

	t.Run("LoggingLevels", func(t *testing.T) {
		logger := NewFakeLogger()
		logger.SetLevel(DebugLevel)

		// Test all logging levels
		logger.Debug("Debug message")
		logger.Info("Info message")
		logger.Warn("Warning message")
		logger.Error("Error message")

		entries := logger.GetEntries()
		if len(entries) != 4 {
			t.Errorf("Expected 4 entries, got %d", len(entries))
		}

		t.Log("LoggingLevels test passed")
	})

	t.Run("LoggingWithFields", func(t *testing.T) {
		logger := NewFakeLogger()

		logger.Info("Test message", "string_field", "value", "int_field", 123, "bool_field", true)

		entries := logger.GetEntries()
		if len(entries) != 1 {
			t.Fatalf("Expected 1 entry, got %d", len(entries))
		}

		entry := entries[0]
		if entry.Message != "Test message" {
			t.Errorf("Expected message 'Test message', got '%s'", entry.Message)
		}

		if entry.Fields["string_field"] != "value" {
			t.Errorf("Expected string_field 'value', got '%v'", entry.Fields["string_field"])
		}

		if entry.Fields["int_field"] != 123 {
			t.Errorf("Expected int_field 123, got %v", entry.Fields["int_field"])
		}

		if entry.Fields["bool_field"] != true {
			t.Errorf("Expected bool_field true, got %v", entry.Fields["bool_field"])
		}

		t.Log("LoggingWithFields test passed")
	})

	t.Run("TracingInformation", func(t *testing.T) {
		traceID := "trace-789"
		spanID := "span-012"
		logger := NewFakeLoggerWithTracing(traceID, spanID)

		logger.Info("Message with tracing")

		entries := logger.GetEntries()
		if len(entries) != 1 {
			t.Fatalf("Expected 1 entry, got %d", len(entries))
		}

		entry := entries[0]
		if entry.TraceID != traceID {
			t.Errorf("Expected trace ID '%s', got '%s'", traceID, entry.TraceID)
		}

		if entry.SpanID != spanID {
			t.Errorf("Expected span ID '%s', got '%s'", spanID, entry.SpanID)
		}

		// Test entries with trace
		tracedEntries := logger.GetEntriesWithTrace()
		if len(tracedEntries) != 1 {
			t.Errorf("Expected 1 traced entry, got %d", len(tracedEntries))
		}

		t.Log("TracingInformation test passed")
	})

	t.Run("LevelFiltering", func(t *testing.T) {
		logger := NewFakeLogger()
		logger.SetLevel(WarnLevel) // Only warn and above

		logger.Debug("Debug message") // Should be filtered out
		logger.Info("Info message")   // Should be filtered out
		logger.Warn("Warning message")
		logger.Error("Error message")

		entries := logger.GetEntries()
		if len(entries) != 2 {
			t.Errorf("Expected 2 entries (warn and error), got %d", len(entries))
		}

		// Check only warn and error entries exist
		for _, entry := range entries {
			if entry.Level != WarnLevel && entry.Level != ErrorLevel {
				t.Errorf("Unexpected entry level: %v", entry.Level)
			}
		}

		t.Log("LevelFiltering test passed")
	})

	t.Run("GetEntriesByLevel", func(t *testing.T) {
		logger := NewFakeLogger()
		logger.SetLevel(DebugLevel)

		// Add multiple entries of different levels
		logger.Info("Info 1")
		logger.Error("Error 1")
		logger.Info("Info 2")
		logger.Warn("Warning 1")
		logger.Error("Error 2")

		// Test filtering by level
		infoEntries := logger.GetEntriesByLevel(InfoLevel)
		if len(infoEntries) != 2 {
			t.Errorf("Expected 2 info entries, got %d", len(infoEntries))
		}

		errorEntries := logger.GetEntriesByLevel(ErrorLevel)
		if len(errorEntries) != 2 {
			t.Errorf("Expected 2 error entries, got %d", len(errorEntries))
		}

		warnEntries := logger.GetEntriesByLevel(WarnLevel)
		if len(warnEntries) != 1 {
			t.Errorf("Expected 1 warn entry, got %d", len(warnEntries))
		}

		t.Log("GetEntriesByLevel test passed")
	})

	t.Run("CountMethods", func(t *testing.T) {
		logger := NewFakeLogger()

		logger.Info("Info message")
		logger.Error("Error message")
		logger.Info("Another info message")

		if logger.CountEntries() != 3 {
			t.Errorf("Expected 3 total entries, got %d", logger.CountEntries())
		}

		if logger.CountEntriesByLevel(InfoLevel) != 2 {
			t.Errorf("Expected 2 info entries, got %d", logger.CountEntriesByLevel(InfoLevel))
		}

		if logger.CountEntriesByLevel(ErrorLevel) != 1 {
			t.Errorf("Expected 1 error entry, got %d", logger.CountEntriesByLevel(ErrorLevel))
		}

		t.Log("CountMethods test passed")
	})

	t.Run("HasEntryMethods", func(t *testing.T) {
		logger := NewFakeLogger()

		logger.Info("Test message", "user_id", 123, "action", "login")

		if !logger.HasEntry("Test message") {
			t.Error("Should have entry with message 'Test message'")
		}

		if logger.HasEntry("Non-existent message") {
			t.Error("Should not have entry with non-existent message")
		}

		if !logger.HasEntryWithField("user_id", 123) {
			t.Error("Should have entry with field user_id=123")
		}

		if !logger.HasEntryWithField("action", "login") {
			t.Error("Should have entry with field action=login")
		}

		if logger.HasEntryWithField("user_id", 456) {
			t.Error("Should not have entry with field user_id=456")
		}

		t.Log("HasEntryMethods test passed")
	})

	t.Run("DisableEnable", func(t *testing.T) {
		logger := NewFakeLogger()

		// Log when enabled
		logger.Info("Message 1")

		// Disable and log
		logger.Disable()
		logger.Info("Message 2") // Should not be recorded

		// Enable and log
		logger.Enable()
		logger.Info("Message 3")

		entries := logger.GetEntries()
		if len(entries) != 2 {
			t.Errorf("Expected 2 entries, got %d", len(entries))
		}

		if entries[0].Message != "Message 1" {
			t.Errorf("Expected first message 'Message 1', got '%s'", entries[0].Message)
		}

		if entries[1].Message != "Message 3" {
			t.Errorf("Expected second message 'Message 3', got '%s'", entries[1].Message)
		}

		t.Log("DisableEnable test passed")
	})

	t.Run("OutputWriter", func(t *testing.T) {
		var buffer bytes.Buffer
		logger := NewFakeLogger()
		logger.SetOutput(&buffer)

		logger.Info("Test output message", "key", "value")

		output := buffer.String()
		if output == "" {
			t.Error("Expected output to be written to buffer")
		}

		// Check that output contains expected elements
		if !bytes.Contains(buffer.Bytes(), []byte("INFO:")) {
			t.Error("Output should contain log level")
		}

		if !bytes.Contains(buffer.Bytes(), []byte("Test output message")) {
			t.Error("Output should contain log message")
		}

		if !bytes.Contains(buffer.Bytes(), []byte("key=value")) {
			t.Error("Output should contain log fields")
		}

		t.Log("OutputWriter test passed")
	})

	t.Run("WithTracing", func(t *testing.T) {
		logger := NewFakeLogger()
		tracedLogger := logger.WithTracing("new-trace", "new-span")

		// Original logger should not have tracing
		logger.Info("Original message")
		originalEntries := logger.GetEntries()
		if len(originalEntries) != 1 {
			t.Fatalf("Expected 1 entry in original logger, got %d", len(originalEntries))
		}
		if originalEntries[0].TraceID != "" {
			t.Error("Original logger should not have trace ID")
		}

		// Traced logger should have tracing
		tracedLogger.Info("Traced message")
		tracedEntries := tracedLogger.GetEntries()
		if len(tracedEntries) != 1 {
			t.Fatalf("Expected 1 entry in traced logger, got %d", len(tracedEntries))
		}
		if tracedEntries[0].TraceID != "new-trace" {
			t.Errorf("Expected trace ID 'new-trace', got '%s'", tracedEntries[0].TraceID)
		}
		if tracedEntries[0].SpanID != "new-span" {
			t.Errorf("Expected span ID 'new-span', got '%s'", tracedEntries[0].SpanID)
		}

		t.Log("WithTracing test passed")
	})

	t.Run("Clear", func(t *testing.T) {
		logger := NewFakeLogger()

		logger.Info("Message 1")
		logger.Error("Message 2")

		if logger.CountEntries() != 2 {
			t.Errorf("Expected 2 entries before clear, got %d", logger.CountEntries())
		}

		logger.Clear()

		if logger.CountEntries() != 0 {
			t.Errorf("Expected 0 entries after clear, got %d", logger.CountEntries())
		}

		t.Log("Clear test passed")
	})

	t.Run("Reset", func(t *testing.T) {
		logger := NewFakeLoggerWithTracing("trace", "span")
		logger.SetLevel(ErrorLevel)
		logger.Disable()
		logger.Info("Test message")

		// Reset logger
		logger.Reset()

		// Check all properties are reset
		if logger.level != InfoLevel {
			t.Errorf("Expected level InfoLevel after reset, got %v", logger.level)
		}

		if logger.traceID != "" {
			t.Errorf("Expected empty trace ID after reset, got '%s'", logger.traceID)
		}

		if logger.spanID != "" {
			t.Errorf("Expected empty span ID after reset, got '%s'", logger.spanID)
		}

		if logger.disabled {
			t.Error("Logger should not be disabled after reset")
		}

		if logger.CountEntries() != 0 {
			t.Errorf("Expected 0 entries after reset, got %d", logger.CountEntries())
		}

		t.Log("Reset test passed")
	})

	t.Run("GetStats", func(t *testing.T) {
		logger := NewFakeLoggerWithTracing("trace", "span")
		logger.SetLevel(DebugLevel) // Enable debug logging

		// Add various entries with tracing
		logger.Debug("Debug message")
		logger.Info("Info message")
		logger.Info("Another info message")
		logger.Warn("Warning message")
		logger.Error("Error message")

		// Add entry without tracing by temporarily clearing trace info
		originalTraceID := logger.traceID
		originalSpanID := logger.spanID
		logger.SetTracing("", "")
		logger.Info("Plain message")
		logger.SetTracing(originalTraceID, originalSpanID)

		// Add final traced message
		logger.Info("Final traced message")

		stats := logger.GetStats()

		if stats["total"] != 7 {
			t.Errorf("Expected 7 total entries, got %d", stats["total"])
		}

		if stats["info"] != 4 {
			t.Errorf("Expected 4 info entries, got %d", stats["info"])
		}

		if stats["with_trace"] != 6 {
			t.Errorf("Expected 6 entries with trace, got %d", stats["with_trace"])
		}

		t.Log("GetStats test passed")
	})
}

// TestFakeLoggerFromContext tests creating logger from context.
// zh: TestFakeLoggerFromContext 測試從上下文建立日誌記錄器。
func TestFakeLoggerFromContext(t *testing.T) {
	t.Run("ContextWithTracing", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "trace_id", "ctx-trace-123")
		ctx = context.WithValue(ctx, "span_id", "ctx-span-456")

		logger := FakeLoggerFromContext(ctx)

		if logger.traceID != "ctx-trace-123" {
			t.Errorf("Expected trace ID 'ctx-trace-123', got '%s'", logger.traceID)
		}

		if logger.spanID != "ctx-span-456" {
			t.Errorf("Expected span ID 'ctx-span-456', got '%s'", logger.spanID)
		}

		// Test logging with context tracing
		logger.Info("Context traced message")

		entries := logger.GetEntries()
		if len(entries) != 1 {
			t.Fatalf("Expected 1 entry, got %d", len(entries))
		}

		entry := entries[0]
		if entry.TraceID != "ctx-trace-123" {
			t.Errorf("Expected entry trace ID 'ctx-trace-123', got '%s'", entry.TraceID)
		}

		if entry.SpanID != "ctx-span-456" {
			t.Errorf("Expected entry span ID 'ctx-span-456', got '%s'", entry.SpanID)
		}

		t.Log("ContextWithTracing test passed")
	})

	t.Run("ContextWithoutTracing", func(t *testing.T) {
		ctx := context.Background()
		logger := FakeLoggerFromContext(ctx)

		if logger.traceID != "" {
			t.Errorf("Expected empty trace ID, got '%s'", logger.traceID)
		}

		if logger.spanID != "" {
			t.Errorf("Expected empty span ID, got '%s'", logger.spanID)
		}

		t.Log("ContextWithoutTracing test passed")
	})

	t.Run("ContextWithPartialTracing", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "trace_id", "partial-trace")
		// No span_id in context

		logger := FakeLoggerFromContext(ctx)

		if logger.traceID != "partial-trace" {
			t.Errorf("Expected trace ID 'partial-trace', got '%s'", logger.traceID)
		}

		if logger.spanID != "" {
			t.Errorf("Expected empty span ID, got '%s'", logger.spanID)
		}

		t.Log("ContextWithPartialTracing test passed")
	})
}

// TestLogEntry tests the LogEntry structure.
// zh: TestLogEntry 測試 LogEntry 結構。
func TestLogEntry(t *testing.T) {
	t.Run("LogEntryFields", func(t *testing.T) {
		logger := NewFakeLoggerWithTracing("test-trace", "test-span")

		before := time.Now()
		logger.Info("Test entry", "field1", "value1", "field2", 42)
		after := time.Now()

		entries := logger.GetEntries()
		if len(entries) != 1 {
			t.Fatalf("Expected 1 entry, got %d", len(entries))
		}

		entry := entries[0]

		// Check timestamp is reasonable
		if entry.Timestamp.Before(before) || entry.Timestamp.After(after) {
			t.Error("Entry timestamp is not within expected range")
		}

		// Check level
		if entry.Level != InfoLevel {
			t.Errorf("Expected level InfoLevel, got %v", entry.Level)
		}

		// Check message
		if entry.Message != "Test entry" {
			t.Errorf("Expected message 'Test entry', got '%s'", entry.Message)
		}

		// Check fields
		if len(entry.Fields) != 2 {
			t.Errorf("Expected 2 fields, got %d", len(entry.Fields))
		}

		if entry.Fields["field1"] != "value1" {
			t.Errorf("Expected field1 'value1', got '%v'", entry.Fields["field1"])
		}

		if entry.Fields["field2"] != 42 {
			t.Errorf("Expected field2 42, got %v", entry.Fields["field2"])
		}

		// Check tracing
		if entry.TraceID != "test-trace" {
			t.Errorf("Expected trace ID 'test-trace', got '%s'", entry.TraceID)
		}

		if entry.SpanID != "test-span" {
			t.Errorf("Expected span ID 'test-span', got '%s'", entry.SpanID)
		}

		t.Log("LogEntryFields test passed")
	})
}
