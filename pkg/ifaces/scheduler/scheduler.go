package scheduler

import "context"

// Job represents a task that can be scheduled and executed.
// zh: Job 代表一個可排程執行的任務。
type Job interface {
	// Run executes the job logic.
	// zh: 執行任務邏輯。
	Run(ctx context.Context) error

	// Name returns the job name.
	// zh: 回傳任務名稱（用於識別與日誌）。
	Name() string

	// Spec returns the cron-style schedule string for this job.
	// zh: 回傳此任務的排程時間字串（cron 格式）。
	Spec() string
}

// Scheduler defines the interface for a job scheduler.
// zh: Scheduler 定義排程器的操作介面。
type Scheduler interface {
	// Register registers a job with the scheduler.
	// zh: 註冊任務至排程器中。
	Register(job Job)

	// Start initiates the job scheduler.
	// zh: 啟動排程器。
	Start(ctx context.Context) error

	// Stop gracefully stops the scheduler and all running jobs.
	// zh: 停止排程器與所有正在執行的任務。
	Stop(ctx context.Context) error
}
