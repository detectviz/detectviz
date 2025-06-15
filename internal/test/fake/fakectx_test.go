package fake

import (
	"context"
	"testing"
	"time"
)

// TestFakeContext tests the fake context functionality.
func TestFakeContext(t *testing.T) {
	t.Run("NewFakeContext", func(t *testing.T) {
		ctx := NewFakeContext()
		if ctx == nil {
			t.Fatal("NewFakeContext should not return nil")
		}

		if ctx.parent == nil {
			t.Error("Parent should be set to background context")
		}

		if ctx.IsCanceled() {
			t.Error("Context should not be canceled initially")
		}

		t.Log("NewFakeContext test passed")
	})

	t.Run("ContextWithTracing", func(t *testing.T) {
		traceID := "test-trace-123"
		spanID := "test-span-456"
		ctx := NewFakeContextWithTracing(traceID, spanID)

		if ctx.GetTraceID() != traceID {
			t.Errorf("Expected trace ID '%s', got '%s'", traceID, ctx.GetTraceID())
		}

		if ctx.GetSpanID() != spanID {
			t.Errorf("Expected span ID '%s', got '%s'", spanID, ctx.GetSpanID())
		}

		// Test Value method
		if ctx.Value("trace_id") != traceID {
			t.Errorf("Expected trace_id value '%s', got '%v'", traceID, ctx.Value("trace_id"))
		}

		if ctx.Value("span_id") != spanID {
			t.Errorf("Expected span_id value '%s', got '%v'", spanID, ctx.Value("span_id"))
		}

		t.Log("ContextWithTracing test passed")
	})

	t.Run("SetValue", func(t *testing.T) {
		ctx := NewFakeContext()

		ctx.SetValue("key1", "value1")
		ctx.SetValue("key2", 42)
		ctx.SetValue("key3", true)

		if ctx.Value("key1") != "value1" {
			t.Errorf("Expected key1 'value1', got '%v'", ctx.Value("key1"))
		}

		if ctx.Value("key2") != 42 {
			t.Errorf("Expected key2 42, got %v", ctx.Value("key2"))
		}

		if ctx.Value("key3") != true {
			t.Errorf("Expected key3 true, got %v", ctx.Value("key3"))
		}

		if ctx.Value("nonexistent") != nil {
			t.Errorf("Expected nil for nonexistent key, got %v", ctx.Value("nonexistent"))
		}

		t.Log("SetValue test passed")
	})

	t.Run("HasValue", func(t *testing.T) {
		ctx := NewFakeContext()
		ctx.SetValue("existing_key", "value")

		if !ctx.HasValue("existing_key") {
			t.Error("Should have existing_key")
		}

		if ctx.HasValue("nonexistent_key") {
			t.Error("Should not have nonexistent_key")
		}

		t.Log("HasValue test passed")
	})

	t.Run("WithValue", func(t *testing.T) {
		ctx := NewFakeContext()
		ctx.SetValue("original", "value")

		newCtx := ctx.WithValue("new_key", "new_value")

		// Original context should not have new key
		if ctx.Value("new_key") != nil {
			t.Error("Original context should not have new_key")
		}

		// New context should have both keys
		if newCtx.Value("original") != "value" {
			t.Error("New context should have original key")
		}

		if newCtx.Value("new_key") != "new_value" {
			t.Error("New context should have new_key")
		}

		t.Log("WithValue test passed")
	})

	t.Run("WithTracing", func(t *testing.T) {
		ctx := NewFakeContext()
		tracedCtx := ctx.WithTracing("new-trace", "new-span")

		// Original context should not have tracing
		if ctx.GetTraceID() != "" {
			t.Error("Original context should not have trace ID")
		}

		// New context should have tracing
		if tracedCtx.GetTraceID() != "new-trace" {
			t.Errorf("Expected trace ID 'new-trace', got '%s'", tracedCtx.GetTraceID())
		}

		if tracedCtx.GetSpanID() != "new-span" {
			t.Errorf("Expected span ID 'new-span', got '%s'", tracedCtx.GetSpanID())
		}

		t.Log("WithTracing test passed")
	})

	t.Run("Cancel", func(t *testing.T) {
		ctx := NewFakeContext()

		// Check initial state
		if ctx.IsCanceled() {
			t.Error("Context should not be canceled initially")
		}

		if ctx.Err() != nil {
			t.Error("Context should not have error initially")
		}

		// Cancel context
		testErr := context.Canceled
		ctx.Cancel(testErr)

		// Check canceled state
		if !ctx.IsCanceled() {
			t.Error("Context should be canceled after Cancel()")
		}

		if ctx.Err() != testErr {
			t.Errorf("Expected error %v, got %v", testErr, ctx.Err())
		}

		// Check Done channel
		select {
		case <-ctx.Done():
			// Expected
		default:
			t.Error("Done channel should be closed after cancel")
		}

		t.Log("Cancel test passed")
	})

	t.Run("WithTimeout", func(t *testing.T) {
		timeout := 50 * time.Millisecond
		ctx := NewFakeContextWithTimeout(timeout)

		// Check deadline
		deadline, ok := ctx.Deadline()
		if !ok {
			t.Error("Context should have deadline")
		}

		expectedDeadline := time.Now().Add(timeout)
		if deadline.Before(expectedDeadline.Add(-10*time.Millisecond)) || deadline.After(expectedDeadline.Add(10*time.Millisecond)) {
			t.Error("Deadline is not within expected range")
		}

		// Wait for timeout
		select {
		case <-ctx.Done():
			// Expected
		case <-time.After(100 * time.Millisecond):
			t.Error("Context should have timed out")
		}

		// Check error is deadline exceeded
		if ctx.Err() != context.DeadlineExceeded {
			t.Errorf("Expected DeadlineExceeded error, got %v", ctx.Err())
		}

		t.Log("WithTimeout test passed")
	})

	t.Run("Clone", func(t *testing.T) {
		ctx := NewFakeContextWithTracing("trace", "span")
		ctx.SetValue("key", "value")

		cloned := ctx.Clone()

		// Check values are copied
		if cloned.GetTraceID() != "trace" {
			t.Error("Cloned context should have same trace ID")
		}

		if cloned.Value("key") != "value" {
			t.Error("Cloned context should have same values")
		}

		// Check independence
		cloned.SetValue("new_key", "new_value")
		if ctx.Value("new_key") != nil {
			t.Error("Original context should not be affected by changes to clone")
		}

		t.Log("Clone test passed")
	})

	t.Run("GetAllValues", func(t *testing.T) {
		ctx := NewFakeContext()
		ctx.SetValue("key1", "value1")
		ctx.SetValue("key2", 42)
		ctx.SetValue("key3", true)

		values := ctx.GetAllValues()

		if len(values) != 3 {
			t.Errorf("Expected 3 values, got %d", len(values))
		}

		if values["key1"] != "value1" {
			t.Error("Values should contain key1")
		}

		if values["key2"] != 42 {
			t.Error("Values should contain key2")
		}

		if values["key3"] != true {
			t.Error("Values should contain key3")
		}

		t.Log("GetAllValues test passed")
	})

	t.Run("Reset", func(t *testing.T) {
		ctx := NewFakeContextWithTracing("trace", "span")
		ctx.SetValue("key", "value")
		ctx.Cancel(context.Canceled)

		// Check initial state
		if !ctx.IsCanceled() {
			t.Error("Context should be canceled before reset")
		}

		// Reset context
		ctx.Reset()

		// Check reset state
		if ctx.IsCanceled() {
			t.Error("Context should not be canceled after reset")
		}

		if ctx.GetTraceID() != "" {
			t.Error("Context should not have trace ID after reset")
		}

		if ctx.Value("key") != nil {
			t.Error("Context should not have values after reset")
		}

		if ctx.Err() != nil {
			t.Error("Context should not have error after reset")
		}

		t.Log("Reset test passed")
	})
}

// TestFakeContextBuilder tests the context builder functionality.
func TestFakeContextBuilder(t *testing.T) {
	t.Run("BuilderBasic", func(t *testing.T) {
		builder := NewFakeContextBuilder()
		ctx := builder.
			WithValue("key1", "value1").
			WithValue("key2", 42).
			WithTracing("trace-123", "span-456").
			Build()

		if ctx.Value("key1") != "value1" {
			t.Error("Built context should have key1")
		}

		if ctx.Value("key2") != 42 {
			t.Error("Built context should have key2")
		}

		if ctx.GetTraceID() != "trace-123" {
			t.Error("Built context should have trace ID")
		}

		if ctx.GetSpanID() != "span-456" {
			t.Error("Built context should have span ID")
		}

		t.Log("BuilderBasic test passed")
	})

	t.Run("BuilderWithTimeout", func(t *testing.T) {
		timeout := 50 * time.Millisecond
		builder := NewFakeContextBuilder()
		ctx := builder.
			WithTimeout(timeout).
			WithValue("test", "value").
			Build()

		// Check deadline
		deadline, ok := ctx.Deadline()
		if !ok {
			t.Error("Built context should have deadline")
		}

		expectedDeadline := time.Now().Add(timeout)
		if deadline.Before(expectedDeadline.Add(-10*time.Millisecond)) || deadline.After(expectedDeadline.Add(10*time.Millisecond)) {
			t.Error("Deadline is not within expected range")
		}

		// Check value is preserved
		if ctx.Value("test") != "value" {
			t.Error("Built context should have test value")
		}

		t.Log("BuilderWithTimeout test passed")
	})

	t.Run("BuilderWithParent", func(t *testing.T) {
		parent := context.WithValue(context.Background(), "parent_key", "parent_value")

		builder := NewFakeContextBuilder()
		ctx := builder.
			WithParent(parent).
			WithValue("child_key", "child_value").
			Build()

		// Should have parent value
		if ctx.Value("parent_key") != "parent_value" {
			t.Error("Built context should inherit parent values")
		}

		// Should have child value
		if ctx.Value("child_key") != "child_value" {
			t.Error("Built context should have child values")
		}

		t.Log("BuilderWithParent test passed")
	})
}

// TestHelperFunctions tests the helper functions.
func TestHelperFunctions(t *testing.T) {
	t.Run("ContextWithTracing", func(t *testing.T) {
		ctx := ContextWithTracing("helper-trace", "helper-span")

		if ctx.Value("trace_id") != "helper-trace" {
			t.Error("Helper context should have trace ID")
		}

		if ctx.Value("span_id") != "helper-span" {
			t.Error("Helper context should have span ID")
		}

		t.Log("ContextWithTracing test passed")
	})

	t.Run("ContextWithTimeout", func(t *testing.T) {
		timeout := 50 * time.Millisecond
		ctx := ContextWithTimeout(timeout)

		deadline, ok := ctx.Deadline()
		if !ok {
			t.Error("Helper context should have deadline")
		}

		expectedDeadline := time.Now().Add(timeout)
		if deadline.Before(expectedDeadline.Add(-10*time.Millisecond)) || deadline.After(expectedDeadline.Add(10*time.Millisecond)) {
			t.Error("Deadline is not within expected range")
		}

		t.Log("ContextWithTimeout test passed")
	})

	t.Run("ContextWithValues", func(t *testing.T) {
		values := map[interface{}]interface{}{
			"key1": "value1",
			"key2": 42,
			"key3": true,
		}

		ctx := ContextWithValues(values)

		for k, v := range values {
			if ctx.Value(k) != v {
				t.Errorf("Helper context should have value for key %v", k)
			}
		}

		t.Log("ContextWithValues test passed")
	})
}
