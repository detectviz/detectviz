package manager_test

import (
	"testing"

	"github.com/detectviz/detectviz/internal/plugins/manager"
)

type testPlugin struct {
	name    string
	version string
}

func (p *testPlugin) Name() string    { return p.name }
func (p *testPlugin) Version() string { return p.version }
func (p *testPlugin) Init() error     { return nil }
func (p *testPlugin) Close() error    { return nil }

// TestManagerRegistry_RegisterAndGet 測試插件註冊與查詢。
// zh: 驗證 ManagerRegistry 是否能正確儲存並回傳 Plugin。
func TestManagerRegistry_RegisterAndGet(t *testing.T) {
	r := manager.NewManagerRegistry()
	p := &testPlugin{name: "modA", version: "v1.0.0"}

	if err := r.Register(p); err != nil {
		t.Fatalf("unexpected register error: %v", err)
	}

	got, ok := r.Get("modA")
	if !ok {
		t.Fatal("expected plugin to be found")
	}
	if got.Name() != "modA" {
		t.Errorf("expected plugin name 'modA', got %s", got.Name())
	}
	if got.Version() != "v1.0.0" {
		t.Errorf("expected version 'v1.0.0', got %s", got.Version())
	}
}

// TestManagerRegistry_Duplicate 測試註冊重複名稱插件時回傳錯誤。
// zh: 插件名稱不可重複，否則應報錯。
func TestManagerRegistry_Duplicate(t *testing.T) {
	r := manager.NewManagerRegistry()
	_ = r.Register(&testPlugin{name: "dup", version: "v1.0.0"})

	err := r.Register(&testPlugin{name: "dup", version: "v1.1.0"})
	if err == nil {
		t.Error("expected duplicate registration error, got nil")
	}
}
