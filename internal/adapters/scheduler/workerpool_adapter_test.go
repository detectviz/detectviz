package scheduleradapter_test

import (
	"context"
	"testing"
	"time"

	scheduleradapter "github.com/detectviz/detectviz/internal/adapters/scheduler"
	"github.com/detectviz/detectviz/internal/test/testutil"
	ifacescheduler "github.com/detectviz/detectviz/pkg/ifaces/scheduler"
	"github.com/stretchr/testify/assert"
)

// mockJobWithRetry implements scheduler.Job and simulates a retryable task.
// zh: 模擬會失敗後成功的任務，驗證 retry 機制。
type mockJobWithRetry struct {
	name      string
	spec      string
	calls     int
	failUntil int
}

func (m *mockJobWithRetry) Run(ctx context.Context) error {
	m.calls++
	if m.calls <= m.failUntil {
		return assert.AnError
	}
	return nil
}

func (m *mockJobWithRetry) Name() string {
	return m.name
}

func (m *mockJobWithRetry) Spec() string {
	return m.spec
}

// TestWorkerPoolSchedulerIntegration 驗證 WorkerPoolScheduler 能夠依據排程與重試邏輯執行工作。
// zh: 測試 worker pool 型排程器是否能正確執行失敗後可重試的任務。
func TestWorkerPoolSchedulerIntegration(t *testing.T) {
	ctx := context.Background()
	job := &mockJobWithRetry{
		name:      "retry-job",
		spec:      "@every 5s",
		failUntil: 1,
	}

	log := testutil.NewTestLogger()
	s := scheduleradapter.NewWorkerPoolScheduler(1, log)

	var _ ifacescheduler.Scheduler = s // 型別符合驗證

	s.Register(job)

	err := s.Start(ctx)
	assert.NoError(t, err)

	time.Sleep(7 * time.Second) // 等待 worker 排程與重試

	err = s.Stop(ctx)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, job.calls, 2)
}
