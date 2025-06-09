package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/detectviz/detectviz/internal/test/fakes"
)

// TestServer_RunAndShutdown tests the basic lifecycle of a Server implementation.
// zh: 測試 Server 實作是否能正確執行 Run 和 Shutdown。
func TestServer_RunAndShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := &fakes.FakeServer{}

	go func() {
		_ = server.Run(ctx)
	}()

	// 模擬短時間後觸發關閉
	time.Sleep(10 * time.Millisecond)
	_ = server.Shutdown(context.Background())

	if !server.RunCalled {
		t.Errorf("Run was not called")
	}
	if !server.ShutdownCalled {
		t.Errorf("Shutdown was not called")
	}
}
