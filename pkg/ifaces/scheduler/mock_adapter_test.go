package scheduler_test

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/pkg/ifaces/scheduler"
	"github.com/stretchr/testify/assert"
)

// mockJob implements the scheduler.Job interface for testing.
// zh: mockJob 是模擬用的排程任務，實作 scheduler.Job 介面。
type mockJob struct {
	name  string // zh: 任務名稱
	spec  string // zh: 排程時間（cron 格式）
	calls int    // zh: 被執行次數
}

// Run increments the call counter.
// zh: Run 被呼叫時會遞增 calls 計數。
func (m *mockJob) Run(ctx context.Context) error {
	m.calls++
	return nil
}

// Name returns the job name.
// zh: 回傳任務名稱。
func (m *mockJob) Name() string {
	return m.name
}

// Spec returns the job's schedule.
// zh: 回傳排程規則。
func (m *mockJob) Spec() string {
	return m.spec
}

// mockScheduler is a lightweight mock implementation of Scheduler.
// zh: mockScheduler 是模擬用的 Scheduler 實作，不會實際執行任務。
type mockScheduler struct {
	Jobs []scheduler.Job // Registered jobs. zh: 已註冊的任務清單
}

// Register adds a job to the scheduler.
// zh: 註冊一個任務。
func (m *mockScheduler) Register(job scheduler.Job) {
	m.Jobs = append(m.Jobs, job)
}

// Start is a no-op.
// zh: 啟動排程器（模擬，無實際行為）。
func (m *mockScheduler) Start(ctx context.Context) error {
	return nil
}

// Stop is a no-op.
// zh: 停止排程器（模擬，無實際行為）。
func (m *mockScheduler) Stop(ctx context.Context) error {
	return nil
}

// TestMockSchedulerIntegration tests the registration and lifecycle flow of the mock scheduler.
// zh: 測試 mockScheduler 的任務註冊與啟停流程。
func TestMockSchedulerIntegration(t *testing.T) {
	job := &mockJob{name: "test-job", spec: "@every 1m"}
	mockSched := &mockScheduler{}

	mockSched.Register(job)

	assert.Len(t, mockSched.Jobs, 1)
	assert.Equal(t, "test-job", mockSched.Jobs[0].Name())

	err := mockSched.Start(context.Background())
	assert.NoError(t, err)

	err = mockSched.Stop(context.Background())
	assert.NoError(t, err)
}
