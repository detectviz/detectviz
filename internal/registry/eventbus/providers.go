package eventbus

import (
	"fmt"
	"sync"

	"github.com/detectviz/detectviz/pkg/ifaces/eventbus"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[string]eventbus.DispatcherProvider)
)

// RegisterProvider 註冊 provider，例如 in-memory / kafka / nats
func RegisterProvider(p eventbus.DispatcherProvider) {
	providersMu.Lock()
	defer providersMu.Unlock()
	providers[p.Name()] = p
}

// GetProvider 依據 provider 名稱取得 EventDispatcher
func GetProvider(name string) (eventbus.EventDispatcher, error) {
	providersMu.RLock()
	defer providersMu.RUnlock()
	p, ok := providers[name]
	if !ok {
		return nil, fmt.Errorf("eventbus provider '%s' not found", name)
	}
	return p.Build(), nil
}
