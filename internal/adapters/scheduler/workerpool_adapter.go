package scheduleradapter

import (
	"context"
	"sync"
	"time"

	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	"github.com/detectviz/detectviz/pkg/ifaces/scheduler"
)

// WorkerPoolScheduler executes jobs using a fixed pool of worker goroutines.
// zh: WorkerPoolScheduler 使用固定數量的 worker goroutine 執行排程任務。
type WorkerPoolScheduler struct {
	jobs       []scheduler.Job
	workerSize int
	started    bool
	mu         sync.Mutex
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	log        logger.Logger
}

// NewWorkerPoolScheduler creates a new scheduler with a given worker count and logger.
// zh: 建立一個新的 WorkerPoolScheduler，可指定 worker 數量與注入 logger。
func NewWorkerPoolScheduler(workerSize int, log logger.Logger) *WorkerPoolScheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPoolScheduler{
		workerSize: workerSize,
		ctx:        ctx,
		cancel:     cancel,
		log:        log,
	}
}

// Register adds a job to the queue.
// zh: 註冊任務至排程器中。
func (s *WorkerPoolScheduler) Register(job scheduler.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs = append(s.jobs, job)
}

// Start runs the worker pool and schedules all jobs in round-robin.
// zh: 啟動 worker pool 並以輪詢方式執行所有註冊任務。
func (s *WorkerPoolScheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		return nil
	}
	s.started = true
	s.mu.Unlock()

	for i := 0; i < s.workerSize; i++ {
		s.wg.Add(1)
		go s.worker()
	}
	return nil
}

// Stop cancels execution and waits for all workers to complete.
// zh: 停止排程器並等待所有 worker 結束。
func (s *WorkerPoolScheduler) Stop(ctx context.Context) error {
	s.cancel()
	s.wg.Wait()
	return nil
}

// worker is the function executed by each worker goroutine.
// zh: worker 是每個 worker goroutine 執行的函式。
func (s *WorkerPoolScheduler) worker() {
	defer s.wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.mu.Lock()
			jobs := make([]scheduler.Job, len(s.jobs))
			copy(jobs, s.jobs)
			s.mu.Unlock()
			for _, job := range jobs {
				go s.runWithRetry(job)
			}
		}
	}
}

// runWithRetry executes a job and retries up to 3 times if it fails.
// zh: 嘗試執行任務，若失敗則最多重試 3 次。
func (s *WorkerPoolScheduler) runWithRetry(job scheduler.Job) {
	const maxRetry = 3
	for i := 0; i < maxRetry; i++ {
		err := job.Run(s.ctx)
		if err == nil {
			return
		}
		s.log.Warn("job failed", "name", job.Name(), "attempt", i+1, "error", err)
		time.Sleep(1 * time.Second)
	}
	s.log.Error("job permanently failed after retries", "name", job.Name())
}
