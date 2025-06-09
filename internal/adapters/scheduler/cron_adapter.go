package scheduleradapter

import (
	"context"
	"sync"

	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	"github.com/detectviz/detectviz/pkg/ifaces/scheduler"
	"github.com/robfig/cron/v3"
)

// CronScheduler is an implementation of the Scheduler interface using robfig/cron.
// zh: CronScheduler 是使用 robfig/cron 套件實作的排程器。
type CronScheduler struct {
	c     *cron.Cron
	jobs  []scheduler.Job
	mu    sync.Mutex
	start sync.Once
	log   logger.Logger
}

// NewCronScheduler creates a new CronScheduler instance with logger support.
// zh: 建立一個新的 CronScheduler 實例，支援注入 logger。
func NewCronScheduler(log logger.Logger) *CronScheduler {
	return &CronScheduler{
		c:   cron.New(),
		log: log,
	}
}

// Register adds a job to the scheduler. The job must implement Spec() string to define its schedule.
// zh: 註冊一個任務到排程器，排程時間由 job.Spec() 提供。
func (s *CronScheduler) Register(job scheduler.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs = append(s.jobs, job)
}

// Start schedules all registered jobs using their Spec() value.
// zh: 根據每個註冊任務的 Spec() 結果進行排程並啟動排程器。
func (s *CronScheduler) Start(ctx context.Context) error {
	s.start.Do(func() {
		for _, j := range s.jobs {
			spec := j.Spec()
			_, err := s.c.AddFunc(spec, func() {
				if err := j.Run(ctx); err != nil && s.log != nil {
					s.log.Error("job run failed", "name", j.Name(), "error", err)
				}
			})
			if err != nil && s.log != nil {
				s.log.Error("failed to add cron job", "spec", spec, "name", j.Name(), "error", err)
			}
		}
		s.c.Start()
	})
	return nil
}

// Stop stops the cron engine gracefully.
// zh: 優雅地停止排程引擎。
func (s *CronScheduler) Stop(ctx context.Context) error {
	s.c.Stop()
	return nil
}

/*
範例：如何使用 CronScheduler 註冊任務

log := logger.NewZapLogger(...)
sched := NewCronScheduler(log)

job := &MyJob{} // 實作 scheduler.Job 介面，並實作 Spec() string 方法
sched.Register(job)

ctx := context.Background()
_ = sched.Start(ctx)
*/
