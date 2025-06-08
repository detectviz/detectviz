package scheduler_test

import (
	"context"
	"testing"
	"time"

	adapters "github.com/detectviz/detectviz/internal/adapters/scheduler"
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

func TestCronScheduler_Run(t *testing.T) {
	ctx := context.Background()
	job := &mockJob{name: "cron-job", spec: "@every 2s"}

	logger := &testLogger{}
	sched := adapters.NewCronScheduler(logger)
	sched.Register(job)

	err := sched.Start(ctx)
	assert.NoError(t, err)

	time.Sleep(5 * time.Second) // 至少執行兩次

	err = sched.Stop(ctx)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, job.calls, 2)
}
