package alert

import (
	"fmt"
	"sync"

	"github.com/detectviz/detectviz/internal/adapters/alert/flux"
	"github.com/detectviz/detectviz/internal/adapters/alert/prom"
	"github.com/detectviz/detectviz/pkg/ifaces/alert"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// AlertEvaluatorRegistry 管理所有註冊的 AlertEvaluator 實例
// zh: 負責註冊與查詢可用的告警評估器實作。
type AlertEvaluatorRegistry struct {
	evaluators map[string]alert.AlertEvaluator
	mu         sync.RWMutex
}

// NewAlertEvaluatorRegistry 建立註冊中心實例
// zh: 初始化並回傳告警評估器註冊中心。
func NewAlertEvaluatorRegistry() *AlertEvaluatorRegistry {
	return &AlertEvaluatorRegistry{
		evaluators: make(map[string]alert.AlertEvaluator),
	}
}

// NewDefaultAlertEvaluatorRegistry 建立包含預設實作的註冊中心
// zh: 注入 PromEvaluator 與 FluxEvaluator 供預設使用。
func NewDefaultAlertEvaluatorRegistry(log logger.Logger) *AlertEvaluatorRegistry {
	r := NewAlertEvaluatorRegistry()
	r.Register("prometheus", prom.NewEvaluator(log))
	r.Register("flux", flux.NewEvaluator(log))
	return r
}

// Register 註冊一個新的 AlertEvaluator
// zh: 以指定名稱註冊一個告警評估器實作。
func (r *AlertEvaluatorRegistry) Register(name string, evaluator alert.AlertEvaluator) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.evaluators[name] = evaluator
}

// Get 根據名稱取得已註冊的 AlertEvaluator
// zh: 回傳對應名稱的告警評估器，若不存在則回傳錯誤。
func (r *AlertEvaluatorRegistry) Get(name string) (alert.AlertEvaluator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	evaluator, ok := r.evaluators[name]
	if !ok {
		return nil, fmt.Errorf("alert evaluator '%s' not found", name)
	}
	return evaluator, nil
}
