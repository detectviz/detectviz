package scheduler

import (
	"context"
	"fmt"

	"github.com/detectviz/detectviz/pkg/ifaces/scheduler"
)

// MockJob is a simple mock implementation of a scheduled job.
// zh: MockJob 是模擬用的任務實作，用於測試排程器功能。
type MockJob struct {
	ID   string
	Logs *[]string // zh: 用於收集執行記錄
}

// Run appends a message to Logs when executed.
// zh: Run 在執行時將訊息寫入 Logs。
func (j *MockJob) Run(ctx context.Context) error {
	if j.Logs != nil {
		*j.Logs = append(*j.Logs, fmt.Sprintf("job %s executed", j.ID))
	}
	return nil
}

// Name returns the job ID.
// zh: 回傳任務名稱。
func (j *MockJob) Name() string {
	return j.ID
}

// MockScheduler is a no-op implementation for testing.
// zh: MockScheduler 是模擬用排程器實作，不會實際執行任務。
type MockScheduler struct {
	Jobs []scheduler.Job
}

// NewMockScheduler returns a new instance.
// zh: 建立一個新的 MockScheduler 實例。
func NewMockScheduler() *MockScheduler {
	return &MockScheduler{}
}

// Register adds a job to internal list.
// zh: 註冊任務至模擬排程器。
func (s *MockScheduler) Register(job scheduler.Job) {
	s.Jobs = append(s.Jobs, job)
}

// Start does nothing.
// zh: 模擬啟動，實際不執行任何任務。
func (s *MockScheduler) Start(ctx context.Context) error {
	return nil
}

// Stop does nothing.
// zh: 模擬停止，實際不執行任何操作。
func (s *MockScheduler) Stop(ctx context.Context) error {
	return nil
}
