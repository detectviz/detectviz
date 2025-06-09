package plugins

import (
	"log"
	"os"
	"path/filepath"
	"plugin"

	"github.com/detectviz/detectviz/pkg/ifaces/plugins"
	pluginregistry "github.com/detectviz/detectviz/pkg/registry/apis/plugin"
)

var registry = pluginregistry.New()

// LoadPlugins 掃描 plugins 目錄並載入所有符合條件的插件。
// zh: 掃描 ./plugins 目錄並註冊每個動態插件。
func LoadPlugins(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".so" {
			continue
		}

		pPath := filepath.Join(dir, f.Name())
		p, err := plugin.Open(pPath)
		if err != nil {
			log.Printf("failed to open plugin: %v", err)
			continue
		}

		sym, err := p.Lookup("Plugin")
		if err != nil {
			log.Printf("failed to find 'Plugin' symbol: %v", err)
			continue
		}

		plug, ok := sym.(plugins.Plugin)
		if !ok {
			log.Printf("invalid plugin type: %s", f.Name())
			continue
		}

		if err := plug.Init(); err != nil {
			log.Printf("plugin init error [%s]: %v", plug.Name(), err)
			continue
		}

		if err := registry.Register(plug); err != nil {
			log.Printf("plugin register error: %v", err)
			continue
		}

		log.Printf("loaded plugin: %s (%s)", plug.Name(), plug.Version())
	}

	return nil
}

// Registry 回傳全域 plugin 註冊表。
// zh: 提供供外部查詢與擴充的 plugin registry。
func Registry() *pluginregistry.Registry {
	return registry
}
