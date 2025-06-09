package fakes

import (
	"errors"
	"sync"
)

// FakeEngine 為 engine.Interface 的測試假件實作。
// zh: 用於模擬 Engine 行為以供模組測試注入使用。
type FakeEngine struct {
	mu      sync.Mutex
	started bool
	stopped bool
	status  string
}

// NewFakeEngine 建立新的 Engine 假件。
// zh: 可設定初始狀態或注入特定行為模擬。
func NewFakeEngine() *FakeEngine {
	return &FakeEngine{
		status: "initialized",
	}
}

// Start 模擬啟動行為。
func (e *FakeEngine) Start() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.started {
		return errors.New("engine already started")
	}
	e.started = true
	e.status = "running"
	return nil
}

// Stop 模擬停止行為。
func (e *FakeEngine) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.stopped {
		return errors.New("engine already stopped")
	}
	e.stopped = true
	e.status = "stopped"
	return nil
}

// Status 回傳當前狀態。
func (e *FakeEngine) Status() string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.status
}
