package notifier_test

import (
	"testing"

	"github.com/detectviz/detectviz/internal/registry/notifier"
	"github.com/detectviz/detectviz/internal/test/mock/logger"
	"github.com/detectviz/detectviz/pkg/configtypes"
)

// TestNewNotifierRegistry 測試 notifier 註冊邏輯。
// zh: 檢查是否依據設定正確建立 notifier 清單。
func TestNewNotifierRegistry(t *testing.T) {
	cfgs := []configtypes.NotifierConfig{
		{Name: "email", Type: "email", Target: "noreply@example.com", Enable: true},
		{Name: "slack", Type: "slack", Target: "https://hooks.slack.com/xxx", Enable: true},
		{Name: "webhook", Type: "webhook", Target: "https://example.com/webhook", Enable: false}, // zh: 停用
		{Name: "invalid", Type: "foo", Target: "xxx", Enable: true},                              // zh: 無效類型
	}

	log := logger.NewTestLogger()
	notifiers := notifier.NewNotifierRegistry(cfgs, log)

	if len(notifiers) != 2 {
		t.Errorf("expected 2 notifiers, got %d", len(notifiers))
	}

	names := map[string]bool{}
	for _, n := range notifiers {
		names[n.Name()] = true
	}

	if !names["email"] {
		t.Error("email notifier not registered")
	}
	if !names["slack"] {
		t.Error("slack notifier not registered")
	}
	if names["webhook"] {
		t.Error("disabled notifier should not be registered")
	}
	if names["invalid"] {
		t.Error("invalid type should not be registered")
	}
}
