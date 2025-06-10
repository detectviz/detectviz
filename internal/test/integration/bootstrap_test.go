package integration_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/detectviz/detectviz/internal/bootstrap"
)

// fakeComponent 用於測試注入的模組，紀錄啟動與停止日誌
type fakeComponent struct {
	started bool
	stopped bool
	log     *bytes.Buffer
}

func (f *fakeComponent) Run(ctx context.Context) error {
	f.log.WriteString("start\n")
	f.started = true
	<-ctx.Done()
	f.stopped = true
	f.log.WriteString("stop\n")
	return nil
}

func (f *fakeComponent) Shutdown(ctx context.Context) error {
	f.stopped = true
	f.log.WriteString("shutdown\n")
	return nil
}

func (f *fakeComponent) Name() string {
	return "fakeComponent"
}

// TestBootstrapFullFlow 測試完整 Init -> BuildServer -> Run 流程
// zh: 驗證使用 fakeComponent 注入模組，模組是否按序啟動與關閉
func TestBootstrapFullFlow(t *testing.T) {
	logBuf := &bytes.Buffer{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 初始化 bootstrap，注入 fakeComponent
	comp := &fakeComponent{log: logBuf}
	bs := bootstrap.NewBuilder()
	bs.InjectModule(comp)

	// Build Server，假設此階段會初始化 Engine 等核心服務
	srv, err := bs.BuildServer()
	if err != nil {
		t.Fatalf("BuildServer failed: %v", err)
	}

	// 以 goroutine 執行 Run，並在測試結束時取消 context 停止
	go func() {
		if err := srv.Run(ctx); err != nil {
			t.Errorf("Server Run failed: %v", err)
		}
	}()

	// 取消 context 以觸發 Shutdown
	cancel()

	// 驗證 fakeComponent 的啟動與停止狀態
	if !comp.started {
		t.Error("expected component to be started")
	}
	if !comp.stopped {
		t.Error("expected component to be stopped")
	}

	logStr := logBuf.String()
	expectedStart := "start"
	expectedStop := "stop"

	if !strings.Contains(logStr, expectedStart) {
		t.Errorf("expected log to contain %q", expectedStart)
	}
	if !strings.Contains(logStr, expectedStop) {
		t.Errorf("expected log to contain %q", expectedStop)
	}
}
