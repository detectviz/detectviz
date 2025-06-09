// Package prom 提供基於 Prometheus 查詢語言的 AlertEvaluator 實作。
package promadapter

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/alert"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// PromAlertEvaluator 是基於 PromQL 的告警評估器實作。
// zh: 使用 Prometheus 查詢語法進行評估與比對。
type PromAlertEvaluator struct {
	Log logger.Logger // zh: 日誌記錄器，用於觀察查詢與比對過程
}

// NewEvaluator 建立一個 PromAlertEvaluator。
// zh: 可傳入 logger 實例以觀察觸發情況。
func NewEvaluator(log logger.Logger) alert.AlertEvaluator {
	return &PromAlertEvaluator{
		Log: log.Named("prom"),
	}
}

// Evaluate 根據 AlertCondition 執行查詢與閾值比對。
// zh: 實際查詢與閾值判斷尚未實作。
func (p *PromAlertEvaluator) Evaluate(ctx context.Context, cond alert.AlertCondition) (alert.AlertResult, error) {
	p.Log.Debug("prom evaluator not implemented", "expr", cond.Expr)

	// TODO: 實作 Prometheus 查詢、解析結果、進行閾值比對
	return alert.AlertResult{
		Firing:  false,
		Message: "prom evaluation not implemented",
		Value:   0,
	}, nil
}
