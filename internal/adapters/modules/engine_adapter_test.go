package modulesadapter_test

import (
	"context"
	"errors"
	"testing"

	modulesadapter "github.com/detectviz/detectviz/internal/adapters/modules"
	core "github.com/detectviz/detectviz/internal/modules"
)

// fakeEngine 是模擬 core.Engine interface 的假件，方便測試。
// zh: 提供模擬核心 Engine 行為以測試 Adapter。
type fakeEngine struct {
	registerCalled    bool
	runAllCalled      bool
	shutdownAllCalled bool
	runAllErr         error
	shutdownAllErr    error
}

func (f *fakeEngine) Register(m core.Module) {
	f.registerCalled = true
}

func (f *fakeEngine) RunAll(ctx context.Context) error {
	f.runAllCalled = true
	return f.runAllErr
}

func (f *fakeEngine) ShutdownAll(ctx context.Context) error {
	f.shutdownAllCalled = true
	return f.shutdownAllErr
}

func TestEngineAdapter_Register(t *testing.T) {
	fake := &fakeEngine{}
	adapter := modulesadapter.NewEngineAdapter(fake)

	module := &fakeLifecycleModule{}

	adapter.Register(module)

	if !fake.registerCalled {
		t.Error("expected Register to call core.Engine.Register")
	}
}

func TestEngineAdapter_RunAll(t *testing.T) {
	fake := &fakeEngine{}
	adapter := modulesadapter.NewEngineAdapter(fake)

	ctx := context.Background()

	fake.runAllErr = nil
	err := adapter.RunAll(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !fake.runAllCalled {
		t.Error("expected RunAll to call core.Engine.RunAll")
	}

	fake.runAllErr = errors.New("run all failed")
	err = adapter.RunAll(ctx)
	if err == nil {
		t.Error("expected error from RunAll, got nil")
	}
}

func TestEngineAdapter_ShutdownAll(t *testing.T) {
	fake := &fakeEngine{}
	adapter := modulesadapter.NewEngineAdapter(fake)

	ctx := context.Background()

	fake.shutdownAllErr = nil
	err := adapter.ShutdownAll(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !fake.shutdownAllCalled {
		t.Error("expected ShutdownAll to call core.Engine.ShutdownAll")
	}

	fake.shutdownAllErr = errors.New("shutdown all failed")
	err = adapter.ShutdownAll(ctx)
	if err == nil {
		t.Error("expected error from ShutdownAll, got nil")
	}
}

// fakeLifecycleModule 是模擬 core.Module interface。
// zh: 提供測試用的生命週期模組假件。
type fakeLifecycleModule struct{}

func (f *fakeLifecycleModule) Run(ctx context.Context) error {
	return nil
}

func (f *fakeLifecycleModule) Shutdown(ctx context.Context) error {
	return nil
}
