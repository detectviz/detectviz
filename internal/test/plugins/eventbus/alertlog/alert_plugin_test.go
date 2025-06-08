package alertlog_test

import (
	"context"
	"strings"
	"testing"

	// 匯入 plugin，觸發 init 註冊處理器
	_ "github.com/detectviz/detectviz/internal/plugins/eventbus/alertlog"

	"github.com/detectviz/detectviz/internal/registry/eventbus"
	testlogger "github.com/detectviz/detectviz/internal/test/testutil"
	"github.com/detectviz/detectviz/pkg/ifaces/event"
)

func TestAlertPluginHandler_Registration(t *testing.T) {
	// 使用可驗證輸出的測試 logger
	log := testlogger.NewTestLogger()

	// 替換 plugin handler 的 logger（需 plugin 內部支援）
	eventbus.OverrideDefaultLogger(log)

	dispatcher, err := eventbus.NewEventDispatcher("in-memory")
	if err != nil {
		t.Fatalf("failed to create event dispatcher: %v", err)
	}

	err = dispatcher.DispatchAlertTriggered(context.Background(), event.AlertTriggeredEvent{
		AlertID:    "test-alert",
		RuleName:   "test-rule",
		Level:      "warning",
		Instance:   "test-instance",
		Metric:     "test-metric",
		Comparison: ">",
		Value:      100,
		Threshold:  100,
		Message:    "plugin alert test",
	})
	if err != nil {
		t.Errorf("unexpected dispatch error: %v", err)
	}

	// 驗證 logger 是否捕捉到 alert 處理訊息
	found := false
	for _, line := range log.Messages() {
		if strings.Contains(line, "plugin alert test") {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected alert message not found in log")
	}
}
