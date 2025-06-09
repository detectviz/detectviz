package manager

import (
	"context"
	"log"
	"os/exec"
	"sync"
)

// PluginProcess 管理單一 plugin 的後端進程。
// zh: 用於啟動與停止外部 plugin backend 的封裝結構。
type PluginProcess struct {
	Name string
	Cmd  *exec.Cmd
}

// ProcessManager 負責管理所有 plugin 的後端進程。
// zh: 可啟動、追蹤並終止所有已註冊插件的執行實體。
type ProcessManager struct {
	mu        sync.Mutex
	processes map[string]*PluginProcess
}

// NewProcessManager 建立新的 ProcessManager。
func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		processes: make(map[string]*PluginProcess),
	}
}

// Start 啟動指定 plugin 名稱的後端程序。
// zh: 可執行指定二進位檔或 shell 指令。
func (pm *ProcessManager) Start(name string, command string, args ...string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.processes[name]; exists {
		return nil // 已啟動
	}

	cmd := exec.CommandContext(context.Background(), command, args...)
	if err := cmd.Start(); err != nil {
		return err
	}

	log.Printf("plugin backend started: %s (pid=%d)", name, cmd.Process.Pid)

	pm.processes[name] = &PluginProcess{
		Name: name,
		Cmd:  cmd,
	}
	return nil
}

// Stop 停止指定 plugin 的後端程序。
// zh: 優雅中止 plugin 程式。
func (pm *ProcessManager) Stop(name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	p, ok := pm.processes[name]
	if !ok {
		return nil
	}

	if err := p.Cmd.Process.Kill(); err != nil {
		return err
	}
	delete(pm.processes, name)
	log.Printf("plugin backend stopped: %s", name)
	return nil
}

// List 回傳目前執行中的 plugin 名稱。
// zh: 回傳目前啟動中的 plugin 清單。
func (pm *ProcessManager) List() []string {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	var names []string
	for name := range pm.processes {
		names = append(names, name)
	}
	return names
}
