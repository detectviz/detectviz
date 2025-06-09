package manager_test

import (
	"testing"

	"github.com/detectviz/detectviz/internal/plugins/manager"
)

type mockPlugin struct {
	name    string
	version string
	inited  bool
	closed  bool
}

func (p *mockPlugin) Name() string    { return p.name }
func (p *mockPlugin) Version() string { return p.version }
func (p *mockPlugin) Init() error {
	p.inited = true
	return nil
}
func (p *mockPlugin) Close() error {
	p.closed = true
	return nil
}

// TestPluginLifecycleManager_Flow 測試完整插件註冊與生命周期流程。
// zh: 驗證註冊、InitAll 與 ShutdownAll 是否正確觸發 Plugin 方法。
func TestPluginLifecycleManager_Flow(t *testing.T) {
	mgr := manager.NewLifecycleManager()

	p1 := &mockPlugin{name: "alpha", version: "v1"}
	p2 := &mockPlugin{name: "beta", version: "v2"}

	_ = mgr.Register(p1)
	_ = mgr.Register(p2)

	mgr.InitAll()

	if !p1.inited || !p2.inited {
		t.Errorf("expected both plugins to be initialized")
	}

	mgr.ShutdownAll()

	if !p1.closed || !p2.closed {
		t.Errorf("expected both plugins to be closed")
	}
}
