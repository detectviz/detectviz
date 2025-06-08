package scheduler_test

import (
	"context"
	"testing"
	"time"

	adapters "github.com/detectviz/detectviz/internal/adapters/scheduler"
	"github.com/detectviz/detectviz/internal/adapters/scheduler/testlogger"
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

func TestWorkerPoolSchedulerIntegration(t *testing.T) {
	ctx := context.Background()
	job := &mockJobWithRetry{
		name:      "retry-job",
		spec:      "@every 5s",
		failUntil: 1,
	}

	log := testlogger.New()
	s := adapters.NewWorkerPoolScheduler(1, log)

	var _ ifacescheduler.Scheduler = s // 型別符合驗證

	s.Register(job)

	err := s.Start(ctx)
	assert.NoError(t, err)

	time.Sleep(7 * time.Second) // 等待 worker 排程與重試

	err = s.Stop(ctx)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, job.calls, 2)
}
