package manager_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/detectviz/detectviz/internal/plugins/manager"
)

// TestScanPlugins_EmptyDir 驗證掃描空目錄不會 panic。
// zh: 測試 ScanPlugins 是否能正常處理空目錄情境。
func TestScanPlugins_EmptyDir(t *testing.T) {
	dir := t.TempDir()

	err := manager.ScanPlugins(dir)
	if err != nil {
		t.Fatalf("ScanPlugins failed on empty dir: %v", err)
	}
}

// TestScanPlugins_InvalidFile 測試非 .so 檔案不應被載入。
// zh: 確保非 plugin 檔案 (.txt 等) 不會造成錯誤。
func TestScanPlugins_InvalidFile(t *testing.T) {
	dir := t.TempDir()
	dummyFile := filepath.Join(dir, "not_a_plugin.txt")

	if err := os.WriteFile(dummyFile, []byte("invalid"), 0644); err != nil {
		t.Fatalf("failed to write dummy file: %v", err)
	}

	err := manager.ScanPlugins(dir)
	if err != nil {
		t.Fatalf("ScanPlugins failed on invalid file: %v", err)
	}
}
