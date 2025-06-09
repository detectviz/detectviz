package manager

import (
	"log"
	"os"
	"path/filepath"
	"plugin"

	iface "github.com/detectviz/detectviz/pkg/ifaces/plugins"
)

var registry = NewManagerRegistry()

// ScanPlugins 掃描 plugins 目錄，動態載入符合 Plugin interface 的模組。
// zh: 掃描本地 plugins 資料夾，載入並初始化所有 plugin 實體。
func ScanPlugins(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".so" {
			continue
		}

		path := filepath.Join(dir, f.Name())
		p, err := plugin.Open(path)
		if err != nil {
			log.Printf("plugin open error: %v", err)
			continue
		}

		sym, err := p.Lookup("Plugin")
		if err != nil {
			log.Printf("plugin lookup error: %v", err)
			continue
		}

		instance, ok := sym.(iface.Plugin)
		if !ok {
			log.Printf("invalid plugin type: %s", f.Name())
			continue
		}

		if err := instance.Init(); err != nil {
			log.Printf("plugin init failed: %v", err)
			continue
		}

		if err := registry.Register(instance); err != nil {
			log.Printf("plugin register failed: %v", err)
			continue
		}

		log.Printf("plugin loaded: %s (%s)", instance.Name(), instance.Version())
	}

	return nil
}

// Registry 回傳 plugin 管理中心實體。
// zh: 提供給外部使用的 plugin 管理註冊表。
func Registry() *ManagerRegistry {
	return registry
}
