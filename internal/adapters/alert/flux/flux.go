// Package flux 提供基於 Flux 查詢語言的 AlertEvaluator 實作。
package fluxadapter

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/alert"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// FluxAlertEvaluator 是基於 Flux 查詢語言的告警評估器實作。
// zh: 透過 InfluxDB 的 Flux 語法進行指標查詢與告警條件比對。
type FluxAlertEvaluator struct {
	Log logger.Logger // zh: 日誌紀錄器，用於除錯與追蹤查詢過程
}

// NewEvaluator 建立一個 FluxAlertEvaluator。
// zh: 傳入 logger 實例以支援除錯與紀錄。
func NewEvaluator(log logger.Logger) alert.AlertEvaluator {
	return &FluxAlertEvaluator{
		Log: log.Named("flux"),
	}
}

// Evaluate 根據 AlertCondition 進行查詢與比對。
// zh: 目前尚未實作實際查詢與比對邏輯。
func (f *FluxAlertEvaluator) Evaluate(ctx context.Context, cond alert.AlertCondition) (alert.AlertResult, error) {
	f.Log.Debug("flux evaluator not implemented", "expr", cond.Expr)

	// TODO: 解析 cond.Expr，執行 Flux 查詢，解析回傳值進行閾值比對
	return alert.AlertResult{
		Firing:  false,
		Message: "flux evaluation not implemented",
		Value:   0,
	}, nil
}
