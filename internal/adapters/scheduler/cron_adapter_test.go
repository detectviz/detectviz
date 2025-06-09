package scheduleradapter_test

import (
	"context"
	"testing"
	"time"

	scheduleradapter "github.com/detectviz/detectviz/internal/adapters/scheduler"
	"github.com/detectviz/detectviz/internal/test/testutil"
	"github.com/stretchr/testify/assert"
)

// mockJob implements scheduler.Job interface.
// zh: 用於測試 CronScheduler 的模擬任務。
type mockJob struct {
	name  string
	spec  string
	calls int
}

func (m *mockJob) Run(ctx context.Context) error {
	m.calls++
	return nil
}

func (m *mockJob) Name() string {
	return m.name
}

func (m *mockJob) Spec() string {
	return m.spec
}

// TestCronScheduler_Run 驗證 CronScheduler 能夠依據時間規則週期執行註冊任務。
// zh: 測試 Cron 型排程器是否能依指定的頻率成功執行任務。
func TestCronScheduler_Run(t *testing.T) {
	ctx := context.Background()
	job := &mockJob{name: "cron-job", spec: "@every 2s"}

	logger := testutil.NewTestLogger()
	sched := scheduleradapter.NewCronScheduler(logger)
	sched.Register(job)

	err := sched.Start(ctx)
	assert.NoError(t, err)

	time.Sleep(5 * time.Second) // 等待任務至少執行兩次

	err = sched.Stop(ctx)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, job.calls, 2)
}
