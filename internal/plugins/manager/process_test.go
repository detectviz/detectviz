package manager_test

import (
	"testing"

	"github.com/detectviz/detectviz/internal/plugins/manager"
)

// TestProcessManager_StartAndStop 測試啟動與停止 plugin 後端程序。
// zh: 驗證 ProcessManager 是否可成功啟動與停止 plugin backend。
func TestProcessManager_StartAndStop(t *testing.T) {
	pm := manager.NewProcessManager()

	err := pm.Start("test_plugin", "/bin/sleep", "2")
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	list := pm.List()
	if len(list) != 1 || list[0] != "test_plugin" {
		t.Errorf("expected 'test_plugin' in list, got %v", list)
	}

	err = pm.Stop("test_plugin")
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	list = pm.List()
	if len(list) != 0 {
		t.Errorf("expected no active plugins, got %v", list)
	}
}

// TestProcessManager_DoubleStart 測試重複啟動同一 plugin 不會報錯。
// zh: 已啟動的 plugin 不應再次啟動，且不得報錯。
func TestProcessManager_DoubleStart(t *testing.T) {
	pm := manager.NewProcessManager()

	err := pm.Start("dup_plugin", "/bin/sleep", "1")
	if err != nil {
		t.Fatalf("first Start failed: %v", err)
	}
	defer pm.Stop("dup_plugin")

	err = pm.Start("dup_plugin", "/bin/sleep", "1")
	if err != nil {
		t.Errorf("second Start should not error, got: %v", err)
	}
}
