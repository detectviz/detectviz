package modules_test

import (
	"context"
	"sync"
	"testing"
	"time"

	core "github.com/detectviz/detectviz/internal/modules"
)

type fakeModule struct {
	name    string
	started bool
	stopped bool
	mu      sync.Mutex
}

// Run 模擬模組執行，啟動後會阻塞，直到 ctx 結束。
// zh: 啟動模組並模擬長時間運行，直到 context 取消。
func (f *fakeModule) Run(ctx context.Context) error {
	f.mu.Lock()
	f.started = true
	f.mu.Unlock()

	select {
	case <-ctx.Done():
		f.mu.Lock()
		f.stopped = true
		f.mu.Unlock()
		return nil
	case <-time.After(5 * time.Second):
		return nil
	}
}

// Shutdown 模擬模組關閉。
// zh: 優雅停止模組，標記已停止狀態。
func (f *fakeModule) Shutdown(ctx context.Context) error {
	f.mu.Lock()
	f.stopped = true
	f.mu.Unlock()
	return nil
}

// Name 回傳模組名稱。
// zh: 取得模組識別名稱。
func (f *fakeModule) Name() string {
	return f.name
}

// TestEngineRegistryRunnerIntegration 測試 Engine、Registry、Runner 的整合流程。
// zh: 模擬模組註冊、啟動、關閉，驗證整合流程正確。
func TestEngineRegistryRunnerIntegration(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	engine := core.NewEngine()
	registry := core.NewRegistry()
	graph := core.NewDependencyGraph()
	runner := core.NewRunner(engine, registry, graph)

	module := &fakeModule{name: "testModule"}
	err := registry.Register(module.Name(), module)
	if err != nil {
		t.Fatalf("failed to register module: %v", err)
	}

	engine.Register(module)

	go func() {
		if err := runner.Start(ctx); err != nil {
			t.Errorf("Runner run failed: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	module.mu.Lock()
	started := module.started
	module.mu.Unlock()
	if !started {
		t.Errorf("expected module to be started")
	}

	cancel()

	time.Sleep(100 * time.Millisecond)

	module.mu.Lock()
	stopped := module.stopped
	module.mu.Unlock()
	if !stopped {
		t.Errorf("expected module to be stopped")
	}
}
