package registry_test

import (
	"testing"

	"github.com/detectviz/detectviz/pkg/registry/apis/plugin"
)

type testPlugin struct {
	name    string
	version string
}

func (p *testPlugin) Name() string    { return p.name }
func (p *testPlugin) Version() string { return p.version }
func (p *testPlugin) Init() error     { return nil }
func (p *testPlugin) Close() error    { return nil }

// TestPluginRegistry_RegisterAndGet 測試 Plugin 註冊與查詢邏輯。
// zh: 驗證 Register() 與 Get() 是否正確註冊並可查詢 Plugin。
func TestPluginRegistry_RegisterAndGet(t *testing.T) {
	r := plugin.New()
	p := &testPlugin{name: "test", version: "v1.0.0"}

	if err := r.Register(p); err != nil {
		t.Fatalf("unexpected error on register: %v", err)
	}

	got, ok := r.Get("test")
	if !ok {
		t.Fatal("expected plugin to be found")
	}
	if got.Name() != "test" {
		t.Errorf("expected name 'test', got %s", got.Name())
	}
	if got.Version() != "v1.0.0" {
		t.Errorf("expected version 'v1.0.0', got %s", got.Version())
	}
}

// TestPluginRegistry_Duplicate 測試重複註冊 Plugin 是否錯誤。
// zh: 檢查 Plugin 名稱重複時是否回傳錯誤。
func TestPluginRegistry_Duplicate(t *testing.T) {
	r := plugin.New()
	p1 := &testPlugin{name: "dup", version: "v1.0.0"}
	p2 := &testPlugin{name: "dup", version: "v1.1.0"}

	if err := r.Register(p1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := r.Register(p2); err == nil {
		t.Error("expected error for duplicate registration, got nil")
	}
}
