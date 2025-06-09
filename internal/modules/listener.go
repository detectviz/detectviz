package modules

import (
	"context"
	"log"
	"time"
)

// HealthCheckable represents modules that can report health status.
// zh: 具備健康狀態回報能力的模組介面。
type HealthCheckable interface {
	Healthy() bool
}

// Listener monitors the health of modules and triggers shutdown if needed.
// zh: Listener 監控所有模組健康狀態，必要時觸發全域停機。
type Listener struct {
	engine   *Engine
	registry *Registry
	interval time.Duration
	cancel   context.CancelFunc
}

// NewListener creates a new health monitoring listener.
// zh: 建立健康狀態監控器。
func NewListener(engine *Engine, registry *Registry, interval time.Duration) *Listener {
	return &Listener{
		engine:   engine,
		registry: registry,
		interval: interval,
	}
}

// Start launches the periodic health check loop.
// zh: 啟動定期健康檢查迴圈。
func (l *Listener) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	l.cancel = cancel

	go func() {
		ticker := time.NewTicker(l.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if !l.checkAll() {
					log.Println("[listener] unhealthy module detected, shutting down all modules")
					_ = l.engine.ShutdownAll(ctx)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Stop stops the health listener.
// zh: 停止健康監聽器。
func (l *Listener) Stop() {
	if l.cancel != nil {
		l.cancel()
	}
}

// checkAll verifies the health status of all registered modules that support HealthCheckable.
// zh: 檢查所有具備健康檢查功能的模組是否健康。
func (l *Listener) checkAll() bool {
	for _, name := range l.registry.List() {
		m, ok := l.registry.Get(name)
		if !ok {
			continue
		}
		hc, ok := m.(HealthCheckable)
		if ok && !hc.Healthy() {
			log.Printf("[listener] module %q is unhealthy", name)
			return false
		}
	}
	return true
}
