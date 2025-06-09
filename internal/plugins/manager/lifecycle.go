package manager

import (
	"log"

	"github.com/detectviz/detectviz/pkg/ifaces/plugins"
)

// PluginLifecycleManager 控制 plugin 的初始化與關閉流程。
// zh: 管理插件的整體生命週期，包括註冊、啟動與關閉。
type PluginLifecycleManager struct {
	registry *ManagerRegistry
	process  *ProcessManager
}

// NewLifecycleManager 建立 PluginLifecycleManager。
// zh: 初始化註冊表與後端管理器。
func NewLifecycleManager() *PluginLifecycleManager {
	return &PluginLifecycleManager{
		registry: NewManagerRegistry(),
		process:  NewProcessManager(),
	}
}

// Register 註冊 plugin。
// zh: 加入註冊表以供查詢與管理。
func (m *PluginLifecycleManager) Register(p plugins.Plugin) error {
	if err := m.registry.Register(p); err != nil {
		return err
	}
	log.Printf("plugin registered: %s (%s)", p.Name(), p.Version())
	return nil
}

// InitAll 初始化所有註冊 plugin。
// zh: 執行 Init() 以完成啟動邏輯。
func (m *PluginLifecycleManager) InitAll() {
	for _, p := range m.registry.List() {
		if err := p.Init(); err != nil {
			log.Printf("plugin init failed: %s - %v", p.Name(), err)
			continue
		}
		log.Printf("plugin initialized: %s", p.Name())
	}
}

// ShutdownAll 關閉所有插件。
// zh: 執行 Close() 並中止後端程序。
func (m *PluginLifecycleManager) ShutdownAll() {
	for _, p := range m.registry.List() {
		if err := p.Close(); err != nil {
			log.Printf("plugin close failed: %s - %v", p.Name(), err)
		} else {
			log.Printf("plugin closed: %s", p.Name())
		}
		_ = m.process.Stop(p.Name()) // 後端若有啟動則同步關閉
	}
}
